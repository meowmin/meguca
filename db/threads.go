package db

import (
	"context"

	"github.com/bakape/pg_util"
	"github.com/jackc/pgx/v4"
)

type ThreadInsertParams struct {
	Subject string

	// Must include between 1 and 3 tags
	Tags []string

	PostInsertParamsCommon `db:"-"`
}

// Insert thread and empty post into DB and return the post ID.
//
// authKey is optional, in case this is a migration of a legacy post
func InsertThread(ctx context.Context, p ThreadInsertParams) (
	id uint64, err error,
) {
	err = InTransaction(context.Background(), func(tx pgx.Tx) (err error) {
		q, args := pg_util.BuildInsert(pg_util.InsertOpts{
			Table:  "threads",
			Data:   p,
			Suffix: "returning id",
		})
		err = tx.QueryRow(ctx, q, args...).Scan(&id)
		if err != nil {
			return
		}

		_, err = InsertPost(ctx, tx, OPInsertparams{
			ID: id,
			ReplyInsertParams: ReplyInsertParams{
				Thread:                 id,
				PostInsertParamsCommon: p.PostInsertParamsCommon,
			},
		})
		return
	})
	return
}

// Check, if thread exists in the database
func ThreadExists(ctx context.Context, id uint64) (exists bool, err error) {
	err = db.
		QueryRow(
			ctx,
			`select exists (
				select
				from threads
				where id = $1
			)`,
			id,
		).
		Scan(&exists)
	return
}

// Read feed data for initializing Pulsar as JSON
func GetFeedData() (buf []byte, err error) {
	err = db.
		QueryRow(
			context.Background(),
			`select
				jsonb_agg(
					jsonb_build_object(
						'thread', t.id,
						'recent_posts', coalesce(r.val, '{}'::jsonb),
						'open_posts', coalesce(o.val, '{}'::jsonb)
					)
				)
			from threads t
			left join (
				select
					r.thread,
					jsonb_object_agg(
						r.id,
						to_unix(r.created_on)
					) val
				from posts r
				where r.created_on > now() - interval '16 minutes'
				group by r.thread
			) r on r.thread = t.id
			left join (
				select
					o.thread,
					jsonb_object_agg(
						o.id,
						jsonb_build_object(
							'has_image', o.image is not null,
							'image_spoilered', o.image_spoilered,
							'created_on', to_unix(o.created_on),
							'thread', o.thread,
							'body', o.body
						)
					) val
				from posts o
				where o.open
				group by o.thread
			) o on o.thread = t.id`,
		).
		Scan(&buf)
	if err != nil {
		return
	}
	if len(buf) == 0 {
		buf = []byte(`[]`)
	}
	return
}
