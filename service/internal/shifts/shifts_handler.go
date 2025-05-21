package shifts

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/afrianjunior/justpayd/internal/pkg"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ShiftHandler struct {
	ShiftService ShiftService
	logger       *zap.SugaredLogger
}

func NewShiftHandler(shiftService ShiftService, logger *zap.SugaredLogger) *ShiftHandler {
	return &ShiftHandler{
		ShiftService: shiftService,
		logger:       logger,
	}
}

func (h *ShiftHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.GetShifts)
	r.Post("/", h.CreateShift) // Assuming admin check is handled by middleware or in service
	r.Get("/{id}", h.GetShiftByID)
	r.Put("/{id}", h.UpdateShift)
	r.Delete("/{id}", h.DeleteShift)
}

// GetShifts godoc
// @Summary List semua shift
// @Description List semua shift
// @Tags shifts
// @Produce json
// @Success 200 {object} pkg.BaseResponse{data=[]ShiftResponse} "Successfully retrieved shifts"
// @Failure 500 {object} pkg.BaseResponse "Internal server error"
// @Router /shifts [get]
func (h *ShiftHandler) GetShifts(w http.ResponseWriter, r *http.Request) {
	shifts, err := h.ShiftService.GetShifts(r.Context())
	if err != nil {
		h.logger.Errorf("Error getting shifts: %v", err)
		pkg.WriteJSON(w, http.StatusInternalServerError, pkg.NewErrorResponse("Failed to retrieve shifts"))
		return
	}
	pkg.WriteJSON(w, http.StatusOK, pkg.SuccessResponse(shifts))
}

// CreateShift godoc
// @Summary Admin creates shift
// @Description Admin creates shift
// @Tags shifts
// @Accept json
// @Produce json
// @Param payload body CreateShiftRequest true "Shift creation payload"
// @Success 201 {object} pkg.BaseResponse{data=ShiftResponse} "Shift created successfully"
// @Failure 400 {object} pkg.BaseResponse "Invalid request payload"
// @Failure 500 {object} pkg.BaseResponse "Internal server error"
// @Router /shifts [post]
func (h *ShiftHandler) CreateShift(w http.ResponseWriter, r *http.Request) {
	var payload CreateShiftRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		pkg.WriteJSON(w, http.StatusBadRequest, pkg.NewErrorResponse("Invalid request payload: "+err.Error()))
		return
	}

	shift, err := h.ShiftService.CreateShift(r.Context(), &payload)
	if err != nil {
		h.logger.Errorf("Error creating shift: %v", err)
		pkg.WriteJSON(w, http.StatusInternalServerError, pkg.NewErrorResponse("Failed to create shift"))
		return
	}
	pkg.WriteJSON(w, http.StatusCreated, pkg.SuccessResponse(shift))
}

// GetShiftByID godoc
// @Summary Detail shift
// @Description Mendapatkan detail shift berdasarkan ID
// @Tags shifts
// @Produce json
// @Param id path int true "Shift ID"
// @Success 200 {object} pkg.BaseResponse{data=ShiftResponse} "Successfully retrieved shift detail"
// @Failure 400 {object} pkg.BaseResponse "Invalid shift ID"
// @Failure 404 {object} pkg.BaseResponse "Shift not found"
// @Failure 500 {object} pkg.BaseResponse "Internal server error"
// @Router /shifts/{id} [get]
func (h *ShiftHandler) GetShiftByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.WriteJSON(w, http.StatusBadRequest, pkg.NewErrorResponse("Invalid shift ID"))
		return
	}

	shift, err := h.ShiftService.GetShiftByID(r.Context(), id)
	if err != nil {
		// TODO: Differentiate between not found and other errors if service layer supports it
		h.logger.Errorf("Error getting shift by ID %d: %v", id, err)
		pkg.WriteJSON(w, http.StatusInternalServerError, pkg.NewErrorResponse("Failed to retrieve shift"))
		return
	}
	if shift == nil { // Assuming service returns nil, nil if not found
		pkg.WriteJSON(w, http.StatusNotFound, pkg.NewErrorResponse("Shift not found"))
		return
	}
	pkg.WriteJSON(w, http.StatusOK, pkg.SuccessResponse(shift))
}

// UpdateShift godoc
// @Summary Edit shift
// @Description Mengedit shift berdasarkan ID
// @Tags shifts
// @Accept json
// @Produce json
// @Param id path int true "Shift ID"
// @Param payload body UpdateShiftRequest true "Shift update payload"
// @Success 200 {object} pkg.BaseResponse{data=ShiftResponse} "Shift updated successfully"
// @Failure 400 {object} pkg.BaseResponse "Invalid request payload or shift ID"
// @Failure 404 {object} pkg.BaseResponse "Shift not found"
// @Failure 500 {object} pkg.BaseResponse "Internal server error"
// @Router /shifts/{id} [put]
func (h *ShiftHandler) UpdateShift(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.WriteJSON(w, http.StatusBadRequest, pkg.NewErrorResponse("Invalid shift ID"))
		return
	}

	var payload UpdateShiftRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		pkg.WriteJSON(w, http.StatusBadRequest, pkg.NewErrorResponse("Invalid request payload: "+err.Error()))
		return
	}

	shift, err := h.ShiftService.UpdateShift(r.Context(), id, &payload)
	if err != nil {
		// TODO: Differentiate between not found and other errors
		h.logger.Errorf("Error updating shift ID %d: %v", id, err)
		pkg.WriteJSON(w, http.StatusInternalServerError, pkg.NewErrorResponse("Failed to update shift"))
		return
	}
	if shift == nil { // Assuming service returns nil, nil if not found and update is not partial
		pkg.WriteJSON(w, http.StatusNotFound, pkg.NewErrorResponse("Shift not found or update failed"))
		return
	}
	pkg.WriteJSON(w, http.StatusOK, pkg.SuccessResponse(shift))
}

// DeleteShift godoc
// @Summary Hapus shift
// @Description Menghapus shift berdasarkan ID
// @Tags shifts
// @Produce json
// @Param id path int true "Shift ID"
// @Success 200 {object} pkg.BaseResponse "Shift deleted successfully"
// @Failure 400 {object} pkg.BaseResponse "Invalid shift ID"
// @Failure 404 {object} pkg.BaseResponse "Shift not found"
// @Failure 500 {object} pkg.BaseResponse "Internal server error"
// @Router /shifts/{id} [delete]
func (h *ShiftHandler) DeleteShift(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.WriteJSON(w, http.StatusBadRequest, pkg.NewErrorResponse("Invalid shift ID"))
		return
	}

	err = h.ShiftService.DeleteShift(r.Context(), id)
	if err != nil {
		// TODO: Differentiate between not found and other errors
		h.logger.Errorf("Error deleting shift ID %d: %v", id, err)
		pkg.WriteJSON(w, http.StatusInternalServerError, pkg.NewErrorResponse("Failed to delete shift"))
		return
	}
	pkg.WriteJSON(w, http.StatusOK, pkg.SuccessResponse(map[string]string{"message": "Shift deleted successfully"}))
}
