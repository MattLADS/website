package main

//import "C"
import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Topic represents a forum topic with related comments.
type Topic struct {
	gorm.Model
	Title    string
	Username string
	Content  string
	Comments []Comment `gorm:"foreignKey:TopicID"`
}

// Comment represents a comment on a forum topic.
type Comment struct {
	gorm.Model
	TopicID  uint
	Content  string
	Username string
}

var forumDB *gorm.DB

// InitializeForumDB sets up the database connection and creates the tables if they don't exist.
func InitializeForumDB() {
	var err error
	forumDB, err = gorm.Open(sqlite.Open("forum.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to forum database:", err)
	}

	// Enable foreign key constraints.
	forumDB.Exec("PRAGMA foreign_keys = ON;")

	// Migrate the User, Topic, and Comment tables.
	forumDB.AutoMigrate(&User{}, &Topic{}, &Comment{})
}

// ForumHandler serves the forum page with a list of topics.
func ForumHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("forum.html")
	if err != nil {
		log.Fatal(err)
	}

	var topics []Topic
	// Retrieve topics from the database
	forumDB.Preload("Comments").Find(&topics)

	// Pass topics data to the template
	tmpl.Execute(w, topics)
}

// TopicHandler serves a specific topic page along with its comments.
func TopicHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")

	// Retrieve the topic by title from forumDB
	var topic Topic
	err := forumDB.Preload("Comments").Where("title = ?", title).First(&topic).Error
	if err != nil {
		http.Error(w, "Topic not found", http.StatusNotFound)
		return
	}

	// Render the topic page with the retrieved topic and its comments
	tmpl, err := template.ParseFiles("topic.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(w, topic)
}

// NewTopicHandler adds a new topic to the database.
func NewTopicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		comments := []Comment{}

		cookie, err := r.Cookie("username")
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		username := cookie.Value

		if title != "" && content != "" && username != "" {
			newTopic := Topic{Title: title, Content: content, Username: username, Comments: comments}
			forumDB.Create(&newTopic)

			http.Redirect(w, r, "/forum/", http.StatusSeeOther)
		} else {
			http.Error(w, "Title, content, and username are required", http.StatusBadRequest)
		}
	}
}

// NewCommentHandler adds a new comment to a topic.
func NewCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		comment := r.FormValue("comment")

		cookie, err := r.Cookie("username")
		if err != nil {
			log.Fatal(err)
		}
		username := cookie.Value

		if title != "" && comment != "" && username != "" {
			// Find the topic by title in forumDB
			var topic Topic
			err := forumDB.Where("title = ?", title).First(&topic).Error
			if err != nil {
				http.Error(w, "Topic not found", http.StatusNotFound)
				return
			}
			// Create and save the new comment associated with the topic's ID
			newComment := Comment{TopicID: topic.ID, Content: comment, Username: username}
			forumDB.Create(&newComment)

			url := fmt.Sprintf("/topic?title=%s", title)
			http.Redirect(w, r, url, http.StatusSeeOther)
		} else {
			http.Error(w, "Comment, title, and username are required", http.StatusBadRequest)
		}
	}
}
