package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	IsSuccess bool   `json:"is_success"`
	Stack     string `json:"stack"`
	Message   string `json:"message"`
}

type BaseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func JsonResponse(w http.ResponseWriter, d any, c int) {
	dj, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)

	_, _ = fmt.Fprintf(w, "%s", dj)
}

// jsonResponseUsingBase write json body response using struct BaseResponse
func JsonResponseUsingBase(w http.ResponseWriter, msg string, payload any, err error, c int) {
	resp := BaseResponse{
		Message: msg,
		Success: err == nil,
		Data:    payload,
	}

	dj, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)

	_, _ = fmt.Fprintf(w, "%s", dj)
}

// WriteJSON writes a JSON response with the given status code
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// SuccessResponse creates a standard success response
func SuccessResponse(data interface{}) BaseResponse {
	return BaseResponse{
		Success: true,
		Message: "Success",
		Data:    data,
	}
}

// NewErrorResponse creates a standard error response with a message
func NewErrorResponse(message string) BaseResponse {
	return BaseResponse{
		Success: false,
		Message: message,
		Data:    nil,
	}
}
