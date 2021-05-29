package db

// BlogPost represents blog post entry in DB
type BlogPost struct {
    BlogID   uint32 `db:"blogid"`
    Author   string `db:"author"`
    Date     string `db:"date"`
    Duration uint32 `db:"duration"`
    URL      []byte `db:"url"`
    Content  []byte `db:"content"`
}
