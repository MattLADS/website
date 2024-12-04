package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
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

// validateEmailDomain checks if the email's domain has valid MX records.
func validateEmailDomain(email string) error {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return fmt.Errorf("invalid email format")
	}

	_, err := net.LookupMX(parts[1])
	if err != nil {
		return fmt.Errorf("invalid email domain: %v", err)
	}
	return nil
}

// validatePassword checks if the password meets the specified requirements.
func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one special character")
	}
	return nil
}

// SignUpHandler handles user registration.
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        t, err := template.ParseFiles("signup.html")
        if err != nil {
            log.Println("Error parsing signup.html:", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }
        t.Execute(w, nil)
        return
    }

    if r.Method == "POST" {
        firstName := r.FormValue("first_name")
        lastName := r.FormValue("last_name")
        username := r.FormValue("username")
        email := r.FormValue("email")
        password := r.FormValue("password")
        confirmPassword := r.FormValue("confirm_password")

        errorMessage := ""

        // Check if passwords match
        if password != confirmPassword {
            errorMessage = "Passwords do not match."
        }

        // Check for empty fields
        if firstName == "" || lastName == "" || username == "" || email == "" || password == "" || confirmPassword == "" {
            errorMessage = "All fields are required."
        } else if errorMessage == "" {
            // Validate the email domain
            if err := validateEmailDomain(email); err != nil {
                errorMessage = "Invalid email address."
            } else if err := validatePassword(password); err != nil {
                errorMessage = err.Error()
            } else {
                // Check if the username or email already exists
                var existingUser User
                if err := forumDB.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
                    errorMessage = "Username or email already exists. Please choose another."
                }
            }
        }

        // If there is an error, re-render the signup page with the error message
        if errorMessage != "" {
            t, err := template.ParseFiles("signup.html")
            if err != nil {
                log.Println("Error parsing signup.html:", err)
                http.Error(w, "Internal server error", http.StatusInternalServerError)
                return
            }

            // Pass the error message to the template
            data := struct {
                ErrorMessage string
                FirstName    string
                LastName     string
                Username     string
                Email        string
            }{
                ErrorMessage: errorMessage,
                FirstName:    firstName,
                LastName:     lastName,
                Username:     username,
                Email:        email,
            }

            t.Execute(w, data)
            return
        }

        // If no errors, create the new user
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

// func SignUpHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "GET" {
// 		t, err := template.ParseFiles("signup.html")
// 		if err != nil {
// 			log.Println("Error parsing signup.html:", err)
// 			http.Error(w, "Internal server error", http.StatusInternalServerError)
// 			return
// 		}
// 		t.Execute(w, nil)
// 		return
// 	}
// 
// 	if r.Method == "POST" {
// 		firstName := r.FormValue("first_name")
// 		lastName := r.FormValue("last_name")
// 		username := r.FormValue("username")
// 		email := r.FormValue("email")
// 		password := r.FormValue("password")
// 
// 		errorMessage := ""
// 
// 		// Check for empty fields
// 		if firstName == "" || lastName == "" || username == "" || email == "" || password == "" {
// 			errorMessage = "All fields are required."
// 		} else {
// 			// Validate the email domain
// 			if err := validateEmailDomain(email); err != nil {
// 				errorMessage = "Invalid email address."
// 			} else if err := validatePassword(password); err != nil {
// 				errorMessage = err.Error()
// 			} else {
// 				// Check if the username or email already exists
// 				var existingUser User
// 				if err := forumDB.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
// 					errorMessage = "Username or email already exists. Please choose another."
// 				}
// 			}
// 		}
// 
// 		// If there is an error, re-render the signup page with the error message
// 		if errorMessage != "" {
// 			t, err := template.ParseFiles("signup.html")
// 			if err != nil {
// 				log.Println("Error parsing signup.html:", err)
// 				http.Error(w, "Internal server error", http.StatusInternalServerError)
// 				return
// 			}
// 
// 			// Pass the error message to the template
// 			data := struct {
// 				ErrorMessage string
// 				FirstName    string
// 				LastName     string
// 				Username     string
// 				Email        string
// 			}{
// 				ErrorMessage: errorMessage,
// 				FirstName:    firstName,
// 				LastName:     lastName,
// 				Username:     username,
// 				Email:        email,
// 			}
// 
// 			t.Execute(w, data)
// 			return
// 		}
// 
// 		// If no errors, create the new user
// 		newUser := User{
// 			FirstName: firstName,
// 			LastName:  lastName,
// 			Username:  username,
// 			Email:     email,
// 			Password:  password,
// 		}
// 		if err := forumDB.Create(&newUser).Error; err != nil {
// 			http.Error(w, "Failed to create user.", http.StatusInternalServerError)
// 			return
// 		}
// 
// 		http.Redirect(w, r, "/", http.StatusFound)
// 	}
// }

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

// ForgotPasswordHandler serves the "Forgot Password" page.
func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("forgot_password.html")
		if err != nil {
			http.Error(w, "Error loading forgot password page", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		email := r.FormValue("email")
		var user User
		err := forumDB.Where("email = ?", email).First(&user).Error
		if err != nil {
			tmpl, _ := template.ParseFiles("forgot_password.html")
			tmpl.Execute(w, struct {
				ErrorMessage   string
				SuccessMessage string
			}{
				ErrorMessage:   "Invalid email address.",
				SuccessMessage: "",
			})
			return
		}

		// Generate a reset token
		token := generateResetToken()

		// Store the token (this can be done in DB or temporary storage)
		resetTokens[email] = token

		// Log the simulated reset link for debugging
		log.Printf("Password reset link for %s: http://localhost:8080/reset-password/?token=%s", email, token)

		// Pass success message to the template
		tmpl, _ := template.ParseFiles("forgot_password.html")
		tmpl.Execute(w, struct {
			ErrorMessage   string
			SuccessMessage string
		}{
			ErrorMessage:   "",
			SuccessMessage: "Sent you an email with a link to reset your password.",
		})
	}
}

// ResetPasswordHandler handles the password reset.
func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		token := r.URL.Query().Get("token")
		if !isValidResetToken(token) {
			http.Error(w, "Invalid or expired reset token", http.StatusBadRequest)
			return
		}

		tmpl, _ := template.ParseFiles("reset_password.html")
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		token := r.FormValue("token")
		newPassword := r.FormValue("password")

		if !isValidResetToken(token) {
			http.Error(w, "Invalid or expired reset token", http.StatusBadRequest)
			return
		}

		// Validate the new password
		if err := validatePassword(newPassword); err != nil {
			tmpl, _ := template.ParseFiles("reset_password.html")
			tmpl.Execute(w, struct{ ErrorMessage string }{ErrorMessage: err.Error()})
			return
		}

		// Find the associated email and user
		email := getEmailByToken(token)
		var user User
		forumDB.Where("email = ?", email).First(&user)

		// Check if the new password is the same as the old one
		if user.Password == newPassword {
			tmpl, _ := template.ParseFiles("reset_password.html")
			tmpl.Execute(w, struct{ ErrorMessage string }{ErrorMessage: "New password cannot be the same as the old password."})
			return
		}

		// Update the password
		user.Password = newPassword
		forumDB.Save(&user)

		// Invalidate the reset token
		delete(resetTokens, email)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Helper functions
var resetTokens = map[string]string{}

func generateResetToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func isValidResetToken(token string) bool {
	for _, v := range resetTokens {
		if v == token {
			return true
		}
	}
	return false
}

func getEmailByToken(token string) string {
	for email, t := range resetTokens {
		if t == token {
			return email
		}
	}
	return ""
}
