package http

import (
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"

	"restaurant-api/internal/core/ports"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductHandler struct {
	productService ports.ProductService
	storageService ports.StorageService
}

func NewProductHandler(productService ports.ProductService, storageService ports.StorageService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		storageService: storageService,
	}
}

func parseIngredients(ingredientsStr string) ([]uuid.UUID, error) {
	if ingredientsStr == "" {
		return []uuid.UUID{}, nil
	}
	parts := strings.Split(ingredientsStr, ",")
	var uuids []uuid.UUID
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := uuid.Parse(part)
		if err != nil {
			return nil, err
		}
		uuids = append(uuids, id)
	}
	return uuids, nil
}

func isImageFile(file *multipart.FileHeader) bool {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp"
}

// Create godoc
// @Summary Create a new product
// @Description Create a new product with name, description, ingredients, price, and image
// @Tags products
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Product Name"
// @Param description formData string false "Product Description"
// @Param ingredients formData string false "Comma-separated UUIDs of ingredients"
// @Param price formData number true "Product Price"
// @Param image formData file true "Product Image"
// @Success 201 {object} domain.Product
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/products [post]
func (h *ProductHandler) Create(c *fiber.Ctx) error {
	name := c.FormValue("name")
	if name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name is required"})
	}

	description := c.FormValue("description")

	priceStr := c.FormValue("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid price"})
	}

	ingredientsStr := c.FormValue("ingredients")
	ingredients, err := parseIngredients(ingredientsStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ingredients format"})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "image is required"})
	}

	if !isImageFile(file) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid image format"})
	}

	f, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to open image"})
	}
	defer f.Close()

	imageURL, err := h.storageService.UploadImage(c.UserContext(), f, file.Filename, file.Header.Get("Content-Type"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to upload image: " + err.Error()})
	}

	product, err := h.productService.Create(name, description, ingredients, price, imageURL)
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
// @Router /api/v1/products/{id} [get]
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
// @Router /api/v1/products [get]
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
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Product ID"
// @Param name formData string false "Product Name"
// @Param description formData string false "Product Description"
// @Param ingredients formData string false "Comma-separated UUIDs of ingredients"
// @Param price formData number false "Product Price"
// @Param image formData file false "Product Image"
// @Success 200 {object} domain.Product
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/products/{id} [put]
func (h *ProductHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid product id"})
	}

	product, err := h.productService.Get(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "product not found"})
	}

	name := c.FormValue("name")
	if name == "" {
		name = product.Name
	}

	description := c.FormValue("description")
	if description == "" {
		description = product.Description
	}

	priceStr := c.FormValue("price")
	price := product.Price
	if priceStr != "" {
		p, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid price"})
		}
		price = p
	}

	ingredientsStr := c.FormValue("ingredients")
	ingredients := product.Ingredients
	if ingredientsStr != "" {
		ings, err := parseIngredients(ingredientsStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ingredients format"})
		}
		ingredients = ings
	}

	imageURL := product.Image
	file, err := c.FormFile("image")
	if err == nil {
		if !isImageFile(file) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid image format"})
		}
		f, err := file.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to open image"})
		}
		defer f.Close()

		imgURL, err := h.storageService.UploadImage(c.UserContext(), f, file.Filename, file.Header.Get("Content-Type"))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to upload image: " + err.Error()})
		}
		imageURL = imgURL
	}

	updatedProduct, err := h.productService.Update(id, name, description, ingredients, price, imageURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(updatedProduct)
}

// Delete godoc
// @Summary Delete a product
// @Description Delete a product by ID
// @Tags products
// @Param id path string true "Product ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/products/{id} [delete]
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
