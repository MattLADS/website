package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Handle routes
	http.HandleFunc("/signin", signInHandler)
	http.HandleFunc("/profile", profileHandler)
	http.HandleFunc("/signoff", signOffHandler)

	// Serve static files
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))

	// Start server
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
