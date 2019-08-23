package config

import(
	"fmt"
	"github.com/go-redis/redis"
)

func GetRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "", //when empty means no password is set in Redis
		DB: 0,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully pinged redis with response of " + pong)

	return client, err
}
