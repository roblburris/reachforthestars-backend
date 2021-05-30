package tests

import (
    "context"
    "github.com/jackc/pgx/v4/pgxpool"
    "github.com/roblburris/reachforthestars-backend/db"
    "testing"
)

// Runs all unit tests that relate to the Db
func TestDB(t *testing.T) {
    ctx := context.Background()
    pool, err := pgxpool.Connect(ctx, "postgres://localhost:5432/rfts-test")
    defer pool.Close()
    if err != nil {
        t.Fatalf("Unable to connect to DB, no tests run. Error: %v\n", err)
    }
    SetupTestDB(t, ctx, pool)
    // Unit tests for testing blog DB
    testGetAllBlogPostsDB(t, ctx, pool)
    testGetSpecificBP(t, ctx, pool)
    testInsertNewBP(t, ctx, pool)
}

func testGetAllBlogPostsDB(t *testing.T, ctx context.Context, conn *pgxpool.Pool) {
    rows := db.GetAllBlogPosts(ctx, conn)
    // check the results of the first row
    res0 := rows[0]
    if res0.BlogID != 1 {
        t.Fatalf("FAILED: expected BlogID 1 but got %d", res0.BlogID)
    }
    if res0.Author != "John Doe" {
        t.Fatalf("FAILED: expected Author `John Doe` but got `%s`", res0.Author)
    }
    if res0.Date != "2020-01-01" {
        t.Fatalf("FAILED: expected Date `2020-01-01` but got `%s`", res0.Date)
    }
    if res0.Duration != 4 {
        t.Fatalf("FAILED: expected Duration 4 but got %d", res0.Duration)
    }
    if string(res0.URL) != "https://www.google.com/" {
        t.Fatalf("FAILED: expected URL `https://www.google.com/` but got `%s`", res0.URL)
    }
    if string(res0.Content) != "my name is John Doe" {
        t.Fatalf("FAILED: expected Content `my name is John Doe` but got `%s`", res0.Content)
    }
    res1 := rows[1]
    if res1.BlogID != 2 {
        t.Fatalf("FAILED: expected BlogID 2 but got %d", res1.BlogID)
    }
    if res1.Author != "Jane Doe" {
        t.Fatalf("FAILED: expected Author `Jane Doe` but got `%s`", res1.Author)
    }
    if res1.Date != "2021-01-01" {
        t.Fatalf("FAILED: expected Date `2021-01-01` but got `%s`", res1.Date)
    }
    if res1.Duration != 100 {
        t.Fatalf("FAILED: expected Duration 4 but got %d", res1.Duration)
    }

    if string(res1.URL) != "https://www.google.com/maps" {
        t.Fatalf("FAILED: expected URL `https://www.google.com/maps` but got `%s`", res1.URL)
    }

    if string(res1.Content) != "i am Jane Doe" {
        t.Fatalf("FAILED: expected Content `i am Jane Doe` but got `%s`", res1.Content)
    }
    t.Log("db.GetAllBlogPosts tests passed\n")
}

func testGetSpecificBP(t *testing.T, ctx context.Context, conn *pgxpool.Pool) {
    res0 := db.GetBlogPostByID(ctx, conn, "john-doe")

    if res0.BlogID != 1 {
        t.Fatalf("FAILED: expected BlogID 1 but got %d", res0.BlogID)
    }
    if res0.Author != "John Doe" {
        t.Fatalf("FAILED: expected Author `John Doe` but got `%s`", res0.Author)
    }
    if res0.Date != "2020-01-01" {
        t.Fatalf("FAILED: expected Date `2020-01-01` but got `%s`", res0.Date)
    }
    if res0.Duration != 4 {
        t.Fatalf("FAILED: expected Duration 4 but got %d", res0.Duration)
    }
    if string(res0.URL) != "https://www.google.com/" {
        t.Fatalf("FAILED: expected URL `https://www.google.com/` but got `%s`", res0.URL)
    }
    if string(res0.Content) != "my name is John Doe" {
        t.Fatalf("FAILED: expected Content `my name is John Doe` but got `%s`", res0.Content)
    }

    t.Logf("db.GetBlogPostByID test passed")
}

func testInsertNewBP(t *testing.T, ctx context.Context, conn *pgxpool.Pool) {
    testInsertData := db.BlogPost{
        BlogID:   0,
        Author:   "Foo Bar",
        Date:     "2021-05-30",
        Duration: 109,
        URL:      []byte("https://mail.google.com/"),
        Content:  []byte("foo bar baz"),
    }
    title := "foo-bar"
    err := db.InsertNewBlogPost(ctx, conn, &testInsertData, title)
    if err != nil {
        t.Fatalf("FAILED: could not insert data, got error: %v\n", err)
    }

    res0 := db.GetBlogPostByID(ctx, conn, title)
    if res0.BlogID != 3 {
        t.Fatalf("FAILED: expected BlogID 3 but got %d", res0.BlogID)
    }
    if res0.Author != "Foo Bar" {
        t.Fatalf("FAILED: expected Author `Foo Bar` but got `%s`", res0.Author)
    }
    if res0.Date != "2021-05-30" {
        t.Fatalf("FAILED: expected Date `2021-05-30` but got `%s`", res0.Date)
    }
    if res0.Duration != 109 {
        t.Fatalf("FAILED: expected Duration 109 but got %d", res0.Duration)
    }
    if string(res0.URL) != "https://mail.google.com/" {
        t.Fatalf("FAILED: expected URL `https://mail.google.com/` but got `%s`", string(res0.URL))
    }
    if string(res0.Content) != "foo bar baz" {
        t.Fatalf("FAILED: expected Content `foo bar baz` but got `%s`", string(res0.Content))
    }

    t.Logf("db.InsertNewBlogPost test passed")
}