package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func Error(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func renderTemplate(w http.ResponseWriter, tmpl string) {
	// Combine the layout and the specific page template
	tmplPath := filepath.Join("templates", tmpl+".html")
	layoutPath := filepath.Join("templates", "layout.html")

	// Log the paths of the templates being parsed
	log.Printf("Parsing templates: layout=%s, page=%s", layoutPath, tmplPath)

	// Parse all templates
	t, err := template.ParseFiles(layoutPath, tmplPath)
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log successful parsing of templates
	log.Printf("Successfully parsed templates: layout=%s, page=%s", layoutPath, tmplPath)

	// Render the layout template, injecting the content from the specific page template
	err = t.ExecuteTemplate(w, "layout.html", nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		// Log successful execution of the template
		log.Printf("Successfully executed template: layout=%s, page=%s", layoutPath, tmplPath)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Rendering home page")
	renderTemplate(w, "home")
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Rendering blog page")
	http.ServeFile(w, r, "blog")
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Rendering projects page")
	renderTemplate(w, "projects")
}

func writingpageHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Rendering a blog article: %v", r.URL.Path)

	renderTemplate(w, "writingpage")

	err := r.Write(w)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Rendering about page")
	renderTemplate(w, "about")
}

func main() {
	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Define routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/projects", projectsHandler)
	http.HandleFunc("/blog", blogHandler)
	http.HandleFunc("/writingpage", writingpageHandler)
	http.HandleFunc("/about", aboutHandler)
	// Start the server on port 8080
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
