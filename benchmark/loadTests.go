package main

import (
	"benchmark/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
)

type TestResult struct {
	X int           `json:"x"`
	Y time.Duration `json:"y"`
}

var unknownFilterType = errors.New("unknown filter type")

func processTestResults(testName string, t map[string][]TestResult) {
	printTestResults(testName, t)
	dumpTestResults(testName, t)
}

func printTestResults(testName string, t map[string][]TestResult) {
	logrus.Infof("**** %s results ****", testName)
	for k, v := range t {
		logrus.Infof("********** %s **********", k)
		for _, res := range v {
			logrus.Infof("x : %d, y : %v", res.X, res.Y.String())
		}
	}
}

func dumpTestResults(testName string, t map[string][]TestResult) {
	file, _ := json.Marshal(t)
	_ = ioutil.WriteFile(fmt.Sprintf("./results/%s.json", testName), file, 0644)
}

func RunLoadTests() error {
	var err error
	var testResult map[string][]TestResult

	//testResult, err = testMAddPerNumberOfParalleledTests()
	//if err != nil {
	//	return err
	//}
	//processTestResults("testMAddPerNumberOfParalleledTests", testResult)

	//testResult, err = testMAddPerNumberOfItems()
	//if err != nil {
	//	return err
	//}
	//processTestResults("testMAddPerNumberOfItems", testResult)

	testResult, err = testExistsPerNumberOfParalleledTests()
	if err != nil {
		return err
	}
	processTestResults("testExistsPerNumberOfParalleledTests", testResult)
	return nil
}

func testMAddPerNumberOfParalleledTests() (map[string][]TestResult, error) {
	res := make(map[string][]TestResult)
	bfResults := make([]TestResult, 0)
	pfResults := make([]TestResult, 0)
	numberOfAdds := 200
	for i := 1; i < 200; i += 10 {
		d, err := testMAddTime("bf", i, numberOfAdds)
		if err != nil {
			return nil, err
		}
		bfResults = append(bfResults, TestResult{
			X: i,
			Y: d,
		})

		d2, err := testMAddTime("pf", i, numberOfAdds)
		if err != nil {
			return nil, err
		}
		pfResults = append(pfResults, TestResult{
			X: i,
			Y: d2,
		})
	}
	res["bf"] = bfResults
	res["pf"] = pfResults
	return res, nil
}

func testMAddPerNumberOfItems() (map[string][]TestResult, error) {
	res := make(map[string][]TestResult)
	bfResults := make([]TestResult, 0)
	pfResults := make([]TestResult, 0)
	parallelTests := 1
	for i := 1; i < 1500; i += 10 {
		d, err := testMAddTime("bf", parallelTests, i)
		if err != nil {
			return nil, err
		}
		bfResults = append(bfResults, TestResult{
			X: i,
			Y: d,
		})
		d2, err := testMAddTime("pf", parallelTests, i)
		if err != nil {
			return nil, err
		}
		pfResults = append(pfResults, TestResult{
			X: i,
			Y: d2,
		})
	}
	res["bf"] = bfResults
	res["pf"] = pfResults
	return res, nil
}

func testMAddTime(filterType string, parallelTests, numberOfAdds int) (time.Duration, error) {
	if filterType == "bf" {
		return utils.RunTest(avgOfIterations, parallelTests, func() error {
			strArr := utils.GetRandomStrings(10, 20, numberOfAdds)
			return redisClient.MAddBF(strArr)
		})
	}
	if filterType == "pf" {
		return utils.RunTest(avgOfIterations, parallelTests, func() error {
			strArr := utils.GetRandomStrings(10, 20, numberOfAdds)
			return redisClient.MAddPF(strArr)
		})
	}
	return 0, unknownFilterType
}

func testExistsPerNumberOfParalleledTests() (map[string][]TestResult, error) {
	res := make(map[string][]TestResult)
	bfResults := make([]TestResult, 0)
	cfResults := make([]TestResult, 0)
	pfResults := make([]TestResult, 0)
	for i := 1; i < 300; i += 20 {
		d, err := testExistsTime("bf", i)
		if err != nil {
			return nil, err
		}
		bfResults = append(bfResults, TestResult{
			X: i,
			Y: d,
		})
		d, err = testExistsTime("cf", i)
		if err != nil {
			return nil, err
		}
		cfResults = append(cfResults, TestResult{
			X: i,
			Y: d,
		})

		d, err = testExistsTime("pf", i)
		if err != nil {
			return nil, err
		}
		pfResults = append(pfResults, TestResult{
			X: i,
			Y: d,
		})
	}
	res["bf"] = bfResults
	res["cf"] = cfResults
	res["pf"] = pfResults
	return res, nil
}

func testExistsTime(filterType string, parallelTests int) (time.Duration, error) {
	if filterType == "bf" {
		return utils.RunTest(avgOfIterations, parallelTests, func() error {
			str := utils.GetRandomString(20)
			_, err := redisClient.ExistsBF(str)
			return err
		})
	}
	if filterType == "cf" {
		return utils.RunTest(avgOfIterations, parallelTests, func() error {
			str := utils.GetRandomString(20)
			_, err := redisClient.ExistsCF(str)
			return err
		})
	}
	if filterType == "pf" {
		return utils.RunTest(avgOfIterations, parallelTests, func() error {
			str := utils.GetRandomString(20)
			_, err := redisClient.ExistsPF(str)
			return err
		})
	}
	return 0, unknownFilterType
}
