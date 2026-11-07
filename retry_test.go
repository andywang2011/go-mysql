package mysql

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryManager_execute_successFirstAttempt(t *testing.T) {
	rm := newRetryManager(DefaultRetryConfig())
	rm.sleep = func(context.Context, time.Duration) error { return nil }
	calls := 0
	err := rm.execute(context.Background(), func(context.Context) error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("calls = %d", calls)
	}
}

func TestRetryManager_execute_exhausted(t *testing.T) {
	rm := newRetryManager(RetryConfig{
		MaxAttempts:  3,
		InitialDelay: time.Millisecond,
		MaxDelay:     time.Millisecond,
		Multiplier:   2,
	})
	rm.sleep = func(context.Context, time.Duration) error { return nil }
	wantErr := errors.New("fail")
	err := rm.execute(context.Background(), func(context.Context) error {
		return wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("err = %v", err)
	}
}

func TestRetryManager_execute_contextCanceled(t *testing.T) {
	rm := newRetryManager(DefaultRetryConfig())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := rm.execute(ctx, func(context.Context) error { return errors.New("x") })
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v", err)
	}
}

func TestRetryManager_jitterRange(t *testing.T) {
	rm := newRetryManager(DefaultRetryConfig())
	delay := 100 * time.Millisecond
	j := rm.jitter(delay)
	if j < 80*time.Millisecond || j > 120*time.Millisecond {
		t.Fatalf("jitter out of range: %s", j)
	}
}
