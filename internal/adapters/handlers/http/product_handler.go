package http

import (
	"restaurant-api/internal/core/ports"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductHandler struct {
	productService ports.ProductService
}

func NewProductHandler(productService ports.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

type createProductReq struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Ingredients []uuid.UUID `json:"ingredients"`
	Price       float64     `json:"price"`
}

// Create godoc
// @Summary Create a new product
// @Description Create a new product with name, description, ingredients and price
// @Tags products
// @Accept json
// @Produce json
// @Param request body createProductReq true "Product details"
// @Success 201 {object} domain.Product
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products [post]
func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var req createProductReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	product, err := h.productService.Create(req.Name, req.Description, req.Ingredients, req.Price)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

// Get godoc
// @Summary Get product by ID
// @Description Get details of a single product by ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} domain.Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func (h *ProductHandler) Get(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid product id"})
	}

	product, err := h.productService.Get(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "product not found"})
	}

	return c.JSON(product)
}

// List godoc
// @Summary List all products
// @Description Get a list of all products
// @Tags products
// @Produce json
// @Success 200 {array} domain.Product
// @Failure 500 {object} map[string]string
// @Router /products [get]
func (h *ProductHandler) List(c *fiber.Ctx) error {
	products, err := h.productService.List()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(products)
}

// Update godoc
// @Summary Update a product
// @Description Update product details by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body createProductReq true "Product details"
// @Success 200 {object} domain.Product
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [put]
func (h *ProductHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid product id"})
	}

	var req createProductReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	product, err := h.productService.Update(id, req.Name, req.Description, req.Ingredients, req.Price)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(product)
}

// Delete godoc
// @Summary Delete a product
// @Description Delete a product by ID
// @Tags products
// @Param id path string true "Product ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [delete]
func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid product id"})
	}

	if err := h.productService.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
