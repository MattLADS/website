package main

//import "C"
import (
	"html/template"
	"net/http"
)

// profile handler
func profileHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve cookies for first name, last name, username, and email
	firstNameCookie, err1 := r.Cookie("first_name")
	lastNameCookie, err2 := r.Cookie("last_name")
	usernameCookie, err3 := r.Cookie("username")
	emailCookie, err4 := r.Cookie("email")

	// Redirect to sign-in page if any cookie is missing
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Create a struct to hold user data for the template
	data := struct {
		FirstName string
		LastName  string
		Username  string
		Email     string
	}{
		FirstName: firstNameCookie.Value,
		LastName:  lastNameCookie.Value,
		Username:  usernameCookie.Value,
		Email:     emailCookie.Value,
	}

	// Parse and execute the profile template
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
