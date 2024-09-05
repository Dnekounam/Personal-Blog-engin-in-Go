package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "blog.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT,
        content TEXT
    )`)
	if err != nil {
		log.Fatal(err)
	}
}

func fetchPost(id string) (Post, error) {
	var post Post
	err := db.QueryRow("SELECT title, content FROM posts WHERE id = ?", id).Scan(&post.Title, &post.Content)
	return post, err
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

type Post struct {
	Title   string
	Content string
}

var posts = map[string]Post{
	"1": {Title: "First Post", Content: "This is the first post."},
	"2": {Title: "Second Post", Content: "This is the second post."},
}

// this will be by a database later
func PostHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Path[len("/post/"):]

	post, err := fetchPost(postID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("templates/post.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, post)

}

func main() {

	initDB()
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/post/", PostHandler)

	fmt.Println("Starting server on :9092")
	http.ListenAndServe(":9092", nil)
}
