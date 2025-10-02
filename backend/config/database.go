package config

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv" // Library to load .env file
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	var err error

	// Load .env file (if it exists)
	// This is where "taco:cat@tcp(mysql.tacocat.com:3306)/tacocat?parseTime=true" gets loaded.
	godotenv.Load()

	// Get database URL from environment
	dsn := os.Getenv("DATABASE_URL")

	// Check if the environment variable is set
	if dsn == "" {
		log.Fatal("❌ FATAL: DATABASE_URL environment variable is not set. Did you create your .env file?")
	}

	// Open database connection
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		// This will print an error if the connection string format is bad or the DB is unreachable
		log.Fatal("❌ Failed to connect to database:", err)
	}

	// Configure connection pool (Keep these for performance)
	DB.SetMaxOpenConns(25)                 // Maximum open connections
	DB.SetMaxIdleConns(25)                 // Maximum idle connections
	DB.SetConnMaxLifetime(5 * time.Minute) // Connection max lifetime

	// Test the connection
	if err = DB.Ping(); err != nil {
		log.Fatal("❌ Failed to ping database:", err)
	}

	log.Println("✅ Database connected successfully")

	// Optionally create tables and views if they don't exist
	createTablesIfNotExist() // Assuming you uncommented this line to fix the "unused" error
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("❌ Error closing database: %v", err)
		} else {
			log.Println("✅ Database connection closed")
		}
	}
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return DB
}

// createTablesIfNotExist creates tables and views if they don't exist (optional)
func createTablesIfNotExist() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
            id INT AUTO_INCREMENT PRIMARY KEY,
            fname VARCHAR(30) NOT NULL,
            lname VARCHAR(60) NOT NULL,
            email VARCHAR(60) NOT NULL UNIQUE,
            pwd CHAR(128) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            INDEX idx_users_email (email)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		`CREATE TABLE IF NOT EXISTS artists (
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(60) NOT NULL,
            codename VARCHAR(60),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		`CREATE TABLE IF NOT EXISTS user_artists (
            user_id INT NOT NULL,
            artist_id INT NOT NULL,
            PRIMARY KEY(user_id, artist_id),
            FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
            FOREIGN KEY(artist_id) REFERENCES artists(id) ON DELETE CASCADE
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		`CREATE TABLE IF NOT EXISTS artworks (
            id INT AUTO_INCREMENT PRIMARY KEY,
            artist_id INT NOT NULL,
            grade VARCHAR(20),
            school VARCHAR(30),
			title VARCAR(100),
            description VARCHAR(500),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY(artist_id) REFERENCES artists(id) ON DELETE CASCADE,
            INDEX idx_artworks_artist_id (artist_id)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		`CREATE TABLE IF NOT EXISTS mediums (
            id INT AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(60) NOT NULL UNIQUE
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		`CREATE TABLE IF NOT EXISTS artworks_mediums (
            artwork_id INT NOT NULL,
            medium_id INT NOT NULL,
            PRIMARY KEY(artwork_id, medium_id),
            FOREIGN KEY(artwork_id) REFERENCES artworks(id) ON DELETE CASCADE,
            FOREIGN KEY(medium_id) REFERENCES mediums(id) ON DELETE CASCADE
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		`CREATE TABLE IF NOT EXISTS images (
            id INT AUTO_INCREMENT PRIMARY KEY,
            artwork_id INT NOT NULL UNIQUE,
            url VARCHAR(255),
			original_mime VARCHAR(50) NOT NULL, -- e.g., 'image/png', 'image/gif'
            thumb BLOB,
            image MEDIUMBLOB NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY(artwork_id) REFERENCES artworks(id) ON DELETE CASCADE,
            UNIQUE INDEX idx_images_artwork_id (artwork_id)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,

		// --- VIEW DEFINITION ---
		`
		-- Collates together all the artist/artwork/medium data + a thumbnail and a link.
		-- Uses CREATE OR REPLACE VIEW to handle existence and updates in a single statement.
		CREATE OR REPLACE VIEW all_artwork_data AS
		SELECT
			a.id AS artwork_id,
			a.grade,
			a.school,
			a.title,
			a.description,
			COALESCE(ar.codename, ar.name) AS artist_name,
			i.url,
			i.thumb, -- BLOB thumbnail
			GROUP_CONCAT(m.name ORDER BY m.name SEPARATOR ', ') AS mediums
		FROM artworks a
		JOIN artists ar ON a.artist_id = ar.id
		LEFT JOIN images i ON a.id = i.artwork_id
		LEFT JOIN artworks_mediums am ON a.id = am.artwork_id
		LEFT JOIN mediums m ON am.medium_id = m.id
		GROUP BY
			a.id,
			a.grade,
			a.school,
			a.title,
			a.description,
			ar.codename,
			ar.name,
			i.url,
			i.thumb
		ORDER BY a.id;
		`,
		// -----------------------
	}

	for _, query := range queries {
		if _, err := DB.Exec(query); err != nil {
			log.Printf("⚠️  Error executing DDL query: %v\nQuery: %s", err, query)
		}
	}

	log.Println("✅ Tables and views verified/created")
}
