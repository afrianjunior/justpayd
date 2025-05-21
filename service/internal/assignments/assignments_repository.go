package assignments

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// AssignmentRepository defines the interface for assignment data operations
type AssignmentRepository interface {
	GetAssignments(ctx context.Context) ([]AssignmentResponse, error)
	GetAssignmentByID(ctx context.Context, id int) (*AssignmentResponse, error)
	UpdateAssignment(ctx context.Context, id int, req *UpdateAssignmentRequest) (*AssignmentResponse, error)
	CreateAssignment(ctx context.Context, req *CreateAssignmentRequest) (*AssignmentResponse, error)
}

type assignmentRepository struct {
	db *sql.DB
}

// NewAssignmentRepository creates a new instance of AssignmentRepository
func NewAssignmentRepository(db *sql.DB) AssignmentRepository {
	return &assignmentRepository{db: db}
}

func (r *assignmentRepository) GetAssignments(ctx context.Context) ([]AssignmentResponse, error) {
	query := `
		SELECT 
			a.id, 
			a.shift_id, 
			a.user_id, 
			u.name as user_name,
			s.date,
			s.start_time,
			s.end_time,
			a.assigned_at
		FROM assignments a
		JOIN users u ON a.user_id = u.id
		JOIN shifts s ON a.shift_id = s.id
		ORDER BY s.date, s.start_time
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []AssignmentResponse
	for rows.Next() {
		var assignment AssignmentResponse
		var dateStr, startTimeStr, endTimeStr string
		var assignedAtStr string

		if err := rows.Scan(
			&assignment.ID,
			&assignment.ShiftID,
			&assignment.UserID,
			&assignment.UserName,
			&dateStr,
			&startTimeStr,
			&endTimeStr,
			&assignedAtStr,
		); err != nil {
			return nil, err
		}

		// Parse date string
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
				// Use the string value directly if parsing fails
				assignment.Date = dateStr
			} else {
				assignment.Date = date.Format("2006-01-02")
			}
		} else {
			assignment.Date = date.Format("2006-01-02")
		}

		// Parse start time string
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
				// Use the string value directly if parsing fails
				assignment.StartTime = startTimeStr
			} else {
				assignment.StartTime = startTime.Format("15:04:05")
			}
		} else {
			assignment.StartTime = startTime.Format("15:04:05")
		}

		// Parse end time string
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
				// Use the string value directly if parsing fails
				assignment.EndTime = endTimeStr
			} else {
				assignment.EndTime = endTime.Format("15:04:05")
			}
		} else {
			assignment.EndTime = endTime.Format("15:04:05")
		}

		// Parse assigned_at timestamp
		assignedAt, err := time.Parse(layout, assignedAtStr)
		if err != nil {
			// Try other common formats
			layouts := []string{"2006-01-02 15:04:05", "2006-01-02T15:04:05Z", "2006-01-02 15:04:05.999999999-07:00"}
			for _, l := range layouts {
				assignedAt, err = time.Parse(l, assignedAtStr)
				if err == nil {
					break
				}
			}
			if err != nil {
				fmt.Printf("Failed to parse assigned_at string %q: %v\n", assignedAtStr, err)
			}
		}
		assignment.AssignedAt = assignedAt

		assignments = append(assignments, assignment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assignments, nil
}

func (r *assignmentRepository) GetAssignmentByID(ctx context.Context, id int) (*AssignmentResponse, error) {
	query := `
		SELECT 
			a.id, 
			a.shift_id, 
			a.user_id,
			u.name as user_name,
			s.date,
			s.start_time,
			s.end_time,
			a.assigned_at
		FROM assignments a
		JOIN users u ON a.user_id = u.id
		JOIN shifts s ON a.shift_id = s.id
		WHERE a.id = ?
	`

	var assignment AssignmentResponse
	var dateStr, startTimeStr, endTimeStr string
	var assignedAtStr string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&assignment.ID,
		&assignment.ShiftID,
		&assignment.UserID,
		&assignment.UserName,
		&dateStr,
		&startTimeStr,
		&endTimeStr,
		&assignedAtStr,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, err
	}

	// Parse date string
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
			// Use the string value directly if parsing fails
			assignment.Date = dateStr
		} else {
			assignment.Date = date.Format("2006-01-02")
		}
	} else {
		assignment.Date = date.Format("2006-01-02")
	}

	// Parse start time string
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
			// Use the string value directly if parsing fails
			assignment.StartTime = startTimeStr
		} else {
			assignment.StartTime = startTime.Format("15:04:05")
		}
	} else {
		assignment.StartTime = startTime.Format("15:04:05")
	}

	// Parse end time string
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
			// Use the string value directly if parsing fails
			assignment.EndTime = endTimeStr
		} else {
			assignment.EndTime = endTime.Format("15:04:05")
		}
	} else {
		assignment.EndTime = endTime.Format("15:04:05")
	}

	// Parse assigned_at timestamp
	assignedAt, err := time.Parse(layout, assignedAtStr)
	if err != nil {
		// Try other common formats
		layouts := []string{"2006-01-02 15:04:05", "2006-01-02T15:04:05Z", "2006-01-02 15:04:05.999999999-07:00"}
		for _, l := range layouts {
			assignedAt, err = time.Parse(l, assignedAtStr)
			if err == nil {
				break
			}
		}
		if err != nil {
			fmt.Printf("Failed to parse assigned_at string %q: %v\n", assignedAtStr, err)
		}
	}
	assignment.AssignedAt = assignedAt

	return &assignment, nil
}

func (r *assignmentRepository) UpdateAssignment(ctx context.Context, id int, req *UpdateAssignmentRequest) (*AssignmentResponse, error) {
	// First, check if the assignment exists
	assignment, err := r.GetAssignmentByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if assignment == nil {
		return nil, nil // Not found
	}

	// Update the assignment with new user ID
	query := `
		UPDATE assignments
		SET user_id = ?, assigned_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err = r.db.ExecContext(ctx, query, req.UserID, id)
	if err != nil {
		return nil, err
	}

	// Get full assignment details
	return r.GetAssignmentByID(ctx, id)
}

func (r *assignmentRepository) CreateAssignment(ctx context.Context, req *CreateAssignmentRequest) (*AssignmentResponse, error) {
	query := `
		INSERT INTO assignments (shift_id, user_id, assigned_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
	`

	result, err := r.db.ExecContext(ctx, query, req.ShiftID, req.UserID)
	if err != nil {
		return nil, err
	}

	// Get the ID of the inserted row
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Get the full assignment details
	return r.GetAssignmentByID(ctx, int(id))
}
