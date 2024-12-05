package main

import (
	
	"log"
	"net/http"
	
)

// Declare a global variable for the chatbot instance.
var chatbot *ChatBot

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
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

// goServer starts the HTTP server and registers all routes.
func goServer() {
	// Initialize chatbot before starting the server.
	InitializeChatBot()

	// Load existing user credentials from the database at startup.
	InitializeForumDB()
	defer func() {
		if sqlDB, err := forumDB.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	// Register routes.
	http.Handle("/signup/", enableCORS(http.HandlerFunc(SignUpHandler)))
	http.Handle("/", enableCORS(http.HandlerFunc(SignInHandler)))
	http.Handle("/signout/", enableCORS(http.HandlerFunc(SignOutHandler)))
	http.Handle("/chatbot", enableCORS(http.HandlerFunc(ChatbotHandler))) // Chatbot route.
	http.Handle("/fetch-users", enableCORS(authMiddleware(http.HandlerFunc(FetchUsersHandler))))
	// Other existing routes.
	http.Handle("/send-message", enableCORS(authMiddleware(http.HandlerFunc(SendMessageHandler))))
	http.Handle("/get-messages", enableCORS(authMiddleware(http.HandlerFunc(GetMessagesHandler))))
	http.Handle("/profile/", enableCORS(authMiddleware(http.HandlerFunc(profileHandler))))
	http.Handle("/topic/", enableCORS(authMiddleware(http.HandlerFunc(TopicHandler))))
	http.Handle("/new-topic/", enableCORS(authMiddleware(http.HandlerFunc(NewTopicHandler))))
	http.Handle("/new-comment/", enableCORS(authMiddleware(http.HandlerFunc(NewCommentHandler))))
	http.Handle("/view/", enableCORS(authMiddleware(http.HandlerFunc(ViewHandler))))
	http.Handle("/edit/", enableCORS(authMiddleware(http.HandlerFunc(EditHandler))))
	http.Handle("/save/", enableCORS(authMiddleware(http.HandlerFunc(SaveHandler))))
	http.Handle("/forum/", enableCORS(http.HandlerFunc(ForumHandler)))

	// Start the HTTP server.
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func main() {
	goServer()
}

