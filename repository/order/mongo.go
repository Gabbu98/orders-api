package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Gabbu98/orders-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepo struct {
	Client *mongo.Client
}

// This method closes mongoDB connection and cancel context. Func exectues first, then cancel().
func close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc){
	defer cancel()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

// returns a mongo.Client, context.Context, 
// context.CancelFunc and error.
// mongo.Client will be used for further database 
// operation. context.Context will be used set 
// deadlines for process. context.CancelFunc will 
// be used to cancel context and resource 
// associated with it.
func connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	return client, ctx, cancel, err
}

func getMongoContext(m *MongoRepo) (context.Context, *mongo.Collection) {
		// get Client, Context, CancelFunc and err from connect method.
    client, ctx, cancel, err := connect("mongodb://localhost:27017")
    if err != nil {
        panic(err)
    }
    
    // Release resource when main function is returned.
    defer close(client, ctx, cancel)

	collection := m.Client.Database("").Collection("")

	return ctx, collection
}

func (m *MongoRepo) Insert(order model.Order) error {
	ctx, collection := getMongoContext(m)

	_, err := collection.InsertOne(ctx, order)

	return err
}

func(m *MongoRepo) FindByID(id uint64) (model.Order, error) {
	ctx, collection := getMongoContext(m)

	var result model.Order
	err := collection.FindOne(ctx, bson.M{"_id":id}).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return model.Order{}, ErrNoExist
	} else if err != nil {
		return model.Order{}, fmt.Errorf("get error: %w", err)
	}

	return result, err
}

func(m *MongoRepo) DeleteByID(id uint64) error {
	ctx, collection := getMongoContext(m)

	_, err := collection.DeleteOne(ctx, bson.M{"_id":id})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return ErrNoExist
	} else if err != nil {
		return fmt.Errorf("failed to remove order: %w", err)
	}
	
	return err
}

func(m *MongoRepo) Update(order model.Order) error {
	ctx, collection := getMongoContext(m)

	_, err := collection.UpdateByID(ctx, order.OrderID, order)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return ErrNoExist
	} else if err != nil {
		return fmt.Errorf("failed to remove order: %w", err)
	} 

	return err
}

// https://dev.to/hackmamba/build-a-rest-api-with-golang-and-mongodb-fiber-version-4la0
// https://www.geeksforgeeks.org/how-to-use-go-with-mongodb/