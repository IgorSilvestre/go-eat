package http

import (
	"restaurant-api/internal/core/ports"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type OrderHandler struct {
	orderService ports.OrderService
}

func NewOrderHandler(orderService ports.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

type createOrderReq struct {
	UserID uuid.UUID              `json:"user_id"`
	Items  []ports.OrderItemInput `json:"items"`
}

// Create godoc
// @Summary Create a new order
// @Description Create a new order with user ID and items
// @Tags orders
// @Accept json
// @Produce json
// @Param request body createOrderReq true "Order details"
// @Success 201 {object} domain.Order
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [post]
func (h *OrderHandler) Create(c *fiber.Ctx) error {
	var req createOrderReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	order, err := h.orderService.CreateOrder(req.UserID, req.Items)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}

// Get godoc
// @Summary Get order by ID
// @Description Get details of a single order by ID
// @Tags orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} domain.Order
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /orders/{id} [get]
func (h *OrderHandler) Get(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid order id"})
	}

	order, err := h.orderService.Get(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "order not found"})
	}

	return c.JSON(order)
}

// List godoc
// @Summary List all orders
// @Description Get a list of all orders
// @Tags orders
// @Produce json
// @Success 200 {array} domain.Order
// @Failure 500 {object} map[string]string
// @Router /orders [get]
func (h *OrderHandler) List(c *fiber.Ctx) error {
	orders, err := h.orderService.List()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(orders)
}
