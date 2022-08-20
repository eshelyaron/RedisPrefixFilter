package main

import (
	"benchmark/redis"
	"github.com/sirupsen/logrus"
	"math"
)

var redisClient *redis.RedisClient

const avgOfIterations = 5

func main() {
	initRedis()
	err := reserveBF()
	if err != nil {
		logrus.WithError(err).Error("reserveBF")
		return
	}
	err = RunLoadTests()
	if err != nil {
		logrus.WithError(err).Error("RunLoadTests")
		return
	}
}

func initRedis() {
	redisClient = redis.NewRedisClient()
}

func reserveBF() error {
	errorRate := float64(1) / math.Pow(10, 6)
	capacity := math.Pow(10, 6)
	err := redisClient.ReserveBF(errorRate, int64(capacity))
	if err != nil {
		return err
	}
	return nil
}
