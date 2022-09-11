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

	testResult, err = testMAddPerNumberOfParalleledTests()
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

	testResult, err = testMExistsPerNumberOfItems()
	if err != nil {
		return err
	}
	processTestResults("testMExistsPerNumberOfItems", testResult)

	logrus.Infof("All load tests finished successfully!")
	return nil
}

func testMAddPerNumberOfParalleledTests() (map[string][]TestResult, error) {
	res := make(map[string][]TestResult)
	bfResults := make([]TestResult, 0)
	pfResults := make([]TestResult, 0)
	numberOfAdds := 200
	for i := 1; i < 300; i += 10 {
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
	strArr := utils.GetRandomStrings(10, 20, numberOfAdds)
	if filterType == "bf" {
		return utils.RunTest(avgOfIterations, parallelTests, func() error {
			return redisClient.MAddBF(strArr)
		})
	}
	if filterType == "pf" {
		return utils.RunTest(avgOfIterations, parallelTests, func() error {
			return redisClient.MAddPF(strArr)
		})
	}
	if filterType == "cf" {
		logrus.Infof("cuckoo filter don't support the MADD command")
		logrus.Exit(0)
	}
	return 0, unknownFilterType
}

func testMExistsPerNumberOfItems() (map[string][]TestResult, error) {
	res := make(map[string][]TestResult)
	bfResults := make([]TestResult, 0)
	cfResults := make([]TestResult, 0)
	pfResults := make([]TestResult, 0)
	parallelTests := 1
	for i := 1; i < 1500; i += 100 {
		d, err := testMExistsTime("bf", parallelTests, i)
		if err != nil {
			return nil, err
		}
		bfResults = append(bfResults, TestResult{
			X: i,
			Y: d,
		})
		d, err = testMExistsTime("cf", parallelTests, i)
		if err != nil {
			return nil, err
		}
		cfResults = append(cfResults, TestResult{
			X: i,
			Y: d,
		})

		d, err = testMExistsTime("pf", parallelTests, i)
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

func testExistsPerNumberOfParalleledTests() (map[string][]TestResult, error) {
	res := make(map[string][]TestResult)
	bfResults := make([]TestResult, 0)
	cfResults := make([]TestResult, 0)
	pfResults := make([]TestResult, 0)
	numberOfItems := 1
	for i := 1; i < 300; i += 20 {
		d, err := testMExistsTime("bf", i, numberOfItems)
		if err != nil {
			return nil, err
		}
		bfResults = append(bfResults, TestResult{
			X: i,
			Y: d,
		})
		d, err = testMExistsTime("cf", i, numberOfItems)
		if err != nil {
			return nil, err
		}
		cfResults = append(cfResults, TestResult{
			X: i,
			Y: d,
		})

		d, err = testMExistsTime("pf", i, numberOfItems)
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

func testMExistsTime(filterType string, parallelTests, numberOfItems int) (time.Duration, error) {
	strArr := utils.GetRandomStrings(10, 20, numberOfItems)
	if filterType == "bf" {
		return utils.RunTest(avgOfIterations, parallelTests, func() error {
			_, err := redisClient.MExistsBF(strArr)
			return err
		})
	}
	if filterType == "cf" {
		return utils.RunTest(avgOfIterations, parallelTests, func() error {
			_, err := redisClient.MExistsCF(strArr)
			return err
		})
	}
	if filterType == "pf" {
		return utils.RunTest(avgOfIterations, parallelTests, func() error {
			_, err := redisClient.MExistsPF(strArr)
			return err
		})
	}
	return 0, unknownFilterType
}
