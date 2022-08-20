package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

const bfKey = "bfKey"
const pfKey = "pfKey"

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6380",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &RedisClient{
		client: rdb,
	}
}

func (r *RedisClient) ReserveBF(errorRate float64, capacity int64) error {
	_, err := bfReserveScript.exec(context.Background(), r.client, []string{bfKey}, errorRate, capacity)
	return err
}

func (r *RedisClient) DeleteBF() error {
	res := r.client.Del(context.Background(), bfKey)
	return res.Err()
}

func (r *RedisClient) MAddBF(items []string) error {
	_, err := bfMAddScript.exec(context.Background(), r.client, []string{bfKey}, items)
	return err
}

func (r *RedisClient) ExistsBF(item string) (bool, error) {
	_, err := bfExistsScript.exec(context.Background(), r.client, []string{bfKey}, item)
	return false, err
}
