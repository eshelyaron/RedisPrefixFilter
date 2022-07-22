package main

import (
	"benchmark/redis"
	"time"
)

var redisClient *redis.RedisClient

type TestResult struct {
	x int
	y time.Duration
}

func main() {
	err := initRedis()
	if err != nil {
		println(err)
		return
	}

}

func initRedis() error {
	redisClient = redis.NewRedisClient()
	err := redisClient.ReserveBF()
	if err != nil {
		return err
	}
	//TODO reserve PF
	return nil
}

//func testMAdd() {
//
//}
//
//func testMAddTime(numberOfItems int) (*TestResult, error) {
//	strArr := utils.GetRandomStrings(10, 20, 300)
//
//	return utils.RunTest(1, func() error {
//		return redisClient.MAddBF(strArr)
//	})
//}

//
