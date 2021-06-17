package endpoints

import "net/http"

type RequestHandler = func(w http.ResponseWriter, r *http.Request)

type BlogPostJSON struct {
    BlogID   uint32 `db:"blogid"`
    Author   string `db:"author"`
    Date     string `db:"date"`
    Duration uint32 `db:"duration"`
    URL      string `db:"url"`
    Content  string `db:"content"`
}