package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("REST_PORT")
	if port == "" {
		port = "8000"
	}

	r := mux.NewRouter()
	r.HandleFunc("/", Index).Methods("GET")
	r.HandleFunc("/items", GetItems).Methods("GET")
	r.HandleFunc("/items/{id}", GetItem).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func Index(w http.ResponseWriter, r *http.Request)    {}
func GetItems(w http.ResponseWriter, r *http.Request) {}
func GetItem(w http.ResponseWriter, r *http.Request)  {}
