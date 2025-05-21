package shifts

import "context"

// ShiftService defines the interface for shift business logic
type ShiftService interface {
	CreateShift(ctx context.Context, req *CreateShiftRequest) (*ShiftResponse, error)
	GetShifts(ctx context.Context) ([]ShiftResponse, error)
	GetShiftByID(ctx context.Context, id int) (*ShiftResponse, error)
	UpdateShift(ctx context.Context, id int, req *UpdateShiftRequest) (*ShiftResponse, error)
	DeleteShift(ctx context.Context, id int) error
}

type shiftService struct {
	shiftRepository ShiftRepository
}

// NewShiftService creates a new instance of ShiftService
func NewShiftService(shiftRepository ShiftRepository) ShiftService {
	return &shiftService{shiftRepository: shiftRepository}
}

func (s *shiftService) CreateShift(ctx context.Context, req *CreateShiftRequest) (*ShiftResponse, error) {
	// TODO: Add validation or other business logic here if needed
	return s.shiftRepository.CreateShift(ctx, req)
}

func (s *shiftService) GetShifts(ctx context.Context) ([]ShiftResponse, error) {
	return s.shiftRepository.GetShifts(ctx)
}

func (s *shiftService) GetShiftByID(ctx context.Context, id int) (*ShiftResponse, error) {
	return s.shiftRepository.GetShiftByID(ctx, id)
}

func (s *shiftService) UpdateShift(ctx context.Context, id int, req *UpdateShiftRequest) (*ShiftResponse, error) {
	// TODO: Add validation or other business logic here if needed
	return s.shiftRepository.UpdateShift(ctx, id, req)
}

func (s *shiftService) DeleteShift(ctx context.Context, id int) error {
	return s.shiftRepository.DeleteShift(ctx, id)
}
