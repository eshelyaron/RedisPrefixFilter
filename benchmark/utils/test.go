package utils

import (
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"time"
)

func RunTest(avgOfIterations, parallelTests int, f func() error) (time.Duration, error) {
	logrus.Infof("%s\n", time.Now().UTC().String())
	totalDuration := time.Duration(0)
	for i := 0; i < avgOfIterations; i++ {
		g, _ := errgroup.WithContext(context.Background())
		start := time.Now()
		for i := 0; i < parallelTests; i++ {
			g.Go(f)
		}
		err := g.Wait()
		if err != nil {
			logrus.Error(err)
			return 0, err
		}
		totalDuration += time.Now().Sub(start)
	}
	return time.Duration(int64(totalDuration) / int64(avgOfIterations)), nil
}
