package main

import (
	"log"
	"net/http"
)

func goServer() {
	InitializeForumDB()
	defer func() {
		// clse the forum database connection
		if sqlDB, err := forumDB.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	http.Handle("/signup/", enableCORS(http.HandlerFunc(SignUpHandler)))
	http.Handle("/", enableCORS(http.HandlerFunc(SignInHandler)))
	http.Handle("/signout/", enableCORS(http.HandlerFunc(SignOutHandler)))
	http.Handle("/forum/", enableCORS(authMiddleware(http.HandlerFunc(ForumHandler))))
	http.Handle("/profile/", enableCORS(authMiddleware(http.HandlerFunc(profileHandler))))
	http.Handle("/topic/", enableCORS(authMiddleware(http.HandlerFunc(TopicHandler))))
	http.Handle("/new-topic/", enableCORS(authMiddleware(http.HandlerFunc(NewTopicHandler))))
	http.Handle("/new-comment/", enableCORS(authMiddleware(http.HandlerFunc(NewCommentHandler))))
	http.Handle("/view/", enableCORS(authMiddleware(http.HandlerFunc(ViewHandler))))
	http.Handle("/edit/", enableCORS(authMiddleware(http.HandlerFunc(EditHandler))))
	http.Handle("/save/", enableCORS(authMiddleware(http.HandlerFunc(SaveHandler))))

	// new routes for assignments
	http.Handle("/upload-assignment/", enableCORS(authMiddleware(http.HandlerFunc(UploadAssignmentHandler))))
	http.Handle("/assignments/", enableCORS(authMiddleware(http.HandlerFunc(ListAssignmentsHandler))))
	http.Handle("/delete-assignment/", enableCORS(authMiddleware(http.HandlerFunc(DeleteAssignmentHandler))))

	//
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
