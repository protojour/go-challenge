package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrReadJson  = errors.New("failed to read json")
	ErrWriteJson = errors.New("failed to write json")
	ErrPretty    = errors.New("failed to pretty print")
)

// SeedsJson
type SeedsJson struct {
	Seeds []string `json:"seeds"`
}

func (sj *SeedsJson) ReadRequest(w http.ResponseWriter, r *http.Request) (err error) {
	defer r.Body.Close()

	if err = json.NewDecoder(r.Body).Decode(&sj); err != nil {
		err = fmt.Errorf("%w, %s", ErrReadJson, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return
}

func (sj SeedsJson) PrettyPrint() (err error) {
	if b, e := json.MarshalIndent(sj, "", "    "); e != nil {
		err = fmt.Errorf("%w, %s", ErrPretty, e)
	} else {
		fmt.Printf("%s\n", b)
	}

	return
}

// HashesJson
type HashesJson struct {
	Hashes []string `json:"hashes"`
}

func (hj *HashesJson) ReadResponse(resp *http.Response) (err error) {
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&hj); err != nil {
		err = fmt.Errorf("%w, %s", ErrReadJson, err)
	}

	return
}

func (hj HashesJson) WriteResponse(w http.ResponseWriter) (err error) {
	if b, e := json.Marshal(hj); e != nil {
		err = fmt.Errorf("%w, %s", ErrWriteJson, e)
		http.Error(w, e.Error(), http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}

	return
}

func (hj HashesJson) PrettyPrint() (err error) {
	if b, e := json.MarshalIndent(hj, "", "    "); e != nil {
		err = fmt.Errorf("%w, %s", ErrPretty, e)
	} else {
		fmt.Printf("%s\n", b)
	}

	return
}
