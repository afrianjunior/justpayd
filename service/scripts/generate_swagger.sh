#!/bin/bash

set -e  # Exit immediately if a command exits with a non-zero status

# Check if swag is installed
if ! command -v swag &> /dev/null; then
    echo "Error: 'swag' command not found."
    echo "Please install it with: go install github.com/swaggo/swag/cmd/swag@latest"
    echo "And make sure your GOPATH/bin is in your PATH environment variable:"
    echo "export PATH=\$PATH:\$(go env GOPATH)/bin"
    exit 1
fi

# Create docs directory if it doesn't exist
mkdir -p docs

# Create a temporary main.go file for Swagger generation
cat > swaggertemp.go << 'EOF'
// @title           JustPayd API
// @version         1.0
// @description     API for JustPayd services
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.

package main

import (
	_ "github.com/afrianjunior/justpayd/internal/auth"
	_ "github.com/afrianjunior/justpayd/internal/subscription"
	_ "github.com/afrianjunior/justpayd/internal/tenant"
	_ "github.com/afrianjunior/justpayd/internal/plan"
	_ "github.com/afrianjunior/justpayd/internal/user"
	_ "github.com/afrianjunior/justpayd/internal/file"
	// Add other handler packages as needed
)

// This is just a placeholder file for Swagger.
// Don't worry, we won't actually run this.
func main() {}
EOF

# Generate Swagger documentation using swag
echo "Generating Swagger documentation..."
swag init -g swaggertemp.go -o docs --parseInternal --parseDependency

# Clean up the temporary file
rm swaggertemp.go

# Verify the file was generated
if [ -f "docs/swagger.json" ]; then
  echo "Swagger documentation generated successfully"
  # Copy to static location as well if needed
  mkdir -p static/swagger
  cp docs/swagger.json static/swagger/swagger.json
else
  echo "Failed to generate Swagger documentation"
  exit 1
fi