package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// entrypoint
func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/hello", basicHandler)

	// create a server instance and storing it as a memory address
	server := &http.Server{
		Addr: ":3000",
		Handler: router, // a handler interface used for when our server receives a request
	}

	// run and listen
	err := server.ListenAndServe()

	if err != nil {
		fmt.Println("Gailed to listen to server", err)
	}
}

/*
w - will allow use to write our http response
r - pointer for inbound request received from the client-side
*/
func basicHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}