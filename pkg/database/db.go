package database

import (
    "database/sql"
    "log"
    "os"

    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Connect initializes the database connection
func Connect() {
    var err error
    dsn := os.Getenv("DATABASE_URL") // Ensure to set this environment variable
    db, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
    }

    if err = db.Ping(); err != nil {
        log.Fatalf("Database is unreachable: %v", err)
    }

    log.Println("Database connection established")
}

// GetDB returns the database connection
func GetDB() *sql.DB {
    return db
}

// Close closes the database connection
func Close() {
    if err := db.Close(); err != nil {
        log.Fatalf("Error closing the database: %v", err)
    }
    log.Println("Database connection closed")
}