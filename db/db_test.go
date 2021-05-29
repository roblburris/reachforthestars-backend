package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"io/ioutil"
	"testing"
)

// Runs all unit tests that relate to the Db
func TestDB(t *testing.T) {
	conn, err := pgx.Connect(context.Background(), "postgres://localhost:5432/rfts-test")
	if err != nil {
		t.Fatalf("Unable to connect to DB, no tests run. Error: %v\n", err)
	}
	defer conn.Close(context.Background())

	// clear DB and run create tables
	deleteTables, err := ioutil.ReadFile("./setup/delete-tables.sql")
	if err != nil {
		t.Fatalf("Unable to read delete-tables file, no tests run. Error: %v\n", err)
	}
	_, err = conn.Exec(context.Background(), string(deleteTables))
	if err != nil {
		t.Fatalf("Unable to wipe existing tables, no tests run. Error: %v\n", err)
	}
	createTables, err := ioutil.ReadFile("./setup/create-tables.sql")
	if err != nil {
		t.Fatalf("Unable to read create-tables file, no tests run. Error: %v\n", err)
	}
	_, err = conn.Exec(context.Background(), string(createTables))
	if err != nil {
		t.Fatalf("Unable create necessary tables, no tests run %v\n", err)
	}

	// insert test data into BLOG_POSTS table
	test1 := BlogPost {
		blogid:   0,
		author:   "John Doe",
		date:     "2020-01-01",
		duration: 4,
		url:      []byte("https://www.google.com/"),
		content:  []byte("my name is John Doe"),
	}

	test2 := BlogPost {
		blogid:   1,
		author:   "Jane Doe",
		date:     "2021-01-01",
		duration: 100,
		url:     []byte("https://www.google.com/maps"),
		content:  []byte("i am Jane Doe"),
	}

	// insert test data
	_, err = conn.Exec(
		context.Background(),
		"INSERT INTO BLOG_POSTS VALUES ($1, $2, $3, $4, $5, $6)",
		test1.blogid,
		test1.author,
		test1.date,
		test1.duration,
		test1.url,
		test1.content)
	if err != nil {
		t.Fatalf("Unable to insert test data, no tests run. Error: %v\n", err)
	}

	_, err = conn.Exec(
		context.Background(),
		"INSERT INTO BLOG_POSTS VALUES ($1, $2, $3, $4, $5, $6)",
		test2.blogid,
		test2.author,
		test2.date,
		test2.duration,
		test2.url,
		test2.content)
	if err != nil {
		t.Fatalf("Unable to insert test data, no tests run. Error: %v\n", err)
	}

	// Unit tests for testing blog DB
	testGetAllBlogPosts(t, context.Background(), conn)
}

func testGetAllBlogPosts(t * testing.T, context context.Context, conn *pgx.Conn) {
	rows := getAllBlogPosts(context, conn)
	// check the results of the first row
	res0 := rows[0]
	error := false
	if res0.blogid != 0 {
		t.Errorf("ERROR: expected blogid 0 but got %d", res0.blogid)
		error = true
	}
	if res0.author != "John Doe" {
		t.Errorf("ERROR: expected author `John Doe` but got `%s`", res0.author)
		error = true
	}
	if res0.date != "2020-01-01" {
		t.Errorf("ERROR: expected date `2020-01-01` but got `%s`", res0.date)
		error = true
	}
	if res0.duration != 4 {
		t.Errorf("ERROR: expected duration 4 but got %d", res0.duration)
		error = true
	}

	if string(res0.url) != "https://www.google.com/" {
		t.Errorf("ERROR: expected url `https://www.google.com/` but got `%s`", res0.url)
		error = true
	}

	if string(res0.content) != "my name is John Doe" {
		t.Errorf("ERROR: expected content `my name is John Doe` but got `%s`", res0.content)
		error = true
	}
	res1 := rows[1]
	if res1.blogid != 1 {
		t.Errorf("ERROR: expected blogid 1 but got %d", res1.blogid)
		error = true
	}
	if res1.author != "Jane Doe" {
		t.Errorf("ERROR: expected author `Jane Doe` but got `%s`", res1.author)
		error = true
	}
	if res1.date != "2021-01-01" {
		t.Errorf("ERROR: expected date `2021-01-01` but got `%s`", res1.date)
		error = true
	}
	if res1.duration != 100 {
		t.Errorf("ERROR: expected duration 4 but got %d", res1.duration)
		error = true
	}

	if string(res1.url) != "https://www.google.com/maps" {
		t.Errorf("ERROR: expected url `https://www.google.com/maps` but got `%s`", res1.url)
		error = true
	}

	if string(res1.content) != "i am Jane Doe" {
		t.Errorf("ERROR: expected content `i am Jane Doe` but got `%s`", res1.content)
		error = true
	}

	if !error {t.Log("GetAllBlogPosts tests passed\n")}
}

