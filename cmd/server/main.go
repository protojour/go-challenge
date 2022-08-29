package main

import (
	"go-challenge/internal/server"
)

func main() {
	server.Serve("localhost", "5000")
}
