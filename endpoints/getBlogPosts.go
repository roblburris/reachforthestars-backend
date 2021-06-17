package endpoints

import (
    "context"
    "encoding/json"
    "github.com/jackc/pgx/v4/pgxpool"
    "log"
    "net/http"
    "strings"

    "github.com/gorilla/mux"
    "github.com/roblburris/reachforthestars-backend/db"
)

// min returns minimum between two ints
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func convertBlogPostDBSlice (dbRes []db.BlogPost) []BlogPostJSON {
    var res []BlogPostJSON
    for i := 0; i < len(dbRes); i++ {
        res = append(res, convertBlogPostDBToJSON(dbRes[i]))
    }
    return res
}

func convertBlogPostDBToJSON(dbRes db.BlogPost) BlogPostJSON {
    var temp BlogPostJSON
    temp.BlogID = dbRes.BlogID
    temp.Duration = dbRes.Duration
    temp.Author = dbRes.Author
    temp.Date = dbRes.Date
    temp.Content = string(dbRes.Content)
    temp.URL = string(dbRes.URL)

    return temp
}

// GetAllBlogPosts Returns all blog posts for the blog page
func GetAllBlogPosts(ctx context.Context, conn *pgxpool.Pool) RequestHandler {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            text := "405: expected GET"
            log.Printf("incorrect request, %s\n", text)
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }

        // get blog posts from the DB and filter the content to create
        // summaries
        res := convertBlogPostDBSlice(db.GetAllBlogPosts(ctx, conn))
        for i := 0; i < len(res); i++ {
            strippedContent := strings.Fields(res[i].Content)
            temp := strings.Join(strippedContent[:min(len(strippedContent), 20)], " ")
            temp += "..."
            res[i].Content = temp
        }

        jsonRes, err := json.Marshal(res)
        if err != nil {
            text := "500: failed to marshal into JSON"
            log.Printf("internal server error, %s\n", text)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write(jsonRes)
    }
}

func GetSpecificBlogPost(ctx context.Context, conn *pgxpool.Pool) RequestHandler {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            text := "405: expected GET"
            log.Printf("incorrect request, %s\n", text)
            w.WriteHeader(http.StatusMethodNotAllowed)
            return
        }

        desiredTitle := mux.Vars(r)["title"]
        if desiredTitle == "" {
            text := "400: bad request"
            log.Printf("unable to parse request, %s\n", text)
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        // title is OK (for now), query DB to get desired blog post
        post := convertBlogPostDBToJSON(db.GetBlogPostByID(ctx, conn, desiredTitle))
        if post.Author == "" {
            text := "404: not found"
            log.Printf("unable to find post, %s\n", text)
            w.WriteHeader(http.StatusNotFound)
            return
        }

        // we have the post, encode to JSON and send back to client
        jsonRes, err := json.Marshal(post)
        if err != nil {
            text := "500: failed to marshal into JSON"
            log.Printf("internal server error, %s\n", text)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write(jsonRes)
    }
}
