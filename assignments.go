package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"
)

// Assignment struct for the database
type Assignment struct {
	ID          int `gorm:"primaryKey"`
	Title       string
	Description string
	FilePath    string
	Username    string
	CreatedAt   time.Time
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

// EditAssignmentHandler renders the edit form for an assignment.
func EditAssignmentHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	// Retrieve the assignment from the database
	var assignment Assignment
	err := forumDB.First(&assignment, id).Error
	if err != nil {
		http.Error(w, "Assignment not found", http.StatusNotFound)
		return
	}

	// Render the edit-assignment template
	tmpl, err := template.ParseFiles("edit-assignment.html")
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, assignment)
}

// UpdateAssignmentHandler processes the form to update an assignment.
func UpdateAssignmentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.FormValue("id")
		title := r.FormValue("title")
		description := r.FormValue("description")
		file, handler, err := r.FormFile("file")

		var assignment Assignment
		err = forumDB.First(&assignment, id).Error
		if err != nil {
			http.Error(w, "Assignment not found", http.StatusNotFound)
			return
		}

		// Update the assignment fields
		assignment.Title = title
		assignment.Description = description

		// If a new file is uploaded, save it
		if err == nil {
			filePath := fmt.Sprintf("uploads/%s", handler.Filename)
			dst, err := os.Create(filePath)
			if err != nil {
				http.Error(w, "Unable to save file", http.StatusInternalServerError)
				return
			}
			defer dst.Close()
			_, err = dst.ReadFrom(file)
			if err != nil {
				http.Error(w, "Unable to save file content", http.StatusInternalServerError)
				return
			}

			// Update the file path in the database
			assignment.FilePath = filePath
		}

		forumDB.Save(&assignment)

		http.Redirect(w, r, "/assignments/", http.StatusSeeOther)
	}

	http.HandleFunc("/upload-assignment", UploadAssignmentHandler)
	http.HandleFunc("/edit-assignment", EditAssignmentHandler)
	http.HandleFunc("/delete-assignment", DeleteAssignmentHandler)
	http.ListenAndServe(":8080", nil)

}
