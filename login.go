package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// User structure represents a user with a unique ID, username, and password.
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string
}

var db *gorm.DB

// InitializeDB sets up the database connection and creates the users table if it doesn't exist.
func InitializeDB() {
	var err error
	db, err = gorm.Open("sqlite3", "users.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Enable foreign key constraints.
	db.Exec("PRAGMA foreign_keys = ON;")

	// Automatically migrate the User struct to create/update the users table.
	db.AutoMigrate(&User{})
}

// authMiddleware is a middleware function that checks if the user is authenticated.
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")

		if err != nil || cookie.Value != "authenticated" {
			http.Redirect(w, r, "/", http.StatusFound) // Redirect to sign-in if not authenticated.
			return
		}

		next(w, r) // Proceed to the next handler if authenticated.
	}
}

// SignUpHandler handles user registration (sign-up).
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("signup.html")
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Check if the username already exists.
		var existingUser User
		if err := db.Where("username = ?", username).First(&existingUser).Error; err == nil {
			http.Error(w, "Username already exists. Please choose another one.", http.StatusConflict)
			return
		}

		// Create a new user record.
		newUser := User{Username: username, Password: password}
		if err := db.Create(&newUser).Error; err != nil {
			http.Error(w, "Failed to create user.", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// SignInHandler handles user sign-in.
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("SignInHandler called")
	log.Printf("Request method: %s", r.Method)

	if r.Method == "GET" {
		t, _ := template.ParseFiles("signin.html")
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		log.Printf("Received username: %s, password: %s", username, password)

		var user User
		if err := db.Where("username = ? AND password = ?", username, password).First(&user).Error; err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Create session if valid
		sessionID := fmt.Sprintf("%d", time.Now().UnixNano())
		log.Println("Login successful, setting session token")

		http.SetCookie(w, &http.Cookie{
			Name:  "session_token",
			Value: "authenticated",
			Path:  "/",
		})

		// Set a cookie with session ID
		http.SetCookie(w, &http.Cookie{
			Name:    "session_id",
			Value:   sessionID,
			Expires: time.Now().Add(24 * time.Hour),
			Path:    "/",
		})

		http.SetCookie(w, &http.Cookie{
			Name:    "username",
			Value:   username,
			Expires: time.Now().Add(24 * time.Hour),
			Path:    "/",
		})

		http.Redirect(w, r, "/view/", http.StatusFound)
	}
}

// SignOutHandler handles user sign-out by clearing the session cookie.
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

	http.Redirect(w, r, "/", http.StatusFound)
}
