package pkg

import "time"

// User represents a user in the system
type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Role      string    `json:"role" db:"role"` // worker, admin
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Shift represents a work shift in the system
type Shift struct {
	ID        int       `json:"id" db:"id"`
	Date      time.Time `json:"date" db:"date"`
	StartTime time.Time `json:"start_time" db:"start_time"`
	EndTime   time.Time `json:"end_time" db:"end_time"`
	Role      string    `json:"role" db:"role"`
	Location  string    `json:"location" db:"location"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ShiftRequest represents a request by a user to work a shift
type ShiftRequest struct {
	ID          int        `json:"id" db:"id"`
	UserID      int        `json:"user_id" db:"user_id"`
	ShiftID     int        `json:"shift_id" db:"shift_id"`
	Status      string     `json:"status" db:"status"` // pending, approved, rejected
	RequestedAt time.Time  `json:"requested_at" db:"requested_at"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty" db:"approved_at"`
}

// Assignment represents an assigned shift to a user
type Assignment struct {
	ID         int       `json:"id" db:"id"`
	ShiftID    int       `json:"shift_id" db:"shift_id"`
	UserID     int       `json:"user_id" db:"user_id"`
	AssignedAt time.Time `json:"assigned_at" db:"assigned_at"`
}
