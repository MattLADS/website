package main

import (
	"log"
	"net/http"
)

func main() {
	// Load existing user credentials from the database at startup.
	InitializeDB()
	defer db.Close()

	// Set up HTTP handlers for different routes.
	http.HandleFunc("/signup/", SignUpHandler)
	http.HandleFunc("/", SignInHandler)
	http.HandleFunc("/signout/", SignOutHandler)

	http.HandleFunc("/forum/", authMiddleware(ForumHandler))
	
	// profile stuff from Abi's branch, using Kunj's authMiddleware method
	http.HandleFunc("/profile/", authMiddleware(profileHandler))

	http.HandleFunc("/topic/", authMiddleware(TopicHandler))
	http.HandleFunc("/new-topic/", authMiddleware(NewTopicHandler))
	http.HandleFunc("/new-comment/", authMiddleware(NewCommentHandler))
	http.HandleFunc("/view/", authMiddleware(ViewHandler))
	http.HandleFunc("/edit/", authMiddleware(EditHandler))
	http.HandleFunc("/save/", authMiddleware(SaveHandler))

	// Start the HTTP server on port 8080.
	log.Println("Starting server on :8090")
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
