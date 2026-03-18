package http

import (
	"restaurant-api/internal/core/ports"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type IngredientHandler struct {
	ingredientService ports.IngredientService
}

func NewIngredientHandler(ingredientService ports.IngredientService) *IngredientHandler {
	return &IngredientHandler{ingredientService: ingredientService}
}

type createIngredientReq struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Create godoc
// @Summary Create a new ingredient
// @Description Create a new ingredient with name and price
// @Tags ingredients
// @Accept json
// @Produce json
// @Param request body createIngredientReq true "Ingredient details"
// @Success 201 {object} domain.Ingredient
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/ingredients [post]
func (h *IngredientHandler) Create(c *fiber.Ctx) error {
	var req createIngredientReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ingredient, err := h.ingredientService.Create(req.Name, req.Price)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(ingredient)
}

// Get godoc
// @Summary Get ingredient by ID
// @Description Get details of a single ingredient by ID
// @Tags ingredients
// @Produce json
// @Param id path string true "Ingredient ID"
// @Success 200 {object} domain.Ingredient
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/ingredients/{id} [get]
func (h *IngredientHandler) Get(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ingredient id"})
	}

	ingredient, err := h.ingredientService.Get(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "ingredient not found"})
	}

	return c.JSON(ingredient)
}

// List godoc
// @Summary List all ingredients
// @Description Get a list of all ingredients
// @Tags ingredients
// @Produce json
// @Success 200 {array} domain.Ingredient
// @Failure 500 {object} map[string]string
// @Router /api/v1/ingredients [get]
func (h *IngredientHandler) List(c *fiber.Ctx) error {
	ingredients, err := h.ingredientService.List()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(ingredients)
}

// Update godoc
// @Summary Update an ingredient
// @Description Update ingredient details by ID
// @Tags ingredients
// @Accept json
// @Produce json
// @Param id path string true "Ingredient ID"
// @Param request body createIngredientReq true "Ingredient details"
// @Success 200 {object} domain.Ingredient
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/ingredients/{id} [put]
func (h *IngredientHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ingredient id"})
	}

	var req createIngredientReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ingredient, err := h.ingredientService.Update(id, req.Name, req.Price)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(ingredient)
}

// Delete godoc
// @Summary Delete an ingredient
// @Description Delete an ingredient by ID
// @Tags ingredients
// @Param id path string true "Ingredient ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/ingredients/{id} [delete]
func (h *IngredientHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ingredient id"})
	}

	if err := h.ingredientService.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
