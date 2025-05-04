package database

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

// DB is a global variable holding the GORM database connection instance.
// It's accessible throughout the application after being initialized by Connect.
var DB *gorm.DB

// Connect establishes the database connection using environment variables for configuration.
// It loads variables from a .env file, constructs the DSN (Data Source Name),
// and initializes the global DB variable with the GORM connection instance.
//
// Returns:
//   - error: An error if loading environment variables or connecting to the database fails, otherwise nil.
func Connect() error {
	// Load environment variables from a .env file in the current directory.
	// This allows for easy configuration without hardcoding credentials.
	if err := godotenv.Load(); err != nil {
		// Return an error if the .env file cannot be loaded.
		return fmt.Errorf("failed to load environment variables: %w", err)
	}

	// Retrieve database connection details from environment variables.
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// Construct the Data Source Name (DSN) string required by the PostgreSQL driver.
	dsn := CreateDSN(host, user, password, db_name, port)

	// Declare err variable to store potential errors from gorm.Open.
	var err error
	// Attempt to open a connection to the PostgreSQL database using the constructed DSN.
	// The result (connection instance or error) is assigned to the global DB variable and the local err variable.
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// If gorm.Open returns an error, wrap it with context and return it.
		// This allows the calling function (e.g., main) to handle the connection failure gracefully.
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Log a success message if the connection is established.
	log.Println("✅ Database connected")
	// Return nil to indicate a successful connection.
	return nil
}

// CreateDSN constructs the Data Source Name (DSN) string for connecting to a PostgreSQL database.
// It takes individual connection parameters and formats them into the required DSN format.
//
// Parameters:
//   - host: The database server hostname or IP address.
//   - user: The username for database authentication.
//   - password: The password for database authentication.
//   - dbName: The name of the database to connect to.
//   - port: The port number the database server is listening on.
//
// Returns:
//   - string: The formatted DSN string.
func CreateDSN(host, user, password, dbName, port string) string {
	// Use fmt.Sprintf to build the DSN string. sslmode=disable is used for simplicity;
	// consider enabling SSL (e.g., sslmode=require) for production environments.
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbName, port,
	)
}
