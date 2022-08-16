package api

import (
	"encoding/json"
	"net/http"

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
}

func (s *Server) convertFromSeedsToHashes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sc SeedCluster
		var hc HashCluster
		if error := json.NewDecoder(r.Body).Decode(&sc); error != nil {
			http.Error(w, error.Error(), http.StatusBadRequest)
			return
		}
		hc.Hashes = sc.Seeds
		s.HashClusters = append(s.HashClusters, hc)

		w.Header().Set("Content-Type", "application/json")
		if error := json.NewEncoder(w).Encode(hc); error != nil {
			http.Error(w, error.Error(), http.StatusInternalServerError)
			return
		}
	}
}
