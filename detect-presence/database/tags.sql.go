// Code generated by sqlc. DO NOT EDIT.
// source: tags.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const listTags = `-- name: ListTags :many
SELECT tag AS name, COUNT(trip_id) AS trip_count
FROM trip_taggings
GROUP BY tag
ORDER BY COUNT(trip_id) DESC
LIMIT $1
`

type ListTagsRow struct {
	Name      string
	TripCount int64
}

func (q *Queries) ListTags(ctx context.Context, limit int32) ([]ListTagsRow, error) {
	rows, err := q.db.QueryContext(ctx, listTags, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListTagsRow
	for rows.Next() {
		var i ListTagsRow
		if err := rows.Scan(&i.Name, &i.TripCount); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const tagTrip = `-- name: TagTrip :exec
INSERT INTO trip_taggings (trip_id, tag)
VALUES ($1, $2)
`

type TagTripParams struct {
	TripID uuid.UUID
	Tag    string
}

func (q *Queries) TagTrip(ctx context.Context, arg TagTripParams) error {
	_, err := q.db.ExecContext(ctx, tagTrip, arg.TripID, arg.Tag)
	return err
}

const untagTrip = `-- name: UntagTrip :exec
DELETE
FROM trip_taggings
WHERE trip_id = $1
  AND tag = $2
`

type UntagTripParams struct {
	TripID uuid.UUID
	Tag    string
}

func (q *Queries) UntagTrip(ctx context.Context, arg UntagTripParams) error {
	_, err := q.db.ExecContext(ctx, untagTrip, arg.TripID, arg.Tag)
	return err
}
