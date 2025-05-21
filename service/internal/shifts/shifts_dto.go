package shifts

import "time"

type CreateShiftRequest struct {
	Date      string `json:"date" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
	Role      string `json:"role" binding:"required"`
	Location  string `json:"location"`
}

type UpdateShiftRequest struct {
	Date      *string `json:"date"`
	StartTime *string `json:"start_time"`
	EndTime   *string `json:"end_time"`
	Role      *string `json:"role"`
	Location  *string `json:"location"`
}

type ShiftResponse struct {
	ID         int       `json:"id"`
	Date       string    `json:"date"`
	StartTime  string    `json:"start_time"`
	EndTime    string    `json:"end_time"`
	Role       string    `json:"role"`
	Assignee   string    `json:"assignee"`
	IsAssigned bool      `json:"is_assigned"`
	Location   string    `json:"location"`
	CreatedAt  time.Time `json:"created_at"`
}
