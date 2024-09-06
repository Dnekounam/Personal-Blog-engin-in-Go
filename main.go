package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

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

const PostPerPage = 5

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}

	offset := (pageNum - 1) * PostPerPage

	rows, err := db.Query("SELECT title, content FROM posts LIMIT ? OFFSET ?", PostPerPage, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.Title, &post.Content); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)

	}
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct {
		Posts       []Post
		CurrentPage int
		NextPage    int
		PrevPage    int
	}{
		Posts:       posts,
		CurrentPage: pageNum,
		NextPage:    pageNum + 1,
		PrevPage:    pageNum - 1,
	}

	tmpl.Execute(w, data)

}

type Post struct {
	Title   string
	Content string
}

func NewPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/new_post.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	}
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")

		_, err := db.Exec("INSERT INTO posts (title, content) VALUES (?, ?)", title, content)
		if err != nil {
			http.Error(w, "Unable to create post", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
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
	http.HandleFunc("/new1", NewPostHandler)
	http.HandleFunc("/new", CreatePostHandler)

	fmt.Println("Starting server on :9092")
	http.ListenAndServe(":9092", nil)
}
