package main

//import "C"
//

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

// User represents a user with a unique ID, username, and password.
type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
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
		email := r.FormValue("email")
		password := r.FormValue("password")

		var existingUser User
		if err := forumDB.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
			http.Error(w, "Username or email already exists. Please choose another one.", http.StatusConflict)
			return
		}

		newUser := User{Username: username, Email: email, Password: password}
		if err := forumDB.Create(&newUser).Error; err != nil {
			http.Error(w, "Failed to create user.", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// SignInHandler handles user sign-in.
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login_or_register.dart")
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var user User
		if err := forumDB.Where("username = ? AND password = ?", username, password).First(&user).Error; err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		sessionID := fmt.Sprintf("%d", time.Now().UnixNano())

		http.SetCookie(w, &http.Cookie{
			Name:  "session_token",
			Value: "authenticated",
			Path:  "/",
		})

		http.SetCookie(w, &http.Cookie{
			Name:    "session_id",
			Value:   sessionID,
			Expires: time.Now().Add(24 * time.Hour),
			Path:    "/",
		})

		http.SetCookie(w, &http.Cookie{
			Name:    "username",
			Value:   user.Username,
			Expires: time.Now().Add(24 * time.Hour),
			Path:    "/",
		})

		http.SetCookie(w, &http.Cookie{
			Name:    "email",
			Value:   user.Email,
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

	http.Redirect(w, r, "/", http.StatusFound)
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
