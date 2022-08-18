package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

// struct definition of the json-objects and the server
type SeedCluster struct {
	Seeds []string `json:"seeds"`
}

type HashCluster struct {
	Hashes []string `json:"hashes"`
}

type Server struct {
	*mux.Router

	HashClusters []HashCluster
}

// initialization function for the server used by main.go
func NewServer() *Server {
	s := &Server{
		Router:       mux.NewRouter(),
		HashClusters: []HashCluster{},
	}
	s.Routes()
	return s
}

// setts up where the client finds the services and what type of http request they are
func (s *Server) Routes() {
	s.HandleFunc("/hash", s.ConvertToHashes()).Methods("POST")
	s.HandleFunc("/hash", s.ListHashes()).Methods("GET")
}

// a worker function used by the Waitgroup in convertToHashes() to convert
// an individual seed to a hash.
// the channel provided (chnl) takes a lists of 2 strings so that it is possible
// to pair the index of the seed with the correct index of the created hash.
// this way it is possible to maintain the correct order of hashes in the final result.
// the order is important to ensure that we know what hash the seed turned into.
func HashWorker(seed string, chnl chan [2]string, index int) {
	var res [2]string                              // used to write the result to the channel
	var sum [32]byte = sha256.Sum256([]byte(seed)) // converting seed to sha256-hash
	res[0] = strconv.Itoa(index)                   // converting index to string and adding it to res
	res[1] = hex.EncodeToString(sum[:])            // converting hash to hex-format and adding it to res
	chnl <- res                                    // writing res to channel
}

// function used for prosessing the seeds and POSTING them to the
// servers collection of HashClusters
func (s *Server) ConvertToHashes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.SetOutput(os.Stdout) // by default logging to stdout
		log.Print("initiating conversion of seeds to sha256-hashes")
		var sc SeedCluster
		var hc HashCluster
		var wg sync.WaitGroup
		channel := make(chan [2]string) // list of two strings to be able to pair hash to the correct index

		// recieving request and checking for network errors
		// the request is recieved with a SeedCluster (sc) and is
		// thus defaulted to "nil" if the input uses a wrong format
		if error := json.NewDecoder(r.Body).Decode(&sc); error != nil {
			http.Error(w, error.Error(), http.StatusBadRequest)
			log.SetOutput(os.Stderr)
			log.Print(error.Error())
			return
		}
		// checking if the input is invalid
		if sc.Seeds == nil {
			err := "error: input uses an invalid format"
			http.Error(w, err, http.StatusUnprocessableEntity)
			log.SetOutput(os.Stderr)
			log.Print(err)
			return
		}
		log.Print("input accepted")
		// initiating the list of hashes to have the same length as the list of seeds
		hc.Hashes = make([]string, len(sc.Seeds))

		// hashing every seed individually and concurrently using a waitGroup as a goroutine
		for i := 0; i < len(sc.Seeds); i++ {
			wg.Add(1)

			index := i

			go func() {
				defer wg.Done()
				go HashWorker(sc.Seeds[index], channel, index)
				data := <-channel                 // reading from channel
				ind, err := strconv.Atoi(data[0]) // converting index back to integer
				if err != nil {
					log.Print(err)
					return
				}
				hc.Hashes[ind] = data[1] // adding hash to HashCluster
			}()
		}
		wg.Wait()
		log.Print("hashing completed")

		s.HashClusters = append(s.HashClusters, hc) // adding hashcluster to the servers list of results

		// setting header expecting a json object
		w.Header().Set("Content-Type", "application/json")
		if error := json.NewEncoder(w).Encode(hc); error != nil {
			http.Error(w, error.Error(), http.StatusInternalServerError)
			log.SetOutput(os.Stderr)
			log.Print(error.Error())
			return
		}
		log.Print("request completed")
	}
}

// function used for GETTING the HashClusters that has allready been processed
// this function is used for testing
func (s *Server) ListHashes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.SetOutput(os.Stdout) // by default logging to stdout
		log.Print("initiating listing of results from previous hashing requests")
		w.Header().Set("Content-Type", "application/json")
		if error := json.NewEncoder(w).Encode(s.HashClusters); error != nil {
			http.Error(w, error.Error(), http.StatusInternalServerError)
			log.SetOutput(os.Stderr)
			log.Print(error.Error())
			return
		}
		log.Print("request completed")
	}
}
