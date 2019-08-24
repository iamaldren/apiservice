package auth

import (
	"apiservice/config"
	"apiservice/customutil"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"io/ioutil"
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

type Data struct {
	Key string `json:"key"`
	Value string `json:"value"`
}

func GetData(w http.ResponseWriter, r *http.Request) {
	client := config.GetRedisClient()

	data := mux.Vars(r)["data"]

	result, err := client.Get(data).Result()
	if err == redis.Nil {
		fmt.Println("Data doesn't exists in Redis")
		var status Status
		status.StatusCode = 404
		status.StatusDesc = "Record Not Found"

		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(status)
	} else if err != nil {
		debug.PrintStack()
		panic(err)
	} else {
		var message Message
		message.PassedData = data
		message.StoredData = result

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(message)
	}
}

func SetData(w http.ResponseWriter, r *http.Request) {
	client := config.GetRedisClient()

	body, err := ioutil.ReadAll(r.Body)
	triggerError(err)

	var data Data
	json.Unmarshal(body, &data)

	err = client.Set(data.Key, data.Value, 0).Err()
	createResponseForRedisInsert(err,w)
}

func ZAddData(w http.ResponseWriter, r *http.Request) {
	client := config.GetRedisClient()

	body, err := ioutil.ReadAll(r.Body)
	triggerError(err)

	table := mux.Vars(r)["table"]
	strData := customutil.FormatJsonForRedis(string(body))

	count, err := client.ZCard(table).Result()
	triggerError(err)

	if count != 0 {
		sets, err := client.ZRevRange(table, 0, 0).Result()
		triggerError(err)

		rank, err := client.ZRank(table, sets[0]).Result()
		triggerError(err)

		count = rank + 2
	} else {
		count = 1
	}

	err = client.ZAdd(table, &redis.Z{
		Score: float64(count),
		Member: strData,
	}).Err()
	createResponseForRedisInsert(err,w)
}

func triggerError(err error) {
	if err != nil {
		debug.PrintStack()
		panic(err)
	}
}

func createResponseForRedisInsert(err error, w http.ResponseWriter) {
	if err != nil {
		fmt.Println("Error encountered in storing data to Redis.")
		var status Status
		status.StatusCode = 400
		status.StatusDesc = "Bad Request"

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(status)
	} else {
		var status Status
		status.StatusCode = 200
		status.StatusDesc = "Success"

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(status)
	}
}