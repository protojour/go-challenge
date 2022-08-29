package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-challenge/internal/server"
	"net/http"
	"os"
)

func main() {
	sj := server.SeedsJson{Seeds: os.Args[1:]}
	hashResult := server.HashesJson{}

	if e := sj.PrettyPrint(); e != nil {
		fmt.Println(e)
	} else if b, e := json.Marshal(sj); e != nil {
		fmt.Println("Failed to marshal seeds to json", e)
	} else if resp, e := http.Post("http://localhost:5000/hash", "application/json", bytes.NewBuffer(b)); e != nil {
		fmt.Println("failed to post seeds to server", e)
	} else if e := hashResult.ReadResponse(resp); e != nil {
		fmt.Println(e)
	} else if e = hashResult.PrettyPrint(); e != nil {
		fmt.Println(e)
	}
}
