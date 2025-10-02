package handlers

import (
	"database/sql"
	"fmt"
	"go-art-api/config"
	"go-art-api/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// The required helper functions (sendSuccessResponse, sendErrorResponse, etc.)
// are assumed to be present in this package.

// Inside handlers/artworks.go

// CreateArtworkAndUploadImage handles creating a new artwork and uploading its initial image simultaneously.
func CreateArtworkAndUploadImage(w http.ResponseWriter, r *http.Request) {
	// 1. Parse the Multipart Form Data (Max 10MB)
	log.Printf("--- START: CreateArtworkAndUploadImage ---")
	log.Printf("1. Attempting to parse multipart form data...")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("ERROR 1.1: Failed to parse form or size exceeded: %v", err)
		sendErrorResponse(w, "File size limit (10MB) exceeded or failed to parse form", http.StatusBadRequest)
		return
	}

	// 2. Extract Artwork Metadata from form values
	title := r.FormValue("title")
	artistIDStr := r.FormValue("artist_id")
	grade := r.FormValue("grade")
	school := r.FormValue("school")

	log.Printf("1.2. Parsed form values: Title='%s', ArtistIDStr='%s'", title, artistIDStr)

	artistID, err := strconv.Atoi(artistIDStr)
	if err != nil || title == "" || artistID == 0 {
		log.Printf("ERROR 1.3: Validation failed. ArtistID conversion error: %v", err)
		sendErrorResponse(w, "Title, Artist ID, and valid data are required", http.StatusBadRequest)
		return
	}

	// 3. Insert Artwork into Database
	log.Printf("2. Starting DB INSERT for new artwork (ArtistID: %d)", artistID)
	result, err := config.DB.Exec(
		"INSERT INTO artworks (artist_id, title, grade, school) VALUES (?, ?, ?, ?)",
		artistID, title, grade, school,
	)
	if err != nil {
		log.Printf("FATAL DB ERROR 2.1: Failed to insert artwork. Check FK/constraints: %v", err)
		sendErrorResponse(w, "Failed to create artwork entry", http.StatusInternalServerError)
		return
	}

	artworkID, _ := result.LastInsertId()
	if artworkID == 0 {
		log.Printf("ERROR 2.2: Failed to retrieve new artwork ID.")
		sendErrorResponse(w, "Failed to retrieve new artwork ID", http.StatusInternalServerError)
		return
	}
	log.Printf("2.3. Artwork created successfully. New Artwork ID: %d", artworkID)

	// 4. Extract and Process the Image File
	log.Printf("3. Attempting to read image file from form field 'image'...")
	file, header, err := r.FormFile("image") // 'image' is the file input field name
	if err != nil {
		log.Printf("ERROR 3.1: Image file missing: %v. Deleting created artwork.", err)
		config.DB.Exec("DELETE FROM artworks WHERE id = ?", artworkID) // Clean up artwork
		sendErrorResponse(w, "No image file provided in the 'image' field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	originalMime := header.Header.Get("Content-Type")
	log.Printf("3.2. File found. Original MIME: %s. Starting image processing...", originalMime)

	// Process the Image (Generates 2 JPEG BLOBs)
	thumbData, imageData, err := utils.ProcessImage(file)
	if err != nil {
		// This is the common hang point if processing is too long or crashes.
		log.Printf("FATAL PROCESSING ERROR 3.3: Image processing failed: %v. Deleting created artwork.", err)
		config.DB.Exec("DELETE FROM artworks WHERE id = ?", artworkID) // Clean up artwork
		sendErrorResponse(w, fmt.Sprintf("Image processing failed: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("4. Image processing complete. Thumb Size: %d bytes, Full Size: %d bytes.", len(thumbData), len(imageData))

	// 5. Insert Image Data into Database (Using the new artworkID)
	query := `
        INSERT INTO images (artwork_id, original_mime, thumb, image, url) 
        VALUES (?, ?, ?, ?, NULL)
    `
	log.Printf("5. Starting final DB INSERT for image data (Artwork ID: %d)...", artworkID)
	result, err = config.DB.Exec(query, artworkID, originalMime, thumbData, imageData)

	if err != nil {
		// This is the other common crash point (e.g., if a BLOB exceeds MySQL size limit)
		log.Printf("FATAL DB ERROR 5.1: Database error during image INSERT: %v. Deleting created artwork.", err)
		config.DB.Exec("DELETE FROM artworks WHERE id = ?", artworkID) // Clean up
		sendErrorResponse(w, "Failed to save image data to the database", http.StatusInternalServerError)
		return
	}

	imageID, _ := result.LastInsertId()
	log.Printf("5.2. Image data saved successfully. New Image ID: %d", imageID)

	// 6. Success Response
	log.Printf("6. Sending final SUCCESS response.")
	sendSuccessResponse(w, map[string]interface{}{
		"artwork_id": artworkID,
		"image_id":   imageID,
		"title":      title,
		"image_size": fmt.Sprintf("%.2f KB", float64(len(imageData))/1024),
	}, "Artwork and image created successfully", http.StatusCreated)
	log.Printf("--- END: CreateArtworkAndUploadImage ---")
}

// UploadImage handles the upload of a new image, processes it, and saves it to the database.
func UploadImage(w http.ResponseWriter, r *http.Request) {
	// 1. Get Artwork ID from URL
	vars := mux.Vars(r)
	artworkIDStr := vars["id"]
	artworkID, err := strconv.Atoi(artworkIDStr)
	if err != nil {
		sendErrorResponse(w, "Invalid artwork ID", http.StatusBadRequest)
		return
	}

	// 2. Parse the Multipart Form Data (Max 10MB)
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		sendErrorResponse(w, "File size limit (10MB) exceeded or failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image") // 'image' is the file input field name
	if err != nil {
		sendErrorResponse(w, "No image file provided in the 'image' field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get original MIME type for metadata storage
	originalMime := header.Header.Get("Content-Type")

	// 3. Process the Image (Generates 2 JPEG BLOBs)
	thumbData, imageData, err := utils.ProcessImage(file)
	if err != nil {
		sendErrorResponse(w, fmt.Sprintf("Image processing failed: %v", err), http.StatusInternalServerError)
		return
	}

	// 4. Database Insertion (UPSERT Logic)
	var existingID int
	err = config.DB.QueryRow("SELECT id FROM images WHERE artwork_id = ?", artworkID).Scan(&existingID)

	var result sql.Result
	if err == sql.ErrNoRows {
		// INSERT (New image)
		query := `
            INSERT INTO images (artwork_id, original_mime, thumb, image, url) 
            VALUES (?, ?, ?, ?, NULL)
        `
		result, err = config.DB.Exec(query, artworkID, originalMime, thumbData, imageData)
	} else if err == nil {
		// UPDATE (Replace existing image)
		query := `
            UPDATE images SET original_mime = ?, thumb = ?, image = ?, url = NULL
            WHERE artwork_id = ?
        `
		result, err = config.DB.Exec(query, originalMime, thumbData, imageData, artworkID)
	}

	if err != nil {
		log.Printf("Database error during image UPSERT: %v", err)
		sendErrorResponse(w, "Failed to save image data to the database", http.StatusInternalServerError)
		return
	}

	imageID, _ := result.LastInsertId()
	if imageID == 0 && existingID != 0 {
		imageID = int64(existingID) // Use existing ID if it was an update
	}

	// 5. Success Response
	sendSuccessResponse(w, map[string]interface{}{
		"image_id":        imageID,
		"artwork_id":      artworkID,
		"thumb_size":      fmt.Sprintf("%.2f KB", float64(len(thumbData))/1024),
		"image_size":      fmt.Sprintf("%.2f KB", float64(len(imageData))/1024),
		"stored_format":   "image/jpeg",
		"original_format": originalMime,
	}, "Image uploaded, processed, and saved successfully", http.StatusCreated)
}

// GetImage retrieves and serves the full-size image (BLOB)
func GetImage(w http.ResponseWriter, r *http.Request) {
	serveImage(w, r, "image")
}

// GetThumbnail retrieves and serves the thumbnail image (BLOB)
func GetThumbnail(w http.ResponseWriter, r *http.Request) {
	serveImage(w, r, "thumb")
}

// serveImage is a helper function to retrieve and serve the requested BLOB column.
func serveImage(w http.ResponseWriter, r *http.Request, column string) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sendErrorResponse(w, "Invalid image ID", http.StatusBadRequest)
		return
	}

	// 1. Fetch the image data (BLOB) from the database
	var imageData []byte
	query := fmt.Sprintf("SELECT %s FROM images WHERE id = ?", column)
	err = config.DB.QueryRow(query, id).Scan(&imageData)

	if err == sql.ErrNoRows {
		sendErrorResponse(w, "Image not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("DB error fetching image: %v", err)
		sendErrorResponse(w, "Failed to retrieve image data", http.StatusInternalServerError)
		return
	}

	// 2. Set the appropriate HTTP headers (Always JPEG since we process it that way)
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "public, max-age=2592000, immutable")

	// 3. Write the raw image data (BLOB) to the response body
	if _, err := w.Write(imageData); err != nil {
		log.Printf("Error writing image to response: %v", err)
		// Note: Can't send an HTTP error here, as headers are already sent.
	}
}
