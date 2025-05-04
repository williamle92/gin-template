# Gin Web API Template with Docker, Traefik, PostgreSQL and GORM
This project is a starter template for building web APIs in Go using the Gin framework, set up for local development with Docker Compose, Traefik as an API Gateway, and GORM for database interactions with PostgreSQL.

## Table of Contents
- [Gin Web API Template with Docker, Traefik, PostgreSQL and GORM](#gin-web-api-template-with-docker-traefik-postgresql-and-gorm)
  - [Table of Contents](#table-of-contents)
  - [Technology Stack and Features](#technology-stack-and-features)
  - [Prerequisites](#prerequisites)
  - [Getting Started](#getting-started)
    - [Cloning the Repository](#cloning-the-repository)
    - [Environment Configuration](#environment-configuration)
    - [Starting the Services](#starting-the-services)
    - [Database Initialization](#database-initialization)
  - [Database Migrations](#database-migrations)
  - [Core Components](#core-components)
    - [Models](#models)
    - [Routes](#routes)


## Technology Stack and Features
* **[Go](https://go.dev/):** The primary programming language.
* **[Gin](https://github.com/gin-gonic/gin):** A fast, minimalist web framework for Go, used for building the API endpoints.
* **[Docker & Docker Compose](https://www.docker.com/):** For containerizing the application and its dependencies (database, Traefik) and orchestrating them for easy local development setup.
* **[Traefik](https://traefik.io/traefik/):** Used as an API Gateway/Reverse Proxy to manage incoming requests and route them to the correct services based on rules (like hostname).
* **[PostgreSQL](https://www.postgresql.org/):** A powerful open-source relational database.
* **[GORM](https://gorm.io/):** An Object-Relational Mapper (ORM) library for Go, simplifying database interactions (defining models, querying, creating, updating, deleting data, and handling migrations).
* **[Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt):** Used for securely hashing user passwords before storing them in the database.
* **[JWT](https://pkg.go.dev/github.com/golang-jwt/jwt/v5) (JSON Web Tokens):** Intended for stateless authentication (though the provided code snippet focuses on user creation, JWT would be used for subsequent login/authorized requests).
* **[Atlas](https://atlasgo.io/):** (Implied by your usage) A database schema management tool for managing and applying database migrations.


## Prerequisites
Before you begin, ensure you have the following installed on your machine:

* **Git:** For cloning the repository (`https://git-scm.com/downloads`).
* **Go:** The Go programming language (version 1.20 or higher recommended) (`https://golang.org/doc/install`).
* **Docker and Docker Compose:** (`https://docs.docker.com/get-docker/`). Docker Compose is usually bundled with Docker Desktop.
* **Atlas CLI:** For managing database migrations (`https://atlasgo.io/getting-started`).

## Getting Started

Follow these steps to get the project running locally:

### Cloning the Repository

Open your terminal and clone the repository:

```bash
git clone <repository_url> # Replace <repository_url> with the actual URL
cd gin-template # Navigate into the project directory
```

### Environment Configuration
Create a file named `.env` at the root of the project directory. This file will hold sensitive information and configuration variables for your services.

Copy the following structure into your .env file and fill in the values. Do not commit this file if it contains production secrets!
```
# Database Configuration (used by both docker-compose and your Go app)
DB_HOST=                 # The name of the database service in docker-compose
DB_PORT=                 # The internal port of the postgres container
DB_NAME=                 # The name of the database to create/connect to
DB_USER=                 # The database user
POSTGRES_PASSWORD=       # The database password

# Application Specific
JWT_SALT=                # Secret key for JWT signing
# Add any other app-specific variables here
```

### Starting the Services
Navigate to the root directory of the project in your terminal and run the following command:

```bash
docker compose up --build
```
- `docker compose up`: Starts the services defined in `docker-compose.yml`.
- `--build`: Builds the api service's Docker image before starting if it's not already built or has changed.

This command will:
1. Build your Go application's Docker image.
2. Start the database (PostgreSQL) container.
3. Start the traefik container.
4. Start the api container. The api container will wait for the database container to report as healthy before starting.

### Database Initialization
With the database service running, you need to:
1. Create the Database: The name of the database should match what you set as `DB_NAME`.
2. Run Migrations: You need to apply your database schema migrations using Atlas. Atlas needs to connect to the database running on your host machine via the exposed port. Make sure your Atlas configuration file (usually atlas.hcl at the root) is set up correctly to connect to the database via localhost:5433. Example atlas.hcl snippet:
```
# atlas.hcl
env "gorm" {
  # Set the URL using environment variables read from your .env file by Atlas
  url = "postgres://${DB_USER}:${DB_PASSWORD}@localhost:5433/${DB_NAME}?sslmode=disable"

  # Specify the directory containing your migration files
  # Ensure you have a migrations directory and generated migration files
  migration {
    dir = "file://migrations"
  }
}
```
3. Apply migrations: run `atlas migrate apply --env gorm`

## Database Migrations

Database migrations are version-controlled scripts that manage changes to your database schema over time (e.g., adding tables, columns, altering types, adding indexes). They provide a structured, trackable, and reliable way to evolve your database alongside your application code.

This project uses [Atlas](https://atlasgo.io/versioned/intro) to manage the database migrations. Atlas inspects the GORM models to understand the desired schema state and automatically generates SQL scripts to transition the database from its current state to the desired state. This allows us to:

* Keep schema changes in source control (`migrations` directory).
* Apply database updates reliably across different environments.
* Simplify collaboration on schema changes.

**How to Use Migrations:**

Follow these steps whenever you modify your GORM models (in the `./models` directory) and need to update the database schema:

1.  **Define the Desired Schema:** Modify your Go structs in the `./models` package to reflect the new desired state of your database schema (e.g., add a new field to a struct, change a field's type).

2.  **Generate the Migration Script:**
    Run the following command from your project root:

    ```bash
    atlas migrate diff --env gorm
    ```
    * This command compares the current state of your database schema (derived from the existing migration files in the `migrations` directory, applied to a temporary development database defined in `atlas.hcl`) with the desired state (derived from your GORM models, as configured by `src` in the `gorm` env).
    * Atlas will generate a new SQL file in the `./migrations` directory with the necessary SQL statements to bridge the difference.

3.  **Review the Generated Script (Critical Step):**
    Open the newly created `.sql` file in the `./migrations` directory. **Carefully review** the generated SQL statements to ensure they accurately represent the changes you intended and do not contain any unexpected or potentially destructive operations. Atlas is smart, but reviewing the generated SQL is essential before applying it to any real database. Edit the file manually if necessary.

4.  **Apply the Migration:**
    Once you are confident that the generated SQL script is correct, apply it to your target database. You will need to provide the URL for the database you want to apply the migration to (e.g., your local development database, a staging database, or production).

    ```bash
    atlas migrate apply --env gorm --url "postgres://user:pass@host:port/db_name?sslmode=disable"
    ```
    * Replace `"postgres://user:pass@host:port/db_name?sslmode=disable"` with the actual connection string for your target database.
    * This command executes the generated migration scripts (those in the `migrations` directory that haven't been applied yet, tracked by the `atlas_schema_revisions` table) against the specified `--url`.

**Important Notes:**

* The `atlas.hcl` file configures how Atlas interacts with your GORM models (`src`) and defines a **temporary development database** (`dev`) using a Docker container. This temporary database is used *only* during the `migrate diff` process for schema comparison and is separate from your actual database(s) used by the application.
* The `--env gorm` flag tells Atlas to use the configuration block named `gorm` in your `atlas.hcl`.
* The `--url` flag is required for `migrate apply` to specify the target database. Atlas will record applied migrations in the `atlas_schema_revisions` table in this target database.

## Core Components
### Models
Database models are defined as Go structs, this can be found in /models. GORM uses struct tags (gorm:"...", json:"...") to map the struct fields to database table columns and define column properties (like primaryKey, unique, not null, index, column names).

```
package models

import (
    "gorm.io/gorm"
)

type User struct {
    gorm.Model // Adds ID, CreatedAt, UpdatedAt, DeletedAt fields automatically
    ID          uint   `gorm:"primaryKey;autoIncrement;index" json:"id"` // Explicitly define ID if needed for json or custom gorm tags
    FirstName   string `json:"first_name" gorm:"not null"`
    LastName    string `json:"last_name" gorm:"not null"`
    Email       string `gorm:"unique;index;not null" json:"email"` // unique and indexed for quick lookups
    PhoneNumber string `gorm:"unique;index;not null" json:"phone_number"` // unique and indexed
    Password    string `json:"-"` // The "-" tag tells json marshaller to ignore this field (good for security)
}
```
- The `gorm.Model` embed adds standard fields automatically managed by GORM.
- `gorm:"..."` tags configure database column mapping and constraints.
- `json:"..."` tags configure how the struct is serialized/deserialized to/from JSON. json:"-" hides the field in JSON output.

### Routes
- Routes define the API endpoints (e.g., POST /user). These are set up using the Gin router in the main.go file.
- Handlers are Go functions that execute when a specific route is matched. They receive a `*gin.Context` which provides access to the request (body, headers, params, query) and methods for writing the response (JSON, status codes, etc.). Handlers are defined in the routes directory.