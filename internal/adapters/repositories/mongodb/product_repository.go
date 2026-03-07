package mongodb

import (
	"context"
	"restaurant-api/internal/core/domain"
	"restaurant-api/internal/core/ports"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type productRepo struct {
	coll *mongo.Collection
}

func NewProductRepository(db *mongo.Database) ports.ProductRepository {
	return &productRepo{
		coll: db.Collection("products"),
	}
}

func (r *productRepo) Create(product *domain.Product) error {
	product.ID = uuid.New()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	// Ensure slice is not nil for BSON serialization
	if product.Ingredients == nil {
		product.Ingredients = []uuid.UUID{}
	}

	_, err := r.coll.InsertOne(context.Background(), product)
	return err
}

func (r *productRepo) GetByID(id uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	err := r.coll.FindOne(context.Background(), bson.M{"_id": id}).Decode(&product)
	if err != nil {
		return nil, err
	}
	if product.Ingredients == nil {
		product.Ingredients = []uuid.UUID{}
	}
	return &product, nil
}

func (r *productRepo) GetAll() ([]domain.Product, error) {
	var products []domain.Product
	cursor, err := r.coll.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &products); err != nil {
		return nil, err
	}

	for i := range products {
		if products[i].Ingredients == nil {
			products[i].Ingredients = []uuid.UUID{}
		}
	}
	return products, nil
}

func (r *productRepo) Update(product *domain.Product) error {
	product.UpdatedAt = time.Now()
	_, err := r.coll.ReplaceOne(
		context.Background(),
		bson.M{"_id": product.ID},
		product,
	)
	return err
}

func (r *productRepo) Delete(id uuid.UUID) error {
	_, err := r.coll.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
