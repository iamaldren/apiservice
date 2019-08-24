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

	router.HandleFunc("/v1/get/{key}", auth.GetData).Methods("GET")
	router.HandleFunc("/v1/set", auth.SetData).Methods("POST")

	router.HandleFunc("/v1/zset/{table}", auth.ZAddData).Methods("POST")
	router.HandleFunc("/v1/zgetall/{table}", auth.ZRangeByScoreGetAll).Methods("GET")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		debug.PrintStack()
		log.Fatal(err)
	}
}
