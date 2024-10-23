package main

import (
	"html/template"
	"log"
	"net/http"
)

// Data structure for topics
type Topic struct {
	Title string
}

// List of topics
var topics []Topic

// ForumHandler serves the forum page
func ForumHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("forum.html")
	if err != nil {
		log.Fatal(err)
	}

	// Pass topics data to the template
	tmpl.Execute(w, topics)
}

// TopicHandler serves a specific topic page
func TopicHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")

	if title != "" {
		tmpl, err := template.ParseFiles("topic.html")
		if err != nil {
			log.Fatal(err)
		}

		tmpl.Execute(w, struct{ Title string }{Title: title})
	} else {
		http.Error(w, "Topic not found", http.StatusNotFound)
	}
}

// NewTopicHandler adds a new topic to the list
func NewTopicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		title := r.FormValue("title")
		if title != "" {
			topics = append(topics, Topic{Title: title})
			http.Redirect(w, r, "/forum/", http.StatusSeeOther)
		} else {
			http.Error(w, "Title is required", http.StatusBadRequest)
		}
	}
}
