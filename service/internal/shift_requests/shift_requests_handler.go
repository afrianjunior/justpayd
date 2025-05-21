package shift_requests

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/afrianjunior/justpayd/internal/pkg"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ShiftRequestHandler struct {
	ShiftRequestService ShiftRequestService
	logger              *zap.SugaredLogger
}

func NewShiftRequestHandler(shiftRequestService ShiftRequestService, logger *zap.SugaredLogger) *ShiftRequestHandler {
	return &ShiftRequestHandler{
		ShiftRequestService: shiftRequestService,
		logger:              logger,
	}
}

func (h *ShiftRequestHandler) RegisterRoutes(r chi.Router) {
	r.Post("/", h.CreateShiftRequest)
	r.Get("/", h.GetShiftRequests)
	r.Put("/approve/{id}", h.ApproveShiftRequest)
	r.Put("/reject/{id}", h.RejectShiftRequest)
}

// CreateShiftRequest godoc
// @Summary User creates shift request
// @Description User with role "user" creates shift request
// @Tags shift-requests
// @Accept json
// @Produce json
// @Param payload body CreateShiftRequestDTO true "Shift request creation payload"
// @Success 201 {object} pkg.BaseResponse{data=ShiftRequestResponse} "Shift request created successfully"
// @Failure 400 {object} pkg.BaseResponse "Invalid request payload"
// @Failure 401 {object} pkg.BaseResponse "Unauthorized"
// @Failure 500 {object} pkg.BaseResponse "Internal server error"
// @Router /shift_requests [post]
func (h *ShiftRequestHandler) CreateShiftRequest(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (assuming it's set by auth middleware)
	user, ok := pkg.GetUserFromContext(r.Context())
	if !ok {
		pkg.WriteJSON(w, http.StatusUnauthorized, pkg.NewErrorResponse("User not authenticated"))
		return
	}
	userID := user.ID

	// Get user role from context (assuming it's set by auth middleware)
	userRole := user.Role
	if userRole != "worker" {
		pkg.WriteJSON(w, http.StatusForbidden, pkg.NewErrorResponse("Only wokrers can create shift requests"))
		return
	}

	var payload CreateShiftRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		pkg.WriteJSON(w, http.StatusBadRequest, pkg.NewErrorResponse("Invalid request payload: "+err.Error()))
		return
	}

	request, err := h.ShiftRequestService.CreateShiftRequest(r.Context(), userID, payload.ShiftID, &payload)
	if err != nil {
		h.logger.Errorf("Error creating shift request: %v", err)
		pkg.WriteJSON(w, http.StatusInternalServerError, pkg.NewErrorResponse(err.Error()))
		return
	}
	pkg.WriteJSON(w, http.StatusCreated, pkg.SuccessResponse(request))
}

// GetShiftRequests godoc
// @Summary List all shift requests
// @Description Admin gets all shift requests, can filter by status, user_id, and shift_id
// @Tags shift-requests
// @Produce json
// @Param status query string false "Filter by status (pending, approved, rejected)"
// @Param user_id query integer false "Filter by user ID"
// @Param shift_id query integer false "Filter by shift ID"
// @Success 200 {object} pkg.BaseResponse{data=[]ShiftRequestResponse} "Successfully retrieved shift requests"
// @Failure 401 {object} pkg.BaseResponse "Unauthorized"
// @Failure 403 {object} pkg.BaseResponse "Forbidden - Not an admin"
// @Failure 500 {object} pkg.BaseResponse "Internal server error"
// @Router /shift_requests [get]
func (h *ShiftRequestHandler) GetShiftRequests(w http.ResponseWriter, r *http.Request) {
	// Get user role from context (assuming it's set by auth middleware)
	user, ok := pkg.GetUserFromContext(r.Context())
	if !ok {
		pkg.WriteJSON(w, http.StatusUnauthorized, pkg.NewErrorResponse("User not authenticated"))
		return
	}
	userRole := user.Role
	if userRole != "admin" {
		pkg.WriteJSON(w, http.StatusForbidden, pkg.NewErrorResponse("Only admins can view all shift requests"))
		return
	}

	// Get filters from query params
	filter := &ShiftRequestFilter{}

	// Get status filter if provided
	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = status
	}

	// Get user_id filter if provided
	if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
		userID, err := strconv.Atoi(userIDStr)
		if err == nil && userID > 0 {
			filter.UserID = userID
		} else if err != nil {
			h.logger.Warnf("Invalid user_id parameter: %s", userIDStr)
		}
	}

	// Get shift_id filter if provided
	if shiftIDStr := r.URL.Query().Get("shift_id"); shiftIDStr != "" {
		shiftID, err := strconv.Atoi(shiftIDStr)
		if err == nil && shiftID > 0 {
			filter.ShiftID = shiftID
		} else if err != nil {
			h.logger.Warnf("Invalid shift_id parameter: %s", shiftIDStr)
		}
	}

	requests, err := h.ShiftRequestService.GetShiftRequests(r.Context(), filter)
	if err != nil {
		h.logger.Errorf("Error getting shift requests: %v", err)
		pkg.WriteJSON(w, http.StatusInternalServerError, pkg.NewErrorResponse("Failed to retrieve shift requests"))
		return
	}
	pkg.WriteJSON(w, http.StatusOK, pkg.SuccessResponse(requests))
}

