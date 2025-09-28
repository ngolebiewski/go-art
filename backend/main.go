package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// --- API route ---
	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"message":"Hello from Go!"}`)
	})

	// --- Static file server (production build folder) ---
	buildDir := "./dist"

	// Check if dist directory exists
	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		log.Printf("WARNING: dist directory not found at %s", buildDir)
		// List current directory contents for debugging
		files, _ := os.ReadDir(".")
		log.Printf("Current directory contents:")
		for _, file := range files {
			log.Printf("  %s", file.Name())
		}
	}

	fileServer := http.FileServer(http.Dir(buildDir))

	// --- SPA fallback handler ---
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Trim leading "/" to avoid filepath.Join bug
		reqPath := strings.TrimPrefix(r.URL.Path, "/")
		fsPath := filepath.Join(buildDir, reqPath)

		// Ensure path is inside buildDir (prevents ../ traversal)
		absFSPath, err1 := filepath.Abs(fsPath)
		absBuildDir, err2 := filepath.Abs(buildDir)
		if err1 != nil || err2 != nil || !strings.HasPrefix(absFSPath, absBuildDir) {
			log.Printf("Path traversal attempt or error: %s", reqPath)
			http.NotFound(w, r)
			return
		}

		info, err := os.Stat(absFSPath)
		if err != nil || info.IsDir() {
			// File missing → fallback to index.html
			indexPath := filepath.Join(buildDir, "index.html")
			if _, err := os.Stat(indexPath); os.IsNotExist(err) {
				log.Printf("index.html not found at %s", indexPath)
				http.NotFound(w, r)
				return
			}
			log.Printf("Serving index.html for: %s", reqPath)
			http.ServeFile(w, r, indexPath)
			return
		}

		// Add caching for static assets
		if strings.HasPrefix(reqPath, "static/") || strings.HasPrefix(reqPath, "assets/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else if strings.HasSuffix(reqPath, ".html") {
			w.Header().Set("Cache-Control", "no-cache")
		}

		log.Printf("Serving static file: %s", reqPath)
		fileServer.ServeHTTP(w, r)
	})

	// --- Listen address (env PORT or fallback to 8080) ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	log.Printf("Starting server on http://localhost%s (API at /api/...)\n", addr)
	log.Printf("Serving static files from: %s", buildDir)

	// Log absolute path for debugging
	if absPath, err := filepath.Abs(buildDir); err == nil {
		log.Printf("Absolute path to dist: %s", absPath)
	}

	srv := &http.Server{
		Addr:         addr,
		Handler:      nil, // use DefaultServeMux
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
