package tests

import (
    "context"
    "github.com/jackc/pgx/v4"
    "github.com/roblburris/reachforthestars-backend/db"
    "io/ioutil"
    "testing"
)

func SetupTestDB(t *testing.T, ctx context.Context, conn *pgx.Conn) {
    // clear DB and run create tables
    deleteTables, err := ioutil.ReadFile("../db/setup/delete-tables.sql")
    if err != nil {
        t.Fatalf("Unable to read delete-tables file, no tests run. Error: %v\n", err)
    }
    _, err = conn.Exec(ctx, string(deleteTables))
    if err != nil {
        t.Fatalf("Unable to wipe existing tables, no tests run. Error: %v\n", err)
    }
    createTables, err := ioutil.ReadFile("../db/setup/create-tables.sql")
    if err != nil {
        t.Fatalf("Unable to read create-tables file, no tests run. Error: %v\n", err)
    }
    _, err = conn.Exec(ctx, string(createTables))
    if err != nil {
        t.Fatalf("Unable create necessary tables, no tests run %v\n", err)
    }

    // insert test data into BLOG_POSTS table
    test1 := db.BlogPost{
        BlogID:   1,
        Author:   "John Doe",
        Date:     "2020-01-01",
        Duration: 4,
        URL:      []byte("https://www.google.com/"),
        Content:  []byte("my name is John Doe"),
    }

    test2 := db.BlogPost{
        BlogID:   2,
        Author:   "Jane Doe",
        Date:     "2021-01-01",
        Duration: 100,
        URL:      []byte("https://www.google.com/maps"),
        Content:  []byte("i am Jane Doe"),
    }

    // insert test data
    _, err = conn.Exec(
        ctx,
        "INSERT INTO BLOG_POSTS VALUES ($1, $2, $3, $4, $5, $6)",
        test1.BlogID,
        test1.Author,
        test1.Date,
        test1.Duration,
        test1.URL,
        test1.Content)
    if err != nil {
        t.Fatalf("Unable to insert test data, no tests run. Error: %v\n", err)
    }

    _, err = conn.Exec(
        ctx,
        "INSERT INTO BLOG_POSTS VALUES ($1, $2, $3, $4, $5, $6)",
        test2.BlogID,
        test2.Author,
        test2.Date,
        test2.Duration,
        test2.URL,
        test2.Content)
    if err != nil {
        t.Fatalf("Unable to insert test data, no tests run. Error: %v\n", err)
    }

    _, err = conn.Exec(
        ctx,
        "INSERT INTO BLOG_POST_TITLES VALUES ($1, $2)",
        test1.BlogID,
        "john-doe")
    if err != nil {
        t.Fatalf("Unable to insert test data, no tests run. Error: %v\n", err)
    }

    _, err = conn.Exec(
        ctx,
        "INSERT INTO BLOG_POST_TITLES VALUES ($1, $2)",
        test2.BlogID,
        "jane-doe")
    if err != nil {
        t.Fatalf("Unable to insert test data, no tests run. Error: %v\n", err)
    }
}
