{
  description = "A Nix-flake-based Node.js development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }:

    utils.lib.eachDefaultSystem (system:
      let
        goVersion = 22;
        dbName = "agentdb";
        dbUser = "agentuser";
        dbPassword = "agentpassword";
        dbLocation = "sqlite3://data/main.db";
        backendPath = "";

      overlays = [
        (final: prev: {
          go = prev."go_1_${toString goVersion}";
        })
      ];

      pkgs = import nixpkgs { inherit overlays system; };
      
      # Define swag package
      swag = pkgs.buildGoModule rec {
        pname = "swag";
        version = "1.16.4";
        
        src = pkgs.fetchFromGitHub {
          owner = "swaggo";
          repo = "swag";
          rev = "v${version}";
          sha256 = "sha256-wqBT7uan5XL51HHDGINRH9NTb1tybF44d/rWRxl6Lak=";
        };
        
        vendorHash = "sha256-6L5LzXtYjrA/YKmNEC/9dyiHpY/8gkH/CvW0JTo+Bwc=";
        
        subPackages = [ "cmd/swag" ];
        
        meta = with pkgs.lib; {
          description = "Automatically generate RESTful API documentation with Swagger 2.0 for Go";
          homepage = "https://github.com/swaggo/swag";
          license = licenses.mit;
          maintainers = with maintainers; [ ];
        };
      };
      
      scripts = with pkgs; [
        (writeScriptBin "create-migration" ''
          migrate create -ext sql -dir $BACKEND_MIGRATION_DIR -seq $@
        '')

        (writeScriptBin "run-migration" ''
          migrate -path $BACKEND_MIGRATION_DIR -database ${dbLocation} up
        '')

        (writeScriptBin "drop-migration" ''
          migrate -path $BACKEND_MIGRATION_DIR -database ${dbLocation} drop
        '')

        (writeScriptBin "rollback" ''
          migrate -path $BACKEND_MIGRATION_DIR -database ${dbLocation} down
        '')
        
        (writeShellApplication {
          name = "generate-swagger";
          runtimeInputs = [ swag ];
          text = ''
            # Create docs directory if it doesn't exist
            mkdir -p docs

            # Create a temporary main.go file for Swagger generation
            cat > swaggertemp.go << 'EOFSWAGGER'
            // @title           Univ Core API
            // @version         1.0
            // @description     API for Univ Core services
            // @termsOfService  http://swagger.io/terms/

            // @contact.name   API Support
            // @contact.url    http://www.example.com/support
            // @contact.email  support@example.com

            // @license.name  Apache 2.0
            // @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

            // @host      localhost:8000
            // @BasePath  /api

            // @securityDefinitions.apikey BearerAuth
            // @in header
            // @name Authorization
            // @description Type "Bearer" followed by a space and the JWT token.

            package main

            import (
                _ "github.com/afrianjunior/justpayd/internal/auth"
                _ "github.com/afrianjunior/justpayd/internal/shift"
                _ "github.com/afrianjunior/justpayd/internal/shift_requests"
                _ "github.com/afrianjunior/justpayd/internal/users"
                _ "github.com/afrianjunior/justpayd/internal/assignments"
                // Add other handler packages as needed
            )

            // This is just a placeholder file for Swagger.
            // Don't worry, we won't actually run this.
            EOFSWAGGER

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
          '';
        })
      ];

      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            go-migrate
            gopls
            gotools
            golangci-lint
            swag # Add swag to the shell
          ] ++ scripts;

          shellHook = with pkgs;''
            export GOPATH="$(${go}/bin/go env GOPATH)";
            export PATH="$PATH:$GOPATH/bin";
            export PROJECT_DIR=$(pwd);
            export BACKEND_DIR="$PROJECT_DIR/${backendPath}";
            export BACKEND_MIGRATION_DIR="$BACKEND_DIR/migrations";
            
            echo ""
            echo ""
            echo "[info]===================================================================="
            echo ""
            echo "database location:"
            echo "${dbLocation}"
            echo ""

            echo "packages:"
            echo "- `${go}/bin/go version`"
            echo "- swag: $(swag --version || echo 'not available')"

            echo ""
            echo "commands:"
            echo "- create-migration: Create a new migration file"
            echo "- run-migration: Run migrations"
            echo "- rollback: Rollback the last migration"
            echo "- drop-migration: Drop all migrations"
            echo "- generate-swagger: Generate Swagger documentation"
            echo ""
            echo "[info]===================================================================="
            echo ""
            echo ""
          '';
        };

      });
}
