// Package imager handles image, video, etc. upload requests and processing
package imager

import (
	"bytes"
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"gopkg.in/vansante/go-ffprobe.v2"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bakape/meguca/auth"
	"github.com/bakape/meguca/common"
	"github.com/bakape/meguca/config"
	"github.com/bakape/meguca/db"
	"github.com/bakape/meguca/util"
	"github.com/bakape/thumbnailer/v2"
	"github.com/go-playground/log"
)

// Minimal capacity of large buffers in the pool
const largeBufCap = 12 << 10

var (
	// Map of MIME types to the constants used internally
	mimeTypes = map[string]uint8{
		"image/jpeg":                    common.JPEG,
		"image/png":                     common.PNG,
		"image/gif":                     common.GIF,
		"image/webp":                    common.WEBP,
		"image/avif":                    common.AVIF,
		mimePDF:                         common.PDF,
		"video/webm":                    common.WEBM,
		"application/ogg":               common.OGG,
		"video/mp4":                     common.MP4,
		"video/quicktime":               common.MP4,
		"audio/mpeg":                    common.MP3,
		mime7Zip:                        common.SevenZip,
		mimeTarGZ:                       common.TGZ,
		mimeTarXZ:                       common.TXZ,
		mimeZip:                         common.ZIP,
		"audio/x-flac":                  common.FLAC,
		mimeText:                        common.TXT,
		"application/x-rar-compressed":  common.RAR,
		"application/vnd.comicbook+zip": common.CBZ,
		"application/vnd.comicbook-rar": common.CBR,
	}

	// MIME types from thumbnailer to accept
	allowedMimeTypes map[string]bool

	errTooLarge = errors.New("file too large")

	// Large buffer pool of length=0 capacity=12+KB
	largeBufPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 0, largeBufCap)
		},
	}
)

func init() {
	allowedMimeTypes = make(map[string]bool, len(mimeTypes))
	for t := range mimeTypes {
		allowedMimeTypes[t] = true
	}
}

// Return large buffer pool, if eligable
func returnLargeBuf(buf []byte) {
	if cap(buf) >= largeBufCap {
		largeBufPool.Put(buf[:0])
	}
}

// NewImageUpload  handles the clients' image (or other file) upload request
func NewImageUpload(w http.ResponseWriter, r *http.Request) {
	var id string
	err := func() (err error) {
		bypass, err := validateUploader(w, r)
		if err != nil {
			return
		}

		// Limit data received to the maximum uploaded file size limit
		r.Body = http.MaxBytesReader(w, r.Body, int64(config.Get().MaxSize<<20))

		id, err = ParseUpload(r)
		switch err {
		case nil:
			if !bypass {
				err = incrementSpamScore(w, r)
			}
			return
		case io.EOF:
			return common.StatusError{
				Err:  err,
				Code: 400,
			}
		default:
			return
		}
	}()
	if err != nil {
		LogError(w, r, err)
	}

	w.Write([]byte(id))
}

// Apply security restrictions to uploader
func validateUploader(
	w http.ResponseWriter,
	r *http.Request,
) (bypass bool, err error) {
	if s := r.Header.Get("Authorization"); s != "" &&
		s == "Bearer "+config.Get().Salt {
		// Internal upload bypass
		bypass = true
		return
	}

	ip, err := auth.GetIP(r)
	if err != nil {
		return
	}
	_, err = db.IsBanned("all", ip)
	if err != nil {
		return
	}

	var session auth.Base64Token
	err = session.EnsureCookie(w, r)
	if err != nil {
		return
	}
	need, err := db.NeedCaptcha(session, ip)
	if err != nil {
		return
	}
	if need {
		err = common.StatusError{errors.New("captcha required"), 403}
		return
	}

	return
}

