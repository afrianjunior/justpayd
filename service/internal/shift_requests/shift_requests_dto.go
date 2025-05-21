package shift_requests

import "time"

// Possible status values for shift requests
const (
	StatusPending  = "pending"
	StatusApproved = "approved"
	StatusRejected = "rejected"
)

type CreateShiftRequestDTO struct {
	ShiftID int `json:"shift_id" binding:"required"`
}

type ShiftRequestResponse struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	UserName    string    `json:"user_name"`
	ShiftID     int       `json:"shift_id"`
	Status      string    `json:"status"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Date        time.Time `json:"date"`
	RequestedAt time.Time `json:"requested_at"`
}

type ShiftRequestFilter struct {
	Status  string `json:"status"`
	UserID  int    `json:"user_id"`
	ShiftID int    `json:"shift_id"`
}
