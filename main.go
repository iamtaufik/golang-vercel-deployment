package main

import (
	"log"
	"net/http"
	"os"

	// Import your Vercel function's package
	// The path here depends on your project structure.
	// If api/index.go is at the root, it might be "your-module-name/api"
	// If your go.mod is at the root and api/index.go is there, it's just "api"
	// For this example, let's assume your go.mod is at the root and
	// your Vercel function is in `api/index.go` with package `handler`.
	// The import path will be like this:
	api "github.com/iamtaufik/golang-vercel-deployment/api" // <--- IMPORTANT: Replace with your actual module name
	"github.com/iamtaufik/golang-vercel-deployment/internals/db"
)

func main() {
	// Set the port to listen on. Default to 8080 if not specified.
	db.ConnectDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Register your Vercel function handler
	// If your Vercel function is in `api/index.go` and its package is `handler`,
	// and the function itself is named `Handler`, then you'd use `api.Handler`.
	http.HandleFunc("/api/index", api.Handler) // Make sure the path matches your intended access

	// You can also handle the root path if you want
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Redirect or serve a simple message
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/api/index", http.StatusFound)
			return
		}
		http.NotFound(w, r)
	})

	log.Printf("Server listening on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}