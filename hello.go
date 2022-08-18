package main

import (
	"net/http"

	"example.com/m/api"
)

func main() {
	srv := api.NewServer()            // Initiates the server used by the client
	http.ListenAndServe(":5000", srv) // uses an arbitrary available port and the server as arguments
}
