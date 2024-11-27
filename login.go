package main

//import "C"

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// User represents a user with a unique ID, username, and password.
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// authMiddleware checks if the user is authenticated.
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil || cookie.Value != "authenticated" {
			http.Redirect(w, r, "/", http.StatusFound) // Redirect to sign-in if not authenticated.
			return
		}
		next(w, r)
	}
}

// SignUpHandler handles user registration.
func SignUpHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("SignUpHandler called")
	if r.Method == "POST" {
		var request User
		//username := r.FormValue("username")
		//password := r.FormValue("password")

		// Decode JSON body
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
			return
		}
		// Log received username and password for verification
		log.Printf("Received login request - Username: %s, Password: %s\n", request.Username, request.Password)

		//checking user info in database
		var existingUser User
		if err := forumDB.Where("username = ?", request.Username).First(&existingUser).Error; err == nil {
			//http.Error(w, "Username already exists. Please choose another one.", http.StatusConflict)
			// Send JSON response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(`{"error": "Username already exists. Please choose another one."}`))
			return
		}

		newUser := User{Username: request.Username, Password: request.Password}
		if err := forumDB.Create(&newUser).Error; err != nil {
			http.Error(w, "Failed to create user.", http.StatusInternalServerError)
			return
		}
		//changed redirect to send information in JSON format
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": " Created user successfully"}`))
	}
}

// SignInHandler handles user sign-in.
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("SignInHandler called")
	log.Printf("Request method: %s", r.Method)

	if r.Method == "POST" {
		var request struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		// Decode JSON body
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Println("Error decoding JSON:", err)
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
			return
		}

		log.Printf("Received username: %s, password: %s", request.Username, request.Password)

		//checking user info in database
		var user User
		if err := forumDB.Where("username = ? AND password = ?", request.Username, request.Password).First(&user).Error; err != nil {
			log.Println("Invalid username or password")
			// Send JSON response for invalid credentials
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			//w.Write([]byte(`{"error": "Invalid username and password. Try again."}`))
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid username and password. Try again."})

			return

		}

		// Create session if valid
		sessionID := fmt.Sprintf("%d", time.Now().UnixNano())
		log.Println("Login successful, setting session token")

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "authenticated",
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode, // Allow cross-origin requests
			Secure:   false,                 // Set to true if using HTTPS
		})

		// Set a cookie with session ID
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Expires:  time.Now().Add(24 * time.Hour),
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
			Secure:   false,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "username",
			Value:    request.Username,
			Expires:  time.Now().Add(24 * time.Hour),
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
			Secure:   false,
		})
		log.Printf("Set-Cookie header for username: %s", request.Username)

		log.Println("SignUpHandler: Sending response")
		//changed redirect to send information in JSON format
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Login was successful"}`))
	}
}

// SignOutHandler handles user sign-out.
func SignOutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	// Set a cookie with session ID
	http.SetCookie(w, &http.Cookie{
		Name:  "session_id",
		Value: "",
		Path:  "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:  "username",
		Value: "",
		Path:  "/",
	})

	//changed redirect to send information in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Sign out was successful"}`))
}

// AuthStatusHandler checks if the user is logged in by verifying the session token.
func AuthStatusHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("AuthStatusHandler called")

	// Check if session_token cookie exists and is set to "authenticated"
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value != "authenticated" {
		// If no valid session, return 401 Unauthorized
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// If valid session, return 200 OK
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Session is active"}`))
}
