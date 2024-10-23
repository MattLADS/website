package main

import (
	"log"
	"net/http"
)

func main() {
	InitializeDB()
	defer db.Close()

	// define routes
	http.HandleFunc("/signup", SignUpHandler)
	http.HandleFunc("/signin", SignInHandler)
	http.HandleFunc("/signout", SignOutHandler)

	// profile stuff from my branch, using Kunj's authMiddleware method
	http.HandleFunc("/profile", authMiddleware(profileHandler))

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
