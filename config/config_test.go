package config

import (
	"errors"
	"github.com/elliotchance/redismock"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newTestRedis() *redismock.ClientMock {
	return redismock.NewMock()
}

func RedisIsAvailable(client *redismock.ClientMock) bool {
	return client.Ping().Err() == nil
}

func TestRedisCannotBePinged(t *testing.T) {
	r := newTestRedis()
	r.On("Ping").Return(redis.NewStatusResult("", errors.New("server not available")))
	assert.False(t, RedisIsAvailable(r))
}

func TestRedisCanBePinged(t *testing.T) {
	r := newTestRedis()
	r.On("Ping").Return(redis.NewStatusResult("PONG", nil))
	assert.True(t, RedisIsAvailable(r))
}

func TestGetRedisClient(t *testing.T) {
	redisCli := GetRedisClient()
	assert.True(t, redisCli != nil)
}
