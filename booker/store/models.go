// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package store

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Book struct {
	ID        int32
	Title     string
	Author    string
	Genre     pgtype.Text
	Price     pgtype.Numeric
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}
