package main

import (
    "context"
    "github.com/jackc/pgx/v4"
    "log"
)

func main() {
    // create context
    ctx := context.Background()
    // create DB connection
    conn, err := pgx.Connect(ctx, URL)
    if err != nil {
       log.Fatalf("Unable to connect to DB %v\n", err)
    }
    defer conn.Close(ctx)

    // clear DB and run create tables
}
