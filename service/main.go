package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/afrianjunior/justpayd/cmd"
	"github.com/afrianjunior/justpayd/internal/pkg"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"go.uber.org/zap"
)

type App struct {
	db     *sql.DB
	logger *zap.SugaredLogger
	config *pkg.Config
}

func loadConfigFromEnv() *pkg.Config {
	config := &pkg.Config{}

	config.JWT.Secret = os.Getenv("JWT_SECRET")
	expStr := os.Getenv("JWT_EXPIRATION")
	if expStr != "" {
		exp, err := strconv.Atoi(expStr)
		if err == nil {
			config.JWT.Expiration = exp
		}
	}

	if config.JWT.Secret == "" {
		config.JWT.Secret = "secret"
	}
	if config.JWT.Expiration == 0 {
		config.JWT.Expiration = 3600
	}

	config.StoragePath = os.Getenv("STORAGE_PATH")
	if config.StoragePath == "" {
		config.StoragePath = "./data"
	}
	config.ServerPort = os.Getenv("SERVER_PORT")
	if config.ServerPort == "" {
		config.ServerPort = "8080"
	}

	config.LogLevel = os.Getenv("LOG_LEVEL")
	if config.LogLevel == "" {
		config.LogLevel = "info"
	}

	return config
}

func setupLogger(level string) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	if level == "debug" {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("error creating logger: %v", err)
	}
	return logger.Sugar(), nil
}

func NewApp(ctx context.Context, config *pkg.Config) (*App, error) {
	logger, err := setupLogger(config.LogLevel)
	if err != nil {
		return nil, err
	}

	// Ensure data directory exists
	if err := os.MkdirAll(config.StoragePath, 0755); err != nil {
		fmt.Println(config.StoragePath)
		return nil, fmt.Errorf("error creating data directory: %v", err)
	}

	// SQLite database file path
	sqliteDBPath := fmt.Sprintf("%s/main.db", config.StoragePath)

	// Initialize SQLite connection
	db, err := sql.Open("sqlite3", sqliteDBPath)
	if err != nil {
		return nil, fmt.Errorf("error opening sqlite database: %v", err)
	}

	// Ping the database to ensure the connection is live
	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error pinging sqlite database: %v", err)
	}

	return &App{
		db:     db,
		logger: logger,
		config: config,
	}, nil
}

func main() {
	ctx := context.Background()
	fmt.Println("Starting application...")
	config := loadConfigFromEnv()

	app, err := NewApp(ctx, config)
	if err != nil {
		log.Fatalf("Error creating application: %v", err)
	}

	defer app.db.Close()

	// Create a REST server with both API and documentation capabilities
	restServer := cmd.NewRest(
		app.db,
		app.logger,
		app.config,
	)

	// Generate Swagger documentation first
	log.Println("Generating API documentation...")
	err = generateSwaggerDocs()
	if err != nil {
		log.Printf("Warning: Could not generate API documentation: %v", err)
	} else {
		// Verify that the swagger.json file was created
		if _, err := os.Stat("docs/swagger.json"); err == nil {
			log.Println("API documentation generated successfully")
		} else {
			log.Printf("Warning: swagger.json file was not found after generation: %v", err)
		}
	}

	// Start the server with both API and documentation
	log.Printf("Starting server on port %s...", app.config.ServerPort)
	log.Printf("API documentation available at http://localhost:%s/reference", app.config.ServerPort)
	restServer.Start(app.config.ServerPort)
}

// generateSwaggerDocs runs the script to generate Swagger documentation
func generateSwaggerDocs() error {
	// Ensure the docs directory exists
	if err := os.MkdirAll("docs", 0755); err != nil {
		return fmt.Errorf("error creating docs directory: %v", err)
	}

	// Run the script to generate Swagger documentation
	scriptPath := "./scripts/generate_swagger.sh"
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("documentation script not found: %s", scriptPath)
	}

	cmd := exec.Command("/bin/sh", scriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error generating documentation: %v - %s", err, string(output))
	}

	return nil
}
