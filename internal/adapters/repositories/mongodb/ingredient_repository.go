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

type ingredientRepo struct {
	coll *mongo.Collection
}

func NewIngredientRepository(db *mongo.Database) ports.IngredientRepository {
	return &ingredientRepo{
		coll: db.Collection("ingredients"),
	}
}

func (r *ingredientRepo) Create(ingredient *domain.Ingredient) error {
	ingredient.ID = uuid.New()
	ingredient.CreatedAt = time.Now()
	ingredient.UpdatedAt = time.Now()

	_, err := r.coll.InsertOne(context.Background(), ingredient)
	return err
}

func (r *ingredientRepo) GetByID(id uuid.UUID) (*domain.Ingredient, error) {
	var ingredient domain.Ingredient
	err := r.coll.FindOne(context.Background(), bson.M{"_id": id}).Decode(&ingredient)
	if err != nil {
		return nil, err
	}
	return &ingredient, nil
}

func (r *ingredientRepo) GetAll() ([]domain.Ingredient, error) {
	var ingredients []domain.Ingredient
	cursor, err := r.coll.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &ingredients); err != nil {
		return nil, err
	}
	return ingredients, nil
}

func (r *ingredientRepo) Update(ingredient *domain.Ingredient) error {
	ingredient.UpdatedAt = time.Now()
	_, err := r.coll.ReplaceOne(
		context.Background(),
		bson.M{"_id": ingredient.ID},
		ingredient,
	)
	return err
}

func (r *ingredientRepo) Delete(id uuid.UUID) error {
	_, err := r.coll.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}

func (r *ingredientRepo) GetManyByIDs(ids []uuid.UUID) ([]domain.Ingredient, error) {
	var ingredients []domain.Ingredient
	if len(ids) == 0 {
		return ingredients, nil
	}
	cursor, err := r.coll.Find(context.Background(), bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &ingredients); err != nil {
		return nil, err
	}
	return ingredients, nil
}
