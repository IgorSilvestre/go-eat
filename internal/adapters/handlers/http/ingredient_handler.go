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

func (h *IngredientHandler) List(c *fiber.Ctx) error {
	ingredients, err := h.ingredientService.List()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(ingredients)
}

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
