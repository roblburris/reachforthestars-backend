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
}

func testGetAllBlogPostsDB(t *testing.T, ctx context.Context, conn *pgx.Conn) {
	rows := db.GetAllBlogPostsDB(ctx, conn)
	// check the results of the first row
	res0 := rows[0]
	error := false
	if res0.BlogID != 0 {
		t.Errorf("ERROR: expected BlogID 0 but got %d", res0.BlogID)
		error = true
	}
	if res0.Author != "John Doe" {
		t.Errorf("ERROR: expected Author `John Doe` but got `%s`", res0.Author)
		error = true
	}
	if res0.Date != "2020-01-01" {
		t.Errorf("ERROR: expected Date `2020-01-01` but got `%s`", res0.Date)
		error = true
	}
	if res0.Duration != 4 {
		t.Errorf("ERROR: expected Duration 4 but got %d", res0.Duration)
		error = true
	}

	if string(res0.URL) != "https://www.google.com/" {
		t.Errorf("ERROR: expected URL `https://www.google.com/` but got `%s`", res0.URL)
		error = true
	}

	if string(res0.Content) != "my name is John Doe" {
		t.Errorf("ERROR: expected Content `my name is John Doe` but got `%s`", res0.Content)
		error = true
	}
	res1 := rows[1]
	if res1.BlogID != 1 {
		t.Errorf("ERROR: expected BlogID 1 but got %d", res1.BlogID)
		error = true
	}
	if res1.Author != "Jane Doe" {
		t.Errorf("ERROR: expected Author `Jane Doe` but got `%s`", res1.Author)
		error = true
	}
	if res1.Date != "2021-01-01" {
		t.Errorf("ERROR: expected Date `2021-01-01` but got `%s`", res1.Date)
		error = true
	}
	if res1.Duration != 100 {
		t.Errorf("ERROR: expected Duration 4 but got %d", res1.Duration)
		error = true
	}

	if string(res1.URL) != "https://www.google.com/maps" {
		t.Errorf("ERROR: expected URL `https://www.google.com/maps` but got `%s`", res1.URL)
		error = true
	}

	if string(res1.Content) != "i am Jane Doe" {
		t.Errorf("ERROR: expected Content `i am Jane Doe` but got `%s`", res1.Content)
		error = true
	}

	if !error {t.Log("GetAllBlogPostsDB tests passed\n")}
}

