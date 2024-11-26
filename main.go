package main

// import "C"

import (
	"log"
	"net/http"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin") // Get the Origin from the request
		log.Println("Request Origin:", origin)

		// Allow any localhost origin
		//if origin == "http://localhost:8080" || (len(origin) > 16 && origin[:17] == "http://localhost:") {
		//	w.Header().Set("Access-Control-Allow-Origin", origin)
		//}
		//w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", origin)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

//export goServer
func goServer() {
	// Load existing user credentials from the database at startup.
	InitializeForumDB()
	defer func() {
		// Close the forum database connection
		if sqlDB, err := forumDB.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	// Set up HTTP handlers for different routes (EDIT: enabling CORS).
	http.Handle("/signup/", enableCORS(http.HandlerFunc(SignUpHandler)))
	http.Handle("/", enableCORS(http.HandlerFunc(SignInHandler)))
	http.Handle("/signout/", enableCORS(http.HandlerFunc(SignOutHandler)))
	//http.Handle("/forum/", enableCORS(authMiddleware(http.HandlerFunc(ForumHandler))))
	http.Handle("/profile/", enableCORS(authMiddleware(http.HandlerFunc(profileHandler))))
	http.Handle("/topic/", enableCORS(authMiddleware(http.HandlerFunc(TopicHandler))))
	//http.Handle("/new-topic/", enableCORS(authMiddleware(http.HandlerFunc(NewTopicHandler))))
	http.Handle("/new-comment/", enableCORS(authMiddleware(http.HandlerFunc(NewCommentHandler))))
	http.Handle("/view/", enableCORS(authMiddleware(http.HandlerFunc(ViewHandler))))
	http.Handle("/edit/", enableCORS(authMiddleware(http.HandlerFunc(EditHandler))))
	http.Handle("/save/", enableCORS(authMiddleware(http.HandlerFunc(SaveHandler))))
	http.Handle("/auth/status", enableCORS(http.HandlerFunc(AuthStatusHandler)))
	http.HandleFunc("/logout/", SignOutHandler)
	//log.Println("Registered /forum/ route")
	http.Handle("/forum/", enableCORS(http.HandlerFunc(ForumHandler)))
	http.Handle("/new-topic/", enableCORS(http.HandlerFunc(NewTopicHandler)))

	//http.HandleFunc("/testtopics", TestTopicsHandler)
	//http.HandleFunc("/testtopics", enableCORS(http.HandlerFunc(TestTopicsHandler)))
	//http.Handle("/testtopics", enableCORS(http.HandlerFunc(TestTopicsHandler)))

	// Start the HTTP server on port 8080.
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func main() {
	goServer()
}
