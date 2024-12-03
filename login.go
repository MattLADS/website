package main

import (
	//"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	"io"
	"bytes"
)

// User represents a user with a unique ID, username, and password.
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string
	//flag	 bool
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("username")
        if err != nil || cookie.Value == "" {
            log.Println("Unauthorized access, redirecting to sign-in")
            http.Redirect(w, r, "/signin/", http.StatusFound)
            return
        }
        log.Println("Middleware passed, username cookie:", cookie.Value)
        next(w, r)
    }
}


// SignUpHandler handles user registration.
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("signup.html")
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var existingUser User
		if err := forumDB.Where("username = ?", username).First(&existingUser).Error; err == nil {
			http.Error(w, "Username already exists. Please choose another one.", http.StatusConflict)
			return
		}

		newUser := User{Username: username, Password: password}
		if err := forumDB.Create(&newUser).Error; err != nil {
			http.Error(w, "Failed to create user.", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("SignInHandler called")
    log.Printf("Request method: %s", r.Method)

    if r.Method == "GET" {
        log.Println("Handling GET request for SignInHandler")
        t, err := template.ParseFiles("signin.html")
        if err != nil {
            log.Println("Error parsing signin.html:", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }
        t.Execute(w, nil)
        return
    }

    if r.Method == "POST" {
        // Log the raw request body for debugging
        bodyBytes, err := io.ReadAll(r.Body)
        if err != nil {
            log.Println("Error reading request body:", err)
        } else {
            log.Println("Raw Request Body:", string(bodyBytes))
        }
        r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Reset the body so it can be re-read

        // Parse the form data
        if err := r.ParseForm(); err != nil {
            log.Println("Error parsing form data:", err)
            http.Error(w, "Invalid form submission", http.StatusBadRequest)
            return
        }

        username := r.FormValue("username")
        password := r.FormValue("password")
        log.Printf("Received username: '%s', password: '%s'", username, password)

        if username == "" || password == "" {
            log.Println("Empty username or password")
            http.Error(w, "Username and password are required", http.StatusBadRequest)
            return
        }

        var user User
        if err := forumDB.Where("username = ? AND password = ?", username, password).First(&user).Error; err != nil {
            log.Println("Invalid credentials for username:", username)
            http.Error(w, "Invalid username or password", http.StatusUnauthorized)
            return
        }

        log.Println("Login successful, setting cookies")

        http.SetCookie(w, &http.Cookie{
            Name:  "session_token",
            Value: "authenticated",
            Path:  "/",
        })

        http.SetCookie(w, &http.Cookie{
            Name:    "username",
            Value:   username,
            Expires: time.Now().Add(24 * time.Hour),
            Path:    "/",
        })

        log.Println("Set session_token and username cookies")
        http.Redirect(w, r, "/view/", http.StatusFound)
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

	http.Redirect(w, r, "/", http.StatusFound)
}
