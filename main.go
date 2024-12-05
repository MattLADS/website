package main

// import "C"

import (
	"log"
	"net/http"
)

var chatbot *ChatBot

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
	// Initialize chatbot before starting the server.
	InitializeChatBot()

	// Load existing user credentials from the database at startup.
	InitializeForumDB()
	defer func() {
		// Close the forum database connection
		if sqlDB, err := forumDB.DB(); err == nil {
			sqlDB.Close()
		}
	}()
	/*
		//testing forum handler
		http.Handle("/forum/", enableCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Received request on /forum/")
			ForumHandler(w, r)
		})))
	*/
	// Register the chatbot handler
	http.Handle("/chatbot/openai", enableCORS(http.HandlerFunc(OpenAIHandler)))
	http.Handle("/chatbot", enableCORS(http.HandlerFunc(ChatbotHandler)))
	// Set up HTTP handlers for different routes (EDIT: enabling CORS).
	http.Handle("/signup/", enableCORS(http.HandlerFunc(SignUpHandler)))
	http.Handle("/", enableCORS(http.HandlerFunc(SignInHandler)))
	http.Handle("/signout/", enableCORS(http.HandlerFunc(SignOutHandler)))
	http.Handle("/profile/", enableCORS(authMiddleware(http.HandlerFunc(profileHandler))))
	http.Handle("/topic/", enableCORS(authMiddleware(http.HandlerFunc(TopicHandler))))
	http.Handle("/new-comment/", enableCORS(authMiddleware(http.HandlerFunc(NewCommentHandler))))
	http.Handle("/view/", enableCORS(authMiddleware(http.HandlerFunc(ViewHandler))))
	http.Handle("/edit/", enableCORS(authMiddleware(http.HandlerFunc(EditHandler))))
	http.Handle("/save/", enableCORS(authMiddleware(http.HandlerFunc(SaveHandler))))
	http.Handle("/auth/status", enableCORS(http.HandlerFunc(AuthStatusHandler)))
	http.HandleFunc("/logout/", SignOutHandler)
	http.Handle("/forum/", enableCORS(http.HandlerFunc(ForumHandler)))
	http.Handle("/new-topic/", enableCORS(http.HandlerFunc(NewTopicHandler)))
	http.Handle("/fetch-users", enableCORS(authMiddleware(http.HandlerFunc(FetchUsersHandler))))
	// Other existing routes.
	http.Handle("/send-message", enableCORS(authMiddleware(http.HandlerFunc(SendMessageHandler))))
	http.Handle("/get-messages", enableCORS(authMiddleware(http.HandlerFunc(GetMessagesHandler))))

	// New routes for assignments
	http.Handle("/upload-assignment/", enableCORS(authMiddleware(http.HandlerFunc(UploadAssignmentHandler))))
	http.Handle("/assignments/", enableCORS(authMiddleware(http.HandlerFunc(ListAssignmentsHandler))))
	http.Handle("/delete-assignment/", enableCORS(authMiddleware(http.HandlerFunc(DeleteAssignmentHandler))))

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
