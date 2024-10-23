package main

import (
	"net/http"
)

// profile handler
func profileHandler(w http.ResponseWriter, r *http.Request) {
	// using Kunj's session check method
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value != "authenticated" {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	// render profile page with user's info
	http.ServeFile(w, r, "templates/profile.html")
}
