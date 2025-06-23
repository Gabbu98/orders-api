package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/Gabbu98/orders-api/application"
	"go.mongodb.org/mongo-driver/mongo"
)

// This is a user defined method to close resources.
// This method closes mongoDB connection and cancel context.
func closeMongo(client *mongo.Client, ctx *context.Context,
           cancel context.CancelFunc){
           
    // CancelFunc to cancel to context
    defer cancel()
    
    // client provides a method to close 
    // a mongoDB connection.
    defer func(){
    
        // client.Disconnect method also has deadline.
        // returns error if any,
        if err := client.Disconnect(*ctx); err != nil{
            panic(err)
        }
    }()
}

// entrypoint
func main() {
	config := application.LoadConfig()
	app := application.New(config)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := app.Start(ctx)

	if err != nil {
		fmt.Println("failed to start app: ", err)
	}

	defer closeMongo(app.Mdb.Client, app.Mdb.Ctx, app.Mdb.Cancel)
}
