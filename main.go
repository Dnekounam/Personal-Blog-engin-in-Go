package main

import (
	"fmt"
	"html/template"
	"net/http"
)

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

// this will be by a database later
func PostHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Path[len("/post/"):]

	post := Post{
		Title:   "Post" + postID,
		Content: "This is a content of post" + postID,
	}
	tmpl, err := template.ParseFiles("templates/post.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, post)

}

func main() {
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/post/", PostHandler)

	fmt.Println("Starting server on :9092")
	http.ListenAndServe(":9092", nil)
}
