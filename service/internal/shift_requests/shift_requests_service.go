package shift_requests

import (
	"context"
	"fmt"

	"github.com/afrianjunior/justpayd/internal/assignments"
)

// ShiftRequestService defines the interface for shift request business logic
type ShiftRequestService interface {
	CreateShiftRequest(ctx context.Context, userID int, shiftID int, req *CreateShiftRequestDTO) (*ShiftRequestResponse, error)
	GetShiftRequests(ctx context.Context, filter *ShiftRequestFilter) ([]ShiftRequestResponse, error)
	ApproveShiftRequest(ctx context.Context, id int) (*ShiftRequestResponse, error)
	RejectShiftRequest(ctx context.Context, id int) (*ShiftRequestResponse, error)
}

type shiftRequestService struct {
	shiftRequestRepository ShiftRequestRepository
	assignmentRepository   assignments.AssignmentRepository
}

// NewShiftRequestService creates a new instance of ShiftRequestService
func NewShiftRequestService(
	shiftRequestRepository ShiftRequestRepository,
	assignmentRepository assignments.AssignmentRepository,
) ShiftRequestService {
	return &shiftRequestService{
		shiftRequestRepository: shiftRequestRepository,
		assignmentRepository:   assignmentRepository,
	}
}

func (s *shiftRequestService) CreateShiftRequest(ctx context.Context, userID int, shiftID int, req *CreateShiftRequestDTO) (*ShiftRequestResponse, error) {
	// Add validation or business logic here if needed
	return s.shiftRequestRepository.CreateShiftRequest(ctx, userID, shiftID, req)
}

func (s *shiftRequestService) GetShiftRequests(ctx context.Context, filter *ShiftRequestFilter) ([]ShiftRequestResponse, error) {
	return s.shiftRequestRepository.GetShiftRequests(ctx, filter)
}

func (s *shiftRequestService) ApproveShiftRequest(ctx context.Context, id int) (*ShiftRequestResponse, error) {
	// First, update the shift request status
	updatedRequest, err := s.shiftRequestRepository.UpdateShiftRequestStatus(ctx, id, StatusApproved)
	if err != nil {
		return nil, err
	}

	// If update was successful, create an assignment
	if updatedRequest != nil {
		assignmentReq := &assignments.CreateAssignmentRequest{
			ShiftID: updatedRequest.ShiftID,
			UserID:  updatedRequest.UserID,
		}

		assignment, err := s.assignmentRepository.CreateAssignment(ctx, assignmentReq)
		if err != nil {
			// Log the error but don't fail the request approval
			// In a production app, you might want to handle this differently, maybe with a retry mechanism
			// or by rolling back the status update
			fmt.Printf("Failed to create assignment for shift request %d: %v\n", id, err)
		} else if assignment != nil {
			fmt.Printf("Created assignment ID %d for shift request %d (Shift: %d, User: %d)\n",
				assignment.ID, id, updatedRequest.ShiftID, updatedRequest.UserID)
		} else {
			fmt.Printf("Created assignment for shift request %d, but could not retrieve details\n", id)
		}
	}

	return updatedRequest, nil
}

func (s *shiftRequestService) RejectShiftRequest(ctx context.Context, id int) (*ShiftRequestResponse, error) {
	// You could add additional business logic here before rejecting
	return s.shiftRequestRepository.UpdateShiftRequestStatus(ctx, id, StatusRejected)
}