// UploadImageHash attempts to skip image upload, if the file has already
// been thumbnailed and is stored on the server. The client sends an SHA1 hash
// of the file it wants to upload. The server looks up, if such a file is
// thumbnailed. If yes, generates and sends a new image allocation token to
// the client.
func UploadImageHash(w http.ResponseWriter, r *http.Request) {
	token, err := func() (token string, err error) {
		bypass, err := validateUploader(w, r)
		if err != nil {
			return
		}

		buf, err := ioutil.ReadAll(http.MaxBytesReader(w, r.Body, 40))
		if err != nil {
			return
		}
		sha1 := string(buf)

		err = db.InTransaction(false, func(tx *sql.Tx) (err error) {
			exists, err := db.ImageExists(tx, sha1)
			if err != nil {
				return
			}
			if exists {
				token, err = db.NewImageToken(tx, sha1)
			}
			return
		})
		if err != nil {
			return
		}
		if !bypass {
			err = incrementSpamScore(w, r)
		}
		return
	}()
	if err != nil {
		LogError(w, r, err)
	} else if token != "" {
		w.Write([]byte(token))
	}
}
func UploadMeguHash(w http.ResponseWriter, r *http.Request) {
	token, filename, err := func() (token string, filename string, err error) {
		bypass, err := validateUploader(w, r)
		if err != nil {
			return
		}

		buf, err := ioutil.ReadAll(http.MaxBytesReader(w, r.Body, 40))
		if err != nil {
			return
		}
		sha1 := string(buf)

		err = db.InTransaction(false, func(tx *sql.Tx) (err error) {
			filename, err = db.GetImageFilename(sha1)
			if err != nil {
				return
			}
			token, err = db.NewImageToken(tx, sha1)
			return
		})
		if err != nil {
			return
		}
		if !bypass {
			err = incrementSpamScore(w, r)
		}
		return
	}()
	if err != nil {
		LogError(w, r, err)
	} else if token != "" {
		response := struct {
			Token    string `json:"token"`
			Filename string `json:"name"`
		}{token, filename}
		responseBytes, err := json.Marshal(&response)
		if err != nil {
			LogError(w, r, err)
		} else {
			w.Write(responseBytes)
		}
	}
}

func incrementSpamScore(w http.ResponseWriter, r *http.Request) (err error) {
	ip, err := auth.GetIP(r)
	if err != nil {
		return
	}
	var session auth.Base64Token
	err = session.EnsureCookie(w, r)
	if err != nil {
		return
	}
	db.IncrementSpamScore(session, ip, config.Get().ImageScore)
	return
}

// LogError send the client file upload errors and logs them server-side
func LogError(w http.ResponseWriter, r *http.Request, err error) {
	code := 500
	if err, ok := err.(common.StatusError); ok {
		code = err.Code
	}
	http.Error(w, err.Error(), code)

	if common.IsTest || common.CanIgnoreClientError(err) {
		return
	}
	ip, ipErr := auth.GetIP(r)
	if ipErr != nil {
		ip = "invalid IP"
	}
	log.Errorf("upload error: by %s: %s: %#v", ip, err, err)
}

// ParseUpload parses the upload form. Separate function for cleaner error
// handling and reusability.
// Returns the HTTP status code of the response, the ID of the generated image
// and an error, if any.
func ParseUpload(req *http.Request) (string, error) {
	max := config.Get().MaxSize << 20
	length, err := strconv.ParseUint(req.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return "", common.StatusError{err, 413}
	}
	if uint(length) > max {
		return "", common.StatusError{errTooLarge, 400}
	}
	err = req.ParseMultipartForm(0)
	if err != nil {
		return "", common.StatusError{err, 400}
	}

	file, head, err := req.FormFile("image")
	if err != nil {
		return "", common.StatusError{err, 400}
	}
	defer file.Close()
	if uint(head.Size) > max {
		return "", common.StatusError{errTooLarge, 413}
	}

	res := <-requestThumbnailing(file, head.Filename, int(head.Size), nil)
	return res.imageID, res.err
}

// Create a new thumbnail, commit its resources to the DB and filesystem, and
// pass the image data to the client.
func newThumbnail(f multipart.File, filename string, SHA1 string, tiktokName *string) (token string, err error) {
	var img common.ImageCommon
	img.SHA1 = SHA1

	conf := config.Get()
	thumb, err := processFile(f, filename, &img, tiktokName, thumbnailer.Options{
		MaxSourceDims: thumbnailer.Dims{
			Width:  uint(conf.MaxWidth),
			Height: uint(conf.MaxHeight),
		},
		ThumbDims: thumbnailer.Dims{
			Width:  150,
			Height: 150,
		},
		AcceptedMimeTypes: allowedMimeTypes,
	})
	defer returnLargeBuf(thumb)
	if err != nil {
		switch err.(type) {
		case thumbnailer.ErrUnsupportedMIME, thumbnailer.ErrInvalidImage:
			err = common.StatusError{err, 400}
		}
		return
	}

	// Being done in one transaction prevents the image DB record from getting
	// garbage-collected between the calls
	err = db.InTransaction(false, func(tx *sql.Tx) (err error) {
		var thumbR io.ReadSeeker
		if thumb != nil {
			thumbR = bytes.NewReader(thumb)
		}
		err = db.AllocateImage(tx, f, thumbR, img)
		if err != nil && !db.IsConflictError(err) {
			return
		}
		token, err = db.NewImageToken(tx, img.SHA1)
		return
	})
	return
}

