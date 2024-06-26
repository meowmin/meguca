package db

import (
	"database/sql"
	"io"
	"time"

	"github.com/bakape/meguca/auth"
	"github.com/bakape/meguca/common"
	"github.com/bakape/meguca/imager/assets"
	"github.com/bakape/meguca/util"
	"github.com/lib/pq"
)

const (
	// Time it takes for an image allocation token to expire
	tokenTimeout = time.Minute
)

var (
	// ErrInvalidToken occurs, when trying to retrieve an image with an
	// non-existent token. The token might have expired (60 to 119 seconds) or
	// the client could have provided an invalid token to begin with.
	ErrInvalidToken = common.ErrInvalidInput("invalid image token")
	insertImageStmt *sql.Stmt
)

// Video structure
type Video struct {
	FileType uint8         `json:"file_type"`
	Duration time.Duration `json:"-"`
	SHA1     string        `json:"sha1"`
}

func prepareInsertImageStmt() (err error) {
	insertImageStmt, err = sqlDB.Prepare(`
        select insert_image($1::bigint,
                             $2::char(86),
                             $3::varchar(200),
                             $4::bool)
    `)
	return
}

// WriteImage writes a processed image record to the DB. Only used in tests.
func WriteImage(i common.ImageCommon) error {
	return InTransaction(false, func(tx *sql.Tx) error {
		return writeImageTx(tx, i)
	})
}

func writeImageTx(tx *sql.Tx, i common.ImageCommon) (err error) {
	var codec interface{} = i.Codec
	if codec == "" {
		codec = nil
	}
	_, err = sq.
		Insert("images").
		Columns(
			"audio", "video", "file_type", "thumb_type", "dims", "length",
			"size", "MD5", "SHA1", "Title", "Artist", "Codec",
		).
		Values(
			i.Audio, i.Video, int(i.FileType), int(i.ThumbType),
			pq.GenericArray{A: i.Dims}, i.Length, i.Size, i.MD5, i.SHA1,
			i.Title, i.Artist, codec,
		).
		RunWith(tx).
		Exec()
	return
}

// NewImageToken inserts a new image allocation token into the DB and returns
// its ID
func NewImageToken(tx *sql.Tx, SHA1 string) (token string, err error) {
	expires := time.Now().Add(tokenTimeout).UTC()

	// Loop in case there is a primary key collision
	for {
		token, err = auth.RandomID(64)
		if err != nil {
			return
		}

		_, err = sq.
			Insert("image_tokens").
			Columns("token", "SHA1", "expires").
			Values(token, SHA1, expires).
			RunWith(tx).
			Exec()
		switch {
		case err == nil:
			return
		case IsConflictError(err):
			continue
		default:
			return
		}
	}
}

// ImageExists returns, if image exists
func ImageExists(tx *sql.Tx, sha1 string) (exists bool, err error) {
	err = sq.Select("1").
		From("images").
		Where("sha1 = ?", sha1).
		Scan(&exists)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}
func GetImageFilename(sha1 string) (exists string, err error) {
	err = sq.Select("imagename").
		From("posts").
		InnerJoin("images on posts.sha1 = images.sha1").
		Where("images.sha1 = ?", sha1).
		OrderBy("posts.id").
		Limit(1).
		QueryRow().
		Scan(&exists)
	return
}

