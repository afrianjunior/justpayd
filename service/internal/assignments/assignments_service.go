package assignments

import "context"

// AssignmentService defines the interface for assignment business logic
type AssignmentService interface {
	GetAssignments(ctx context.Context) ([]AssignmentResponse, error)
	UpdateAssignment(ctx context.Context, id int, req *UpdateAssignmentRequest) (*AssignmentResponse, error)
	CreateAssignment(ctx context.Context, req *CreateAssignmentRequest) (*AssignmentResponse, error)
}

type assignmentService struct {
	assignmentRepository AssignmentRepository
}

// NewAssignmentService creates a new instance of AssignmentService
func NewAssignmentService(assignmentRepository AssignmentRepository) AssignmentService {
	return &assignmentService{assignmentRepository: assignmentRepository}
}

func (s *assignmentService) GetAssignments(ctx context.Context) ([]AssignmentResponse, error) {
	return s.assignmentRepository.GetAssignments(ctx)
}

func (s *assignmentService) UpdateAssignment(ctx context.Context, id int, req *UpdateAssignmentRequest) (*AssignmentResponse, error) {
	// Here you could add validation or business logic if needed
	return s.assignmentRepository.UpdateAssignment(ctx, id, req)
}

func (s *assignmentService) CreateAssignment(ctx context.Context, req *CreateAssignmentRequest) (*AssignmentResponse, error) {
	// You could add validation or business logic here if needed
	return s.assignmentRepository.CreateAssignment(ctx, req)
}
