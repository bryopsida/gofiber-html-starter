package increment

import (
	"log/slog"

	"github.com/bryopsida/gofiber-pug-starter/interfaces"
)

// ServiceImpl is the implementation of IncrementServiceServer
type ServiceImpl struct {
	repo   interfaces.INumberRepository
	bucket string
}

// NewIncrementService creates a new ServiceImpl
// - repo: INumberRepository number repository
// - bucket: string bucket name
func NewIncrementService(repo interfaces.INumberRepository, bucket string) interfaces.IIncrementService {
	return &ServiceImpl{
		repo:   repo,
		bucket: bucket,
	}
}

// Increment increments a number
// - id: the ID of the number to increment
// Returns the incremented number if successful, otherwise returns an error
func (s *ServiceImpl) Increment(id string) (uint64, error) {
	number, err := s.repo.FindByID(s.bucket)
	if err != nil {
		slog.Info("Bucket not found, creating new bucket", "bucket", s.bucket)
		number = &interfaces.Number{ID: s.bucket, Number: 0}
	}
	slog.Info("Incrementing number", "number", number.Number)
	number.Number++
	slog.Info("Saving number", "number", number.Number)
	saveErr := s.repo.Save(*number)
	if saveErr != nil {
		slog.Error("Error saving number", "error", saveErr)
		return 0, saveErr
	}

	resp := number.Number
	slog.Info("Returning incremented number", "number", resp)
	return resp, nil
}
