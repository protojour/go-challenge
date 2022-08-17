package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
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

func hashingWorker(str string) string {
	var sum [32]byte = sha256.Sum256([]byte(str))
	st := hex.EncodeToString(sum[:])
	fmt.Println(st)
	return st
}

func (s *Server) convertFromSeedsToHashes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sc SeedCluster
		var hc HashCluster
		var wg sync.WaitGroup
		if error := json.NewDecoder(r.Body).Decode(&sc); error != nil {
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}

		for i := 0; i < len(sc.Seeds); i++ {
			wg.Add(1)

			index := i

			go func() {
				defer wg.Done()
				sc.Seeds[index] = hashingWorker(sc.Seeds[index])
			}()
		}

		wg.Wait()

		hc.Hashes = sc.Seeds
		s.HashClusters = append(s.HashClusters, hc)

		w.Header().Set("Content-Type", "application/json")
		if error := json.NewEncoder(w).Encode(hc); error != nil {
			http.Error(w, error.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) listHashes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if error := json.NewEncoder(w).Encode(s.HashClusters); error != nil {
			http.Error(w, error.Error(), http.StatusInternalServerError)
			return
		}
	}
}
