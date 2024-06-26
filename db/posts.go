// Package db handles all core database interactions of the server
package db

import (
	"bytes"
	"database/sql"
	"encoding/json"

	"github.com/Masterminds/squirrel"
	"github.com/bakape/meguca/common"
)

var insertPostStmt *sql.Stmt

// Post is for writing new posts to a database. It contains the Password
// field, which is never exposed publically through Post.
type Post struct {
	common.StandalonePost
	Password []byte
	IP       string
}

func prepareInsertPostStmt() (err error) {

	insertPostStmt, err = sqlDB.Prepare(`
		WITH inserted_post AS (
		INSERT INTO posts (editing, board, op, body, flag, name, trip, auth, sage, PASSWORD, ip)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			RETURNING
				id, time, moderated)
			SELECT
				ip.id,
				ip.time,
				ip.moderated,
				CASE WHEN COUNT(pm.post_id) = 0 THEN
					'[]'::json
				ELSE
					json_agg(json_build_object('type', pm.type, 'length', pm.length, 'by', pm.by, 'data', pm.data))
				END AS moderations
			FROM
				inserted_post ip
			LEFT JOIN post_moderation pm ON ip.id = pm.post_id
		GROUP BY
			ip.id,
			ip.time,
			ip.moderated;
		`)
	return
}

func selectPost(id uint64, columns ...string) rowScanner {
	return sq.Select(columns...).
		From("posts").
		Where("id = ?", id).
		QueryRow()
}

// GetPostParenthood retrieves the board and OP of a post
func GetPostParenthood(id uint64) (board string, op uint64, err error) {
	err = selectPost(id, "board", "op").Scan(&board, &op)
	return
}

// GetPostBoard retrieves the board of a post by ID
func GetPostBoard(id uint64) (board string, err error) {
	err = selectPost(id, "board").Scan(&board)
	return
}

func getCounter(q squirrel.SelectBuilder) (uint64, error) {
	var c sql.NullInt64
	err := q.QueryRow().Scan(&c)
	return uint64(c.Int64), err
}

// BoardCounter retrieves the progress counter of a board
func BoardCounter(board string) (uint64, error) {
	q := sq.Select("max(update_time) + count(*)").
		From("threads").
		Where("board = ?", board)
	return getCounter(q)
}

// AllBoardCounter retrieves the progress counter of the /all/ board
func AllBoardCounter() (uint64, error) {
	q := sq.Select("max(update_time) + count(*)").
		From("threads")
	return getCounter(q)
}

// WritePost writes a post struct to the database. Only used in tests and
// migrations.
func WritePost(tx *sql.Tx, p Post) (err error) {
	// Don't store empty strings of these in the database. Zero value != NULL.
	var (
		img, ip *string
		imgName string
		spoiler bool
	)
	if p.IP != "" {
		ip = &p.IP
	}
	if p.Image != nil {
		img = &p.Image.SHA1
		imgName = p.Image.Name
		spoiler = p.Image.Spoiler
	}

	_, err = sq.Insert("posts").
		Columns(
			"editing", "spoiler", "id", "board", "op", "time", "body", "flag",
			"name", "trip", "auth", "password", "ip",
			"SHA1", "imageName",
			"commands",
		).
		Values(
			p.Editing, spoiler, p.ID, p.Board, p.OP, p.Time, p.Body, p.Flag,
			p.Name, p.Trip, p.Auth, p.Password, ip,
			img, imgName,
			commandRow(p.Commands),
		).
		RunWith(tx).
		Exec()
	if err != nil {
		return
	}
	err = writeLinks(tx, p.ID, p.Links)
	if err != nil {
		return
	}

	if p.Editing {
		err = SetOpenBody(p.ID, []byte(p.Body))
	}
	return
}

// Insert Post into thread and set its ID and creation time and moderation
// status.
// Thread OPs must have their post ID set to the thread ID.
// Any images are to be inserted in a separate call.
func InsertPost(tx *sql.Tx, p *Post) (err error) {
	if p.ID != 0 { // OP of a thread
		args := make([]interface{}, 0, 12)
		args = append(args,
			p.Editing, p.Board, p.OP, p.Body, p.Flag,
			p.Name, p.Trip, p.Auth, p.Sage,
			p.Password, p.IP)

		q := sq.Insert("posts").
			Columns(
				"editing", "board", "op", "body", "flag",
				"name", "trip", "auth", "sage",
				"password", "ip",
			)

		q = q.Columns("id")
		args = append(args, p.ID)
		err = q.
			Values(args...).
			Suffix("returning id, time, moderated").
			RunWith(tx).
			QueryRow().
			Scan(&p.ID, &p.Time, &p.Moderated)
		if err != nil {
			return
		}
	} else {
		var moderationData []byte
		err = tx.Stmt(insertPostStmt).QueryRow(
			p.Editing, p.Board, p.OP, p.Body, p.Flag,
			p.Name, p.Trip, p.Auth, p.Sage,
			p.Password, p.IP,
		).Scan(&p.ID, &p.Time, &p.Moderated, &moderationData)
		if err != nil {
			return
		}
		if bytes.Equal(moderationData, []byte("[]")) {
			return
		}
		p.Moderation = []common.ModerationEntry{}
		err = json.Unmarshal(moderationData, &p.Moderation)
		return
	}

	if p.Moderated {
		// Read moderation log, if post deleted on insert
		//
		// TODO: Get this in db-side JSON in same query, once we have db-side
		// post JSON generation.
		arr := [...]*common.Post{&p.Post}
		err = injectModeration(arr[:], tx)
		if err != nil {
			return
		}

	}
	return
}
func InsertRegularPost(p *Post) (err error) {
	var moderationData []byte
	err = insertPostStmt.QueryRow(
		p.Editing, p.Board, p.OP, p.Body, p.Flag,
		p.Name, p.Trip, p.Auth, p.Sage,
		p.Password, p.IP,
	).Scan(&p.ID, &p.Time, &p.Moderated, &moderationData)
	if err != nil {
		return
	}
	if bytes.Equal(moderationData, []byte("[]")) {
		return
	}
	p.Moderation = []common.ModerationEntry{}
	err = json.Unmarshal(moderationData, &p.Moderation)
	return
}

// GetPostPassword retrieves a post's modification password
func GetPostPassword(id uint64) (p []byte, err error) {
	err = sq.Select("password").From("posts").Where("id = ?", id).Scan(&p)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

// SetPostCounter sets the post counter.
// Should only be used in tests.
func SetPostCounter(c uint64) error {
	_, err := sqlDB.Exec(`SELECT setval('post_id', $1)`, c)
	return err
}
