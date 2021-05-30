package tests

import (
   "context"
   "encoding/json"
   "github.com/jackc/pgx/v4/pgxpool"
   "github.com/roblburris/reachforthestars-backend/db"
   "github.com/roblburris/reachforthestars-backend/endpoints"
   "io/ioutil"
   "net/http"
   "net/http/httptest"
   "testing"
)

func TestEndpoints(t *testing.T) {
   ctx := context.Background()
   // test the getBlogPosts endpoint
   conn, err := pgxpool.Connect(ctx, "postgres://localhost:5432/rfts-test")
   if err != nil {
       t.Fatalf("Unable to connect to DB, no tests run. Error: %v\n", err)
   }
   defer conn.Close()
   SetupTestDB(t, ctx, conn)
   testGetAllBlogPosts(ctx, conn, t)
}

func testGetAllBlogPosts(ctx context.Context, conn *pgxpool.Pool, t *testing.T) {
   getAllBlogPostsHandler := endpoints.GetAllBlogPosts(ctx, conn)
   req := httptest.NewRequest("GET", "http://example.com/foo", nil)
   w := httptest.NewRecorder()
   getAllBlogPostsHandler(w, req)
   response := w.Result()
   if response.StatusCode != http.StatusOK {
       t.Fatalf("FAILED: expected HTTP Status Code 200 but got: %d\n", response.StatusCode)
       return
   }

   if response.Header.Get("Content-Type") != "application/json" {
       t.Fatalf("FAILED: expected Content-Type `application/json` but got %s\n",
           response.Header.Get("Content-Type"))
       return
   }

   var result []db.BlogPost
   body, err := ioutil.ReadAll(response.Body)
   if err != nil {
       t.Fatalf("FAILED: unable to read body. %v\n", err)
       return
   }
   err = json.Unmarshal(body, &result)
   if err != nil {
       t.Fatalf("FAILED: unable to decode JSON %v\n", err)
       return
   }

   res0 := result[0]
   if res0.BlogID != 1 {
       t.Fatalf("FAILED: expected BlogID 1 but got %d", res0.BlogID)
   }
   if res0.Author != "John Doe" {
       t.Fatalf("FAILED: expected Author `John Doe` but got `%s`", res0.Author)
   }
   if res0.Date != "2020-01-01" {
       t.Fatalf("FAILED: expected Date `2020-01-01` but got `%s`", res0.Date)
   }
   if res0.Duration != 4 {
       t.Fatalf("FAILED: expected Duration 4 but got %d", res0.Duration)
   }
   if string(res0.URL) != "https://www.google.com/" {
       t.Fatalf("FAILED: expected URL `https://www.google.com/` but got `%s`", res0.URL)
   }
   if string(res0.Content) != "my name is John Doe..." {
       t.Fatalf("FAILED: expected Content `my name is John Doe...` but got `%s`", res0.Content)
   }

   res1 := result[1]
   if res1.BlogID != 2 {
       t.Fatalf("FAILED: expected BlogID 2 but got %d", res1.BlogID)
   }
   if res1.Author != "Jane Doe" {
       t.Fatalf("FAILED: expected Author `Jane Doe` but got `%s`", res1.Author)
   }
   if res1.Date != "2021-01-01" {
       t.Fatalf("FAILED: expected Date `2021-01-01` but got `%s`", res1.Date)
   }
   if res1.Duration != 100 {
       t.Fatalf("FAILED: expected Duration 4 but got %d", res1.Duration)
   }
   if string(res1.URL) != "https://www.google.com/maps" {
       t.Fatalf("FAILED: expected URL `https://www.google.com/maps` but got `%s`", res1.URL)
   }
   if string(res1.Content) != "i am Jane Doe..." {
       t.Fatalf("FAILED: expected Content `i am Jane Doe...` but got `%s`", res1.Content)
   }

   t.Log("endpoints.GetAllBlogPosts tests passed\n")
}
