package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	router http.Handler
	rdb		*redis.Client
	Mdb		*MongoClient
	config Config
}



// This is a user defined method that returns mongo.Client, 
// context.Context, context.CancelFunc and error.
// mongo.Client will be used for further database operation.
// context.Context will be used set deadlines for process.
// context.CancelFunc will be used to cancel context and 
// resource associated with it.
func connect(uri string)(*mongo.Client, context.Context, 
                          context.CancelFunc, error) {
                          
    // ctx will be used to set deadline for process, here 
    // deadline will of 30 seconds.
    ctx, cancel := context.WithTimeout(context.Background(), 
                                       10 * time.Second)
    
    // mongo.Connect return mongo.Client method
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    return client, ctx, cancel, err
}

type MongoClient struct {
	Client 	*mongo.Client
	Ctx		*context.Context
	Cancel  context.CancelFunc
}

func New(config Config) *App {
	
	client, ctx, cancel, err := connect("mongodb://root:pass@mongodb.middleware:27017/admin?replicaSet=rs0")
	if err != nil {
        panic(err)
    }

	// Release resource when the main
    // function is returned.
    // defer closeMongo(client, ctx, cancel)
	
	app := &App{
		rdb: redis.NewClient(&redis.Options{
			Addr: config.RedisAddress,
		}),
		Mdb: &MongoClient{
			Client: client,
			Ctx: &ctx,
			Cancel: cancel,
		},
		config: config,
	}

	app.loadRoutes()
	
	return app
}

func (a *App) Start(ctx context.Context) error {

	// create a server instance and store its memory address
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", a.config.ServerPort),
		Handler: a.router,
	}

	// pinging redis to confirm if it is up and running or not
	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	defer func ()  {
		if err = a.rdb.Close(); err != nil {
			fmt.Println("failed to close redis", err)
		}
	}()

	fmt.Println("Starting Server")

	// a channel for go routines to communicate 1 is buffer size
	// buffer sizes come in buffered or unbuffered (determines if channel will be blocked or unblocked)
	ch := make(chan error, 1)

	// goroutine that runs server concurrently
	go func(){
		
		// run the server and have it listen to incoming requests
		err = server.ListenAndServe()
		
		if err != nil {
			// wrap the error inside an Error wrapper
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	// wait on multiple channel operations, runs the first case that is ready
	select {
		case err = <- ch:
			return err
		case <- ctx.Done(): // returns a channel thats closed when word is done of behalf of this context
			timeout, cancel := context.WithTimeout(context.Background(), time.Second*10) // used to allow at most 10 seconds for context to shutdown
			defer cancel()

			return server.Shutdown(timeout)
	}
	// err, open := <-ch
	// if !open {

	// }


	return nil
}