// ApproveShiftRequest godoc
// @Summary Admin approves shift request
// @Description Admin approves shift request by ID
// @Tags shift-requests
// @Produce json
// @Param id path int true "Shift Request ID"
// @Success 200 {object} pkg.BaseResponse{data=ShiftRequestResponse} "Shift request approved successfully"
// @Failure 400 {object} pkg.BaseResponse "Invalid request ID"
// @Failure 401 {object} pkg.BaseResponse "Unauthorized"
// @Failure 403 {object} pkg.BaseResponse "Forbidden - Not an admin"
// @Failure 404 {object} pkg.BaseResponse "Shift request not found"
// @Failure 500 {object} pkg.BaseResponse "Internal server error"
// @Router /shift_requests/approve/{id} [put]
func (h *ShiftRequestHandler) ApproveShiftRequest(w http.ResponseWriter, r *http.Request) {
	// Get user role from context (assuming it's set by auth middleware)
	user, ok := pkg.GetUserFromContext(r.Context())
	if !ok {
		pkg.WriteJSON(w, http.StatusUnauthorized, pkg.NewErrorResponse("User not authenticated"))
		return
	}
	userRole := user.Role
	if !ok || userRole != "admin" {
		pkg.WriteJSON(w, http.StatusForbidden, pkg.NewErrorResponse("Only admins can approve shift requests"))
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.WriteJSON(w, http.StatusBadRequest, pkg.NewErrorResponse("Invalid shift request ID"))
		return
	}

	request, err := h.ShiftRequestService.ApproveShiftRequest(r.Context(), id)
	if err != nil {
		h.logger.Errorf("Error approving shift request %d: %v", id, err)
		pkg.WriteJSON(w, http.StatusInternalServerError, pkg.NewErrorResponse("Failed to approve shift request"))
		return
	}
	if request == nil {
		pkg.WriteJSON(w, http.StatusNotFound, pkg.NewErrorResponse("Shift request not found"))
		return
	}
	pkg.WriteJSON(w, http.StatusOK, pkg.SuccessResponse(request))
}

// RejectShiftRequest godoc
// @Summary Admin rejects shift request
// @Description Admin rejects shift request by ID
// @Tags shift-requests
// @Produce json
// @Param id path int true "Shift Request ID"
// @Success 200 {object} pkg.BaseResponse{data=ShiftRequestResponse} "Shift request rejected successfully"
// @Failure 400 {object} pkg.BaseResponse "Invalid request ID"
// @Failure 401 {object} pkg.BaseResponse "Unauthorized"
// @Failure 403 {object} pkg.BaseResponse "Forbidden - Not an admin"
// @Failure 404 {object} pkg.BaseResponse "Shift request not found"
// @Failure 500 {object} pkg.BaseResponse "Internal server error"
// @Router /shift_requests/reject/{id} [put]
func (h *ShiftRequestHandler) RejectShiftRequest(w http.ResponseWriter, r *http.Request) {
	// Get user role from context (assuming it's set by auth middleware)
	user, ok := pkg.GetUserFromContext(r.Context())
	if !ok {
		pkg.WriteJSON(w, http.StatusUnauthorized, pkg.NewErrorResponse("User not authenticated"))
		return
	}
	userRole := user.Role
	if !ok || userRole != "admin" {
		pkg.WriteJSON(w, http.StatusForbidden, pkg.NewErrorResponse("Only admins can reject shift requests"))
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.WriteJSON(w, http.StatusBadRequest, pkg.NewErrorResponse("Invalid shift request ID"))
		return
	}

	request, err := h.ShiftRequestService.RejectShiftRequest(r.Context(), id)
	if err != nil {
		h.logger.Errorf("Error rejecting shift request %d: %v", id, err)
		pkg.WriteJSON(w, http.StatusInternalServerError, pkg.NewErrorResponse("Failed to reject shift request"))
		return
	}
	if request == nil {
		pkg.WriteJSON(w, http.StatusNotFound, pkg.NewErrorResponse("Shift request not found"))
		return
	}
	pkg.WriteJSON(w, http.StatusOK, pkg.SuccessResponse(request))
}
