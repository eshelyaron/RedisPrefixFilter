package main

import (
	"benchmark/redis"
	"github.com/sirupsen/logrus"
	"math"
)

const avgOfIterations = 10

var (
	redisClient *redis.RedisClient
	capacity    = math.Pow(10, 6)
)

func main() {
	initRedis()
	err := reserveBF()
	if err != nil {
		logrus.WithError(err).Error("reserveBF")
		return
	}
	err = reserveCF()
	if err != nil {
		logrus.WithError(err).Error("reserveCF")
		return
	}
	err = reservePF()
	if err != nil {
		logrus.WithError(err).Error("reservePF")
		return
	}
	err = runCorrectnessTests()
	if err != nil {
		logrus.WithError(err).Error("RunCorrectnessTests")
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
	err := redisClient.ReserveBF(errorRate, int64(capacity))
	if err != nil {
		return err
	}
	return nil
}

func reserveCF() error {
	err := redisClient.ReserveCF(int64(capacity))
	if err != nil {
		return err
	}
	return nil
}

func reservePF() error {
	err := redisClient.ReservePF(int64(capacity))
	if err != nil {
		return err
	}
	return nil
}
