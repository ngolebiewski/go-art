package routes

import (
	"go-art-api/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *mux.Router) {
	// Create API subrouter
	api := r.PathPrefix("/api").Subrouter()

	// Add CORS middleware for all API routes
	api.Use(corsMiddleware)

	// Health check
	api.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	api.HandleFunc("/hello", handlers.Hello).Methods("GET")

	// User routes
	setupUserRoutes(api)

	// Artist routes
	setupArtistRoutes(api)

	// Artwork routes
	setupArtworkRoutes(api)

	// Medium routes
	setupMediumRoutes(api)

	// Special/complex routes
	setupSpecialRoutes(api)
}

// setupUserRoutes defines user-related routes
func setupUserRoutes(api *mux.Router) {
	users := api.PathPrefix("/users").Subrouter()

	users.HandleFunc("", handlers.GetUsers).Methods("GET")
	users.HandleFunc("", handlers.CreateUser).Methods("POST")
	users.HandleFunc("/{id:[0-9]+}", handlers.GetUserByID).Methods("GET")
	users.HandleFunc("/{id:[0-9]+}", handlers.UpdateUser).Methods("PUT")
	users.HandleFunc("/{id:[0-9]+}", handlers.DeleteUser).Methods("DELETE")

	// Authentication routes
	api.HandleFunc("/auth/register", handlers.RegisterUser).Methods("POST")
	api.HandleFunc("/auth/login", handlers.LoginUser).Methods("POST")
}

// setupArtistRoutes defines artist-related routes
func setupArtistRoutes(api *mux.Router) {
	artists := api.PathPrefix("/artists").Subrouter()

	artists.HandleFunc("", handlers.GetArtists).Methods("GET")
	artists.HandleFunc("", handlers.CreateArtist).Methods("POST")
	artists.HandleFunc("/{id:[0-9]+}", handlers.GetArtistByID).Methods("GET")
	artists.HandleFunc("/{id:[0-9]+}", handlers.UpdateArtist).Methods("PUT")
	artists.HandleFunc("/{id:[0-9]+}", handlers.DeleteArtist).Methods("DELETE")

	// Artist-specific routes
	artists.HandleFunc("/{id:[0-9]+}/artworks", handlers.GetArtworksByArtist).Methods("GET")
}

// setupArtworkRoutes defines artwork-related routes
func setupArtworkRoutes(api *mux.Router) {
	artworks := api.PathPrefix("/artworks").Subrouter()

	// NEW ROUTE: Handles POST to /api/artworks for creation + upload
	artworks.HandleFunc("", handlers.CreateArtworkAndUploadImage).Methods("POST")

	artworks.HandleFunc("", handlers.GetArtworks).Methods("GET")
	// artworks.HandleFunc("", handlers.CreateArtwork).Methods("POST")
	artworks.HandleFunc("/{id:[0-9]+}", handlers.GetArtworkByID).Methods("GET")
	artworks.HandleFunc("/{id:[0-9]+}", handlers.UpdateArtwork).Methods("PUT")
	artworks.HandleFunc("/{id:[0-9]+}", handlers.DeleteArtwork).Methods("DELETE")

	// Artwork-specific routes
	artworks.HandleFunc("/view", handlers.GetArtworksView).Methods("GET")
	artworks.HandleFunc("/{id:[0-9]+}/mediums", handlers.GetArtworkMediums).Methods("GET")
	artworks.HandleFunc("/{id:[0-9]+}/mediums", handlers.AddArtworkMedium).Methods("POST")
	artworks.HandleFunc("/{id:[0-9]+}/mediums/{medium_id:[0-9]+}", handlers.RemoveArtworkMedium).Methods("DELETE")

	// Image Upload
	artworks.HandleFunc("/{id:[0-9]+}/image", handlers.UploadImage).Methods("POST")

	// Image Retrieval
	artworks.HandleFunc("/images/{id:[0-9]+}", handlers.GetImage).Methods("GET")
	artworks.HandleFunc("/images/{id:[0-9]+}/thumb", handlers.GetThumbnail).Methods("GET")
}

// setupMediumRoutes defines medium-related routes
func setupMediumRoutes(api *mux.Router) {
	mediums := api.PathPrefix("/mediums").Subrouter()

	mediums.HandleFunc("", handlers.GetMediums).Methods("GET")
	mediums.HandleFunc("", handlers.CreateMedium).Methods("POST")
	mediums.HandleFunc("/{id:[0-9]+}", handlers.GetMediumByID).Methods("GET")
	mediums.HandleFunc("/{id:[0-9]+}", handlers.UpdateMedium).Methods("PUT")
	mediums.HandleFunc("/{id:[0-9]+}", handlers.DeleteMedium).Methods("DELETE")
}

// setupSpecialRoutes defines complex/special routes
func setupSpecialRoutes(api *mux.Router) {
	// Search routes
	api.HandleFunc("/search/artists", handlers.SearchArtists).Methods("GET")
	api.HandleFunc("/search/artworks", handlers.SearchArtworks).Methods("GET")

	// Statistics routes
	api.HandleFunc("/stats/overview", handlers.GetOverviewStats).Methods("GET")

	// User-Artist relationship routes
	api.HandleFunc("/users/{user_id:[0-9]+}/artists", handlers.GetUserArtists).Methods("GET")
	api.HandleFunc("/users/{user_id:[0-9]+}/artists/{artist_id:[0-9]+}", handlers.AddUserArtist).Methods("POST")
	api.HandleFunc("/users/{user_id:[0-9]+}/artists/{artist_id:[0-9]+}", handlers.RemoveUserArtist).Methods("DELETE")
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
