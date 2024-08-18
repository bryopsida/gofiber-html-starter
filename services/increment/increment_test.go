package increment

import (
	"testing"

	"github.com/bryopsida/gofiber-pug-starter/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockNumberRepository is a mock implementation of the INumberRepository interface
type MockNumberRepository struct {
	mock.Mock
}

// DeleteByID implements interfaces.INumberRepository.
func (m *MockNumberRepository) DeleteByID(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockNumberRepository) FindByID(id string) (*interfaces.Number, error) {
	args := m.Called(id)
	return args.Get(0).(*interfaces.Number), args.Error(1)
}

func (m *MockNumberRepository) Save(number interfaces.Number) error {
	args := m.Called(number)
	return args.Error(0)
}

func TestNewIncrementService(t *testing.T) {
	mockRepo := new(MockNumberRepository)
	bucket := "test-bucket"
	service := NewIncrementService(mockRepo, bucket)

	assert.NotNil(t, service)
}

func TestIncrement(t *testing.T) {

	t.Run("successful increment", func(t *testing.T) {
		mockRepo := new(MockNumberRepository)
		bucket := "test-bucket-1"
		service := NewIncrementService(mockRepo, bucket)
		mockNumber := &interfaces.Number{ID: bucket, Number: 1}
		expectedNumber := &interfaces.Number{ID: bucket, Number: 2}
		mockRepo.On("FindByID", bucket).Return(mockNumber, nil)
		mockRepo.On("Save", *expectedNumber).Return(nil)

		resp, err := service.Increment(bucket)

		assert.NoError(t, err)
		assert.Equal(t, expectedNumber.Number, resp)
		mockRepo.AssertExpectations(t)
	})

	t.Run("increment with new number", func(t *testing.T) {
		mockRepo := new(MockNumberRepository)
		bucket := "test-bucket-2"
		service := NewIncrementService(mockRepo, bucket)
		mockRepo.On("FindByID", bucket).Return(&interfaces.Number{}, interfaces.ErrNotFound)
		mockRepo.On("Save", mock.AnythingOfType("interfaces.Number")).Return(nil)

		resp, err := service.Increment(bucket)

		assert.NoError(t, err)
		assert.Equal(t, uint64(1), resp)
		mockRepo.AssertExpectations(t)
	})

	t.Run("save error", func(t *testing.T) {
		mockRepo := new(MockNumberRepository)
		bucket := "test-bucket-3"
		service := NewIncrementService(mockRepo, bucket)
		mockNumber := &interfaces.Number{ID: bucket, Number: 1}
		mockRepo.On("FindByID", bucket).Return(mockNumber, nil)
		mockRepo.On("Save", mock.Anything).Return(interfaces.ErrSaveFailed)

		resp, err := service.Increment(bucket)

		assert.Error(t, err)
		assert.Equal(t, uint64(0), resp)
		mockRepo.AssertExpectations(t)
	})
}
