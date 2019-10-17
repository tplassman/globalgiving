package main

import (
	"net/http"

	"github.com/tplassman/globalgiving/pkg/handlers"
)

func main() {
	// Static assets
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	// Define routes
	http.HandleFunc("/", handlers.ProjectsHandler)
	// Start server
	http.ListenAndServe(":8080", nil)
}
