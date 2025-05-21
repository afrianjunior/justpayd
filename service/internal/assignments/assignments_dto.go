package assignments

import "time"

type AssignmentResponse struct {
	ID         int       `json:"id"`
	ShiftID    int       `json:"shift_id"`
	UserID     int       `json:"user_id"`
	UserName   string    `json:"user_name"`
	Date       string    `json:"date"`
	StartTime  string    `json:"start_time"`
	EndTime    string    `json:"end_time"`
	AssignedAt time.Time `json:"assigned_at"`
}

type UpdateAssignmentRequest struct {
	UserID int `json:"user_id" binding:"required"`
}

type CreateAssignmentRequest struct {
	ShiftID int `json:"shift_id" binding:"required"`
	UserID  int `json:"user_id" binding:"required"`
}
