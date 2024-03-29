package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"math"
	"os"
)

type RedisClient struct {
	client *redis.Client
}

const bfKey = "bfKey"
const cfKey = "cfKey"
const pfKey = "pfKey"

var (
	isPrefixFilterSupported bool
	defaultErrorRate        = float64(1) / math.Pow(10, 6)
)

func NewRedisClient() *RedisClient {
	addr := "localhost:6379"
	if envAddr := os.Getenv("redisURL"); envAddr != "" {
		addr = envAddr
	}
	if isPFSupported := os.Getenv("IS_PF_SUPPORTED"); isPFSupported != "0" {
		isPrefixFilterSupported = true
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &RedisClient{
		client: rdb,
	}
}

//region BF

func (r *RedisClient) ReserveBF(errorRate float64, capacity int64) error {
	_, err := bfReserveScript.exec(context.Background(), r.client, []string{bfKey}, errorRate, capacity)
	return err
}

func (r *RedisClient) DeleteBF() error {
	res := r.client.Del(context.Background(), bfKey)
	return res.Err()
}

func (r *RedisClient) AddBF(item string) error {
	_, err := bfAddScript.exec(context.Background(), r.client, []string{bfKey}, item)
	return err
}

func (r *RedisClient) MAddBF(items []string) error {
	_, err := bfMAddScript.exec(context.Background(), r.client, []string{bfKey}, items)
	return err
}

func (r *RedisClient) ExistsBF(item string) (int64, error) {
	res, err := bfExistsScript.exec(context.Background(), r.client, []string{bfKey}, item)
	if err != nil {
		return 0, err
	}
	return res.(int64), err
}

func (r *RedisClient) MExistsBF(items []string) ([]interface{}, error) {
	res, err := bfMExistsScript.exec(context.Background(), r.client, []string{bfKey}, items)
	if err != nil {
		return nil, err
	}
	return res.([]interface{}), err
}

//endregion

//region CF

func (r *RedisClient) ReserveCF(capacity int64) error {
	_, err := cfReserveScript.exec(context.Background(), r.client, []string{cfKey}, capacity)
	return err
}

func (r *RedisClient) DeleteCF() error {
	res := r.client.Del(context.Background(), cfKey)
	return res.Err()
}

func (r *RedisClient) AddCF(items []string) error {
	_, err := cfAddScript.exec(context.Background(), r.client, []string{cfKey}, items)
	return err
}

func (r *RedisClient) ExistsCF(item string) (int64, error) {
	res, err := cfExistsScript.exec(context.Background(), r.client, []string{cfKey}, item)
	if err != nil {
		return 0, err
	}
	return res.(int64), err
}

func (r *RedisClient) MExistsCF(items []string) ([]interface{}, error) {
	res, err := cfMExistsScript.exec(context.Background(), r.client, []string{bfKey}, items)
	if err != nil {
		return nil, err
	}
	return res.([]interface{}), err
}

//endregion

//region PF

func (r *RedisClient) ReservePF(capacity int64) error {
	if !isPrefixFilterSupported {
		_, err := bfReserveScript.exec(context.Background(), r.client, []string{pfKey}, defaultErrorRate, capacity)
		return err
	}
	_, err := pfReserveScript.exec(context.Background(), r.client, []string{pfKey}, capacity)
	return err
}

func (r *RedisClient) DeletePF() error {
	res := r.client.Del(context.Background(), pfKey)
	return res.Err()
}

func (r *RedisClient) AddPF(item string) error {
	if !isPrefixFilterSupported {
		_, err := bfAddScript.exec(context.Background(), r.client, []string{pfKey}, item)
		return err
	}
	_, err := pfAddScript.exec(context.Background(), r.client, []string{pfKey}, item)
	return err
}

func (r *RedisClient) MAddPF(items []string) error {
	if !isPrefixFilterSupported {
		_, err := bfMAddScript.exec(context.Background(), r.client, []string{pfKey}, items)
		return err
	}
	_, err := pfMAddScript.exec(context.Background(), r.client, []string{pfKey}, items)
	return err
}

func (r *RedisClient) ExistsPF(item string) (int64, error) {
	if !isPrefixFilterSupported {
		res, err := bfExistsScript.exec(context.Background(), r.client, []string{pfKey}, item)
		if err != nil {
			return 0, err
		}
		return res.(int64), err
	}
	res, err := pfExistsScript.exec(context.Background(), r.client, []string{pfKey}, item)
	if err != nil {
		return 0, err
	}
	return res.(int64), err
}

func (r *RedisClient) MExistsPF(items []string) ([]interface{}, error) {
	if !isPrefixFilterSupported {
		res, err := bfMExistsScript.exec(context.Background(), r.client, []string{pfKey}, items)
		if err != nil {
			return nil, err
		}
		return res.([]interface{}), err
	}
	res, err := pfMExistsScript.exec(context.Background(), r.client, []string{pfKey}, items)
	if err != nil {
		return nil, err
	}
	return res.([]interface{}), err
}

//endregion
