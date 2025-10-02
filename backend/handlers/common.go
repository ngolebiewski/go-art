package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"go-art-api/models"
)

// sendJSONResponse sends a JSON response
func sendJSONResponse(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// sendErrorResponse sends an error response
func sendErrorResponse(w http.ResponseWriter, message string, status int) {
	response := models.APIResponse{
		Success: false,
		Error:   message,
	}
	sendJSONResponse(w, response, status)
}

// sendSuccessResponse sends a success response
func sendSuccessResponse(w http.ResponseWriter, data interface{}, message string, status int) {
	response := models.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	sendJSONResponse(w, response, status)
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// nullIfEmpty returns nil if string is empty, otherwise returns the string
func nullIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

// Health check endpoint
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := models.APIResponse{
		Success: true,
		Message: "Go Art API is healthy! 🎨",
		Data: map[string]interface{}{
			"status":  "ok",
			"service": "go-art-api",
		},
	}
	sendJSONResponse(w, response, http.StatusOK)
}

// Hello endpoint for testing
func Hello(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "Hello from Go Art API! 🚀",
		"version": "1.0.0",
	}
	sendJSONResponse(w, response, http.StatusOK)
}

// Placeholder handlers for unimplemented routes
func GetArtists(w http.ResponseWriter, r *http.Request) {
	sendSuccessResponse(w, []interface{}{}, "Artists endpoint - coming soon!", http.StatusOK)
}

func CreateArtist(w http.ResponseWriter, r *http.Request) {
	sendErrorResponse(w, "Create artist endpoint not implemented yet", http.StatusNotImplemented)
}

func GetArtistByID(w http.ResponseWriter, r *http.Request) {
	sendErrorResponse(w, "Get artist by ID endpoint not implemented yet", http.StatusNotImplemented)
}

func UpdateArtist(w http.ResponseWriter, r *http.Request) {
	sendErrorResponse(w, "Update artist endpoint not implemented yet", http.StatusNotImplemented)
}

func DeleteArtist(w http.ResponseWriter, r *http.Request) {
	sendErrorResponse(w, "Delete artist endpoint not implemented yet", http.StatusNotImplemented)
}

func GetArtworks(w http.ResponseWriter, r *http.Request) {
	sendSuccessResponse(w, []interface{}{}, "Artworks endpoint - coming soon!", http.StatusOK)
}

func CreateArtwork(w http.ResponseWriter, r *http.Request) {
	sendErrorResponse(w, "Create artwork endpoint not implemented yet", http.StatusNotImplemented)
}

func GetArtworkByID(w http.ResponseWriter, r *http.Request) {
	sendErrorResponse(w, "Get artwork by ID endpoint not implemented yet", http.StatusNotImplemented)
}

func UpdateArtwork(w http.ResponseWriter, r *http.Request) {
	sendErrorResponse(w, "Update artwork endpoint not implemented yet", http.StatusNotImplemented)
}

func DeleteArtwork(w http.ResponseWriter, r *http.Request) {
	sendErrorResponse(w, "Delete artwork endpoint not implemented yet", http.StatusNotImplemented)
}

func GetMediums(w http.ResponseWriter, r *http.Request) {
	sendSuccessResponse(w, []interface{}{}, "Mediums endpoint - coming soon!", http.StatusOK)
}

func CreateMedium(w http.ResponseWriter, r *http.Request) {
	sendErrorResponse(w, "Create medium endpoint not implemented yet", http.StatusNotImplemented)
}

func GetMediumByID(w http.ResponseWriter, r *http.Request) {
	sendErrorResponse(w, "Get medium by ID endpoint not implemented yet", http.StatusNotImplemented)
}

func UpdateMedium(w http.ResponseWriter, r *http.Request) {
	sendErrorResponse(w, "Update medium endpoint not implemented yet", http.StatusNotImplemented)
}

func DeleteMedium(w http.ResponseWriter, r *http.Request) {
	sendErrorResponse(w, "Delete medium endpoint not implemented yet", http.StatusNotImplemented)
}

// Special route placeholders
func GetArtworksView(w http.ResponseWriter, r *http.Request) {
	sendSuccessResponse(w, []interface{}{}, "Artworks view endpoint - coming soon!", http.StatusOK)
}

func GetArtworksByArtist(w http.ResponseWriter, r *http.Request) {
	sendSuccessResponse(w, []interface{}{}, "Get artworks by artist - coming soon!", http.StatusOK)
}

func GetArtworkMediums(w http.ResponseWriter, r *http.Request) {
	sendSuccessResponse(w, []interface{}{}, "Get artwork mediums - coming soon!", http.StatusOK)
}

func AddArtworkMedium(w http.ResponseWriter, r *http.Request) {
	sendErrorResponse(w, "Add artwork medium endpoint not implemented yet", http.StatusNotImplemented)
}

func RemoveArtworkMedium(w http.ResponseWriter, r *http.Request) {
	sendErrorResponse(w, "Remove artwork medium endpoint not implemented yet", http.StatusNotImplemented)
}

func SearchArtists(w http.ResponseWriter, r *http.Request) {
	sendSuccessResponse(w, []interface{}{}, "Search artists - coming soon!", http.StatusOK)
}

func SearchArtworks(w http.ResponseWriter, r *http.Request) {
	sendSuccessResponse(w, []interface{}{}, "Search artworks - coming soon!", http.StatusOK)
}

func GetOverviewStats(w http.ResponseWriter, r *http.Request) {
	sendSuccessResponse(w, map[string]int{
		"users":    0,
		"artists":  0,
		"artworks": 0,
	}, "Overview stats - coming soon!", http.StatusOK)
}

func GetUserArtists(w http.ResponseWriter, r *http.Request) {
	sendSuccessResponse(w, []interface{}{}, "Get user artists - coming soon!", http.StatusOK)
}

func AddUserArtist(w http.ResponseWriter, r *http.Request) {
	sendErrorResponse(w, "Add user artist endpoint not implemented yet", http.StatusNotImplemented)
}

func RemoveUserArtist(w http.ResponseWriter, r *http.Request) {
	sendErrorResponse(w, "Remove user artist endpoint not implemented yet", http.StatusNotImplemented)
}
