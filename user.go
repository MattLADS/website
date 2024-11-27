package main

//import "C"
import (
	"html/template"
	"log"
	"net/http"
)

// profile handler
func profileHandler(w http.ResponseWriter, r *http.Request) {
	// render profile page with user's info
	tmpl, err := template.ParseFiles("profile.html")
	if err != nil {
		log.Print(err)
	}
	cookie, err := r.Cookie("username")
	if err != nil {
		log.Print(err)
	}

	// Pass topics data to the template
	tmpl.Execute(w, cookie)
}
