# JustPayd Service

A microservice for managing employee shifts, assignments, and shift requests.

## Features

- **Shift Management**: Create, retrieve, update, and delete work shifts
- **User Assignment**: Assign users to shifts and manage assignments
- **Shift Requests**: Allow users to request shifts and approve/reject those requests
- **User Authentication**: Secure API access with JWT authentication
- **Interactive API Documentation**: Swagger UI for exploring and testing API endpoints

## Tech Stack

- **Backend**: Go (Golang)
- **Database**: SQLite
- **API**: RESTful API with Chi router
- **Authentication**: JWT tokens
- **Documentation**: Swagger/OpenAPI
- **Containerization**: Docker

## Prerequisites

- Go 1.23
- Docker (optional, for containerization)

## Getting Started

### Running Locally

1. Clone the repository
   ```bash
   git clone https://github.com/afrianjunior/justpayd.git
   cd justpayd/service
   ```

2. Install dependencies
   ```bash
   go mod download
   ```

3. Set up environment variables (optional)
   ```bash
   export JWT_SECRET=your_secret_key
   export JWT_EXPIRATION=3600
   export STORAGE_PATH=./data
   export SERVER_PORT=8080
   export LOG_LEVEL=info
   ```

4. Run the application
   ```bash
   go run main.go
   ```

5. Access the API at http://localhost:8080/api and the documentation at http://localhost:8080/reference

### Using Docker

1. Build the Docker image
   ```bash
   docker build -t justpayd-service .
   ```

2. Run the container
   ```bash
   docker run -p 8080:8080 justpayd-service
   ```

## API Endpoints

### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration

### Users
- `GET /api/users` - List all users
- `GET /api/users/{id}` - Get user by ID

### Shifts
- `GET /api/shifts` - List all shifts
- `POST /api/shifts` - Create a new shift
- `GET /api/shifts/{id}` - Get shift by ID
- `PUT /api/shifts/{id}` - Update a shift
- `DELETE /api/shifts/{id}` - Delete a shift

### Assignments
- `GET /api/assignments` - List all assignments
- `POST /api/assignments` - Create a new assignment
- `PUT /api/assignments/{id}` - Update an assignment

### Shift Requests
- `GET /api/shift_requests` - List all shift requests
- `POST /api/shift_requests` - Create a new shift request
- `PUT /api/shift_requests/{id}/approve` - Approve a shift request
- `PUT /api/shift_requests/{id}/reject` - Reject a shift request

## Project Structure

```
justpayd/
├── cmd/                # Command line applications
├── internal/           # Internal packages
│   ├── assignments/    # Assignment management
│   ├── auth/           # Authentication
│   ├── pkg/            # Shared packages
│   ├── shift_requests/ # Shift request management
│   ├── shifts/         # Shift management
│   └── users/          # User management
├── data/               # SQLite database storage
├── docs/               # API documentation
├── scripts/            # Utility scripts
├── static/             # Static assets
├── Dockerfile          # Docker configuration
└── main.go             # Application entry point
```

## Development

### Installing Nix

The project uses Nix for reproducible development environments and to run database migrations. To install Nix:

```bash
curl --proto '=https' --tlsv1.2 -sSf -L https://install.determinate.systems/nix | sh -s -- install
```

After installation, you can enter the development environment and run migrations with:

```bash
nix develop .
```

Once in the Nix development environment, you'll have access to the following commands:

- `create-migration [name]`: Create a new migration file with the specified name
- `run-migration`: Run all pending migrations
- `rollback`: Rollback the last migration
- `drop-migration`: Drop all migrations (resets the database)
- `generate-swagger`: Generate Swagger documentation

The Nix environment also provides these development tools:
- Go compiler and standard library
- go-migrate for database migrations
- gopls (Go language server)
- gotools (Go development tools)
- golangci-lint (linter for Go)
- swag (Swagger documentation generator)

### Database Migrations

Database schema migrations are handled automatically on startup. You can also manually run migrations using Nix as mentioned above.

### Generating Swagger Documentation

The API documentation is automatically generated on startup. To manually generate it:

```bash
./scripts/generate_swagger.sh
```