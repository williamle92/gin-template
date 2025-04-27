### Database Migrations

Database migrations are version-controlled scripts that manage changes to your database schema over time (e.g., adding tables, columns, altering types, adding indexes). They provide a structured, trackable, and reliable way to evolve your database alongside your application code.

We use [Atlas](https://atlasgo.io/) to manage our database migrations. Atlas inspects our GORM models to understand the desired schema state and automatically generates SQL scripts to transition the database from its current state to the desired state. This allows us to:

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

---

This documentation provides a good overview of migrations, why you're using Atlas, and a clear step-by-step guide for your development workflow, including the crucial review step and the distinction between the development and target databases.