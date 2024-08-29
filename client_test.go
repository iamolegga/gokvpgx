package gokvpgx_test

import (
	"context"
	"testing"

	"github.com/iamolegga/gokvpgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/philippgille/gokv/encoding"
	"github.com/philippgille/gokv/test"
)

func TestClient(t *testing.T) {
	// Test with JSON
	t.Run("JSON", func(t *testing.T) {
		client := createClient(t, encoding.JSON)
		defer client.Close()
		test.TestStore(client, t)
	})

	// Test with gob
	t.Run("gob", func(t *testing.T) {
		client := createClient(t, encoding.Gob)
		defer client.Close()
		test.TestStore(client, t)
	})
}

func createClient(t *testing.T, codec encoding.Codec) *gokvpgx.Client {
	pool, err := pgxpool.Connect(context.Background(), "postgres://postgres:secret@localhost:5432/gokv?sslmode=disable")
	if err != nil {
		t.Fatalf("Unable to connect to database: %v", err)
	}
	options := gokvpgx.Options{
		Pool:      pool,
		Codec:     codec,
		TableName: "test_table",
	}
	client, err := gokvpgx.NewClient(options)
	if err != nil {
		t.Fatal(err)
	}
	return client
}
