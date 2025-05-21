package users

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/afrianjunior/justpayd/internal/pkg"
)

type UserHandler struct {
	UserService UserService
	logger      *zap.SugaredLogger
}

func NewUserHandler(userService UserService, logger *zap.SugaredLogger) *UserHandler {
	return &UserHandler{
		UserService: userService,
		logger:      logger,
	}
}

func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Post("/", h.CreateUser)
}

// @Summary Create a new user
// @Description Creates a new user in the system
// @Tags users
// @Accept json
// @Produce json
// @Param payload body CreateUserRequest true "User information"
// @Success 200 {object} pkg.BaseResponse "User created successfully"
// @Failure 400 {object} pkg.BaseResponse "Invalid request payload"
// @Failure 500 {object} pkg.BaseResponse "Internal server error"
// @Router /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var payload CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		pkg.WriteJSON(w, http.StatusBadRequest, pkg.NewErrorResponse("Invalid request payload"))
		return
	}

	err := h.UserService.CreateUser(&payload)
	if err != nil {
		pkg.WriteJSON(w, http.StatusInternalServerError, pkg.NewErrorResponse("Failed to create user"))
		return
	}

	pkg.WriteJSON(w, http.StatusCreated, pkg.SuccessResponse(payload))
}
