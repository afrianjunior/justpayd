package cmd

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/afrianjunior/justpayd/internal/assignments"
	"github.com/afrianjunior/justpayd/internal/auth"
	"github.com/afrianjunior/justpayd/internal/pkg"
	"github.com/afrianjunior/justpayd/internal/shift_requests"
	"github.com/afrianjunior/justpayd/internal/shifts"
	"github.com/afrianjunior/justpayd/internal/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	// Import Swagger docs to register documentation
	_ "github.com/afrianjunior/justpayd/docs"
)

type APIError struct {
	Error string `json:"error"`
}

type Rest interface {
	Start(port string)
}

type rest struct {
	db     *sql.DB
	logger *zap.SugaredLogger
	config *pkg.Config
}

func NewRest(
	db *sql.DB,
	logger *zap.SugaredLogger,
	config *pkg.Config,
) Rest {
	return &rest{
		db:     db,
		logger: logger,
		config: config,
	}
}

func (s *rest) Start(port string) {
	router := s.setupRouter()
	s.logger.Infow("Starting server", "port", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		s.logger.Fatalf("Server error: %v", err)
	}
}

func (s *rest) setupRouter() *chi.Mux {
	r := chi.NewRouter()

	// Initialize repositories
	userRepository := users.NewUserRepository(s.db)
	shiftRepository := shifts.NewShiftRepository(s.db)
	shiftRequestRepository := shift_requests.NewShiftRequestRepository(s.db)
	authRepository := auth.NewAuthRepository(s.db)
	assignmentRepository := assignments.NewAssignmentRepository(s.db)

	// Initialize services
	userService := users.NewUserService(userRepository)
	shiftService := shifts.NewShiftService(shiftRepository)
	shiftRequestService := shift_requests.NewShiftRequestService(shiftRequestRepository, assignmentRepository)
	authService := auth.NewAuthService(authRepository, s.config)
	assignmentService := assignments.NewAssignmentService(assignmentRepository)

	// Initialize handlers
	userHandler := users.NewUserHandler(userService, s.logger)
	shiftHandler := shifts.NewShiftHandler(shiftService, s.logger)
	shiftRequestHandler := shift_requests.NewShiftRequestHandler(shiftRequestService, s.logger)
	authHandler := auth.NewAuthHandler(authService, s.logger, s.config)
	assignmentHandler := assignments.NewAssignmentHandler(assignmentService, s.logger)

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Static file server for documentation assets
	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// API Documentation routes
	r.Get("/reference", s.serveDocumentationPage)
	r.Get("/api/swagger.json", s.serveSwaggerJSON)

	// Root redirect to documentation
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/reference", http.StatusFound)
	})

	// API Routes
	r.Route("/api", func(r chi.Router) {
		// Public routes (no authentication required)
		r.Route("/auth", func(r chi.Router) {
			authHandler.RegisterRoutes(r)
		})

		// Protected routes (authentication required)
		r.Group(func(r chi.Router) {
			// Apply JWT middleware to all routes in this group
			r.Use(pkg.RequireAuth(s.config, s.db))

			r.Route("/users", func(r chi.Router) {
				userHandler.RegisterRoutes(r)
			})
			r.Route("/shifts", func(r chi.Router) {
				shiftHandler.RegisterRoutes(r)
			})
			r.Route("/shift_requests", func(r chi.Router) {
				shiftRequestHandler.RegisterRoutes(r)
			})
			r.Route("/assignments", func(r chi.Router) {
				assignmentHandler.RegisterRoutes(r)
			})
		})
	})

	return r
}

// serveSwaggerJSON serves the Swagger JSON file
func (s *rest) serveSwaggerJSON(w http.ResponseWriter, r *http.Request) {
	// Check in docs directory first
	swaggerFile := filepath.Join("docs", "swagger.json")
	if _, err := os.Stat(swaggerFile); os.IsNotExist(err) {
		// Fall back to static directory
		swaggerFile = filepath.Join("static", "swagger", "swagger.json")
		if _, err := os.Stat(swaggerFile); os.IsNotExist(err) {
			s.logger.Errorw("Swagger JSON file not found in either location",
				"docsPath", filepath.Join("docs", "swagger.json"),
				"staticPath", filepath.Join("static", "swagger", "swagger.json"))
			http.Error(w, "Swagger file not found", http.StatusNotFound)
			return
		}
	}

	data, err := os.ReadFile(swaggerFile)
	if err != nil {
		http.Error(w, "Failed to read Swagger file", http.StatusInternalServerError)
		s.logger.Errorw("Failed to read swagger file", "error", err, "path", swaggerFile)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// serveDocumentationPage serves the documentation page with Swagger UI using scalar library
func (s *rest) serveDocumentationPage(w http.ResponseWriter, r *http.Request) {
	// Use relative path to current host instead of absolute file path
	specUrl := "./docs/swagger.json"

	// For local development, verify the file exists first
	swaggerFile := filepath.Join("docs", "swagger.json")
	if _, err := os.Stat(swaggerFile); os.IsNotExist(err) {
		s.logger.Warnw("Swagger file not found, documentation may not display correctly", "path", swaggerFile)
	}

	htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecURL: specUrl,
		CustomOptions: scalar.CustomOptions{
			PageTitle: "API Documentation",
		},
		DarkMode: true,
	})

	if err != nil {
		s.logger.Errorw("Failed to generate API documentation page", "error", err)
		http.Error(w, "Error generating documentation page", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, htmlContent)
}
