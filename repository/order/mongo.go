package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/Gabbu98/orders-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepo struct {
	Client *mongo.Client
}

func getMongoContext(m *MongoRepo) (*mongo.Collection) {
	collection := m.Client.Database("orders").Collection("orders")

	return collection
}

func (m *MongoRepo) Insert(ctx context.Context, order model.Order) error {
	collection := getMongoContext(m)

	_, err := collection.InsertOne(ctx, order)

	return err
}

func(m *MongoRepo) FindByID(ctx context.Context, id uint64) (model.Order, error) {
	collection := getMongoContext(m)

	var result model.Order
	err := collection.FindOne(ctx, bson.M{"_id":id}).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.Order{}, ErrNoExist
	} else if err != nil {
		return model.Order{}, fmt.Errorf("get error: %w", err)
	}

	return result, err
}

func(m *MongoRepo) DeleteByID(ctx context.Context, id uint64) error {
	collection := getMongoContext(m)

	_, err := collection.DeleteOne(ctx, bson.M{"_id":id})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return ErrNoExist
	} else if err != nil {
		return fmt.Errorf("failed to remove order: %w", err)
	}
	
	return err
}

func(m *MongoRepo) Update(ctx context.Context, order model.Order) error {
	collection := getMongoContext(m)

	_, err := collection.UpdateByID(ctx, order.OrderID, order)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return ErrNoExist
	} else if err != nil {
		return fmt.Errorf("failed to remove order: %w", err)
	} 

	return err
}

func (m *MongoRepo) FindAll(ctx context.Context, page FindAllPage) (FindResult, error) {
	collection := getMongoContext(m)
	var orders []model.Order

	l := int64(page.Size)
	skip := int64(page.Offset)
	
	result, err := collection.Find(ctx, options.FindOptions{Limit: &l, Skip: &skip})
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get orders: %w", err)
	}
	
	if err := result.All(ctx, &orders); err != nil {
		return FindResult{}, err
	}

	return FindResult{
		Orders: orders,
		Cursor: uint64(skip),
	}, nil
}

// https://dev.to/hackmamba/build-a-rest-api-with-golang-and-mongodb-fiber-version-4la0
// https://www.geeksforgeeks.org/how-to-use-go-with-mongodb/