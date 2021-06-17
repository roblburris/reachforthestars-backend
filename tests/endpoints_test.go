package tests

import (
    "bytes"
    "context"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/jackc/pgx/v4/pgxpool"
    "github.com/roblburris/reachforthestars-backend/endpoints"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "testing"
)

// TODO: FIX encoding with endpoints
func TestEndpoints(t *testing.T) {
   ctx := context.Background()
   // test the getBlogPosts endpoint
   conn, err := pgxpool.Connect(ctx, "postgres://localhost:5432/rfts-test")
   if err != nil {
       t.Fatalf("Unable to connect to DB, no tests run. Error: %v\n", err)
   }
   defer conn.Close()
   SetupTestDB(t, ctx, conn)
   testGetAllBlogPostsEndpoint(t, ctx, conn)
   testGetSpecificBlogPostEndpoint(t, ctx, conn)
   testNewBlogPostEndpoint(t, ctx, conn)
}

func testGetAllBlogPostsEndpoint(t *testing.T, ctx context.Context, conn *pgxpool.Pool) {
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

   var result []endpoints.BlogPostJSON
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
   if res0.URL != "https://www.google.com/" {
       t.Fatalf("FAILED: expected URL `https://www.google.com/` but got `%s`", res0.URL)
   }
   if res0.Content != "my name is John Doe..." {
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
   if res1.URL != "https://www.google.com/maps" {
       t.Fatalf("FAILED: expected URL `https://www.google.com/maps` but got `%s`", res1.URL)
   }
   if res1.Content != "i am Jane Doe..." {
       t.Fatalf("FAILED: expected Content `i am Jane Doe...` but got `%s`", res1.Content)
   }

   t.Log("endpoints.GetAllBlogPosts test passed\n")
}

func testGetSpecificBlogPostEndpoint(t *testing.T, ctx context.Context, conn *pgxpool.Pool) {
    getSpecificBPHandler := endpoints.GetSpecificBlogPost(ctx, conn)
    r := mux.NewRouter()
    r.HandleFunc("/blog/{title}", getSpecificBPHandler)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, httptest.NewRequest("GET", "/blog/john-doe", nil))
    response := w.Result()
    if response.StatusCode != http.StatusOK {
        t.Fatalf("FAILED: expected status code 200 but got %d\n", w.Code)
    }

    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        t.Fatalf("FAILED: unable to read body. %v\n", err)
        return
    }
    var res0 endpoints.BlogPostJSON
    err = json.Unmarshal(body, &res0)
    if err != nil {
        t.Fatalf("FAILED: unable to decode JSON %v\n", err)
        return
    }

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
    if res0.URL != "https://www.google.com/" {
        t.Fatalf("FAILED: expected URL `https://www.google.com/` but got `%s`", res0.URL)
    }
    if res0.Content != "my name is John Doe" {
        t.Fatalf("FAILED: expected Content `my name is John Doe` but got `%s`", res0.Content)
    }

    t.Log("endpoints.GetSpecificBlogPost test passed\n")
}

func testNewBlogPostEndpoint(t *testing.T, ctx context.Context, conn *pgxpool.Pool) {
    insertNewBPHandler := endpoints.InsertNewBlogPost(ctx, conn)
    testInsert := InsertPost{
        Title:    "foo-bar",
        Author:   "Foo Bar",
        Date:     "2021-05-30",
        Duration: 109,
        URL:      "https://mail.google.com/",
        Content:  "foo bar baz",
    }

    jsonBody, err := json.Marshal(testInsert)
    if err != nil {
        t.Fatalf("ERROR: unable to encode to json. Test failed. %v\n", err)
        return
    }
    req := httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonBody))
    w := httptest.NewRecorder()
    insertNewBPHandler(w, req)
    response := w.Result()
    if response.StatusCode != http.StatusOK {
        t.Fatalf("FAILED: expected HTTP Status Code 200 but got: %d\n", response.StatusCode)
        return
    }

    // ensure post got inserted into DB
    getSpecificBPHandler := endpoints.GetSpecificBlogPost(ctx, conn)
    r := mux.NewRouter()
    r.HandleFunc("/blog/{title}", getSpecificBPHandler)
    w = httptest.NewRecorder()
    r.ServeHTTP(w, httptest.NewRequest("GET", "/blog/foo-bar", nil))
    response = w.Result()
    if response.StatusCode != http.StatusOK {
        t.Fatalf("FAILED: expected status code 200 but got %d\n", w.Code)
    }

    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        t.Fatalf("FAILED: unable to read body. %v\n", err)
        return
    }
    var res0 endpoints.BlogPostJSON
    err = json.Unmarshal(body, &res0)
    if err != nil {
        t.Fatalf("FAILED: unable to decode JSON %v\n", err)
        return
    }

    if res0.BlogID != 3 {
        t.Fatalf("FAILED: expected BlogID 3 but got %d", res0.BlogID)
    }
    if res0.Author != "Foo Bar" {
        t.Fatalf("FAILED: expected Author `Foo Bar` but got `%s`", res0.Author)
    }
    if res0.Date != "2021-05-30" {
        t.Fatalf("FAILED: expected Date `2021-05-30` but got `%s`", res0.Date)
    }
    if res0.Duration != 109 {
        t.Fatalf("FAILED: expected Duration 109 but got %d", res0.Duration)
    }
    if res0.URL != "https://mail.google.com/" {
        t.Fatalf("FAILED: expected URL `https://mail.google.com/` but got `%s`", string(res0.URL))
    }
    if res0.Content != "foo bar baz" {
        t.Fatalf("FAILED: expected Content `foo bar baz` but got `%s`", string(res0.Content))
    }

    t.Logf("endpoints.InsertNewBlogPost test passed")
}