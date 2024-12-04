package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

// User represents a user with a unique ID, username, and password.
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique"`
	Email     string `gorm:"unique"`
	Password  string
	FirstName string
	LastName  string
	//flag	 bool
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
		t, _ := template.ParseFiles("signup.html")
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		var existingUser User
		if err := forumDB.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
			http.Error(w, "Username or email already exists. Please choose another one.", http.StatusConflict)
			return
		}

		newUser := User{
			FirstName: firstName,
			LastName:  lastName,
			Username:  username,
			Email:     email,
			Password:  password,
		}
		if err := forumDB.Create(&newUser).Error; err != nil {
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
		if err := forumDB.Where("username = ? AND password = ?", username, password).First(&user).Error; err != nil {
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

		http.SetCookie(w, &http.Cookie{
			Name:    "email",
			Value:   user.Email,
			Expires: time.Now().Add(24 * time.Hour),
			Path:    "/",
		})

		http.SetCookie(w, &http.Cookie{
			Name:    "first_name",
			Value:   user.FirstName,
			Expires: time.Now().Add(24 * time.Hour),
			Path:    "/",
		})

		http.SetCookie(w, &http.Cookie{
			Name:    "last_name",
			Value:   user.LastName,
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

func EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	firstNameCookie, err1 := r.Cookie("first_name")
	lastNameCookie, err2 := r.Cookie("last_name")
	usernameCookie, err3 := r.Cookie("username")
	emailCookie, err4 := r.Cookie("email")

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	data := struct {
		FirstName string
		LastName  string
		Username  string
		Email     string
	}{
		FirstName: firstNameCookie.Value,
		LastName:  lastNameCookie.Value,
		Username:  usernameCookie.Value,
		Email:     emailCookie.Value,
	}

	t, err := template.ParseFiles("edit_profile.html")
	if err != nil {
		http.Error(w, "Error loading edit profile page", http.StatusInternalServerError)
		return
	}

	t.Execute(w, data)
}

func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	currentUsername, _ := r.Cookie("username")
	currentEmail, _ := r.Cookie("email")

	newUsername := r.FormValue("username")
	newEmail := r.FormValue("email")
	newFirstName := r.FormValue("first_name")
	newLastName := r.FormValue("last_name")

	var user User
	if err := forumDB.Where("username = ? AND email = ?", currentUsername.Value, currentEmail.Value).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Update the user's username and email
	user.Username = newUsername
	user.Email = newEmail
	user.FirstName = newFirstName
	user.LastName = newLastName

	if err := forumDB.Save(&user).Error; err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	// Update cookies
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

	http.SetCookie(w, &http.Cookie{
		Name:    "first_name",
		Value:   newFirstName,
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "last_name",
		Value:   newLastName,
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
	})

	http.Redirect(w, r, "/profile/", http.StatusFound)
}
