// Code generated by sqlc. DO NOT EDIT.
// source: trip_messages.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const getTripForMessage = `-- name: GetTripForMessage :one
SELECT trip_id
FROM trip_messages
WHERE message_id = $1
LIMIT 1
`

func (q *Queries) GetTripForMessage(ctx context.Context, messageID int64) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getTripForMessage, messageID)
	var trip_id uuid.UUID
	err := row.Scan(&trip_id)
	return trip_id, err
}

const setMessageForTrip = `-- name: SetMessageForTrip :exec
INSERT INTO trip_messages (trip_id, message_id)
VALUES ($1, $2)
`

type SetMessageForTripParams struct {
	TripID    uuid.UUID
	MessageID int64
}

func (q *Queries) SetMessageForTrip(ctx context.Context, arg SetMessageForTripParams) error {
	_, err := q.db.ExecContext(ctx, setMessageForTrip, arg.TripID, arg.MessageID)
	return err
}
