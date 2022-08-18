package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

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

func NewServer() *Server {
	s := &Server{
		Router:       mux.NewRouter(),
		HashClusters: []HashCluster{},
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.HandleFunc("/hash", s.convertFromSeedsToHashes()).Methods("POST")
	s.HandleFunc("/hash", s.listHashes()).Methods("GET")
}

func hashingWorker(str string, chnl chan [2]string, index int) {
	var sum [32]byte = sha256.Sum256([]byte(str))
	var res [2]string
	res[0] = strconv.Itoa(index)
	res[1] = hex.EncodeToString(sum[:])
	chnl <- res
}

// function used for prosessing the seeds and POSTING them to the
// servers collection of HashClusters
func (s *Server) convertFromSeedsToHashes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sc SeedCluster
		var hc HashCluster
		var wg sync.WaitGroup
		hash := make(chan [2]string)

		// recieving request and checking for network errors
		if error := json.NewDecoder(r.Body).Decode(&sc); error != nil {
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}
		// checking if the data is invalid or in a wrong format
		if sc.Seeds == nil {
			fmt.Println("Input has wrong format")
			return
		}
		hc.Hashes = make([]string, len(sc.Seeds))

		// hashing every seed individually and concurrently using a waitGroup as a goroutine
		for i := 0; i < len(sc.Seeds); i++ {
			wg.Add(1)

			index := i

			go func() {
				defer wg.Done()
				go hashingWorker(sc.Seeds[index], hash, index)
				fmt.Println("inne")
				data := <-hash
				num, error := strconv.Atoi(data[0])
				if error != nil {
					fmt.Println(error)
				}
				hc.Hashes[num] = data[1]
				//TODO append data
			}()
		}
		wg.Wait()
		fmt.Println("ute")
		// TODO add seeds to a channel instead of manipulating the original data structure
		// adding hashes to the HashCluster structure
		fmt.Println("ute2")

		fmt.Println("ute3")
		//hc.Hashes = sc.Seeds
		s.HashClusters = append(s.HashClusters, hc)

		// setting header expecting a json object
		w.Header().Set("Content-Type", "application/json")
		if error := json.NewEncoder(w).Encode(hc); error != nil {
			http.Error(w, error.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// function used for GETTING the HashClusters that has allready been processed
func (s *Server) listHashes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if error := json.NewEncoder(w).Encode(s.HashClusters); error != nil {
			http.Error(w, error.Error(), http.StatusInternalServerError)
			return
		}
	}
}
