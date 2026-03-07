package main

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "restaurant-api/docs"
	"restaurant-api/internal/adapters/handlers/http"
	"restaurant-api/internal/adapters/repositories/mongodb"
	"restaurant-api/internal/core/services"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/yokeTH/gofiber-scalar/scalar/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @title Restaurant API
// @version 1.0
// @description REST API for a restaurant backend
// @host localhost:7000
// @BasePath /api/v1
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	mongoURI := os.Getenv("DB_MONGO_URL")
	if mongoURI == "" {
		log.Fatal("DB_MONGO_URL environment variable is not set")
	}

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	db := client.Database("restaurant")

	// Repositories
	userRepo := mongodb.NewUserRepository(db)
	ingredientRepo := mongodb.NewIngredientRepository(db)
	productRepo := mongodb.NewProductRepository(db)
	orderRepo := mongodb.NewOrderRepository(db)

	// Services
	userService := services.NewUserService(userRepo)
	ingredientService := services.NewIngredientService(ingredientRepo)
	productService := services.NewProductService(productRepo)
	orderService := services.NewOrderService(orderRepo, userRepo, productRepo, ingredientRepo)

	// Handlers
	userHandler := http.NewUserHandler(userService)
	ingredientHandler := http.NewIngredientHandler(ingredientService)
	productHandler := http.NewProductHandler(productService)
	orderHandler := http.NewOrderHandler(orderService)

	app := fiber.New()

	app.Get("/docs/*", scalar.New())

	app.Get("/health", HealthCheck)

	api := app.Group("/api/v1")

	// Users routes
	users := api.Group("/users")
	users.Post("/", userHandler.Create)
	users.Get("/", userHandler.List)
	users.Get("/:id", userHandler.Get)
	users.Put("/:id", userHandler.Update)
	users.Delete("/:id", userHandler.Delete)

	// Ingredients routes
	ingredients := api.Group("/ingredients")
	ingredients.Post("/", ingredientHandler.Create)
	ingredients.Get("/", ingredientHandler.List)
	ingredients.Get("/:id", ingredientHandler.Get)
	ingredients.Put("/:id", ingredientHandler.Update)
	ingredients.Delete("/:id", ingredientHandler.Delete)

	// Products routes
	products := api.Group("/products")
	products.Post("/", productHandler.Create)
	products.Get("/", productHandler.List)
	products.Get("/:id", productHandler.Get)
	products.Put("/:id", productHandler.Update)
	products.Delete("/:id", productHandler.Delete)

	// Orders routes
	orders := api.Group("/orders")
	orders.Post("/", orderHandler.Create)
	orders.Get("/", orderHandler.List)
	orders.Get("/:id", orderHandler.Get)

	log.Fatal(app.Listen(":7000"))
}

// HealthCheck godoc
// @Summary Show the status of server
// @Description get the status of server
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}
