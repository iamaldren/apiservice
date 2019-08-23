package auth

import(
	"apiservice/config"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"runtime/debug"
)

type Message struct {
	PassedData string `json:"passedData"`
	StoredData string `json:"storedData"`
}

type Status struct {
	StatusCode int `json:"statusCode"`
	StatusDesc string `json:"statusDesc"`
}

func CheckData(w http.ResponseWriter, r *http.Request) {
	client, err := config.GetRedisClient()
	if err != nil {
		log.Fatal(err)
		debug.PrintStack()
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer client.Close()

	w.Header().Set("Content-Type", "application/json")
	data := mux.Vars(r)["data"]

	flag := true

	result, err := client.Get(data).Result()
	if err == redis.Nil {
		fmt.Println("Data doesn't exists in Redis")
		flag  = false
	} else if err != nil {
		debug.PrintStack()
		panic(err)
	}

	var message Message
	message.PassedData = data
	message.StoredData = result

	sendResponse(flag, w, message)
}

func sendResponse(flag bool, w http.ResponseWriter, message Message) {
	if flag {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(message)
	} else {
		var status Status
		status.StatusCode = 404
		status.StatusDesc = "Record Not Found"

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(status)
	}
}