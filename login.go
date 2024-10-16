package main

import (
    "net/http"
    "log"
    "os"
    "bufio"
    "strings"
)

// User structure represents a user with a username and password.
type User struct {
    Username string
    Password string
}

var users = make(map[string]string) // In-memory user storage to hold username-password pairs

// Path to the user file where user credentials are stored.
const userFile = "users.txt"

// LoadUsers loads existing user credentials from the users.txt file into the in-memory map.
func LoadUsers() error {
    file, err := os.Open(userFile) // Attempt to open the user file
    if err != nil {
        return err // Return an error if the file does not exist or can't be opened.
    }
    defer file.Close() // Ensure the file is closed when the function exits.

    scanner := bufio.NewScanner(file) // Create a scanner to read the file line by line.
    for scanner.Scan() {
        line := scanner.Text() // Get the current line.
        parts := strings.SplitN(line, ":", 2) // Split the line into username and password at the first colon.
        if len(parts) == 2 { // Ensure we have both username and password.
            username := parts[0]
            password := parts[1]
            users[username] = password // Add the username and password to the in-memory storage.
        }
    }

    return scanner.Err() // Return any error encountered during scanning.
}

// saveUser saves a new user's credentials to the users.txt file.
func saveUser(username, password string) error {
    file, err := os.OpenFile(userFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600) // Open the file for appending.
    if err != nil {
        return err // Return error if file can't be opened.
    }
    defer file.Close() // Ensure the file is closed when the function exits.

    // Write the username and password to the file, separated by a colon, followed by a newline.
    _, err = file.WriteString(username + ":" + password + "\n")
    return err // Return any error encountered during the write operation.
}

// authMiddleware is a middleware function that checks if the user is authenticated before proceeding to the next handler.
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Check for the session cookie in the user's browser.
        cookie, err := r.Cookie("session_token")

        // If no cookie is found or the cookie value does not indicate authentication,
        // redirect the user to the sign-in page.
        if err != nil || cookie.Value != "authenticated" {
            http.Redirect(w, r, "/", http.StatusFound) // Redirect to sign-in.
            return // Stop further processing.
        }

        // If the user is authenticated, proceed to the next handler.
        next(w, r)
    }
}

// SignUpHandler handles user registration (sign-up).
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        // Render the sign-up page when a GET request is made.
        t, _ := template.ParseFiles("signup.html")
        t.Execute(w, nil)
    } else if r.Method == "POST" {
        // Handle form submission for user registration.
        username := r.FormValue("username") // Get the username from the form.
        password := r.FormValue("password") // Get the password from the form.

        // Check if the username already exists in the in-memory storage.
        if _, exists := users[username]; exists {
            http.Error(w, "Username already exists. Please choose another one.", http.StatusConflict) // Inform user of conflict.
            return
        }

        // Store the new user's credentials in memory.
        users[username] = password

        // Save the new user credentials to the users.txt file.
        err := saveUser(username, password) 
        if err != nil {
            http.Error(w, "Failed to save user credentials.", http.StatusInternalServerError) // Handle file save error.
            return
        }

        // Set a session cookie to mark the user as authenticated.
        http.SetCookie(w, &http.Cookie{
            Name:  "session_token",
            Value: "authenticated", // Mark user as authenticated.
            Path:  "/",
        })

        // Redirect to the home view after successful sign-up.
        http.Redirect(w, r, "/view/home", http.StatusFound)
    }
}

// SignInHandler handles user sign-in.
func SignInHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        // Render the sign-in page for GET requests.
        t, _ := template.ParseFiles("signin.html")
        t.Execute(w, nil)
    } else if r.Method == "POST" {
        // Handle form submission for user sign-in.
        username := r.FormValue("username") // Get the username from the form.
        password := r.FormValue("password") // Get the password from the form.

        // Open the user credentials file to verify the user's credentials.
        file, err := os.Open(userFile)
        if err != nil {
            http.Error(w, "Failed to open user credentials file.", http.StatusInternalServerError) // Handle file open error.
            return
        }
        defer file.Close() // Ensure the file is closed when the function exits.

        scanner := bufio.NewScanner(file) // Create a scanner to read the file line by line.
        var authenticated bool // Flag to track if the user is authenticated.

        // Check each line for matching username and password.
        for scanner.Scan() {
            line := scanner.Text() // Get the current line.
            parts := strings.SplitN(line, ":", 2) // Split the line into username and password.
            if len(parts) == 2 {
                storedUsername := parts[0]
                storedPassword := parts[1]

                // Check if the stored username and password match the input.
                if storedUsername == username && storedPassword == password {
                    authenticated = true // Set authenticated flag if credentials match.
                    break
                }
            }
        }

        if err := scanner.Err(); err != nil {
            http.Error(w, "Error reading user credentials.", http.StatusInternalServerError) // Handle scanning error.
            return
        }

        if authenticated {
            // Set a session cookie to mark the user as authenticated.
            http.SetCookie(w, &http.Cookie{
                Name:  "session_token",
                Value: "authenticated", // Mark user as authenticated.
                Path:  "/",
            })
            // Redirect to the home view upon successful sign-in.
            http.Redirect(w, r, "/view/home", http.StatusFound)
        } else {
            http.Error(w, "Invalid username or password", http.StatusUnauthorized) // Inform user of invalid credentials.
        }
    }
}

// SignOutHandler handles user sign-out by clearing the session cookie.
func SignOutHandler(w http.ResponseWriter, r *http.Request) {
    // Invalidate the session cookie by setting its MaxAge to -1, which deletes it.
    http.SetCookie(w, &http.Cookie{
        Name:   "session_token",
        Value:  "", // Clear the value.
        Path:   "/",
        MaxAge: -1, // Deletes the cookie.
    })

    // Redirect to the sign-in page or home page after logging out.
    http.Redirect(w, r, "/signin", http.StatusFound)
}
