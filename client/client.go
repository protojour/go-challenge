package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"example.com/m/api"
)

func main() {
	// creating data
	sc := api.SeedCluster{Seeds: []string{"abc", "def", "xyz"}}
	// creating bodies for http request
	postBody, _ := json.Marshal(sc)
	responseBody := bytes.NewBuffer(postBody)
	// sending postrequest to server
	resp, err := http.Post("http://localhost:5000/hash", "application/json", responseBody)
	// handle request error
	if err != nil {
		log.Fatalf("A fatal error occured (%s)", err)
	}
	defer resp.Body.Close()
	//reading the response from the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.SetOutput(os.Stdout)
	// logging answer to stdout
	log.Print(sb)
}
