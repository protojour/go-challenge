package server

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
)

type job struct {
	index      int
	hashResult string
}

func Serve(host, port string) {
	r := mux.NewRouter()

	r.HandleFunc("/hash", hashEndpoint).Methods("POST")

	srv := http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("%s:%s", host, port),
	}

	log.Printf("Hosting 'go-challenge' server on %s\n", srv.Addr)
	srv.ListenAndServe()
}

func hashEndpoint(w http.ResponseWriter, r *http.Request) {
	jsonIn := SeedsJson{}

	if e := jsonIn.ReadRequest(w, r); e != nil {
		fmt.Fprint(os.Stderr, e.Error())
	} else {
		jobs := make(chan job)
		wg := sync.WaitGroup{}

		go dispatchJobs(jsonIn, jobs, &wg)

		jsonOut := waitJobs(len(jsonIn.Seeds), &wg, jobs)

		if e := jsonOut.WriteResponse(w); e != nil {
			fmt.Fprint(os.Stderr, e.Error())
		}
	}
}

func dispatchJobs(jsonIn SeedsJson, jobs chan<- job, wg *sync.WaitGroup) {
	for i, seed := range jsonIn.Seeds {
		wg.Add(1)

		go func(jobIndex int, seed string) {
			defer wg.Done()
			sum256 := sha256.Sum256([]byte(seed))

			jobs <- job{
				index:      jobIndex,
				hashResult: hex.EncodeToString(sum256[:]),
			}
		}(i, seed)

	}
}

func waitJobs(jobCount int, wg *sync.WaitGroup, jobs <-chan job) (jsonOut HashesJson) {
	wg.Wait()

	jsonOut.Hashes = make([]string, jobCount)

	for i := 0; i < jobCount; i++ {
		j := <-jobs
		jsonOut.Hashes[j.index] = j.hashResult
	}

	return
}
