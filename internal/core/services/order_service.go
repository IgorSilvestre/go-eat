package services

import (
	"errors"
	"restaurant-api/internal/core/domain"
	"restaurant-api/internal/core/ports"

	"github.com/google/uuid"
)

type orderService struct {
	orderRepo      ports.OrderRepository
	userRepo       ports.UserRepository
	productRepo    ports.ProductRepository
	ingredientRepo ports.IngredientRepository
}

func NewOrderService(orderRepo ports.OrderRepository, userRepo ports.UserRepository, productRepo ports.ProductRepository, ingredientRepo ports.IngredientRepository) ports.OrderService {
	return &orderService{
		orderRepo:      orderRepo,
		userRepo:       userRepo,
		productRepo:    productRepo,
		ingredientRepo: ingredientRepo,
	}
}

func (s *orderService) CreateOrder(userID uuid.UUID, items []ports.OrderItemInput) (*domain.Order, error) {
	// Validate User
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	order := &domain.Order{
		UserID:      user.ID,
		PhoneNumber: user.PhoneNumber,
	}

	var orderProducts []domain.OrderProduct
	var totalPrice float64

	for _, item := range items {
		// Fetch base Product
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			return nil, errors.New("product not found: " + item.ProductID.String())
		}

		itemPrice := product.Price

		// Fetch prices for Adicionais
		if len(item.Adicionais) > 0 {
			adicionais, err := s.ingredientRepo.GetManyByIDs(item.Adicionais)
			if err != nil {
				return nil, errors.New("failed to fetch adicionais")
			}

			// Map ingredient prices by ID for quick lookup
			priceMap := make(map[uuid.UUID]float64)
			for _, ing := range adicionais {
				priceMap[ing.ID] = ing.Price
			}

			// Add price for *every* occurrence of the added ingredient ID
			for _, addID := range item.Adicionais {
				if price, exists := priceMap[addID]; exists {
					itemPrice += price
				}
			}
		}

		// (Removed items do not change the base price according to rules)

		orderProduct := domain.OrderProduct{
			ProductID:   product.ID,
			Name:        product.Name,
			Description: product.Description,
			Ingredients: product.Ingredients,
			Adicionais:  item.Adicionais,
			Removed:     item.Removed,
			Price:       itemPrice,
		}

		orderProducts = append(orderProducts, orderProduct)
		totalPrice += itemPrice
	}

	order.TotalPrice = totalPrice

	// Save Order first
	err = s.orderRepo.Create(order)
	if err != nil {
		return nil, err
	}

	// Save OrderProducts and attach OrderID
	for i := range orderProducts {
		orderProducts[i].OrderID = order.ID
		err = s.orderRepo.CreateOrderProduct(&orderProducts[i])
		if err != nil {
			return nil, err
		}
	}

	// Attach items to order struct for the response
	order.Items = orderProducts

	return order, nil
}

func (s *orderService) Get(id uuid.UUID) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	items, err := s.orderRepo.GetOrderProductsByOrderID(id)
	if err != nil {
		return nil, err
	}

	order.Items = items
	return order, nil
}

func (s *orderService) List() ([]domain.Order, error) {
	orders, err := s.orderRepo.GetAll()
	if err != nil {
		return nil, err
	}

	for i := range orders {
		items, err := s.orderRepo.GetOrderProductsByOrderID(orders[i].ID)
		if err != nil {
			return nil, err
		}
		orders[i].Items = items
	}

	return orders, nil
}
