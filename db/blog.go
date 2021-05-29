package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
)

// GetAllBlogPostsDB returns all blog posts for display on the blog page
func GetAllBlogPostsDB(ctx context.Context, conn *pgx.Conn) []BlogPost {
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.ReadUncommitted,
	})

	if err != nil {
		log.Println("ERROR: Unable to get blog posts")
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
		log.Println("ERROR: Unable to get blog posts")
		return nil
	}

	defer rows.Close()

	var blogPosts []BlogPost
	for rows.Next() {
		var temp BlogPost
		err = rows.Scan(&temp.BlogID, &temp.Author, &temp.Date, &temp.Duration, &temp.URL, &temp.Content)
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