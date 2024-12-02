package main

import (
	"html/template"
	"log"
	"net/http"
	"sync"
	"encoding/json"
)

// Message represents a DM between two users.
type Message struct {
	Sender   string
	Receiver string
	Content  string
}

// In-memory message storage
var messages = struct {
	sync.RWMutex
	data map[string][]Message
}{
	data: make(map[string][]Message),
}

// UsersHandler fetches the list of all users.
func UsersHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var users []User
    if err := forumDB.Find(&users).Error; err != nil {
        http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}


// SendMessageHandler allows a user to send a message to another user.
func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		// Retrieve sender's username from the cookie
		senderCookie, err := r.Cookie("username")
		if err != nil || senderCookie.Value == "" {
			http.Error(w, "User not logged in or invalid session", http.StatusUnauthorized)
			log.Println("Error: Missing or invalid username cookie")
			return
		}
		sender := senderCookie.Value

		// Get receiver's username and message content from the form
		receiver := r.FormValue("receiver")
		content := r.FormValue("content")

		if receiver == "" || content == "" {
			http.Error(w, "Receiver and content are required", http.StatusBadRequest)
			log.Println("Error: Missing receiver or content")
			return
		}

		// Save the message in memory
		message := Message{Sender: sender, Receiver: receiver, Content: content}
		messages.Lock()
		messages.data[sender+"-"+receiver] = append(messages.data[sender+"-"+receiver], message)
		messages.data[receiver+"-"+sender] = append(messages.data[receiver+"-"+sender], message)
		messages.Unlock()

		// Redirect to the messages page
		http.Redirect(w, r, "/messages?receiver="+receiver, http.StatusSeeOther)
	}
}

// ViewMessagesHandler displays messages between two users.
func ViewMessagesHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve sender's username from the cookie
	senderCookie, err := r.Cookie("username")
	if err != nil || senderCookie.Value == "" {
		http.Error(w, "User not logged in or invalid session", http.StatusUnauthorized)
		log.Println("Error: Missing or invalid username cookie")
		return
	}
	sender := senderCookie.Value

	// Get receiver's username from the query parameters
	receiver := r.URL.Query().Get("receiver")
	if receiver == "" {
		http.Error(w, "Receiver username is missing", http.StatusBadRequest)
		log.Println("Error: Missing receiver username in query")
		return
	}

	// Fetch messages between the sender and receiver from memory
	messages.RLock()
	chatMessages := messages.data[sender+"-"+receiver]
	messages.RUnlock()

	// Render the messages page
	tmpl, err := template.ParseFiles("messages.html")
	if err != nil {
		http.Error(w, "Failed to load template", http.StatusInternalServerError)
		log.Println("Error: Failed to parse messages.html template")
		return
	}
	tmpl.Execute(w, struct {
		Messages        []Message
		ReceiverUsername string
	}{Messages: chatMessages, ReceiverUsername: receiver})
}
