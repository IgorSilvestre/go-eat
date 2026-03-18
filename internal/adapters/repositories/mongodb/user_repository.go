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

type userRepo struct {
	coll *mongo.Collection
}

func NewUserRepository(db *mongo.Database) ports.UserRepository {
	return &userRepo{
		coll: db.Collection("users"),
	}
}

func (r *userRepo) Create(user *domain.User) error {
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := r.coll.InsertOne(context.Background(), user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (r *userRepo) GetByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.coll.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.coll.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByClerkID(clerkID string) (*domain.User, error) {
	var user domain.User
	err := r.coll.FindOne(context.Background(), bson.M{"clerk_id": clerkID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetAll() ([]domain.User, error) {
	var users []domain.User
	cursor, err := r.coll.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepo) Update(user *domain.User) error {
	user.UpdatedAt = time.Now()
	_, err := r.coll.ReplaceOne(
		context.Background(),
		bson.M{"_id": user.ID},
		user,
	)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (r *userRepo) Delete(id uuid.UUID) error {
	_, err := r.coll.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
