package main

//import "C"
import (
	"encoding/json"
	"fmt"
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
	log.Println("ForumHandler called")
	log.Printf("Request method: %s", r.Method)
	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	var topics []Topic
	// Retrieve topics from the database
	if err := forumDB.Preload("Comments").Find(&topics).Error; err != nil {
		log.Printf("Error fetching topics: %v", err)
		http.Error(w, "Error fetching topics", http.StatusInternalServerError)
		return
	}

	log.Printf("Number of topics retrieved: %d", len(topics)) // Log the count

	if len(topics) == 0 {
		log.Println("No topics found, returning an empty array.")
	}

	// Encode topics to JSON and send as response
	if err := json.NewEncoder(w).Encode(topics); err != nil {
		http.Error(w, "Error encoding topics to JSON", http.StatusInternalServerError)
	}

}

// TopicHandler serves a specific topic page along with its comments.
func TopicHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	title := r.URL.Query().Get("title")

	// Retrieve the topic by title from forumDB
	var topic Topic
	err := forumDB.Preload("Comments").Where("title = ?", title).First(&topic).Error
	if err != nil {
		http.Error(w, "Topic not found", http.StatusNotFound)
		return
	}

	// Convert topic to JSON and send response
	if err := json.NewEncoder(w).Encode(topic); err != nil {
		http.Error(w, `{"error": "Failed to encode topic"}`, http.StatusInternalServerError)
	}
}

// NewTopicHandler adds a new topic to the database.
func NewTopicHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("NewTopicHandler called")
	log.Printf("Request method: %s", r.Method)

	//checking that request is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("Parsing JSON body...")

	// Decode the JSON body
	var parseJSON struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	log.Println("JSON decoded successfully")

	if err := json.NewDecoder(r.Body).Decode(&parseJSON); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Get title and content from the request
	title := parseJSON.Title
	content := parseJSON.Content

	// Get username from the cookie
	cookie, err := r.Cookie("username")
	if err != nil {
		log.Println("Failed to retrieve username cookie:", err)
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}
	username := cookie.Value

	log.Printf("NewTopicHandler called with title: %s, content: %s, username: %s", parseJSON.Title, parseJSON.Content, username)

	// Check if title, content, and username are valid
	if title == "" || content == "" || username == "" {
		http.Error(w, "Title, content, and username are required", http.StatusBadRequest)
		return
	}

	// Create and save the new topic
	newTopic := Topic{Title: title, Content: content, Username: username}
	log.Printf("Creating topic with Title: %s, Content: %s, Username: %s", title, content, username)
	if err := forumDB.Create(&newTopic).Error; err != nil {
		log.Printf("Error saving topic to database: %v", err)
		http.Error(w, "Failed to create topic", http.StatusInternalServerError)
		return
	}
	log.Println("Topic saved to database successfully")

	// Send a success response
	log.Println("Topic created successfully, sending response")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Topic created successfully"}`))
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
