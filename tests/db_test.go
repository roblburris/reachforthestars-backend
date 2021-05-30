package tests

import (
    "context"
    "github.com/jackc/pgx/v4"
    "github.com/roblburris/reachforthestars-backend/db"
    "testing"
)

// Runs all unit tests that relate to the Db
func TestDB(t *testing.T) {
    ctx := context.Background()
    conn, err := pgx.Connect(ctx, "postgres://localhost:5432/rfts-test")
    if err != nil {
        t.Fatalf("Unable to connect to DB, no tests run. Error: %v\n", err)
    }
    defer conn.Close(ctx)
    SetupTestDB(t, ctx, conn)
    // Unit tests for testing blog DB
    testGetAllBlogPostsDB(t, ctx, conn)
    testGetSpecificBP(t, ctx, conn)
}

func testGetAllBlogPostsDB(t *testing.T, ctx context.Context, conn *pgx.Conn) {
    rows := db.GetAllBlogPosts(ctx, conn)
    // check the results of the first row
    res0 := rows[0]
    if res0.BlogID != 1 {
        t.Fatalf("ERROR: expected BlogID 1 but got %d", res0.BlogID)
    }
    if res0.Author != "John Doe" {
        t.Fatalf("ERROR: expected Author `John Doe` but got `%s`", res0.Author)
    }
    if res0.Date != "2020-01-01" {
        t.Fatalf("ERROR: expected Date `2020-01-01` but got `%s`", res0.Date)
    }
    if res0.Duration != 4 {
        t.Fatalf("ERROR: expected Duration 4 but got %d", res0.Duration)
    }
    if string(res0.URL) != "https://www.google.com/" {
        t.Fatalf("ERROR: expected URL `https://www.google.com/` but got `%s`", res0.URL)
    }
    if string(res0.Content) != "my name is John Doe" {
        t.Fatalf("ERROR: expected Content `my name is John Doe` but got `%s`", res0.Content)
    }
    res1 := rows[1]
    if res1.BlogID != 2 {
        t.Fatalf("ERROR: expected BlogID 2 but got %d", res1.BlogID)
    }
    if res1.Author != "Jane Doe" {
        t.Fatalf("ERROR: expected Author `Jane Doe` but got `%s`", res1.Author)
    }
    if res1.Date != "2021-01-01" {
        t.Fatalf("ERROR: expected Date `2021-01-01` but got `%s`", res1.Date)
    }
    if res1.Duration != 100 {
        t.Fatalf("ERROR: expected Duration 4 but got %d", res1.Duration)
    }

    if string(res1.URL) != "https://www.google.com/maps" {
        t.Fatalf("ERROR: expected URL `https://www.google.com/maps` but got `%s`", res1.URL)
    }

    if string(res1.Content) != "i am Jane Doe" {
        t.Fatalf("ERROR: expected Content `i am Jane Doe` but got `%s`", res1.Content)
    }
    t.Log("db.GetAllBlogPosts tests passed\n")
}

func testGetSpecificBP(t *testing.T, ctx context.Context, conn *pgx.Conn) {
    res0 := db.GetBlogPostByID(ctx, conn, "john-doe")

    if res0.BlogID != 1 {
        t.Fatalf("ERROR: expected BlogID 1 but got %d", res0.BlogID)
    }
    if res0.Author != "John Doe" {
        t.Fatalf("ERROR: expected Author `John Doe` but got `%s`", res0.Author)
    }
    if res0.Date != "2020-01-01" {
        t.Fatalf("ERROR: expected Date `2020-01-01` but got `%s`", res0.Date)
    }
    if res0.Duration != 4 {
        t.Fatalf("ERROR: expected Duration 4 but got %d", res0.Duration)
    }
    if string(res0.URL) != "https://www.google.com/" {
        t.Fatalf("ERROR: expected URL `https://www.google.com/` but got `%s`", res0.URL)
    }
    if string(res0.Content) != "my name is John Doe" {
        t.Fatalf("ERROR: expected Content `my name is John Doe` but got `%s`", res0.Content)
    }

    t.Logf("db.GetBlogPostByID test passed")
}
