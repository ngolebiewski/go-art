package models

import (
	"time"
)

// --- Users and Authentication Models ---

// User represents a user in the system (Table: users)
type User struct {
	ID        int       `json:"id"`
	FName     string    `json:"fname" db:"fname"`
	LName     string    `json:"lname" db:"lname"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	// Pwd (password hash) is stored in the DB but not included in this public struct
}

// UserCreate is used for user registration
type UserCreate struct {
	FName    string `json:"fname" validate:"required,min=1,max=30"`
	LName    string `json:"lname" validate:"required,min=1,max=60"`
	Email    string `json:"email" validate:"required,email,max=60"`
	Password string `json:"password" validate:"required,min=8"` // will be hashed to CHAR(60)
}

// UserLogin is used for authentication
type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// --- Artist and Artwork Models ---

// Artist represents an artist (Table: artists)
type Artist struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" validate:"required,min=1,max=60" db:"name"`
	Codename  string    `json:"codename,omitempty" validate:"max=60" db:"codename"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}

// Artwork represents a piece of art (Table: artworks)
type Artwork struct {
	ID          int       `json:"id"`
	ArtistID    int       `json:"artist_id" validate:"required" db:"artist_id"`
	Grade       string    `json:"grade,omitempty" validate:"max=20" db:"grade"`
	School      string    `json:"school,omitempty" validate:"max=30" db:"school"`
	Title       string    `json:"title,omitempty" validate:"max=100" db:"title"`
	Description string    `json:"description,omitempty" validate:"max=500" db:"description"`
	CreatedAt   time.Time `json:"created_at,omitempty" db:"created_at"`
}

// Image represents the image data for an artwork (Table: images)
type Image struct {
	ID           int       `json:"id"`
	ArtworkID    int       `json:"artwork_id" validate:"required" db:"artwork_id"`
	URL          string    `json:"url,omitempty" validate:"max=255" db:"url"`
	OriginalMime string    `json:"mime" validate:"max=50" db:"original_mime"`
	Thumb        []byte    `json:"thumb,omitempty" db:"thumb"`           // BLOB, max 64KB
	Image        []byte    `json:"image" validate:"required" db:"image"` // MEDIUMBLOB, full image
	CreatedAt    time.Time `json:"created_at,omitempty" db:"created_at"`
}

// Medium represents an art medium (Table: mediums)
type Medium struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"required,min=1,max=60" db:"name"`
}

// --- Relationship Models ---

// UserArtist represents the many-to-many relationship between users and artists (Table: user_artists)
type UserArtist struct {
	UserID   int `json:"user_id" db:"user_id"`
	ArtistID int `json:"artist_id" db:"artist_id"`
}

// ArtworkMedium represents the many-to-many relationship between artworks and mediums (Table: artworks_mediums)
type ArtworkMedium struct {
	ArtworkID int `json:"artwork_id" db:"artwork_id"`
	MediumID  int `json:"medium_id" db:"medium_id"`
}

// --- View Model ---

// ArtworkView represents the rich view of artwork data from the all_artwork_data VIEW
type ArtworkView struct {
	ArtworkID   int    `json:"artwork_id" db:"artwork_id"`
	Grade       string `json:"grade,omitempty" db:"grade"`
	School      string `json:"school,omitempty" db:"school"`
	Title       string `json:"title,omitempty" db:"title"`
	Description string `json:"description,omitempty" db:"description"`
	ArtistName  string `json:"artist_name" db:"artist_name"` // COALESCE(ar.codename, ar.name)
	URL         string `json:"url,omitempty" db:"url"`
	Thumb       []byte `json:"thumb,omitempty" db:"thumb"` // BLOB thumbnail
	Mediums     string `json:"mediums,omitempty" db:"mediums"`
}

// --- API Utility Models ---

// APIResponse is a standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginatedResponse wraps responses with pagination info
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	TotalPages int         `json:"total_pages"`
}
