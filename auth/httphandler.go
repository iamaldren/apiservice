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

type Status struct {
	StatusCode int    `json:"statusCode"`
	StatusDesc string `json:"statusDesc"`
}

type Data struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func GetData(w http.ResponseWriter, r *http.Request) {
	client := config.GetRedisClient()

	key := mux.Vars(r)["key"]

	result, err := client.Get(key).Result()
	if err == redis.Nil {
		recordIsNil(err, w)
	} else if err != nil {
		triggerErr(err)
	} else {
		var data Data
		data.Key = key
		data.Value = result

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	}
}

func SetData(w http.ResponseWriter, r *http.Request) {
	client := config.GetRedisClient()

	body, err := ioutil.ReadAll(r.Body)
	triggerErrorIfNotNil(err)

	var data Data
	json.Unmarshal(body, &data)

	err = client.Set(data.Key, data.Value, 0).Err()
	createResponseForRedisInsert(err, w)
}

func ZAddData(w http.ResponseWriter, r *http.Request) {
	client := config.GetRedisClient()

	body, err := ioutil.ReadAll(r.Body)
	triggerErrorIfNotNil(err)

	table := mux.Vars(r)["table"]
	strData := customutil.FormatJsonForRedis(string(body))

	/**
	* This block is use to get the current count of the Sorted Set in Redis.
	* What we're doing here, is that the score of each entry that we're adding inside the
	* sorted set will be incremental of the previous one.
	*
	* You will notice that if the count is not zero, we are getting the rank (index/position)
	* of the latest entry in the Sorted Set, then we will add 2 to it to get the current score.
	* Why 2? It's because the rank always starts from Zero, and in this sample we started
	* the count of the score as 1.
	 */
	count, err := client.ZCard(table).Result()
	triggerErrorIfNotNil(err)

	if count != 0 {
		sets, err := client.ZRevRange(table, 0, 0).Result()
		triggerErrorIfNotNil(err)

		rank, err := client.ZRank(table, sets[0]).Result()
		triggerErrorIfNotNil(err)

		count = rank + 2
	} else {
		count = 1
	}
	// End of current count logic

	err = client.ZAdd(table, redis.Z{
		Score:  float64(count),
		Member: strData,
	}).Err()
	createResponseForRedisInsert(err, w)
}

func ZRangeByScoreGetAll(w http.ResponseWriter, r *http.Request) {
	client := config.GetRedisClient()

	table := mux.Vars(r)["table"]

	set, err := client.ZRangeByScore(table, redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Offset: 0,
		Count:  0,
	}).Result()
	if err == redis.Nil {
		recordIsNil(err, w)
	} else if err != nil {
		triggerErr(err)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(set)
	}
}

func ZRangeByScoreGet(w http.ResponseWriter, r *http.Request) {
	client := config.GetRedisClient()

	table := mux.Vars(r)["table"]
	score := mux.Vars(r)["score"]

	set, err := client.ZRangeByScore(table, redis.ZRangeBy{
		Min:    score,
		Max:    score,
		Offset: 0,
		Count:  0,
	}).Result()
	if err == redis.Nil {
		recordIsNil(err, w)
	} else if err != nil {
		triggerErr(err)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(set[0])
	}
}

func recordIsNil(err error, w http.ResponseWriter) {
	fmt.Println("Data doesn't exists in Redis")
	var status Status
	status.StatusCode = 404
	status.StatusDesc = "Record Not Found"

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(status)
}

func triggerErrorIfNotNil(err error) {
	if err != nil {
		triggerErr(err)
	}
}

func triggerErr(err error) {
	debug.PrintStack()
	panic(err)
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
