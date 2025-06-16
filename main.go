package main

import (
	"context"
	"fmt"

	"github.com/Gabbu98/orders-api/application"
)

// entrypoint
func main() {
	app := application.New()

	err := app.Start(context.TODO())

	if err != nil {
		fmt.Println("failed to start app: ", err)
	}
}