// ImageVisible returns if the image is attached to any non-deleted and unspoilered posts on the board
func ImageVisible(sha1, board string) (visible bool, err error) {
	err = sq.Select("1").
		From("posts").
		Where("sha1 = ? and board = ? and not is_deleted(id) and not spoiler", sha1, board).
		Scan(&visible)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

// AllocateImage allocates an image's file resources to their respective served
// directories and write its data to the database
func AllocateImage(tx *sql.Tx, src, thumb io.ReadSeeker, img common.ImageCommon,
) (
	err error,
) {
	err = writeImageTx(tx, img)
	if err != nil {
		return err
	}

	err = assets.Write(img.SHA1, img.FileType, img.ThumbType, src, thumb)
	if err != nil {
		return cleanUpFailedAllocation(img, err)
	}
	return nil
}

// Delete any dangling image files in case of a failed image allocation
func cleanUpFailedAllocation(img common.ImageCommon, err error) error {
	delErr := assets.Delete(img.SHA1, img.FileType, img.ThumbType)
	if delErr != nil {
		err = util.WrapError(err.Error(), delErr)
	}
	return err
}

// HasImage returns, if the post has an image allocated. Only used in tests.
func HasImage(id uint64) (has bool, err error) {
	err = sq.Select("true").
		From("posts").
		Where("id = ? and SHA1 IS NOT NULL", id).
		QueryRow().
		Scan(&has)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

func BumpThread(tx *sql.Tx, thread uint64) error {
	_, err := tx.Exec(`SELECT bump_thread($1, $2)`, thread, false)
	return err
}

// InsertImage insert an image into an existing open post and return image
// JSON
func InsertImage(tx *sql.Tx, postID uint64, token, name string, spoiler bool,
) (
	json []byte, err error,
) {
	stmt := tx.Stmt(insertImageStmt)
	err = stmt.QueryRow(postID, token, name, spoiler).Scan(&json)
	if extractException(err) == "invalid image token" {
		err = ErrInvalidToken
	}
	return
}

// GetImage retrieves a thumbnailed image record from the DB.
//
// Only used in tests.
func GetImage(sha1 string) (img common.ImageCommon, err error) {
	var scanner imageScanner
	err = sq.Select("*").
		From("images").
		Where("SHA1 = ?", sha1).
		QueryRow().
		Scan(scanner.ScanArgs()...)
	if err != nil {
		return
	}
	return scanner.Val().ImageCommon, nil
}

func GetImageByPost(id uint64) (img common.ImageCommon, err error) {
	// Define a scanner for the image
	var scanner imageScanner

	// Execute the query with join, selecting all columns from the images table
	err = sq.Select("images.*").
		From("images").
		Join("posts ON posts.sha1 = images.sha1").
		Where("posts.id = ?", id).
		QueryRow().
		Scan(scanner.ScanArgs()...)

	if err != nil {
		return
	}

	return scanner.Val().ImageCommon, nil
}

// SpoilerImage spoilers an already allocated image
func SpoilerImage(id, op uint64) error {
	_, err := sq.Update("posts").
		Set("spoiler", true).
		Where("id = ?", id).
		Exec()
	return err
}

// Try to transfer image from one post to another. Return image, if anything was
// transferred
func TransferImage(fromPost, toPost, thread uint64) (
	transferred *common.Image,
	err error,
) {
	err = InTransaction(false, func(tx *sql.Tx) (err error) {
		var scanner imageScanner
		err = tx.
			QueryRow(
				`select p.imageName, p.spoiler, i.*
				from posts p
				join images i on i.sha1 = p.sha1
				join threads t on t.id = p.op
				where
					p.id = $1
					and p.op = $2
					and p.id != p.op
				for update of p`,
				fromPost,
				thread,
			).
			Scan(
				append(
					[]interface{}{
						&scanner.Name,
						&scanner.Spoiler,
					},
					scanner.ScanArgs()...,
				)...,
			)
		if err != nil {
			return
		}
		transferred = scanner.Val()

		res, err := tx.Exec(
			`update posts
			set
				sha1 = $3,
				imageName = $4,
				spoiler = $5
			where
				id = $1
				and op = $2`,
			toPost,
			thread,

			transferred.SHA1,
			transferred.Name,
			transferred.Spoiler,
		)
		if err != nil {
			return
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return
		}
		if affected == 0 {
			return sql.ErrNoRows
		}

		_, err = tx.Exec(
			`update posts
			set
				imageName = '',
				spoiler = false,
				sha1 = null
			where
				id = $1
				and op = $2`,
			fromPost,
			thread,
		)
		return
	})
	if err == sql.ErrNoRows {
		err = nil
		transferred = nil
	}
	return
}

// VideoPlaylist returns a video playlist for a board
func VideoPlaylist(board string) (videos []Video, err error) {

	// Prepare the query
	query := "SELECT sha1, file_type, length FROM get_megu_playlist($1)"

	// Execute the query
	rows, err := sqlDB.Query(query, board)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process the result set
	for rows.Next() {
		var v Video
		var dur int64

		err := rows.Scan(&v.SHA1, &v.FileType, &dur)
		if err != nil {
			return nil, err
		}

		v.Duration = time.Duration(dur) * time.Second
		videos = append(videos, v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return
}

// Delete images not used in any posts
func deleteUnusedImages() (err error) {
	r, err := sqlDB.Query(`select * from cleanup_images()`)
	if err != nil {
		return
	}
	defer r.Close()

	for r.Next() {
		var (
			sha1                string
			fileType, thumbType uint8
		)
		err = r.Scan(&sha1, &fileType, &thumbType)
		if err != nil {
			return
		}
		err = assets.Delete(sha1, fileType, thumbType)
		if err != nil {
			return
		}
	}

	return r.Err()
}
