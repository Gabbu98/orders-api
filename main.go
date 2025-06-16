package main

import (
	"fmt"
	"net/http"
)

// entrypoint
func main() {
	// create a server instance and storing it as a memory address
	server := &http.Server{
		Addr: ":3000",
		Handler: http.HandlerFunc(basicHandler), // a handler interface used for when our server receives a request
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