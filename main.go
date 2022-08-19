package main

import (
	"net/http"

	"example.com/m/api"
)

func main() {
	srv := api.NewServer()            // initiates the server
	http.ListenAndServe(":5000", srv) // sets up the server on an arbitrary available port
}
