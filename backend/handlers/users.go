package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"go-art-api/config"
	"go-art-api/models"
	"go-art-api/utils"

	"github.com/gorilla/mux"
)

// GetUsers retrieves all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query("SELECT id, fname, lname, email, created_at FROM users ORDER BY id")
	if err != nil {
		sendErrorResponse(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.FName, &u.LName, &u.Email, &u.CreatedAt); err != nil {
			sendErrorResponse(w, "Failed to scan user data", http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	sendJSONResponse(w, users, http.StatusOK)
}

// CreateUser creates a new user (admin only - basic version)
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var userCreate models.UserCreate
	if err := json.NewDecoder(r.Body).Decode(&userCreate); err != nil {
		sendErrorResponse(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Basic validation
	if userCreate.FName == "" || userCreate.LName == "" || userCreate.Email == "" || userCreate.Password == "" {
		sendErrorResponse(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Hash the password using Argon2
	hashedPassword, err := utils.HashPassword(userCreate.Password)
	if err != nil {
		sendErrorResponse(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Insert user into database
	result, err := config.DB.Exec(
		"INSERT INTO users (fname, lname, email, pwd) VALUES (?, ?, ?, ?)",
		userCreate.FName, userCreate.LName, userCreate.Email, hashedPassword,
	)
	if err != nil {
		// Check for duplicate email
		if contains(err.Error(), "Duplicate entry") {
			sendErrorResponse(w, "Email already exists", http.StatusConflict)
			return
		}
		sendErrorResponse(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Get the created user ID
	id, _ := result.LastInsertId()

	// Return the created user (without password)
	user := models.User{
		ID:    int(id),
		FName: userCreate.FName,
		LName: userCreate.LName,
		Email: userCreate.Email,
	}

	sendJSONResponse(w, user, http.StatusCreated)
}

// RegisterUser handles user registration
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var userCreate models.UserCreate
	if err := json.NewDecoder(r.Body).Decode(&userCreate); err != nil {
		sendErrorResponse(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Validate input
	if err := validateUserCreate(userCreate); err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(userCreate.Password)
	if err != nil {
		sendErrorResponse(w, "Failed to process registration", http.StatusInternalServerError)
		return
	}

	// Insert user
	result, err := config.DB.Exec(
		"INSERT INTO users (fname, lname, email, pwd) VALUES (?, ?, ?, ?)",
		userCreate.FName, userCreate.LName, userCreate.Email, hashedPassword,
	)
	if err != nil {
		if contains(err.Error(), "Duplicate entry") {
			sendErrorResponse(w, "Email already registered", http.StatusConflict)
			return
		}
		sendErrorResponse(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()

	response := models.APIResponse{
		Success: true,
		Message: "Registration successful",
		Data: map[string]interface{}{
			"user_id": int(id),
			"email":   userCreate.Email,
		},
	}

	sendJSONResponse(w, response, http.StatusCreated)
}

// LoginUser handles user authentication
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var login models.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		sendErrorResponse(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if login.Email == "" || login.Password == "" {
		sendErrorResponse(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Get user from database
	var user models.User
	var hashedPassword string
	err := config.DB.QueryRow(
		"SELECT id, fname, lname, email, pwd FROM users WHERE email = ?",
		login.Email,
	).Scan(&user.ID, &user.FName, &user.LName, &user.Email, &hashedPassword)

	if err == sql.ErrNoRows {
		sendErrorResponse(w, "Invalid email or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		sendErrorResponse(w, "Login failed", http.StatusInternalServerError)
		return
	}

	// Verify password
	isValid, err := utils.VerifyPassword(login.Password, hashedPassword)
	if err != nil {
		sendErrorResponse(w, "Login failed", http.StatusInternalServerError)
		return
	}

	if !isValid {
		sendErrorResponse(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// In a real app, you'd generate a JWT token here
	response := models.APIResponse{
		Success: true,
		Message: "Login successful",
		Data: map[string]interface{}{
			"user":  user,
			"token": "jwt-token-would-go-here",
		},
	}

	sendJSONResponse(w, response, http.StatusOK)
}

// GetUserByID retrieves a user by ID
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var u models.User
	err = config.DB.QueryRow("SELECT id, fname, lname, email FROM users WHERE id = ?", id).
		Scan(&u.ID, &u.FName, &u.LName, &u.Email)

	if err == sql.ErrNoRows {
		sendErrorResponse(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		sendErrorResponse(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, u, http.StatusOK)
}

// UpdateUser updates an existing user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var userUpdate models.User
	if err := json.NewDecoder(r.Body).Decode(&userUpdate); err != nil {
		sendErrorResponse(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Basic validation
	if userUpdate.FName == "" || userUpdate.LName == "" || userUpdate.Email == "" {
		sendErrorResponse(w, "First name, last name, and email are required", http.StatusBadRequest)
		return
	}

	// Update user
	_, err = config.DB.Exec(
		"UPDATE users SET fname = ?, lname = ?, email = ? WHERE id = ?",
		userUpdate.FName, userUpdate.LName, userUpdate.Email, id,
	)
	if err != nil {
		if contains(err.Error(), "Duplicate entry") {
			sendErrorResponse(w, "Email already exists", http.StatusConflict)
			return
		}
		sendErrorResponse(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// Return updated user
	userUpdate.ID = id
	sendJSONResponse(w, userUpdate, http.StatusOK)
}

// DeleteUser deletes a user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErrorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Check if user exists
	var exists bool
	err = config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		sendErrorResponse(w, "Failed to check user", http.StatusInternalServerError)
		return
	}

	if !exists {
		sendErrorResponse(w, "User not found", http.StatusNotFound)
		return
	}

	// Delete user
	_, err = config.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		sendErrorResponse(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper functions
func validateUserCreate(u models.UserCreate) error {
	if len(u.FName) == 0 || len(u.FName) > 30 {
		return errors.New("first name must be 1-30 characters")
	}
	if len(u.LName) == 0 || len(u.LName) > 60 {
		return errors.New("last name must be 1-60 characters")
	}
	if len(u.Email) == 0 || len(u.Email) > 60 {
		return errors.New("email must be 1-60 characters")
	}
	if len(u.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	// Add email format validation here if needed
	return nil
}
