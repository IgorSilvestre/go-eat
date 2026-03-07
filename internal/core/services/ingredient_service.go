package services

import (
	"restaurant-api/internal/core/domain"
	"restaurant-api/internal/core/ports"

	"github.com/google/uuid"
)

type ingredientService struct {
	repo ports.IngredientRepository
}

func NewIngredientService(repo ports.IngredientRepository) ports.IngredientService {
	return &ingredientService{repo: repo}
}

func (s *ingredientService) Create(name string, price float64) (*domain.Ingredient, error) {
	ingredient := &domain.Ingredient{
		Name:  name,
		Price: price,
	}
	err := s.repo.Create(ingredient)
	if err != nil {
		return nil, err
	}
	return ingredient, nil
}

func (s *ingredientService) Get(id uuid.UUID) (*domain.Ingredient, error) {
	return s.repo.GetByID(id)
}

func (s *ingredientService) List() ([]domain.Ingredient, error) {
	return s.repo.GetAll()
}

func (s *ingredientService) Update(id uuid.UUID, name string, price float64) (*domain.Ingredient, error) {
	ingredient, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	ingredient.Name = name
	ingredient.Price = price

	err = s.repo.Update(ingredient)
	if err != nil {
		return nil, err
	}
	return ingredient, nil
}

func (s *ingredientService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
