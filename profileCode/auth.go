package main

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// place holder user database integate w/ Kunj code
var users = map[string]string{
	"user1": "$2a$10$Vb.kkGr6LCHeUvg..g6xaexM8Bw65yNUs32k.bS6rNrT5YXSmLlw6", // password: password123
}

// Validate user credentials
func validateCredentials(username, password string) bool {
	storedPassword, ok := users[username]
	if !ok {
		return false
	}

	// Compare hashed password
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	return err == nil
}

// Sign-in handler with validation
func signInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Render sign-in form
		http.ServeFile(w, r, "templates/signin.html")
		return
	}

	// Get username and password from form values
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Validate credentials
	if !validateCredentials(username, password) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Create session if valid
	sessionID := fmt.Sprintf("%d", time.Now().UnixNano())
	createSession(sessionID, username)

	// Set a cookie with session ID
	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
	})

	http.Redirect(w, r, "/profile", http.StatusFound)
}
