package imager

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"hash"
	"io"
	"mime/multipart"
	"runtime"
	"sync"

	"github.com/bakape/meguca/db"
)

var (
	scheduleJob      = make(chan jobRequest, 128)
	scheduleSmallJob = make(chan jobRequest, 128)

	// Pool of temp buffers used for hashing
	buf512Pool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 512)
		},
	}
)

type jobRequest struct {
	file       multipart.File
	filename   string
	size       int
	res        chan<- thumbnailingResponse
	tiktokName *string
}

type thumbnailingResponse struct {
	imageID string
	err     error
}

// Queues upload processing to prevent resource overuse
func requestThumbnailing(file multipart.File, filename string, size int, tiktokName *string) <-chan thumbnailingResponse {
	// 2 separate queues - one for small and one for bigger files.
	// Allows for some degree of concurrent thumbnailing without exhausting
	// server resources.
	ch := make(chan thumbnailingResponse)
	req := jobRequest{file, filename, size, ch, tiktokName}
	if size <= 4<<20 {
		scheduleSmallJob <- req
	} else {
		scheduleJob <- req
	}
	return ch
}

// Queue thumbnailing jobs to reduce resource contention and prevent OOM
func init() {
	for _, ch := range [...]<-chan jobRequest{scheduleJob, scheduleSmallJob} {
		go func(queue <-chan jobRequest) {
			runtime.LockOSThread()
			for {
				req := <-queue
				id, err := processRequest(req.file, req.filename, req.size, req.tiktokName)
				req.res <- thumbnailingResponse{id, err}
			}
		}(ch)
	}
}

// Hash file to string
func hashFile(rs io.ReadSeeker, h hash.Hash, encode func([]byte) string,
) (
	hash string, read int, err error,
) {

	_, err = rs.Seek(0, 0)
	if err != nil {
		return
	}
	buf := buf512Pool.Get().([]byte)
	defer buf512Pool.Put(buf)

	for {
		buf = buf[:512] // Reset slicing

		var n int
		n, err = rs.Read(buf)
		buf = buf[:n]
		read += n
		switch err {
		case nil:
			h.Write(buf)
		case io.EOF:
			err = nil
			hash = encode(h.Sum(buf))
			return
		default:
			return
		}
	}
}

func processRequest(file multipart.File, filename string, size int, tiktokName *string) (token string, err error) {
	SHA1, _, err := hashFile(file, sha1.New(), hex.EncodeToString)
	if err != nil {
		return
	}
	var exists bool
	err = db.InTransaction(false, func(tx *sql.Tx) (err error) {
		exists, err = db.ImageExists(tx, SHA1)
		if err != nil {
			return
		}
		if exists { // Already have a thumbnail
			token, err = db.NewImageToken(tx, SHA1)
		}
		return
	})
	if err != nil {
		return
	}
	if !exists {
		token, err = newThumbnail(file, filename, SHA1, tiktokName)
	}
	return
}
