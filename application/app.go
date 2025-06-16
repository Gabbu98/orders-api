package application

import (
	"context"
	"fmt"
	"net/http"
)

type App struct {
	router http.Handler
}

func New() *App {
	app := &App{
		router: loadRoutes(),
	}

	return app
}

func (a *App) Start(ctx context.Context) error {

	// create a server instance and store its memory address
	server := &http.Server{
		Addr: ":3000",
		Handler: a.router,
	}

	// run the server and have it listen to incoming requests
	err := server.ListenAndServe()

	if err != nil {
		// wrap the error inside an Error wrapper
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}