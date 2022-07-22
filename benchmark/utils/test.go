package utils

import (
	"context"
	"golang.org/x/sync/errgroup"
	"time"
)

func RunTest(parallelTests int, f func() error) (time.Duration, error) {
	g, _ := errgroup.WithContext(context.Background())
	start := time.Now()
	for i := 0; i < parallelTests; i++ {
		g.Go(f)
	}
	err := g.Wait()
	if err != nil {
		return 0, err
	}
	return time.Now().Sub(start), nil
}
