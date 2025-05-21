package shift_requests

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ShiftRequestRepository defines the interface for shift request data operations
type ShiftRequestRepository interface {
	CreateShiftRequest(ctx context.Context, userID int, shiftID int, req *CreateShiftRequestDTO) (*ShiftRequestResponse, error)
	GetShiftRequests(ctx context.Context, filter *ShiftRequestFilter) ([]ShiftRequestResponse, error)
	GetShiftRequestByID(ctx context.Context, id int) (*ShiftRequestResponse, error)
	UpdateShiftRequestStatus(ctx context.Context, id int, status string) (*ShiftRequestResponse, error)
}

type shiftRequestRepository struct {
	db *sql.DB
}

// NewShiftRequestRepository creates a new instance of ShiftRequestRepository
func NewShiftRequestRepository(db *sql.DB) ShiftRequestRepository {
	return &shiftRequestRepository{db: db}
}

func (r *shiftRequestRepository) CreateShiftRequest(ctx context.Context, userID int, shiftID int, req *CreateShiftRequestDTO) (*ShiftRequestResponse, error) {
	query := `
		INSERT INTO shift_requests (shift_id, user_id, status)
		VALUES (?, ?, ?)
	`

	fmt.Println("Creating shift request for shift ID:", shiftID, "and user ID:", userID)

	result, err := r.db.ExecContext(
		ctx,
		query,
		shiftID,
		userID,
		StatusPending, // Default status is pending
	)

	if err != nil {
		return nil, err
	}

	// Get the ID of the inserted row
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Create a response
	var response ShiftRequestResponse
	response.ID = int(id)
	response.ShiftID = shiftID
	response.UserID = userID
	response.Status = StatusPending
	// Note: requested_at is set by the database, we'll get it when we fetch the full request

	// Get user name
	userQuery := `SELECT name FROM users WHERE id = ?`
	err = r.db.QueryRowContext(ctx, userQuery, userID).Scan(&response.UserName)
	if err != nil {
		// If we can't get the user name, just return what we have
		return &response, nil
	}

	// Get shift date, start time, and end time
	shiftQuery := `
		SELECT date, start_time, end_time
		FROM shifts
		WHERE id = ?
	`

	var dateStr, startTimeStr, endTimeStr string
	err = r.db.QueryRowContext(ctx, shiftQuery, shiftID).Scan(
		&dateStr,
		&startTimeStr,
		&endTimeStr,
	)

	if err != nil {
		// If we can't get the shift details, just return what we have
		return &response, nil
	}

	// Parse the date and time strings
	layout := time.RFC3339
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		// Try other common formats
		layouts := []string{"2006-01-02", "2006/01/02", "2006-01-02T15:04:05Z"}
		for _, l := range layouts {
			date, err = time.Parse(l, dateStr)
			if err == nil {
				break
			}
		}
	}
	response.Date = date

	// Parse start time
	startTime, err := time.Parse(layout, startTimeStr)
	if err != nil {
		// Try common time formats
		timeLayouts := []string{"15:04:05", "15:04", "15:04:05Z"}
		for _, l := range timeLayouts {
			startTime, err = time.Parse(l, startTimeStr)
			if err == nil {
				break
			}
		}
	}
	response.StartTime = startTime

	// Parse end time
	endTime, err := time.Parse(layout, endTimeStr)
	if err != nil {
		// Try common time formats
		timeLayouts := []string{"15:04:05", "15:04", "15:04:05Z"}
		for _, l := range timeLayouts {
			endTime, err = time.Parse(l, endTimeStr)
			if err == nil {
				break
			}
		}
	}
	response.EndTime = endTime

	return &response, nil
}

