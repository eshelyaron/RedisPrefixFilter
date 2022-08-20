package main

import (
	"benchmark/redis"
	"benchmark/utils"
	"github.com/sirupsen/logrus"
	"math"
	"time"
)

type TestResult struct {
	x int
	y time.Duration
}

func printTestResults(t []TestResult) {
	for _, res := range t {
		logrus.Infof("x : %d, y : %v", res.x, res.y.String())
	}
}

var redisClient *redis.RedisClient

const avgOfIterations = 5

func main() {
	initRedis()
	err := reserveBF()
	if err != nil {
		println(err)
		return
	}
	err = runLoadTests()
	if err != nil {
		println(err)
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

func runLoadTests() error {
	testResult, err := testMAddPerNumberOfParalleledTests()
	if err != nil {
		return err
	}
	logrus.Infof("**** testMAddPerNumberOfParalleledTests results ****")
	printTestResults(testResult)
	testResult, err = testMAddPerNumberOfItems()
	if err != nil {
		return err
	}
	logrus.Infof("**** testMAddPerNumberOfItems results ****")
	printTestResults(testResult)
	testResult, err = testExistsPerNumberOfParalleledTests()
	if err != nil {
		return err
	}
	logrus.Infof("**** testExistsPerNumberOfParalleledTests results ****")
	printTestResults(testResult)
	return nil
}

func testMAddPerNumberOfParalleledTests() ([]TestResult, error) {
	results := make([]TestResult, 0)
	numberOfAdds := 500
	for i := 1; i < 200; i += 10 {
		d, err := testMAddTime(i, numberOfAdds)
		if err != nil {
			return nil, err
		}
		results = append(results, TestResult{
			x: i,
			y: d,
		})
	}
	return results, nil
}

func testMAddPerNumberOfItems() ([]TestResult, error) {
	results := make([]TestResult, 0)
	parallelTests := 5
	for i := 1; i < 200; i += 10 {
		d, err := testMAddTime(parallelTests, i)
		if err != nil {
			return nil, err
		}
		results = append(results, TestResult{
			x: i,
			y: d,
		})
	}
	return results, nil
}

func testMAddTime(parallelTests, numberOfAdds int) (time.Duration, error) {
	return utils.RunTest(avgOfIterations, parallelTests, func() error {
		strArr := utils.GetRandomStrings(10, 20, numberOfAdds)
		return redisClient.MAddBF(strArr)
	})
}

func testExistsPerNumberOfParalleledTests() ([]TestResult, error) {
	results := make([]TestResult, 0)
	for i := 1; i < 200; i += 10 {
		d, err := testExistsTime(i)
		if err != nil {
			return nil, err
		}
		results = append(results, TestResult{
			x: i,
			y: d,
		})
	}
	return results, nil
}

func testExistsTime(parallelTests int) (time.Duration, error) {
	return utils.RunTest(avgOfIterations, parallelTests, func() error {
		str := utils.GetRandomString(20)
		_, err := redisClient.ExistsBF(str)
		return err
	})
}
