package db

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"log"
)

// Prepared SQL statements
const GET_BLOG_POST_STATEMENT = `SELECT 
							FROM BLOG_POSTS as b, BLOG_POST_TITLES as t 
							WHERE t.blogid = b.blogID
							AND t.blogTitle = $1`
const INSERT_NEW_POST = `INSERT INTO BLOG_POSTS VALUES ($1, $2, $3, $4, $5, $6)`
const INSERT_NEW_TITLE = `INSERT INTO BLOG_POST_TITLES VALUES($1, $2)`
const COUNT_TITLE_OCCURENCES = `SELECT COUNT(*) FROM BLOG_POST_TITLES WHERE blogTitle = $1`
const FIND_MAX_BLOGID = `SELECT MAX(blogid) FROM BLOG_POSTS`

// GetAllBlogPostsDB returns all blog posts for display on the blog page
func GetAllBlogPostsDB(ctx context.Context, conn *pgx.Conn) []BlogPost {
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.ReadUncommitted,
	})
	if err != nil {
		log.Printf("ERROR: Unable to set transaction level. %v\n", err)
		return nil
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			tx.Rollback(ctx)
		}
	}(tx, ctx)

	rows, err:= tx.Query(ctx, "SELECT * FROM BLOG_POSTS")
	if err != nil {
		log.Printf("ERROR: Unable to get blog posts from DB. %v\n", err)
		return nil
	}

	defer rows.Close()

	var blogPosts []BlogPost
	for rows.Next() {
		var temp BlogPost
		err = rows.Scan(&temp.BlogID,
		&temp.Author,
		&temp.Date,
		&temp.Duration,
		&temp.URL,
		&temp.Content)

		if err != nil {
			log.Printf("ERROR: Unable to parse SQL data. %v\n", err)
			return nil
		}
		blogPosts = append(blogPosts, temp)
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
	}
	return blogPosts
}

// GetBlogPostByIDDB gets blog post by specific blog 'title'
func GetBlogPostByIDDB(ctx context.Context, conn *pgx.Conn, path string) BlogPost {
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadUncommitted,
	})
	if err != nil {
		log.Printf("ERROR: Unable to set transaction level. %v\n", err)
		return BlogPost{}
	}

	queryRes, err := tx.Query(ctx, GET_BLOG_POST_STATEMENT, path)
	if err != nil {
		log.Printf("ERROR: Unable to get specified blog post from DB. %v\n", err)
		return BlogPost{}
	}
	var postRes BlogPost
	for queryRes.Next() {
		err = queryRes.Scan(&postRes.BlogID,
			&postRes.Author,
			&postRes.Date,
			&postRes.Duration,
			&postRes.URL,
			&postRes.Content)

		if err != nil {
			log.Printf("ERROR: Unable to parse SQL data. %v\n", err)
			return BlogPost{}
		}
	}
	return postRes
}

func InsertNewBlogPost(ctx context.Context, conn *pgx.Conn, postInfo *BlogPost, title string) error {
	// TODO: finish InsertNewBlogPost function that inserts new blog post into DB
	return errors.New("new error")
}