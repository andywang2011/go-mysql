package mysql

import (
	"strings"
	"testing"
	"time"
)

func TestConfig_Validate(t *testing.T) {
	cfg := Config{Address: "127.0.0.1:3306", Database: "app", Username: "u", Password: "p"}
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
	if cfg.TimeoutSec != defaultTimeoutSec {
		t.Fatal(cfg.TimeoutSec)
	}
	if cfg.SlowThreshold != defaultSlowThreshold {
		t.Fatal(cfg.SlowThreshold)
	}
	if cfg.Retry.MaxAttempts <= 0 {
		t.Fatal("retry defaults not applied")
	}
	if cfg.MaxOpenConns != defaultMaxOpenConns || cfg.MaxIdleConns != defaultMaxIdleConns {
		t.Fatalf("pool defaults: open=%d idle=%d", cfg.MaxOpenConns, cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime != defaultConnMaxLifetime {
		t.Fatalf("ConnMaxLifetime = %s", cfg.ConnMaxLifetime)
	}
}

func TestConfigDefault(t *testing.T) {
	c := ConfigDefault()
	if c.MaxOpenConns != defaultMaxOpenConns {
		t.Fatal(c.MaxOpenConns)
	}
}

func TestConfig_Validate_errors(t *testing.T) {
	cfg := Config{}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestConfig_Validate_missingDatabase(t *testing.T) {
	cfg := Config{Address: "127.0.0.1:3306"}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestConfig_Validate_defaultsNegative(t *testing.T) {
	cfg := Config{
		Address:       "127.0.0.1:3306",
		Database:      "app",
		TimeoutSec:    -1,
		SlowThreshold: -1,
		Retry: RetryConfig{
			MaxAttempts:  -1,
			InitialDelay: -1,
			MaxDelay:     -1,
			Multiplier:   -1,
		},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
	if cfg.TimeoutSec != defaultTimeoutSec {
		t.Fatalf("TimeoutSec = %d", cfg.TimeoutSec)
	}
	if cfg.SlowThreshold != defaultSlowThreshold {
		t.Fatalf("SlowThreshold = %s", cfg.SlowThreshold)
	}
}

func TestConfig_Validate_retryMaxDelayAdjusted(t *testing.T) {
	cfg := Config{
		Address:  "127.0.0.1:3306",
		Database: "app",
		Retry: RetryConfig{
			MaxAttempts:  3,
			InitialDelay: 3 * time.Second,
			MaxDelay:     time.Second,
			Multiplier:   2,
		},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
	if cfg.Retry.MaxDelay != cfg.Retry.InitialDelay {
		t.Fatalf("MaxDelay=%s InitialDelay=%s", cfg.Retry.MaxDelay, cfg.Retry.InitialDelay)
	}
}

func TestConfig_DSN(t *testing.T) {
	cfg := Config{Address: "h:3306", Database: "db", Username: "u", Password: "p", TimeoutSec: 10}
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
	dsn := BuildDSN(cfg)
	if dsn == "" {
		t.Fatal("empty dsn")
	}
	if !strings.Contains(dsn, "@tcp(h:3306)/db") {
		t.Fatalf("unexpected dsn: %s", dsn)
	}
	if !strings.Contains(dsn, "timeout=10s") {
		t.Fatalf("unexpected timeout: %s", dsn)
	}
}
