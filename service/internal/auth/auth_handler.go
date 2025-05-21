package auth

import (
	"encoding/json"
	"net/http"

	"github.com/afrianjunior/justpayd/internal/pkg"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService AuthService
	logger      *zap.SugaredLogger
	config      *pkg.Config
}

func NewAuthHandler(authService AuthService, logger *zap.SugaredLogger, config *pkg.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
		config:      config,
	}
}

func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	r.Post("/login", h.Login)
}

// @Summary User login
// @Description Authenticates a user and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body LoginRequest true "Login credentials"
// @Success 200 {object} pkg.BaseResponse{data=LoginResponse} "Successfully authenticated"
// @Failure 400 {object} pkg.BaseResponse "Invalid request format"
// @Failure 401 {object} pkg.BaseResponse "Invalid credentials"
// @Failure 500 {object} pkg.BaseResponse "Authentication failed"
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorw("Failed to decode login request", "error", err)
		pkg.JsonResponse(w, pkg.NewErrorResponse("Invalid request format"), http.StatusBadRequest)
		return
	}

	// Verify user credentials
	user, err := h.authService.VerifyCredentials(r.Context(), req.Email)
	if err != nil {
		h.logger.Errorw("Invalid credentials", "email", req.Email, "error", err)
		pkg.JsonResponse(w, pkg.NewErrorResponse("Invalid credentials"), http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := h.authService.GenerateToken(user.ID)
	if err != nil {
		h.logger.Errorw("Failed to generate token", "error", err)
		pkg.JsonResponse(w, pkg.NewErrorResponse("Authentication failed"), http.StatusInternalServerError)
		return
	}

	// Return the token and user ID
	response := LoginResponse{
		Token:  token,
		UserID: user.ID,
	}

	pkg.JsonResponse(w, pkg.SuccessResponse(response), http.StatusOK)
}
