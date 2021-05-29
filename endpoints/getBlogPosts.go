package endpoints

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4"
	"github.com/roblburris/reachforthestars-backend/db"
	"log"
	"net/http"
	"strings"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
// GetAllBlogPosts Returns all blog posts for the blog page
func GetAllBlogPosts(ctx context.Context, conn *pgx.Conn) RequestHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			text := "405: expected GET"
			log.Printf("incorrect request, %s\n", text)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// get blog posts from the DB and filter the content to create
		// summaries
		res := db.GetAllBlogPostsDB(ctx, conn)
		for i := 0; i < len(res); i++ {
			strippedContent := strings.Fields(string(res[i].Content))
			desLen := min(len(strippedContent), 20)
			temp := strings.Join(strippedContent[:desLen], " ")
			temp += "..."
			res[i].Content = []byte(temp)
		}

		jsonRes, err := json.Marshal(res)
		if err != nil {
			text := "400: failed to marshal into JSON"
			log.Printf("internal server error, %s\n", text)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonRes)
	}

}


