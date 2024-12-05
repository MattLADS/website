package main

//import "C"

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"net/smtp"
	"regexp"
	"strings"
	"time"
)

// User represents a user with a unique ID, username, and password.
type User struct {
	ID          uint   `gorm:"primaryKey"`
	Username    string `json:"username" gorm:"unique"`
	Email       string `json:"email" gorm:"unique"`
	Password    string `json:"password"`
	UserContext string `json:"user_context"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
}

const (
	fromEmail    = "klads723@gmail.com" // Your Gmail address
	fromPassword = "yvdmramzqvjttvdl"   // App password generated earlier
	smtpHost     = "smtp.gmail.com"
	smtpPort     = "587"
)

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

	log.Println("SignUpHandler called")
	if r.Method == "POST" {
		var request User

		// Decode JSON body
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
			return
		}
		// Log received username and password for verification
		log.Printf("Received login request - Username: %s, Password: %s\n", request.Username, request.Password)
		if err := validateEmailDomain(request.Email); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(`{"error": "Invalid email. Please enter a valid email address."}`))
			return
		} else if err := validatePassword(request.Password); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(`{"error": "Password must be at least 8 characters long, contain at least one uppercase letter, and one special character."}`))
			return
		}
		//checking user info in database
		var existingUser User
		if err := forumDB.Where("username = ?", request.Username).First(&existingUser).Error; err == nil {
			// Send JSON response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(`{"error": "Username already exists. Please choose another one."}`))
			return
		}
		// Set username cookie for the session
		http.SetCookie(w, &http.Cookie{
			Name:    "username",
			Value:   request.Username,
			Path:    "/",
			Expires: time.Now().Add(24 * time.Hour), // Setting cookie expiration to 24hrs - we can change later.
		})
		log.Printf("Set-Cookie header for username: %s", request.Username)

		//create new user
		newUser := User{Username: request.Username, Password: request.Password, UserContext: ""}
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
			Name:  "session_token",
			Value: "authenticated",
			Path:  "/",
			//SameSite: http.SameSiteNoneMode, // Allow cross-origin requests
		})

		// Set a cookie with session ID
		http.SetCookie(w, &http.Cookie{
			Name:    "session_id",
			Value:   sessionID,
			Expires: time.Now().Add(24 * time.Hour),
			Path:    "/",
			//SameSite: http.SameSiteNoneMode,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "username",
			Value:    request.Username,
			Expires:  time.Now().Add(24 * time.Hour),
			Path:     "/",
			HttpOnly: true,
			//SameSite: http.SameSiteNoneMode,
			//Secure:   false,
		})
		log.Printf("Set-Cookie header for username: %s", request.Username)

		log.Println("SignInHandler: Sending response")
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

		// Compose the email
		to := []string{email}
		subject := "Password Reset Request"
		body := fmt.Sprintf("Click the link below to reset your password:\n\nhttp://localhost:8080/reset-password/?token=%s", token)
		message := fmt.Sprintf("Subject: %s\n\n%s", subject, body)

		// Send the email
		auth := smtp.PlainAuth("", fromEmail, fromPassword, smtpHost)
		err = smtp.SendMail(smtpHost+":"+smtpPort, auth, fromEmail, to, []byte(message))
		if err != nil {
			log.Printf("Failed to send email: %v", err)
			http.Error(w, "Failed to send reset email.", http.StatusInternalServerError)
			return
		}

		// Log for debugging
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
			tmpl, _ := template.ParseFiles("reset_password.html")
			tmpl.Execute(w, struct {
				ErrorMessage string
				Token        string
			}{
				ErrorMessage: "Invalid or expired reset token.",
				Token:        "",
			})
			return
		}

		tmpl, _ := template.ParseFiles("reset_password.html")
		tmpl.Execute(w, struct {
			ErrorMessage string
			Token        string
		}{
			ErrorMessage: "",
			Token:        token,
		})
	} else if r.Method == http.MethodPost {
		token := r.FormValue("token")
		newPassword := r.FormValue("password")

		if !isValidResetToken(token) {
			tmpl, _ := template.ParseFiles("reset_password.html")
			tmpl.Execute(w, struct {
				ErrorMessage string
				Token        string
			}{
				ErrorMessage: "Invalid or expired reset token.",
				Token:        token,
			})
			return
		}

		// Validate the new password
		if err := validatePassword(newPassword); err != nil {
			tmpl, _ := template.ParseFiles("reset_password.html")
			tmpl.Execute(w, struct {
				ErrorMessage string
				Token        string
			}{
				ErrorMessage: err.Error(),
				Token:        token,
			})
			return
		}

		// Find the associated email and user
		email := getEmailByToken(token)
		var user User
		if err := forumDB.Where("email = ?", email).First(&user).Error; err != nil {
			tmpl, _ := template.ParseFiles("reset_password.html")
			tmpl.Execute(w, struct {
				ErrorMessage string
				Token        string
			}{
				ErrorMessage: "User not found.",
				Token:        token,
			})
			return
		}

		// Check if the new password is the same as the old one
		if user.Password == newPassword {
			tmpl, _ := template.ParseFiles("reset_password.html")
			tmpl.Execute(w, struct {
				ErrorMessage string
				Token        string
			}{
				ErrorMessage: "New password cannot be the same as the old password.",
				Token:        token,
			})
			return
		}

		// Update the password
		user.Password = newPassword
		if err := forumDB.Save(&user).Error; err != nil {
			tmpl, _ := template.ParseFiles("reset_password.html")
			tmpl.Execute(w, struct {
				ErrorMessage string
				Token        string
			}{
				ErrorMessage: "Failed to update the password. Please try again later.",
				Token:        token,
			})
			return
		}

		// Invalidate the reset token
		delete(resetTokens, email)

		tmpl, _ := template.ParseFiles("reset_password.html")
		tmpl.Execute(w, struct {
			ErrorMessage   string
			SuccessMessage string
			Token          string
		}{
			ErrorMessage:   "",
			SuccessMessage: "Password successfully updated. You can now sign in.",
			Token:          "",
		})
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
