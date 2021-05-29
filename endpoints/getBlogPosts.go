package endpoints

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v4"
	"github.com/roblburris/reachforthestars-backend/db"
	"log"
	"net/http"
)

func GetBlogPosts(ctx context.Context, conn *pgx.Conn) RequestHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			text := "405: expected GET"
			log.Printf("incorrect request, %s\n", text)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// get blog posts from the DB
		res := db.GetAllBlogPosts(ctx, conn)
		json_res, err := json.Marshal(res)
		if err != nil {
			text := "400: failed to marshal into JSON"
			log.Printf("internal server error, %s\n", text)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(json_res)
	}

}
