package main

import (
	"benchmark/utils"
	"errors"
	"github.com/sirupsen/logrus"
)

func runCorrectnessTests() error {
	var err error
	err = testBF()
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Infof("All correctness tests passed!")
	return nil
}

func testBF() error {
	record := utils.GetRandomString(80)

	exists, err := redisClient.ExistsBF(record)
	if err != nil {
		return err
	}
	if exists != 0 {
		return errors.New("exists should be false")
	}

	err = redisClient.AddBF(record)
	if err != nil {
		return err
	}

	exists, err = redisClient.ExistsBF(record)
	if err != nil {
		return err
	}
	if exists != 1 {
		return errors.New("exists should be true")
	}

	records := utils.GetRandomStrings(10, 20, 10)

	existsArr, err := redisClient.MExistsBF(records)
	if err != nil {
		return err
	}
	for _, e := range existsArr {
		if e.(int64) != 0 {
			return errors.New("exists should be false")
		}
	}

	err = redisClient.MAddBF(records)
	if err != nil {
		return err
	}

	existsArr, err = redisClient.MExistsBF(records)
	if err != nil {
		return err
	}

	for _, e := range existsArr {
		if e.(int64) != 1 {
			return errors.New("exists should be true")
		}
	}

	return nil
}

func testPF() error {
	record := utils.GetRandomString(80)

	exists, err := redisClient.ExistsPF(record)
	if err != nil {
		return err
	}
	if exists != 0 {
		return errors.New("exists should be false")
	}

	err = redisClient.AddPF(record)
	if err != nil {
		return err
	}

	exists, err = redisClient.ExistsPF(record)
	if err != nil {
		return err
	}
	if exists != 1 {
		return errors.New("exists should be true")
	}

	records := utils.GetRandomStrings(10, 20, 10)

	existsArr, err := redisClient.MExistsPF(records)
	if err != nil {
		return err
	}
	for _, e := range existsArr {
		if e.(int64) != 0 {
			return errors.New("exists should be false")
		}
	}

	err = redisClient.MAddPF(records)
	if err != nil {
		return err
	}

	existsArr, err = redisClient.MExistsPF(records)
	if err != nil {
		return err
	}

	for _, e := range existsArr {
		if e.(int64) != 1 {
			return errors.New("exists should be true")
		}
	}

	return nil
}
