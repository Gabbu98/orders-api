package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/Gabbu98/orders-api/application"
)

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
}
