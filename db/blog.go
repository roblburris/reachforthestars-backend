package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
)

func getAllBlogPosts(context context.Context, conn *pgx.Conn) []BlogPost {
	rows, err:= conn.Query(context, "SELECT * FROM BLOG_POSTS")
	if err != nil {
		log.Println("ERROR: Unable to get blog posts")
		return nil
	}

	defer rows.Close()

	var blogPosts []BlogPost
	for rows.Next() {
		var temp BlogPost
		err = rows.Scan(&temp.blogid, &temp.author, &temp.date, &temp.duration, &temp.url, &temp.content)
		if err != nil {
			log.Printf("ERROR: Unable to parse SQL data. %v\n", err)
			return nil
		}
		blogPosts = append(blogPosts, temp)
	}
	return blogPosts
}