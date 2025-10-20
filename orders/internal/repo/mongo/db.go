// Package mongo is a nice package
package mongo

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/odlev/animal-delivery/orders/internal/domain"
	"github.com/odlev/animal-delivery/orders/internal/repo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	client *mongo.Client
}

func Init(ctx context.Context, uri string) (*Storage, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &Storage{client: client}, err
}

func (s *Storage) Close(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}

func (s *Storage) OrdersCollection() *mongo.Collection {
	return s.client.Database(os.Getenv("MONGO_DB_NAME")).Collection("orders")
}

type OrdersCollection struct {
	coll *mongo.Collection
}

func NewOrdersCollection(storage *Storage) *OrdersCollection {
	return &OrdersCollection{coll: storage.OrdersCollection()}
}

func (c *OrdersCollection) CreateOrder(ctx context.Context, order domain.Order) error {
	const op = "repo.Mongo.Create"
	doc, err := fromDomain(order)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	res, err := c.coll.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	id, ok := res.InsertedID.(string)
	if ok {
		if order.OrderID.String() != id {
			return fmt.Errorf("%s: compare id: %s and inserted id: %s fall", op, order.OrderID.String(), id)
		}
	}

	return nil
}

func (c *OrdersCollection) GetOrder(ctx context.Context, uuid uuid.UUID) (domain.Order, error) {
	const op = "repo.Mongo.Get"

	var doc mongoOrder
	err := c.coll.FindOne(ctx, bson.M{"_id": uuid.String()}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return domain.Order{}, fmt.Errorf("%s: %w: %w", op, repo.ErrNotFound, err)
	}
	return doc.toDomain()
}

func (c *OrdersCollection) DeleteOrder(ctx context.Context, uuid uuid.UUID) error {
	const op = "repo.Mongo.Delete"

	res, err := c.coll.DeleteOne(ctx, bson.M{"_id": uuid.String()})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	delCount := res.DeletedCount
	if delCount == 0 {
		return repo.ErrNotFound
	}
	return nil
}

func fromDomain(o domain.Order) (mongoOrder, error) {
	return mongoOrder{
		ID:           o.OrderID.String(),
		Status:       o.Status,
		CustomerID:   o.CustomerID.String(),
		AnimalType:   o.AnimalType,
		AnimalAge:    o.AnimalAge,
		DeleteReason: o.DeleteReason,
		UpdatedAt:    o.UpdatedAt,
	}, nil
}

func (m mongoOrder) toDomain() (domain.Order, error) {
	orderID, err := uuid.Parse(m.ID)
	if err != nil {
		return domain.Order{}, err
	}
	customerID, err := uuid.Parse(m.CustomerID)
	if err != nil {
		return domain.Order{}, err
	}
	return domain.Order{
		OrderID:      orderID,
		Status:       m.Status,
		CustomerID:   customerID,
		AnimalType:   m.AnimalType,
		AnimalAge:    m.AnimalAge,
		DeleteReason: m.DeleteReason,
		UpdatedAt:    m.UpdatedAt,
	}, nil
}
