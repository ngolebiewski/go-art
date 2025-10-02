package main

import (
	"embed"
	"log"
	"net/http"
	"os"
	"time"

	"go-art-api/config"
	"go-art-api/routes"
	"go-art-api/static"

	"github.com/gorilla/mux"
)

//go:embed dist/*
var staticFiles embed.FS

func main() {
	// Initialize database
	config.InitDB()
	defer config.CloseDB()

	// Setup router
	r := mux.NewRouter()

	// Setup API routes
	routes.SetupRoutes(r)

	// Setup static file serving (hybrid approach)
	static.SetupStaticFiles(r, staticFiles)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 Go Art API Server starting on :%s", port)
	log.Printf("📊 API available at /api/...")
	log.Fatal(srv.ListenAndServe())
}
