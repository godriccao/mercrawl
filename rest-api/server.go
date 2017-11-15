package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/godriccao/mercrawl"
	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("REST_PORT")
	if port == "" {
		port = "8000"
	}

	r := mux.NewRouter()
	r.HandleFunc("/items", GetItems).Methods("GET")
	r.HandleFunc("/items/{id}", GetItem).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func GetItems(w http.ResponseWriter, r *http.Request) {
	items := mercrawl.GetAllItems()
	json.NewEncoder(w).Encode(items)
}
func GetItem(w http.ResponseWriter, r *http.Request) {

}
