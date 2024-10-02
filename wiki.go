// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
package main

import (
    "html/template"
    "log"
    "net/http"
    "os"
    "bufio"
    "strings"
)

// Page structure represents a wiki page with a title and body.
type Page struct {
    Title string
    Body  []byte
}

// save method saves the Page to a text file.
func (p *Page) save() error {
    filename := p.Title + ".txt"
    // Write the page body to a file with read-write permissions for the owner.
    return os.WriteFile(filename, p.Body, 0600)
}

// loadPage loads a Page from a text file with the given title.
// Returns an error if the file doesn't exist.
func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    // Read the page body from the corresponding file.
    body, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    // Return a Page with the title and body.
    return &Page{Title: title, Body: body}, nil
}

// viewHandler renders a wiki page. If the page doesn't exist, it redirects to the edit page.
func viewHandler(w http.ResponseWriter, r *http.Request) {
    // Extract the title from the URL (e.g., "/view/SomePage" -> "SomePage").
    title := r.URL.Path[len("/view/"):]
    p, err := loadPage(title)
    if err != nil {
        // If the page is not found, redirect to the edit page so the user can create it.
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    // Render the "view" template to display the page.
    renderTemplate(w, "view", p)
}

// editHandler renders an edit form to modify a wiki page.
// If the page doesn't exist, it creates a new one.
func editHandler(w http.ResponseWriter, r *http.Request) {
    // Extract the title from the URL (e.g., "/edit/SomePage" -> "SomePage").
    title := r.URL.Path[len("/edit/"):]
    p, err := loadPage(title)
    if err != nil {
        // If the page doesn't exist, create an empty page with the given title.
        p = &Page{Title: title}
    }
    // Render the "edit" template to allow the user to modify the page.
    renderTemplate(w, "edit", p)
}

// saveHandler saves the contents of an edited page.
// The content is taken from the POST form and saved to a text file.
func saveHandler(w http.ResponseWriter, r *http.Request) {
    // Extract the title from the URL (e.g., "/save/SomePage" -> "SomePage").
    title := r.URL.Path[len("/save/"):]
    // Get the body content from the form submission.
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    // Save the page's content to a text file.
    err := p.save()
    if err != nil {
        // If saving fails, return an internal server error.
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    // After saving, redirect the user to view the saved page.
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// renderTemplate is a helper function to parse and execute HTML templates.
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    t, err := template.ParseFiles(tmpl + ".html")
    if err != nil {
        // Return an error if the template cannot be loaded or parsed.
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    // Execute the template with the given page data.
    err = t.Execute(w, p)
    if err != nil {
        // Return an error if rendering the template fails.
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}



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
            http.Redirect(w, r, "/signin", http.StatusFound) // Redirect to sign-in.
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


func main() {
    // Load existing user credentials from the users.txt file at startup.
    err := LoadUsers()
    if err != nil {
        log.Fatalf("Error loading users: %v", err)
    }

    // Set up HTTP handlers for different routes.
    http.HandleFunc("/signup", SignUpHandler)
    http.HandleFunc("/signin", SignInHandler)
    http.HandleFunc("/signout", SignOutHandler)

    http.HandleFunc("/view/", authMiddleware(viewHandler))
    http.HandleFunc("/edit/", authMiddleware(editHandler))
    http.HandleFunc("/save/", authMiddleware(saveHandler))

    // Start the HTTP server on port 8080.
    log.Println("Starting server on :8080")
    err = http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
}


//hi
