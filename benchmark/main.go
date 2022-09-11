package main

import (
	"benchmark/redis"
	"github.com/sirupsen/logrus"
	"math"
	"os"
	"strconv"
	"time"
)

const avgOfIterations = 32

var (
	redisClient *redis.RedisClient
	capacity    = math.Pow(10, 6)
)

func main() {
	initRedis()
	initFilters()
	cmd, filterType, parallelTests, records := getCMD()

	var d time.Duration
	var err error

	switch cmd {
	case "generalTests":
		generalTests()
		return
	case "madd":
		d, err = testMAddTime(filterType, parallelTests, records)
	case "mexists":
		d, err = testMExistsTime(filterType, parallelTests, records)
	default:
		logrus.Fatalf("unsupported cmd %s", cmd)
	}

	if err != nil {
		logrus.WithError(err).Fatal()
	}
	logrus.Infof("total duration - %v", d)
}

func generalTests() {
	err := runCorrectnessTests()
	if err != nil {
		logrus.WithError(err).Fatal("RunCorrectnessTests")
	}
	err = RunLoadTests()
	if err != nil {
		logrus.WithError(err).Fatal("RunLoadTests")
	}
}

func getCMD() (string, string, int, int) {
	args := os.Args[1:]
	if len(args) == 0 {
		return "generalTests", "", 0, 0
	}
	if len(args) != 4 {
		logrus.Fatalf("not enough arguments were passed. expected 4, got %d. \n", len(args))
	}
	parallelTests, err := strconv.Atoi(args[2])
	if err != nil {
		logrus.Fatalf("expected # of parallel tests to be integer, got %s.\n", args[2])
	}
	records, err := strconv.Atoi(args[3])
	if err != nil {
		logrus.Fatalf("expected # of records to be integer, got %s.\n", args[3])
	}
	return args[0], args[1], parallelTests, records
}

func initRedis() {
	redisClient = redis.NewRedisClient()
}

func initFilters() {
	err := reserveBF()
	if err != nil {
		logrus.WithError(err).Fatal("reserveBF")
	}
	err = reserveCF()
	if err != nil {
		logrus.WithError(err).Fatal("reserveCF")
	}
	err = reservePF()
	if err != nil {
		logrus.WithError(err).Fatal("reservePF")
	}
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
