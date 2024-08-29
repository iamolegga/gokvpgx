# gokvpgx

[![PkgGoDev](https://pkg.go.dev/badge/github.com/iamolegga/gokvpgx)](https://pkg.go.dev/github.com/iamolegga/gokvpgx)

`gokvpgx` is a PostgreSQL client for the `gokv` key-value store interface, implemented using the `pgx` library. This package allows you to use PostgreSQL as a key-value store with efficient prepared statements, ensuring optimal performance and security.

## Features

- **pgx-based Implementation**: Utilizes the `pgx` library for efficient PostgreSQL interactions.
- **Efficient Prepared Statements**: SQL statements are prepared at the database level for faster execution, thanks to `pgx`.
- **Flexible Serialization**: Supports different encoding formats like JSON and gob through the `gokv` interface.
- **Connection Pooling**: Leverages `pgxpool.Pool` for efficient database connection management.
- **Autocreation of Table**: Automatically creates the key-value store table if it does not exist.

## Installation

To use `gokvpgx`, you need to have a Go environment set up. Then, you can get the package using:

```bash
go get github.com/iamolegga/gokvpgx
```

## Usage

Here is an example of how to use the gokvpgx package to store, retrieve, and delete key-value pairs in PostgreSQL:

```go
package main

import (
    "context"
    "log"

    "github.com/iamolegga/gokvpgx"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/philippgille/gokv/encoding"
)

func main() {
    // Create a connection pool
    pool, err := pgxpool.New(context.Background(), "postgres://postgres:secret@localhost:5432/gokv?sslmode=disable")
    if err != nil {
        log.Fatalf("Unable to connect to database: %v", err)
    }
    defer pool.Close()

    // Create a new client
    options := gokvpgx.Options{
        Pool:      pool,
        Codec:     encoding.JSON,
        TableName: "kv_store",
    }
    client, err := gokvpgx.NewClient(options)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    defer client.Close()

    // Set a value
    err = client.Set("my_key", "my_value")
    if err != nil {
        log.Fatalf("Failed to set value: %v", err)
    }

    // Get the value
    var value string
    found, err := client.Get("my_key", &value)
    if err != nil {
        log.Fatalf("Failed to get value: %v", err)
    }
    if found {
        log.Printf("Found value: %s", value)
    } else {
        log.Println("Value not found")
    }

    // Delete the value
    err = client.Delete("my_key")
    if err != nil {
        log.Fatalf("Failed to delete value: %v", err)
    }
}
```

## Options

### Options Struct

- Pool: *pgxpool.Pool - The PostgreSQL connection pool to use.
- TableName: string - The name of the table where key-value pairs are stored. Default is "Item".
- Codec: encoding.Codec - The codec used for serializing and deserializing values. Default is encoding.JSON.

### DefaultOptions

- TableName: "Item"
- Codec: encoding.JSON

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/iamolegga/gokvpgx/blob/main/LICENSE) file for details.

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request if you find a bug or want to add a new feature.
