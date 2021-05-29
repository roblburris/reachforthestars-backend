package db

// BlogPost represents blog post entry in DB
type BlogPost struct {
	blogid uint32 `db:"blogid"`
	author string `db:"author"`
	date string `db:"date"`
	duration uint32 `db:"duration"`
	url []byte `db:"url"`
	content []byte `db:"content"`
}

func main() {
	return
}
