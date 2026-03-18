package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"restaurant-api/internal/core/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) Create(name, description string, ingredients []uuid.UUID, price float64, image string) (*domain.Product, error) {
	args := m.Called(name, description, ingredients, price, image)
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductService) Get(id uuid.UUID) (*domain.Product, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductService) List() ([]domain.Product, error) {
	args := m.Called()
	return args.Get(0).([]domain.Product), args.Error(1)
}

func (m *MockProductService) Update(id uuid.UUID, name, description string, ingredients []uuid.UUID, price float64, image string) (*domain.Product, error) {
	args := m.Called(id, name, description, ingredients, price, image)
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductService) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockStorageService struct {
	mock.Mock
}

func (m *MockStorageService) UploadImage(ctx context.Context, reader io.Reader, filename string, contentType string) (string, error) {
	args := m.Called(ctx, reader, filename, contentType)
	return args.String(0), args.Error(1)
}

func TestProductHandler_Create(t *testing.T) {
	app := fiber.New()
	mockProductService := new(MockProductService)
	mockStorageService := new(MockStorageService)
	handler := NewProductHandler(mockProductService, mockStorageService)

	app.Post("/products", handler.Create)

	// Prepare multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("name", "Test Product")
	_ = writer.WriteField("price", "10.50")
	_ = writer.WriteField("description", "Test Description")

	part, _ := writer.CreateFormFile("image", "test.jpg")
	part.Write([]byte("fake-image-content"))
	writer.Close()

	expectedProduct := &domain.Product{
		ID:          uuid.New(),
		Name:        "Test Product",
		Description: "Test Description",
		Price:       10.50,
		Image:       "http://storage/test.jpg",
	}

	mockStorageService.On("UploadImage", mock.Anything, mock.Anything, "test.jpg", mock.Anything).Return("http://storage/test.jpg", nil)
	mockProductService.On("Create", "Test Product", "Test Description", []uuid.UUID{}, 10.50, "http://storage/test.jpg").Return(expectedProduct, nil)

	req := httptest.NewRequest("POST", "/products", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result domain.Product
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, expectedProduct.Name, result.Name)
	assert.Equal(t, expectedProduct.Image, result.Image)

	mockStorageService.AssertExpectations(t)
	mockProductService.AssertExpectations(t)
}

func TestProductHandler_Update(t *testing.T) {
	app := fiber.New()
	mockProductService := new(MockProductService)
	mockStorageService := new(MockStorageService)
	handler := NewProductHandler(mockProductService, mockStorageService)

	app.Put("/products/:id", handler.Update)

	productID := uuid.New()
	existingProduct := &domain.Product{
		ID:          productID,
		Name:        "Old Name",
		Description: "Old Description",
		Ingredients: []uuid.UUID{},
		Price:       5.00,
		Image:       "http://storage/old.jpg",
	}

	// Prepare multipart form with updated name only
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("name", "Updated Name")
	writer.Close()

	mockProductService.On("Get", productID).Return(existingProduct, nil)
	mockProductService.On("Update", productID, "Updated Name", "Old Description", []uuid.UUID{}, 5.00, "http://storage/old.jpg").Return(&domain.Product{
		ID:          productID,
		Name:        "Updated Name",
		Description: "Old Description",
		Price:       5.00,
		Image:       "http://storage/old.jpg",
	}, nil)

	req := httptest.NewRequest("PUT", "/products/"+productID.String(), body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result domain.Product
	json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "Updated Name", result.Name)

	mockProductService.AssertExpectations(t)
}
