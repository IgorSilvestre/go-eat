package services

import (
	"testing"

	"restaurant-api/internal/core/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByClerkID(clerkID string) (*domain.User, error) {
	args := m.Called(clerkID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetAll() ([]domain.User, error) {
	args := m.Called()
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUserService_Create_DuplicateEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	email := "duplicate@example.com"
	existingUser := &domain.User{Email: email}

	mockRepo.On("GetByEmail", email).Return(existingUser, nil)

	user, err := service.Create("John", email, "123", "clerk123")

	assert.Nil(t, user)
	assert.ErrorIs(t, err, domain.ErrEmailAlreadyExists)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Create_DuplicateClerkID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	clerkID := "duplicate_clerk"
	email := "new@example.com"
	existingUser := &domain.User{ClerkID: clerkID}

	mockRepo.On("GetByEmail", email).Return(nil, nil)
	mockRepo.On("GetByClerkID", clerkID).Return(existingUser, nil)

	user, err := service.Create("John", email, "123", clerkID)

	assert.Nil(t, user)
	assert.ErrorIs(t, err, domain.ErrClerkIDAlreadyExists)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Update_DuplicateEmail(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	oldEmail := "old@example.com"
	newEmail := "new_duplicate@example.com"
	clerkID := "clerk123"

	user := &domain.User{ID: userID, Email: oldEmail, ClerkID: clerkID}
	existingUser := &domain.User{ID: uuid.New(), Email: newEmail}

	mockRepo.On("GetByID", userID).Return(user, nil)
	mockRepo.On("GetByEmail", newEmail).Return(existingUser, nil)

	updatedUser, err := service.Update(userID, "John", newEmail, "123", clerkID)

	assert.Nil(t, updatedUser)
	assert.ErrorIs(t, err, domain.ErrEmailAlreadyExists)
	mockRepo.AssertExpectations(t)
}
