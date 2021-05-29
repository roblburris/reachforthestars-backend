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

// GetAllBlogPosts returns all blog posts for display on the blog page
func GetAllBlogPosts(ctx context.Context, conn *pgx.Conn) []BlogPost {
    tx, err := conn.BeginTx(ctx, pgx.TxOptions{
        IsoLevel: pgx.ReadUncommitted,
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

    rows, err := tx.Query(ctx, "SELECT * FROM BLOG_POSTS")
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

// GetBlogPostByID gets blog post by specific blog 'title'
func GetBlogPostByID(ctx context.Context, conn *pgx.Conn, path string) BlogPost {
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
    tx, err := conn.BeginTx(ctx, pgx.TxOptions{
        IsoLevel: pgx.Serializable,
    })
    defer tx.Rollback(ctx)

    // ensure title is unused
    rows, err := tx.Query(ctx, COUNT_TITLE_OCCURENCES)
    var countTitle int
    for rows.Next() {
        err = rows.Scan(&countTitle)
        if err != nil {
            log.Printf("ERROR: unable count title occurences. %v\n", err)
            return errors.New("sql_error")
        }
    }
    if countTitle != 0 {
        return errors.New("duplicate_title")
    }

    // find the max blogID + 1 to use as next blogID
    rows, err = tx.Query(ctx, FIND_MAX_BLOGID)
    if err != nil {
        log.Printf("ERROR: unable to find max blog ID. %v\n", err)
        return errors.New("sql_error")
    }

    var blogID int
    for rows.Next() {
        err = rows.Scan(&blogID)
        if err != nil {
            log.Printf("ERROR: unable to find max blog ID. %v\n", err)
            return errors.New("sql_error")
        }
    }

    // everything is OK, insert into DB
    postInfo.BlogID = uint32(blogID + 1)
    _, err = tx.Exec(ctx,
        INSERT_NEW_POST,
        postInfo.BlogID,
        postInfo.Author,
        postInfo.Date,
        postInfo.Duration,
        postInfo.URL,
        postInfo.Content)
    if err != nil {
        log.Printf("ERROR: unable to insert post into DB. %v\n", err)
        return errors.New("sql_error")
    }

    // likewise, insert into BLOG_POST_TITLES table
    _, err = tx.Exec(ctx, INSERT_NEW_TITLE, postInfo.BlogID, title)
    if err != nil {
        log.Printf("ERROR: unable to insert post into DB. %v\n", err)
        return errors.New("sql_error")
    }

    err = tx.Commit(ctx)
    if err != nil {
        return errors.New("sql_unable_commit")
    }

    return nil
}
