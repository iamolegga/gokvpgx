package gokvpgx

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/philippgille/gokv/encoding"
	"github.com/philippgille/gokv/util"
)

// Client is a gokv.Store implementation for PostgreSQL using pgx.
type Client struct {
	Pool      *pgxpool.Pool
	Codec     encoding.Codec
	TableName string

	// pgx prepares statements for us, so we don't need to do it manually
	// https://github.com/jackc/pgx/issues/791#issuecomment-660486444
	upsertStmt string
	getStmt    string
	deleteStmt string
}

// NewClient creates a new PostgreSQL client using pgx.
//
// You must call the Close() method on the client when you're done working with it.
func NewClient(options Options) (*Client, error) {
	if options.Pool == nil {
		return nil, errors.New("pgxpool.Pool cannot be nil")
	}
	if options.TableName == "" {
		options.TableName = DefaultOptions.TableName
	}
	if options.Codec == nil {
		options.Codec = DefaultOptions.Codec
	}

	client := &Client{
		Pool:       options.Pool,
		Codec:      options.Codec,
		TableName:  options.TableName,
		upsertStmt: fmt.Sprintf("INSERT INTO %s (k, v) VALUES ($1, $2) ON CONFLICT (k) DO UPDATE SET v = EXCLUDED.v", options.TableName),
		getStmt:    fmt.Sprintf("SELECT v FROM %s WHERE k=$1", options.TableName),
		deleteStmt: fmt.Sprintf("DELETE FROM %s WHERE k=$1", options.TableName),
	}

	// Create table if it doesn't exist yet
	_, err := options.Pool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS `+client.TableName+` (k TEXT PRIMARY KEY, v BYTEA NOT NULL)`)
	if err != nil {
		_ = client.Close()
		return nil, err
	}

	return client, nil
}

// Set stores the given value for the given key.
func (c *Client) Set(k string, v any) error {
	if err := util.CheckKeyAndValue(k, v); err != nil {
		return err
	}

	data, err := c.Codec.Marshal(v)
	if err != nil {
		return err
	}

	_, err = c.Pool.Exec(context.Background(), c.upsertStmt, k, data)
	return err
}

// Get retrieves the stored value for the given key.
func (c *Client) Get(k string, v any) (bool, error) {
	if err := util.CheckKeyAndValue(k, v); err != nil {
		return false, err
	}

	var data []byte
	err := c.Pool.QueryRow(context.Background(), c.getStmt, k).Scan(&data)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, c.Codec.Unmarshal(data, v)
}

// Delete deletes the stored value for the given key.
func (c *Client) Delete(k string) error {
	if err := util.CheckKey(k); err != nil {
		return err
	}

	_, err := c.Pool.Exec(context.Background(), c.deleteStmt, k)
	return err
}

// Close closes the connection pool.
func (c *Client) Close() error {
	c.Pool.Close()
	return nil
}
