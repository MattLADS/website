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
	Email 	string `gorm:"unique"`
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


// // SignUpHandler handles user registration.
// func SignUpHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "GET" {
// 		t, _ := template.ParseFiles("login_or_register.dart")
// 		t.Execute(w, nil)
// 	} else if r.Method == "POST" {
// 		username := r.FormValue("username")
// 		password := r.FormValue("password")
// 
// 		var existingUser User
// 		if err := forumDB.Where("username = ?", username).First(&existingUser).Error; err == nil {
// 			http.Error(w, "Username already exists. Please choose another one.", http.StatusConflict)
// 			return
// 		}
// 
// 		newUser := User{Username: username, Password: password}
// 		if err := forumDB.Create(&newUser).Error; err != nil {
// 			http.Error(w, "Failed to create user.", http.StatusInternalServerError)
// 			return
// 		}
// 
// 		http.Redirect(w, r, "/", http.StatusFound)
// 	}
// }

// SignUpHandler handles user registration.
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        // Serve the sign-up page (HTML file)
        t, err := template.ParseFiles("signup.html")
        if err != nil {
            log.Println("Error parsing signup.html:", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }
        t.Execute(w, nil)
        return
    }

    if r.Method == http.MethodPost {
        // Parse form data for sign-up
        username := r.FormValue("username")
        password := r.FormValue("password")

        if username == "" || password == "" {
            log.Println("Username or password missing during signup")
            http.Error(w, "Username and password are required", http.StatusBadRequest)
            return
        }

        // Check if the username already exists
        var existingUser User
        if err := forumDB.Where("username = ?", username).First(&existingUser).Error; err == nil {
            log.Println("Attempt to register an existing username:", username)
            http.Error(w, "Username already exists", http.StatusConflict)
            return
        }

        // Add the new user to the database
        newUser := User{
            Username: username,
            Password: password, // In production, hash the password before storing it
        }

        if err := forumDB.Create(&newUser).Error; err != nil {
            log.Println("Error saving new user to database:", err)
            http.Error(w, "Failed to create user", http.StatusInternalServerError)
            return
        }

        log.Println("New user registered successfully:", username)

        // Automatically sign in the user after successful registration
        http.SetCookie(w, &http.Cookie{
            Name:    "username",
            Value:   username,
            Path:    "/",
            Expires: time.Now().Add(24 * time.Hour),
        })

        // Redirect the user to the DMs or dashboard
        http.Redirect(w, r, "/view/", http.StatusSeeOther)
    }
}


// SignInHandler handles user sign-in.
func SignInHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("SignInHandler called")
    log.Printf("Request method: %s", r.Method)

    // Handle GET requests
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

    // Handle POST requests
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
