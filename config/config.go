package config

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"runtime/debug"
)

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "aldren", //when empty means no password is set in Redis
		DB:       0,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		debug.PrintStack()
		log.Fatal(err)
	}

	fmt.Println("Successfully pinged redis with response of " + pong)
}

func GetRedisClient() *redis.Client {
	return client
}
