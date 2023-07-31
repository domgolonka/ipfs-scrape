package main

import (
	"domgolonka/blockparty/database"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// return a list of all the ipfs
func ipfs(db *database.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ipfs, err := db.SelectIpfs()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		bodyMarshal(w, http.StatusOK, ipfs)

	}
}

// return specific ipfs by cid
func cidIpfs(db *database.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cid := mux.Vars(r)["cid"]
		ipfs, err := db.SelectIpfsByCid(cid)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid cid")
			return
		}
		bodyMarshal(w, http.StatusOK, ipfs)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	bodyMarshal(w, code, map[string]string{"error": message})
}

func bodyMarshal(w http.ResponseWriter, code int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(payload)
	if err != nil {
		w.Write([]byte(`{"success":false,"message":"Some internal error occurred"}`))
		return errors.New("some internal error occurred")
	}
	w.WriteHeader(code)
	w.Write(resp)
	return nil
}

func healthHandler(db *database.Client) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check the health of the server and return a status code accordingly
		if serverIsHealthy(db) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Server is healthy")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Server is not healthy")
		}
	}
}

// serverIsHealthy checks if the database and server is healthy
func serverIsHealthy(db *database.Client) bool {
	err := db.Ping()
	if err != nil {
		return false
	}
	return true
}

// handleRequests handles all the requests
func handleRequests(db *database.Client) {
	http.HandleFunc("/tokens", ipfs(db))
	http.HandleFunc("/tokens/{cid}", cidIpfs(db))
	http.HandleFunc("/health", healthHandler(db))
	log.Fatal(http.ListenAndServe(":80", nil))
}

func main() {
	db, err := database.NewClient()
	if err != nil {
		panic(err)
	}
	handleRequests(db)
}
