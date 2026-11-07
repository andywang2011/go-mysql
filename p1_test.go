package mysql

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
)

func TestNewDB_ctxCancelledBeforeOpen(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := NewDB(ctx, validConfig()); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	//lint:ignore SA1012 intentional nil context in test
	if _, err := NewDB(nil, Config{}); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestNewDB_ctxCancelClosesLateConnection(t *testing.T) {
	oldOpen := openGorm
	defer func() { openGorm = oldOpen }()

	orm, _, mock := mustTestGorm(t)
	mock.ExpectClose()
	release := make(chan struct{})
	exited := make(chan struct{})
	openGorm = func(Config, glog.Interface) (*gorm.DB, error) {
		defer close(exited)
		<-release
		return orm, nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() { time.Sleep(10 * time.Millisecond); cancel() }()
	cfg := validConfig()
	cfg.Logger = glog.Discard
	if _, err := NewDB(ctx, cfg); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	close(release)
	<-exited
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if mock.ExpectationsWereMet() == nil {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("late connection was not closed")
}

func TestNewDBWithRetry_ctxCancelStopsRetry(t *testing.T) {
	oldNew := newDBWithConfig
	defer func() { newDBWithConfig = oldNew }()
	var calls atomic.Int32
	newDBWithConfig = func(context.Context, Config) (DB, error) { calls.Add(1); return nil, errors.New("down") }
	cfg := validConfig()
	cfg.Retry = RetryConfig{MaxAttempts: 50, InitialDelay: time.Hour, MaxDelay: time.Hour, Multiplier: 1}
	ctx, cancel := context.WithCancel(context.Background())
	go func() { time.Sleep(30 * time.Millisecond); cancel() }()
	start := time.Now()
	_, err := NewDBWithRetry(ctx, cfg)
	if err == nil || calls.Load() < 1 || time.Since(start) > time.Second {
		t.Fatal(err, calls.Load())
	}
	// Cancelled before the first attempt: wrapped context error, no connect.
	calls.Store(0)
	if _, err := NewDBWithRetry(ctx, cfg); !errors.Is(err, context.Canceled) || calls.Load() != 0 {
		t.Fatal(err, calls.Load())
	}
	//lint:ignore SA1012 intentional nil context in test
	if _, err := NewDBWithRetry(nil, Config{}); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestRetryManager_ctxPaths(t *testing.T) {
	rm := newRetryManager(RetryConfig{MaxAttempts: 3, InitialDelay: time.Millisecond, MaxDelay: time.Millisecond, Multiplier: 1})
	ctx, cancel := context.WithCancel(context.Background())
	if err := rm.execute(ctx, func(context.Context) error { cancel(); return errors.New("x") }); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	//lint:ignore SA1012 intentional nil context in test
	if err := rm.execute(nil, func(context.Context) error { return nil }); err != nil {
		t.Fatal(err)
	}
	if err := sleepCtx(context.Background(), 0); err != nil {
		t.Fatal(err)
	}
	if d := rm.jitter(0); d != 0 {
		t.Fatal(d)
	}
	slow := newRetryManager(RetryConfig{MaxAttempts: 5, InitialDelay: time.Hour, MaxDelay: time.Hour, Multiplier: 1})
	ctx2, cancel2 := context.WithCancel(context.Background())
	go func() { time.Sleep(20 * time.Millisecond); cancel2() }()
	if err := slow.execute(ctx2, func(context.Context) error { return errors.New("x") }); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}
