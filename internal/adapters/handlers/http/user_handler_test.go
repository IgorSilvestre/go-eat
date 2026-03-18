package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"restaurant-api/internal/core/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(name, email, phone, clerkID string) (*domain.User, error) {
	args := m.Called(name, email, phone, clerkID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) Get(id uuid.UUID) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) List() ([]domain.User, error) {
	args := m.Called()
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserService) Update(id uuid.UUID, name, email, phone, clerkID string) (*domain.User, error) {
	args := m.Called(id, name, email, phone, clerkID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUserHandler_Create_Conflict(t *testing.T) {
	app := fiber.New()
	mockUserService := new(MockUserService)
	handler := NewUserHandler(mockUserService)

	app.Post("/users", handler.Create)

	reqReq := createUserReq{
		Name:    "John",
		Email:   "duplicate@example.com",
		ClerkID: "clerk123",
	}
	body, _ := json.Marshal(reqReq)

	mockUserService.On("Create", "John", "duplicate@example.com", "", "clerk123").Return(nil, domain.ErrEmailAlreadyExists)

	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusConflict, resp.StatusCode)

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, domain.ErrEmailAlreadyExists.Error(), result["error"])

	mockUserService.AssertExpectations(t)
}
