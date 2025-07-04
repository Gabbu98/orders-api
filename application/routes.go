package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Gabbu98/orders-api/handler"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	/*
	w - will allow use to write our http response
	r - pointer for inbound request received from the client-side
	*/
	router.Get("/", func (w http.ResponseWriter, r *http.Request)  {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/orders", a.loadOrderRoutes)

	a.router = router
}

// logically group routes
func (a *App) loadOrderRoutes(router chi.Router) {
	repo, err := (*a.InitRepositoryStrategies()).GetStrategy("MONGO")
	if err != nil {
		panic(err)
	}

	orderHandler := &handler.Order{
		Repo: repo,
	}

	router.Post("/", orderHandler.Create)
	router.Get("/", orderHandler.List)
	router.Get("/{id}", orderHandler.GetByID)
	router.Put("/{id}", orderHandler.UpdateByID)
	router.Delete("/{id}", orderHandler.DeleteByID)
}