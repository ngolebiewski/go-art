package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed dist/*
var staticFiles embed.FS

func main() {
	// --- API route ---
	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"message":"Hello from Go!"}`)
	})

	// --- Determine if we're in dev or production ---
	isDev := os.Getenv("NODE_ENV") != "production" && os.Getenv("RENDER") == ""

	if isDev {
		log.Println("Running in DEVELOPMENT mode - serving from filesystem")
		setupDevServer()
	} else {
		log.Println("Running in PRODUCTION mode - serving from embedded files")
		setupProdServer()
	}

	// --- Start server ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s", port)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func setupProdServer() {
	// Production: use embedded files
	distFS, err := fs.Sub(staticFiles, "dist")
	if err != nil {
		log.Fatal("Failed to create sub filesystem:", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveSPA(w, r, distFS, true)
	})
}

func setupDevServer() {
	// Development: use filesystem (assumes you ran npm run build manually)
	buildDir := "./dist"

	// Check if dist exists, if not provide helpful message
	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		log.Println("No dist/ folder found. Run 'cd frontend && npm run build' first, or start React dev server separately")
		// Still set up handler for API-only development
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, `
					<h1>Go Backend Running</h1>
					<p>API available at <a href="/api/hello">/api/hello</a></p>
					<p>Run 'cd frontend && npm run dev' for React dev server on port 5173</p>
				`)
				return
			}
			http.NotFound(w, r)
		})
		return
	}

	// Serve from filesystem
	distFS := os.DirFS(buildDir)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveSPA(w, r, distFS, false)
	})
}

func serveSPA(w http.ResponseWriter, r *http.Request, distFS fs.FS, isEmbedded bool) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		path = "index.html"
	}

	// Try to serve the requested file
	if file, err := distFS.Open(path); err == nil {
		file.Close()

		// Add appropriate headers
		if strings.HasPrefix(path, "assets/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else if strings.HasSuffix(path, ".html") {
			w.Header().Set("Cache-Control", "no-cache")
		}

		if isEmbedded {
			http.ServeFileFS(w, r, distFS, path)
		} else {
			http.ServeFile(w, r, filepath.Join("dist", path))
		}
		return
	}

	// File not found, serve index.html for SPA routing
	w.Header().Set("Cache-Control", "no-cache")
	if isEmbedded {
		http.ServeFileFS(w, r, distFS, "index.html")
	} else {
		http.ServeFile(w, r, "dist/index.html")
	}
}