func getVideoCodec(file io.Reader) (string, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	data, err := ffprobe.ProbeReader(ctx, file)
	if err != nil {
		return "", err
	}
	if len(data.Streams) == 0 {
		return "", errors.New("No streams found")
	}
	stream := data.FirstVideoStream()
	if stream != nil {
		return stream.CodecName, nil
	}
	return "", errors.New("no video stream found")
}

func getAudioCodec(file io.Reader) (string, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	data, err := ffprobe.ProbeReader(ctx, file)
	if err != nil {
		return "", err
	}
	if len(data.Streams) == 0 {
		return "", errors.New("No streams found")
	}
	stream := data.FirstAudioStream()
	if stream != nil {
		return stream.CodecName, nil
	}
	return "", errors.New("no audio stream found")
}

// Separate function for easier testability
func processFile(f multipart.File, filename string, img *common.ImageCommon, tiktokName *string, opts thumbnailer.Options) (thumb []byte, err error) {
	jpegThumb := config.Get().JPEGThumbnails

	resultCh := make(chan string, 1)
	errCh := make(chan error, 1)

	if tiktokName == nil {
		go func() {
			username, usernameErr := getTiktokUsername(filename)
			resultCh <- username
			errCh <- usernameErr
		}()
	}

	src, thumbImage, err := thumbnailer.Process(f, opts)

	defer func() {
		// Add image internal buffer to pool
		if thumbImage == nil {
			return
		}
		// Only image type used in thumbnailer by default
		img, ok := thumbImage.(*image.RGBA)
		if ok {
			returnLargeBuf(img.Pix)
		}
	}()
	switch err {
	case nil:
		if jpegThumb {
			img.ThumbType = common.JPEG
		} else {
			img.ThumbType = common.WEBP
		}
	case thumbnailer.ErrCantThumbnail:
		err = nil
		img.ThumbType = common.NoFile
	default:
		return
	}

	img.FileType = mimeTypes[src.Mime]

	img.Audio = src.HasAudio
	img.Video = src.HasVideo
	img.Length = uint32(src.Length / time.Second)
	f.Seek(0, 0)

	// check if src.mime starts with "image"
	isImage := strings.HasPrefix(src.Mime, "image")

	//Certain common codecs are returned by the thumbnailer
	//If one of these codecs are returned, we skip the call to ffprobe
	codecs := []string{"h264", "hevc", "mjpeg", "gif", "png"}
	codecSet := false
	for _, codec := range codecs {
		if src.Codec == codec {
			img.Codec = src.Codec
			codecSet = true
			break
		}
	}

	if !codecSet {
		if img.Video || isImage {
			img.Codec, err = getVideoCodec(f)
		} else if img.Audio {
			img.Codec, err = getAudioCodec(f)
		}
	}

	if isImage && img.Codec == "mjpeg" {
		img.Codec = "jpeg"
	}
	if isImage && img.Codec == "av1" {
		img.Codec = "avif"
	}
	// Some media has retardedly long meta strings. Just truncate them, instead
	// of rejecting. Must ensure it's still valid unicode after trancation,
	// incase a rune was split.
	img.Artist = src.Artist
	img.Title = src.Title
	util.TrimString(&img.Artist, 100)
	util.TrimString(&img.Title, 200)
	//Detect tiktok @ if Artist tag isn't present
	if tiktokName == nil {
		if src.Artist == "" {
			tiktokUsername := <-resultCh
			err := <-errCh
			if err == nil {
				img.Artist = "@" + tiktokUsername
			}
		}
	} else {
		img.Artist = "@" + *tiktokName
	}

	img.Dims = [4]uint16{uint16(src.Width), uint16(src.Height), 0, 0}
	if thumbImage != nil {
		b := thumbImage.Bounds()
		img.Dims[2] = uint16(b.Dx())
		img.Dims[3] = uint16(b.Dy())
	}

	img.MD5, img.Size, err = hashFile(f, md5.New(),
		base64.RawURLEncoding.EncodeToString)
	if err != nil {
		return
	}

	if thumbImage != nil {
		if jpegThumb {
			w := bytes.NewBuffer(largeBufPool.Get().([]byte))
			err = jpeg.Encode(w, thumbImage, &jpeg.Options{
				Quality: 90,
			})
			thumb = w.Bytes()
		} else {
			thumb, err = EncodeWebP(thumbImage, 90)
		}
		if err != nil {
			return
		}
	}

	return
}
