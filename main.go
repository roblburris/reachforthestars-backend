package main

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
)

func main() {
	// create DB connection
	conn, err := pgx.Connect(context.Background(), URL)
	if err != nil {
		log.Fatalf("Unable to connect to DB %v\n", err)
	}
	defer conn.Close(context.Background())
}
