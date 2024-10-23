package main

import (
	"html/template"
	"log"
	"net/http"
	"fmt"
)

// Data structure for topics
type Topic struct {
	Title string
	Content string
	Comments []string
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

	// Search for the topic by title and pass both title and content
	for _, topic := range topics {
		if topic.Title == title {
			tmpl, err := template.ParseFiles("topic.html")
			if err != nil {
				log.Fatal(err)
			}
			tmpl.Execute(w, topic)
			return
		}
	}

	http.Error(w, "Topic not found", http.StatusNotFound)
}

// NewTopicHandler adds a new topic to the list
func NewTopicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// Get title and content from the form
		title := r.FormValue("title")
		content := r.FormValue("content")

		// Check if title and content are not empty
		if title != "" && content != "" {
			// Append the new topic to the topics list
			topics = append(topics, Topic{Title: title, Content: content})

			// Redirect back to the homepage
			http.Redirect(w, r, "/forum/", http.StatusSeeOther)
		} else {
			http.Error(w, "Title and content are required", http.StatusBadRequest)
		}
	}
}

// NewCommentHandler adds a new comment to a post
func NewCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// Get title and comment from the form
		title := r.FormValue("title")
		comment := r.FormValue("comment")

		// Check if title and comment are not empty
		if title != "" && comment != "" {
			// Append the new comment to the comments list
			for _, topic := range topics {
				tmpl, err := template.ParseFiles("topic.html")
				if err != nil {
					log.Fatal(err)
				}
				if topic.Title == title {
					topic.Comments = append(topic.Comments, comment)
				}
				tmpl.Execute(w, topic.Comments)
			}
			url := fmt.Sprintf("/topic?title=%s", title)
			http.Redirect(w, r, url, http.StatusSeeOther)
		} else {
			http.Error(w, "Comment is empty", http.StatusBadRequest)
		}
	}
}
