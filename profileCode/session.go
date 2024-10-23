package main

import (
	"net/http"
	"time"
)

var sessions = map[string]string{}

// create session for a user
func createSession(sessionID, username string) {
	sessions[sessionID] = username
}

// get the username associated with a session
func getSessionUser(r *http.Request) (string, bool) {
	// get session cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", false
	}

	username, exists := sessions[cookie.Value]
	return username, exists
}

// sign-off (logout) handler
func signOffHandler(w http.ResponseWriter, r *http.Request) {
	// get session cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// ddelete session from memory
	delete(sessions, cookie.Value)

	// clealear the cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour), // expire the cookie
		Path:    "/",
	})

	http.Redirect(w, r, "/", http.StatusFound)
}
