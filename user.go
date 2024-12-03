package main

//import "C"
import (
	"html/template"
	"net/http"
)

// profile handler
func profileHandler(w http.ResponseWriter, r *http.Request) {
	usernameCookie, err1 := r.Cookie("username")
	emailCookie, err2 := r.Cookie("email")

	if err1 != nil || err2 != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	data := struct {
		Username string
		Email    string
	}{
		Username: usernameCookie.Value,
		Email:    emailCookie.Value,
	}

	t, err := template.ParseFiles("profile.html")
	if err != nil {
		http.Error(w, "Error loading profile page", http.StatusInternalServerError)
		return
	}

	t.Execute(w, data)
}

func EditHandlers(w http.ResponseWriter, r *http.Request) {
	usernameCookie, err1 := r.Cookie("username")
	emailCookie, err2 := r.Cookie("email")

	if err1 != nil || err2 != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	data := struct {
		Username string
		Email    string
	}{
		Username: usernameCookie.Value,
		Email:    emailCookie.Value,
	}

	t, err := template.ParseFiles("edit.html")
	if err != nil {
		http.Error(w, "Error loading edit page", http.StatusInternalServerError)
		return
	}

	t.Execute(w, data)
}
