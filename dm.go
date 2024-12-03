package main

import (
	"encoding/json" // For JSON encoding
	"net/http"      // For HTTP handling
	"time"          // For timestamps
	"log"

	//"gorm.io/gorm" // For GORM tags and database operations
)

// Message struct represents a DM between users
type Message struct {
	ID        uint      `gorm:"primaryKey"`
	Sender    string    // Username of the sender
	Recipient string    // Username of the recipient
	Content   string    // Message content
	Timestamp time.Time // When the message was sent
}


func FetchUsersHandler(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("username")
    if err != nil {
        log.Println("Error fetching username cookie:", err)
        http.Error(w, "User not authenticated", http.StatusUnauthorized)
        return
    }
    log.Println("Fetched username from cookie:", cookie.Value)

    var users []User
    if err := forumDB.Where("username != ?", cookie.Value).Find(&users).Error; err != nil {
        log.Println("Error fetching users from database:", err)
        http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
        return
    }

    log.Println("Users fetched for DMs:", users)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}


func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Parse sender's username from the cookie
    senderCookie, err := r.Cookie("username")
    if err != nil {
        log.Println("Error fetching username cookie:", err)
        http.Error(w, "User not authenticated", http.StatusUnauthorized)
        return
    }
    sender := senderCookie.Value
    log.Println("Sender username:", sender)

    // Parse request body for recipient and content
    var messageData struct {
        Recipient string `json:"recipient"`
        Content   string `json:"content"`
    }
    if err := json.NewDecoder(r.Body).Decode(&messageData); err != nil {
        log.Println("Error decoding request body:", err)
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    log.Printf("Received message data: %+v\n", messageData)

    // Validate the input
    if messageData.Recipient == "" || messageData.Content == "" {
        log.Println("Recipient or content missing")
        http.Error(w, "Recipient and content are required", http.StatusBadRequest)
        return
    }

    // Create and save the message
    message := Message{
        Sender:    sender,
        Recipient: messageData.Recipient,
        Content:   messageData.Content,
        Timestamp: time.Now(),
    }
    if err := forumDB.Create(&message).Error; err != nil {
        log.Println("Error saving message to database:", err)
        http.Error(w, "Failed to send message", http.StatusInternalServerError)
        return
    }

    log.Println("Message sent successfully:", message)
    w.WriteHeader(http.StatusOK)
}

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	sender, err := r.Cookie("username")
	if err != nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	recipient := r.URL.Query().Get("recipient")
	if recipient == "" {
		http.Error(w, "Recipient is required", http.StatusBadRequest)
		return
	}

	var messages []Message
	if err := forumDB.
		Where("(sender = ? AND recipient = ?) OR (sender = ? AND recipient = ?)", sender.Value, recipient, recipient, sender.Value).
		Order("timestamp").
		Find(&messages).Error; err != nil {
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(messages)
}
