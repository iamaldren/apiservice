package main

import (
	"apiservice/auth"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"runtime/debug"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/v1/{data}", auth.GetData).Methods("GET")
	router.HandleFunc("/api/v1", auth.SetData).Methods("POST")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		debug.PrintStack()
		log.Fatal(err)
	}
}
