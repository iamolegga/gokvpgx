package gokvpgx

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/philippgille/gokv/encoding"
)

const defaultTableName = "Item"

// Options are the options for the PostgreSQL client.
type Options struct {
	Pool      *pgxpool.Pool
	TableName string
	Codec     encoding.Codec
}

// DefaultOptions is an Options object with default values.
var DefaultOptions = Options{
	TableName: defaultTableName,
	Codec:     encoding.JSON,
}
