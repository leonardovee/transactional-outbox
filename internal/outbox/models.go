// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package outbox

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Outbox struct {
	ID            string
	AggregateID   string
	AggregateType string
	Type          string
	Payload       []byte
	CreatedAt     pgtype.Timestamp
}