func (r *shiftRequestRepository) GetShiftRequests(ctx context.Context, filter *ShiftRequestFilter) ([]ShiftRequestResponse, error) {
	query := `
		SELECT
			sr.id,
			sr.user_id,
			sr.shift_id,
			u.name as user_name,
			sr.status,
			sr.requested_at,
			s.date,
			s.start_time,
			s.end_time
		FROM shift_requests sr
		JOIN shifts s ON sr.shift_id = s.id
		JOIN users u ON sr.user_id = u.id
	`

	// Add filtering if status or userID is provided
	var args []interface{}
	where := []string{}

	if filter != nil {
		if filter.Status != "" {
			where = append(where, "sr.status = ?")
			args = append(args, filter.Status)
		}

		if filter.UserID > 0 {
			where = append(where, "sr.user_id = ?")
			args = append(args, filter.UserID)
		}

		if filter.ShiftID > 0 {
			where = append(where, "sr.shift_id = ?")
			args = append(args, filter.ShiftID)
		}
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	query += " ORDER BY sr.requested_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []ShiftRequestResponse
	for rows.Next() {
		var request ShiftRequestResponse
		var dateStr, startTimeStr, endTimeStr string

		if err := rows.Scan(
			&request.ID,
			&request.UserID,
			&request.ShiftID,
			&request.UserName,
			&request.Status,
			&request.RequestedAt,
			&dateStr,
			&startTimeStr,
			&endTimeStr,
		); err != nil {
			return nil, err
		}

		// Parse string dates/times into time.Time objects
		// Try standard ISO 8601 format first (which SQLite may use)
		layout := time.RFC3339
		date, err := time.Parse(layout, dateStr)
		if err != nil {
			// Try other common formats
			layouts := []string{"2006-01-02", "2006/01/02", "2006-01-02T15:04:05Z"}
			for _, l := range layouts {
				date, err = time.Parse(l, dateStr)
				if err == nil {
					break
				}
			}
			if err != nil {
				fmt.Printf("Failed to parse date string %q: %v\n", dateStr, err)
			}
		}
		request.Date = date

		// Parse time strings - try with complete ISO format first
		timeLayout := time.RFC3339
		startTime, err := time.Parse(timeLayout, startTimeStr)
		if err != nil {
			// Try common time formats
			timeLayouts := []string{"15:04:05", "15:04", "15:04:05Z"}
			for _, l := range timeLayouts {
				startTime, err = time.Parse(l, startTimeStr)
				if err == nil {
					break
				}
			}
			if err != nil {
				fmt.Printf("Failed to parse start time string %q: %v\n", startTimeStr, err)
			}
		}
		request.StartTime = startTime

		endTime, err := time.Parse(timeLayout, endTimeStr)
		if err != nil {
			// Try common time formats
			timeLayouts := []string{"15:04:05", "15:04", "15:04:05Z"}
			for _, l := range timeLayouts {
				endTime, err = time.Parse(l, endTimeStr)
				if err == nil {
					break
				}
			}
			if err != nil {
				fmt.Printf("Failed to parse end time string %q: %v\n", endTimeStr, err)
			}
		}
		request.EndTime = endTime

		requests = append(requests, request)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *shiftRequestRepository) GetShiftRequestByID(ctx context.Context, id int) (*ShiftRequestResponse, error) {
	query := `
		SELECT
			sr.id,
			sr.user_id,
			sr.shift_id,
			u.name as user_name,
			sr.status,
			sr.requested_at,
			s.date,
			s.start_time,
			s.end_time
		FROM shift_requests sr
		JOIN shifts s ON sr.shift_id = s.id
		JOIN users u ON sr.user_id = u.id
		WHERE sr.id = ?
	`

	var request ShiftRequestResponse
	var dateStr, startTimeStr, endTimeStr string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&request.ID,
		&request.UserID,
		&request.ShiftID,
		&request.UserName,
		&request.Status,
		&request.RequestedAt,
		&dateStr,
		&startTimeStr,
		&endTimeStr,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, err
	}

	// Parse string dates/times into time.Time objects
	// Try standard ISO 8601 format first (which SQLite may use)
	layout := time.RFC3339
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		// Try other common formats
		layouts := []string{"2006-01-02", "2006/01/02", "2006-01-02T15:04:05Z"}
		for _, l := range layouts {
			date, err = time.Parse(l, dateStr)
			if err == nil {
				break
			}
		}
		if err != nil {
			fmt.Printf("Failed to parse date string %q: %v\n", dateStr, err)
		}
	}
	request.Date = date

	// Parse time strings - try with complete ISO format first
	timeLayout := time.RFC3339
	startTime, err := time.Parse(timeLayout, startTimeStr)
	if err != nil {
		// Try common time formats
		timeLayouts := []string{"15:04:05", "15:04", "15:04:05Z"}
		for _, l := range timeLayouts {
			startTime, err = time.Parse(l, startTimeStr)
			if err == nil {
				break
			}
		}
		if err != nil {
			fmt.Printf("Failed to parse start time string %q: %v\n", startTimeStr, err)
		}
	}
	request.StartTime = startTime

	endTime, err := time.Parse(timeLayout, endTimeStr)
	if err != nil {
		// Try common time formats
		timeLayouts := []string{"15:04:05", "15:04", "15:04:05Z"}
		for _, l := range timeLayouts {
			endTime, err = time.Parse(l, endTimeStr)
			if err == nil {
				break
			}
		}
		if err != nil {
			fmt.Printf("Failed to parse end time string %q: %v\n", endTimeStr, err)
		}
	}
	request.EndTime = endTime

	return &request, nil
}

func (r *shiftRequestRepository) UpdateShiftRequestStatus(ctx context.Context, id int, status string) (*ShiftRequestResponse, error) {
	// First check if the request exists
	request, err := r.GetShiftRequestByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if request == nil {
		return nil, nil // Not found
	}

	// Update the status - remove the RETURNING clause for SQLite compatibility
	query := `
		UPDATE shift_requests
		SET status = ?
		WHERE id = ?
	`

	_, err = r.db.ExecContext(
		ctx,
		query,
		status,
		id,
	)

	if err != nil {
		return nil, err
	}

	// Get the updated request
	return r.GetShiftRequestByID(ctx, id)
}
