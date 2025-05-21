package assignments

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/afrianjunior/justpayd/internal/pkg"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type AssignmentHandler struct {
	AssignmentService AssignmentService
	logger            *zap.SugaredLogger
}

func NewAssignmentHandler(assignmentService AssignmentService, logger *zap.SugaredLogger) *AssignmentHandler {
	return &AssignmentHandler{
		AssignmentService: assignmentService,
		logger:            logger,
	}
}

func (h *AssignmentHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.GetAssignments)
	r.Put("/{id}", h.UpdateAssignment)
	r.Post("/", h.CreateAssignment)
}

// GetAssignments godoc
// @Summary List all assignments
// @Description Get all shift assignments
// @Tags assignments
// @Produce json
// @Success 200 {object} pkg.BaseResponse{data=[]AssignmentResponse} "Successfully retrieved assignments"
// @Failure 500 {object} pkg.BaseResponse "Internal server error"
// @Router /assignments [get]
func (h *AssignmentHandler) GetAssignments(w http.ResponseWriter, r *http.Request) {
	assignments, err := h.AssignmentService.GetAssignments(r.Context())
	if err != nil {
		h.logger.Errorf("Error getting assignments: %v", err)
		pkg.WriteJSON(w, http.StatusInternalServerError, pkg.NewErrorResponse("Failed to retrieve assignments"))
		return
	}
	pkg.WriteJSON(w, http.StatusOK, pkg.SuccessResponse(assignments))
}

// UpdateAssignment godoc
// @Summary Update assignment
// @Description Change the user assigned to a shift
// @Tags assignments
// @Accept json
// @Produce json
// @Param id path int true "Assignment ID"
// @Param payload body UpdateAssignmentRequest true "Assignment update payload"
// @Success 200 {object} pkg.BaseResponse{data=AssignmentResponse} "Assignment updated successfully"
// @Failure 400 {object} pkg.BaseResponse "Invalid request payload or assignment ID"
// @Failure 404 {object} pkg.BaseResponse "Assignment not found"
// @Failure 500 {object} pkg.BaseResponse "Internal server error"
// @Router /assignments/{id} [put]
func (h *AssignmentHandler) UpdateAssignment(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.WriteJSON(w, http.StatusBadRequest, pkg.NewErrorResponse("Invalid assignment ID"))
		return
	}

	var payload UpdateAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		pkg.WriteJSON(w, http.StatusBadRequest, pkg.NewErrorResponse("Invalid request payload: "+err.Error()))
		return
	}

	assignment, err := h.AssignmentService.UpdateAssignment(r.Context(), id, &payload)
	if err != nil {
		h.logger.Errorf("Error updating assignment ID %d: %v", id, err)
		pkg.WriteJSON(w, http.StatusInternalServerError, pkg.NewErrorResponse("Failed to update assignment"))
		return
	}
	if assignment == nil {
		pkg.WriteJSON(w, http.StatusNotFound, pkg.NewErrorResponse("Assignment not found"))
		return
	}
	pkg.WriteJSON(w, http.StatusOK, pkg.SuccessResponse(assignment))
}

// CreateAssignment godoc
// @Summary Create a new assignment
// @Description Assign a user to a shift
// @Tags assignments
// @Accept json
// @Produce json
// @Param payload body CreateAssignmentRequest true "Assignment creation payload"
// @Success 201 {object} pkg.BaseResponse{data=AssignmentResponse} "Assignment created successfully"
// @Failure 400 {object} pkg.BaseResponse "Invalid request payload"
// @Failure 500 {object} pkg.BaseResponse "Internal server error"
// @Router /assignments [post]
func (h *AssignmentHandler) CreateAssignment(w http.ResponseWriter, r *http.Request) {
	var payload CreateAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		pkg.WriteJSON(w, http.StatusBadRequest, pkg.NewErrorResponse("Invalid request payload: "+err.Error()))
		return
	}

	assignment, err := h.AssignmentService.CreateAssignment(r.Context(), &payload)
	if err != nil {
		h.logger.Errorf("Error creating assignment: %v", err)
		pkg.WriteJSON(w, http.StatusInternalServerError, pkg.NewErrorResponse("Failed to create assignment"))
		return
	}

	pkg.WriteJSON(w, http.StatusCreated, pkg.SuccessResponse(assignment))
}
