package main

//import "C"

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// User represents a user with a unique ID, username, and password.
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string
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
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login_or_register.dart")
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

	if r.Method == "GET" {
		t, _ := template.ParseFiles("login_or_register.dart")
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		log.Printf("Received username: %s, password: %s", username, password)

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

		sessionID := fmt.Sprintf("%d", time.Now().UnixNano())

		http.SetCookie(w, &http.Cookie{
			Name:  "session_token",
			Value: "authenticated",
			Path:  "/",
			//SameSite: http.SameSiteNoneMode, // Allow cross-origin requests
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "is_teacher",
			Value:    request.IsTeacher,
			Expires:  time.Now().Add(24 * time.Hour),
			Path:     "/",
			HttpOnly: true,
			//SameSite: http.SameSiteNoneMode,
			//Secure:   false,
		})

		http.SetCookie(w, &http.Cookie{
			Name:    "session_id",
			Value:   sessionID,
			Expires: time.Now().Add(24 * time.Hour),
			Path:    "/",
			//SameSite: http.SameSiteNoneMode,
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

	http.SetCookie(w, &http.Cookie{
		Name:  "is_teacher",
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

// ProfileHandler handles user profile.
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	usernameCookie, err1 := r.Cookie("username")
	emailCookie, err2 := r.Cookie("email")

	if err1 != nil || err2 != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	data := struct {
		Username string
		Email    string
	}{
		Username: usernameCookie.Value,
		Email:    emailCookie.Value,
	}

	t, _ := template.ParseFiles("profile.html")
	t.Execute(w, data)
}

// EditProfileHandler renders the edit profile page.
func EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	usernameCookie, err1 := r.Cookie("username")
	emailCookie, err2 := r.Cookie("email")

	if err1 != nil || err2 != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	data := struct {
		Username string
		Email    string
	}{
		Username: usernameCookie.Value,
		Email:    emailCookie.Value,
	}

	t, err := template.ParseFiles("edit_profile.html")
	if err != nil {
		http.Error(w, "Error loading edit profile page", http.StatusInternalServerError)
		return
	}

	t.Execute(w, data)

}

// UpdateProfileHandler updates user profile.
func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	usernameCookie, err1 := r.Cookie("username")
	emailCookie, err2 := r.Cookie("email")

	if err1 != nil || err2 != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	currentUsername := usernameCookie.Value
	currentEmail := emailCookie.Value

	newUsername := r.FormValue("username")
	newEmail := r.FormValue("email")

	var user User
	if err := forumDB.Where("username = ? AND email = ?", currentUsername, currentEmail).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if the new username or email already exists
	var existingUser User
	if err := forumDB.Where("username = ? OR email = ?", newUsername, newEmail).First(&existingUser).Error; err == nil && existingUser.ID != user.ID {
		http.Error(w, "Username or email already in use. Please choose another.", http.StatusConflict)
		return
	}

	// Update user information
	user.Username = newUsername
	user.Email = newEmail

	if err := forumDB.Save(&user).Error; err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	// Update cookies with new username and email
	http.SetCookie(w, &http.Cookie{
		Name:    "username",
		Value:   newUsername,
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "email",
		Value:   newEmail,
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
	})

	http.Redirect(w, r, "/view/", http.StatusFound)
}
