package main

import (
	"fmt"
	"net/http"
)

// profile handler
func profileHandler(w http.ResponseWriter, r *http.Request) {
	// get the username from session
	username, exists := getSessionUser(r)
	if !exists {
		http.Error(w, "You need to log in to view your profile.", http.StatusUnauthorized)
		return
	}

	// render profile page with user's info
	http.ServeFile(w, r, "templates/profile.html")
	fmt.Fprintf(w, "Welcome %s! This is your profile.", username)
}
