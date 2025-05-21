package shifts

import (
	"context"
	"database/sql"
	"errors"
)

// ShiftRepository defines the interface for shift data operations
type ShiftRepository interface {
	CreateShift(ctx context.Context, shift *CreateShiftRequest) (*ShiftResponse, error)
	GetShifts(ctx context.Context) ([]ShiftResponse, error)
	GetShiftByID(ctx context.Context, id int) (*ShiftResponse, error)
	UpdateShift(ctx context.Context, id int, shift *UpdateShiftRequest) (*ShiftResponse, error)
	DeleteShift(ctx context.Context, id int) error
}

type shiftRepository struct {
	db *sql.DB
}

// NewShiftRepository creates a new instance of ShiftRepository
func NewShiftRepository(db *sql.DB) ShiftRepository {
	return &shiftRepository{db: db}
}

func (r *shiftRepository) CreateShift(ctx context.Context, shift *CreateShiftRequest) (*ShiftResponse, error) {
	query := `
		INSERT INTO shifts (date, start_time, end_time, role, location)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		shift.Date,
		shift.StartTime,
		shift.EndTime,
		shift.Role,
		shift.Location,
	)

	if err != nil {
		return nil, err
	}

	// Get the ID of the inserted row
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Get the full shift details
	return r.GetShiftByID(ctx, int(id))
}

func (r *shiftRepository) GetShifts(ctx context.Context) ([]ShiftResponse, error) {
	query := `
		SELECT
			s.id,
			s.date,
			s.start_time,
			s.end_time,
			s.role,
			s.location,
			s.created_at,
			u.name as assignee,
			a.user_id IS NOT NULL as is_assigned
		FROM shifts s
		LEFT JOIN assignments a ON s.id = a.shift_id
		LEFT JOIN users u ON a.user_id = u.id
		ORDER BY s.date DESC, s.start_time ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shifts []ShiftResponse
	for rows.Next() {
		var shift ShiftResponse
		var assigneeNullable sql.NullString // Use NullString to handle NULL values

		if err := rows.Scan(
			&shift.ID,
			&shift.Date,
			&shift.StartTime,
			&shift.EndTime,
			&shift.Role,
			&shift.Location,
			&shift.CreatedAt,
			&assigneeNullable, // Scan into nullable string
			&shift.IsAssigned,
		); err != nil {
			return nil, err
		}

		// Convert nullable string to regular string
		if assigneeNullable.Valid {
			shift.Assignee = assigneeNullable.String
		} else {
			shift.Assignee = "" // Empty string for NULL
		}

		shifts = append(shifts, shift)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return shifts, nil
}

func (r *shiftRepository) GetShiftByID(ctx context.Context, id int) (*ShiftResponse, error) {
	query := `
		SELECT id, date, start_time, end_time, role, location, created_at
		FROM shifts
		WHERE id = ?
	`

	var shift ShiftResponse
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&shift.ID,
		&shift.Date,
		&shift.StartTime,
		&shift.EndTime,
		&shift.Role,
		&shift.Location,
		&shift.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, err
	}

	return &shift, nil
}

func (r *shiftRepository) UpdateShift(ctx context.Context, id int, shift *UpdateShiftRequest) (*ShiftResponse, error) {
	// First, get the current shift to use existing values for fields not being updated
	current, err := r.GetShiftByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, nil // Shift not found
	}

	// Apply updates only to fields that are provided
	date := current.Date
	startTime := current.StartTime
	endTime := current.EndTime
	role := current.Role
	location := current.Location

	if shift.Date != nil {
		date = *shift.Date
	}
	if shift.StartTime != nil {
		startTime = *shift.StartTime
	}
	if shift.EndTime != nil {
		endTime = *shift.EndTime
	}
	if shift.Role != nil {
		role = *shift.Role
	}
	if shift.Location != nil {
		location = *shift.Location
	}

	// Update the shift
	query := `
		UPDATE shifts
		SET date = ?, start_time = ?, end_time = ?, role = ?, location = ?
		WHERE id = ?
	`

	_, err = r.db.ExecContext(
		ctx,
		query,
		date,
		startTime,
		endTime,
		role,
		location,
		id,
	)

	if err != nil {
		return nil, err
	}

	// Get the updated shift
	return r.GetShiftByID(ctx, id)
}

func (r *shiftRepository) DeleteShift(ctx context.Context, id int) error {
	// Check if the shift exists first
	exists, err := r.GetShiftByID(ctx, id)
	if err != nil {
		return err
	}
	if exists == nil {
		return errors.New("shift not found")
	}

	// Delete the shift
	query := "DELETE FROM shifts WHERE id = ?"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("shift not found or could not be deleted")
	}

	return nil
}
