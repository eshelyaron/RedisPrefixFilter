package main

import (
	"benchmark/utils"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
)

type TestResult struct {
	X int           `json:"x"`
	Y time.Duration `json:"y"`
}

func processTestResults(testName string, t []TestResult) {
	printTestResults(testName, t)
	dumpTestResults(testName, t)

}

func printTestResults(testName string, t []TestResult) {
	logrus.Infof("**** %s results ****", testName)
	for _, res := range t {
		logrus.Infof("x : %d, y : %v", res.X, res.Y.String())
	}
}

func dumpTestResults(testName string, t []TestResult) {
	file, _ := json.Marshal(t)
	_ = ioutil.WriteFile(fmt.Sprintf("./results/%s.json", testName), file, 0644)
}

func RunLoadTests() error {
	testResult, err := testMAddPerNumberOfParalleledTests()
	if err != nil {
		return err
	}
	processTestResults("testMAddPerNumberOfParalleledTests", testResult)
	testResult, err = testMAddPerNumberOfItems()
	if err != nil {
		return err
	}
	processTestResults("testMAddPerNumberOfItems", testResult)
	testResult, err = testExistsPerNumberOfParalleledTests()
	if err != nil {
		return err
	}
	processTestResults("testExistsPerNumberOfParalleledTests", testResult)
	return nil
}

func testMAddPerNumberOfParalleledTests() ([]TestResult, error) {
	results := make([]TestResult, 0)
	numberOfAdds := 500
	for i := 1; i < 50; i += 10 {
		d, err := testMAddTime(i, numberOfAdds)
		if err != nil {
			return nil, err
		}
		results = append(results, TestResult{
			X: i,
			Y: d,
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
			X: i,
			Y: d,
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
			X: i,
			Y: d,
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
