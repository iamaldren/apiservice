package main

import (
	"apiservice/auth"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/v1/{data}", auth.CheckData).Methods("GET")
	log.Println(http.ListenAndServe(":8080", router), "Service started listening to port 8080")
}
