package endpoints

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/jackc/pgx/v4/pgxpool"
    "github.com/roblburris/reachforthestars-backend/db"
    "io/ioutil"
    "log"
    "net/http"
)

func InsertNewBlogPost(ctx context.Context, conn *pgxpool.Pool) RequestHandler {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            text := "405: expected Post"
            log.Printf("incorrect request, %s\n", text)
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }

        // decode body of request into map before putting into BlogPost struct
        var result map[string]interface{}
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            text := "400: bad request"
            log.Printf("ERROR: unable to decode body. %s\n", text)
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        err = json.Unmarshal(body, &result)
        if err != nil {
            text := "400: bad request"
            log.Printf("ERROR: unable to decode json. %s\n", text)
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        if result["title"] == nil {
            text := "400: bad request"
            log.Printf("ERROR: could not read \"title\" field. %s\n", text)
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        title := result["title"].(string)
        if title == "" {
            text := "400: bad request"
            log.Printf("ERROR: could not read \"title\" field. %s\n", text)
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        if result["author"] == nil {
            text := "400: bad request"
            log.Printf("ERROR: could not read \"author\" field. %s\n", text)
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        author := result["author"].(string)
        if author == "" {
            text := "400: bad request"
            log.Printf("ERROR: could not read \"author\" field. %s\n", text)
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        if result["date"] == nil {
            text := "400: bad request"
            log.Printf("ERROR: could not read \"date\" field. %s\n", text)
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        date := result["date"].(string)
        if date == "" {
            text := "400: bad request"
            log.Printf("ERROR: could not read \"date\" field. %s\n", text)
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        if result["duration"] == nil {
            text := "400: bad request"
            log.Printf("ERROR: could not read \"duration\" field. %s\n", text)
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        duration := result["duration"].(float64)
        if duration == 0 {
            text := "400: bad request"
            log.Printf("ERROR: could not read \"duration\" field. %s\n", text)
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        if result["url"] == nil {
            text := "400: bad request"
            log.Printf("ERROR: could not read \"url\" field. %s\n", text)
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        url := fmt.Sprintf("%v", result["url"])
        if url == "" {
            text := "400: bad request"
            log.Printf("ERROR: could not read \"url\" field. %s\n", text)
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        if result["content"] == nil {
            text := "400: bad request"
            log.Printf("ERROR: could not read \"content\" field. %s\n", text)
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        content := fmt.Sprintf("%v", result["content"])
        if content == "" {
            text := "400: bad request"
            log.Printf("ERROR: could not read \"url\" field. %s\n", text)
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        // insert into DB being careful to handle the possible error
        blogPost := db.BlogPost{
            BlogID:   0,
            Author:   author,
            Date:     date,
            Duration: uint32(duration),
            URL: []byte(url),
            Content: []byte(content),
        }
        err = db.InsertNewBlogPost(ctx, conn, &blogPost, title)
        if err != nil {
            text := "500: internal error"
            log.Printf("ERROR: unable to insert post. %v. %s", err, text)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
    }
}
