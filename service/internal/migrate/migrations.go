package migrate

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// Migration defines a database migration
type Migration struct {
	Name string
	Up   string
	Down string
}

// Migrations is the list of all migrations for the application
var Migrations = []Migration{
	{
		Name: "01_create_tenants_table",
		Up: `
			CREATE TABLE IF NOT EXISTS tenants (
				id VARCHAR(36) PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				description TEXT,
				created_at TIMESTAMP WITH TIME ZONE NOT NULL,
				updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
				status VARCHAR(50) NOT NULL DEFAULT 'active',
				plan_type VARCHAR(50) NOT NULL DEFAULT 'basic',
				domain VARCHAR(255) UNIQUE,
				settings JSONB NOT NULL DEFAULT '{}'::JSONB
			);
			
			CREATE INDEX idx_tenant_status ON tenants(status);
			CREATE INDEX idx_tenant_domain ON tenants(domain);
		`,
		Down: `
			DROP TABLE IF EXISTS tenants;
		`,
	},
	{
		Name: "02_create_users_table",
		Up: `
			CREATE TABLE IF NOT EXISTS users (
				id VARCHAR(36) PRIMARY KEY,
				tenant_id VARCHAR(36) NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
				email VARCHAR(255) NOT NULL UNIQUE,
				username VARCHAR(255) NOT NULL UNIQUE,
				first_name VARCHAR(255),
				last_name VARCHAR(255),
				hashed_password VARCHAR(255) NOT NULL,
				role VARCHAR(50) NOT NULL DEFAULT 'user',
				status VARCHAR(50) NOT NULL DEFAULT 'active',
				created_at TIMESTAMP WITH TIME ZONE NOT NULL,
				updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
				last_login_at TIMESTAMP WITH TIME ZONE
			);
			
			CREATE INDEX idx_user_tenant_id ON users(tenant_id);
			CREATE INDEX idx_user_email ON users(email);
			CREATE INDEX idx_user_username ON users(username);
			CREATE INDEX idx_user_status ON users(status);
		`,
		Down: `
			DROP TABLE IF EXISTS users;
		`,
	},
	{
		Name: "03_create_files_table",
		Up: `
			CREATE TABLE IF NOT EXISTS files (
				id VARCHAR(36) PRIMARY KEY,
				tenant_id VARCHAR(36) NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
				uploader_id VARCHAR(36) NOT NULL REFERENCES users(id),
				file_name VARCHAR(255) NOT NULL,
				file_size BIGINT NOT NULL,
				content_type VARCHAR(255) NOT NULL,
				storage_path VARCHAR(512) NOT NULL,
				created_at TIMESTAMP WITH TIME ZONE NOT NULL,
				updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
				tags TEXT[] DEFAULT '{}',
				is_public BOOLEAN NOT NULL DEFAULT FALSE,
				expires_at TIMESTAMP WITH TIME ZONE
			);
			
			CREATE INDEX idx_file_tenant_id ON files(tenant_id);
			CREATE INDEX idx_file_uploader_id ON files(uploader_id);
			CREATE INDEX idx_file_created_at ON files(created_at);
			CREATE INDEX idx_file_is_public ON files(is_public);
			CREATE INDEX idx_file_expired ON files(expires_at) WHERE expires_at IS NOT NULL;
		`,
		Down: `
			DROP TABLE IF EXISTS files;
		`,
	},
}

// RunMigrations executes all migrations
func RunMigrations(ctx context.Context, db *pgx.Conn) error {
	// Create migrations table if it doesn't exist
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		);
	`)

	if err != nil {
		return fmt.Errorf("error creating migrations table: %w", err)
	}

	// Get applied migrations
	rows, err := db.Query(ctx, "SELECT name FROM migrations ORDER BY id")
	if err != nil {
		return fmt.Errorf("error getting applied migrations: %w", err)
	}
	defer rows.Close()

	appliedMigrations := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("error scanning migration name: %w", err)
		}
		appliedMigrations[name] = true
	}

	// Run new migrations
	for _, migration := range Migrations {
		if appliedMigrations[migration.Name] {
			continue
		}

		// Start transaction
		tx, err := db.Begin(ctx)
		if err != nil {
			return fmt.Errorf("error starting transaction: %w", err)
		}

		// Run migration
		_, err = tx.Exec(ctx, migration.Up)
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("error running migration %s: %w", migration.Name, err)
		}

		// Record migration
		_, err = tx.Exec(ctx, "INSERT INTO migrations (name) VALUES ($1)", migration.Name)
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("error recording migration %s: %w", migration.Name, err)
		}

		// Commit transaction
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("error committing migration %s: %w", migration.Name, err)
		}

		fmt.Printf("Applied migration: %s\n", migration.Name)
	}

	return nil
}

// RollbackMigration rolls back the last applied migration
func RollbackMigration(ctx context.Context, db *pgx.Conn) error {
	// Get the last applied migration
	var id int
	var name string
	err := db.QueryRow(ctx, "SELECT id, name FROM migrations ORDER BY id DESC LIMIT 1").Scan(&id, &name)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("no migrations to rollback")
		}
		return fmt.Errorf("error getting last migration: %w", err)
	}

	// Find the migration
	var migration *Migration
	for i := range Migrations {
		if Migrations[i].Name == name {
			migration = &Migrations[i]
			break
		}
	}

	if migration == nil {
		return fmt.Errorf("migration %s not found", name)
	}

	// Start transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	// Run rollback
	_, err = tx.Exec(ctx, migration.Down)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("error rolling back migration %s: %w", name, err)
	}

	// Remove migration record
	_, err = tx.Exec(ctx, "DELETE FROM migrations WHERE id = $1", id)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("error removing migration record %s: %w", name, err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing rollback %s: %w", name, err)
	}

	fmt.Printf("Rolled back migration: %s\n", name)
	return nil
}
