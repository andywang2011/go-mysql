package mysql

import (
	"context"
	"math/rand"
	"time"
)

type retryManager struct {
	config RetryConfig
	rand   *rand.Rand
	sleep  func(ctx context.Context, d time.Duration) error
}

func newRetryManager(cfg RetryConfig) *retryManager {
	return &retryManager{
		config: cfg,
		rand:   rand.New(rand.NewSource(time.Now().UnixNano())),
		sleep: func(ctx context.Context, d time.Duration) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(d):
				return nil
			}
		},
	}
}

func (rm *retryManager) execute(ctx context.Context, op func(ctx context.Context) error) error {
	var lastErr error
	delay := rm.config.InitialDelay
	for attempt := 0; attempt < rm.config.MaxAttempts; attempt++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err := op(ctx); err == nil {
			return nil
		} else {
			lastErr = err
		}
		if attempt < rm.config.MaxAttempts-1 {
			sleep := rm.jitter(delay)
			if err := rm.sleep(ctx, sleep); err != nil {
				return err
			}
			delay = time.Duration(float64(delay) * rm.config.Multiplier)
			if delay > rm.config.MaxDelay {
				delay = rm.config.MaxDelay
			}
		}
	}
	return lastErr
}

func (rm *retryManager) jitter(delay time.Duration) time.Duration {
	jitterRange := delay / 5
	if jitterRange == 0 {
		jitterRange = 1
	}
	j := time.Duration(rm.rand.Int63n(int64(jitterRange * 2)))
	return delay - jitterRange + j
}
