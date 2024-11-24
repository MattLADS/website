package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"gorm.io/gorm"
)

// Assignment struct for the database
type Assignment struct {
	gorm.Model
	Title       string
	Description string
	FilePath    string
	Username    string
}

// UploadAssignmentHandler handles file uploads and stores assignment data
func UploadAssignmentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(10 << 20) // 10 MB limit
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		description := r.FormValue("description")
		username := r.FormValue("username")

		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Unable to get file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Save the file locally
		filePath := fmt.Sprintf("uploads/%s", handler.Filename)
		dst, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Unable to save file", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = dst.ReadFrom(file)
		if err != nil {
			http.Error(w, "Unable to write file", http.StatusInternalServerError)
			return
		}

		// Save assignment record to database
		newAssignment := Assignment{
			Title:       title,
			Description: description,
			FilePath:    filePath,
			Username:    username,
		}
		forumDB.Create(&newAssignment)

		http.Redirect(w, r, "/assignments/", http.StatusSeeOther)
	}
}

// ListAssignmentsHandler lists all assignments for the Assignments page
func ListAssignmentsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("assignments.html")
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	var assignments []Assignment
	forumDB.Find(&assignments)

	tmpl.Execute(w, assignments)
}

// DeleteAssignmentHandler deletes an assignment by ID
func DeleteAssignmentHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var assignment Assignment
	err := forumDB.First(&assignment, id).Error
	if err != nil {
		http.Error(w, "Assignment not found", http.StatusNotFound)
		return
	}

	// Remove file from storage
	os.Remove(assignment.FilePath)

	// Remove from database
	forumDB.Delete(&assignment)

	http.Redirect(w, r, "/assignments/", http.StatusSeeOther)
}
