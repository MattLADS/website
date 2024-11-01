package main

import (
	"html/template"
	"net/http"
	"os"
)

// Page structure represents a wiki page with a title and body.
type Page struct {
	Title string
	Body  []byte
}

// save method saves the Page to a text file.
func (p *Page) save() error {
	filename := p.Title + ".txt"
	// Write the page body to a file with read-write permissions for the owner.
	return os.WriteFile(filename, p.Body, 0600)
}

// loadPage loads a Page from a text file with the given title.
// Returns an error if the file doesn't exist.
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	// Read the page body from the corresponding file.
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	// Return a Page with the title and body.
	return &Page{Title: title, Body: body}, nil
}

// viewHandler renders a wiki page. If the page doesn't exist, it redirects to the edit page.
func ViewHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the title from the URL (e.g., "/view/SomePage" -> "SomePage").
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		// If the page is not found, redirect to the edit page so the user can create it.
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	// Render the "view" template to display the page.
	renderTemplate(w, "view", p)
}

// editHandler renders an edit form to modify a wiki page.
// If the page doesn't exist, it creates a new one.
func EditHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the title from the URL (e.g., "/edit/SomePage" -> "SomePage").
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		// If the page doesn't exist, create an empty page with the given title.
		p = &Page{Title: title}
	}
	// Render the "edit" template to allow the user to modify the page.
	renderTemplate(w, "edit", p)
}

// saveHandler saves the contents of an edited page.
// The content is taken from the POST form and saved to a text file.
func SaveHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the title from the URL (e.g., "/save/SomePage" -> "SomePage").
	title := r.URL.Path[len("/save/"):]
	// Get the body content from the form submission.
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	// Save the page's content to a text file.
	err := p.save()
	if err != nil {
		// If saving fails, return an internal server error.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// After saving, redirect the user to view the saved page.
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// renderTemplate is a helper function to parse and execute HTML templates.
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		// Return an error if the template cannot be loaded or parsed.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Execute the template with the given page data.
	err = t.Execute(w, p)
	if err != nil {
		// Return an error if rendering the template fails.
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
