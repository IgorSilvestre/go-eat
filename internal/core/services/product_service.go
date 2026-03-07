package services

import (
	"restaurant-api/internal/core/domain"
	"restaurant-api/internal/core/ports"

	"github.com/google/uuid"
)

type productService struct {
	repo ports.ProductRepository
}

func NewProductService(repo ports.ProductRepository) ports.ProductService {
	return &productService{repo: repo}
}

func (s *productService) Create(name, description string, ingredients []uuid.UUID, price float64) (*domain.Product, error) {
	product := &domain.Product{
		Name:        name,
		Description: description,
		Ingredients: ingredients,
		Price:       price,
	}
	err := s.repo.Create(product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) Get(id uuid.UUID) (*domain.Product, error) {
	return s.repo.GetByID(id)
}

func (s *productService) List() ([]domain.Product, error) {
	return s.repo.GetAll()
}

func (s *productService) Update(id uuid.UUID, name, description string, ingredients []uuid.UUID, price float64) (*domain.Product, error) {
	product, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	product.Name = name
	product.Description = description
	product.Ingredients = ingredients
	product.Price = price

	err = s.repo.Update(product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
