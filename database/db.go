package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect establishes the database connection and assigns it to the global DB variable.
// It now returns an error if the connection fails.
func Connect() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("failed to load environment variables: %w", err)
	}
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := CreateDSN(host, user, password, db_name, port)
	var err error                                           // Declare err here
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}) // Assign to the global DB and check for error
	if err != nil {
		// Instead of log.Fatal, return the error so main.go can handle it.
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("✅ Database connected")
	return nil // Return nil if connection is successful
}

func CreateDSN(host, user, password, dbName, port string) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbName, port,
	)
}
