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

type orderRepo struct {
	orderColl        *mongo.Collection
	orderProductColl *mongo.Collection
}

func NewOrderRepository(db *mongo.Database) ports.OrderRepository {
	return &orderRepo{
		orderColl:        db.Collection("orders"),
		orderProductColl: db.Collection("order_products"),
	}
}

func (r *orderRepo) Create(order *domain.Order) error {
	order.ID = uuid.New()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	_, err := r.orderColl.InsertOne(context.Background(), order)
	return err
}

func (r *orderRepo) GetByID(id uuid.UUID) (*domain.Order, error) {
	var order domain.Order
	err := r.orderColl.FindOne(context.Background(), bson.M{"_id": id}).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepo) GetAll() ([]domain.Order, error) {
	var orders []domain.Order
	cursor, err := r.orderColl.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepo) CreateOrderProduct(op *domain.OrderProduct) error {
	op.ID = uuid.New()
	op.CreatedAt = time.Now()
	op.UpdatedAt = time.Now()

	// Ensure empty slices aren't nil
	if op.Ingredients == nil {
		op.Ingredients = []uuid.UUID{}
	}
	if op.Adicionais == nil {
		op.Adicionais = []uuid.UUID{}
	}
	if op.Removed == nil {
		op.Removed = []uuid.UUID{}
	}

	_, err := r.orderProductColl.InsertOne(context.Background(), op)
	return err
}

func (r *orderRepo) GetOrderProductsByOrderID(orderID uuid.UUID) ([]domain.OrderProduct, error) {
	var ops []domain.OrderProduct
	cursor, err := r.orderProductColl.Find(context.Background(), bson.M{"order_id": orderID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &ops); err != nil {
		return nil, err
	}

	for i := range ops {
		if ops[i].Ingredients == nil {
			ops[i].Ingredients = []uuid.UUID{}
		}
		if ops[i].Adicionais == nil {
			ops[i].Adicionais = []uuid.UUID{}
		}
		if ops[i].Removed == nil {
			ops[i].Removed = []uuid.UUID{}
		}
	}

	return ops, nil
}
