package mysql

import (
	"fmt"
	"time"

	glog "gorm.io/gorm/logger"
)

// Config holds MySQL connection and pool settings.
type Config struct {
	// Address is host:port (often a proxy VIP). Ignored for DSN building when DSN is set.
	Address string
	// Username is the MySQL user (ignored when DSN is set).
	Username string
	// Password is the MySQL password (ignored when DSN is set).
	Password string
	// Database is the default schema name (required when DSN is empty).
	Database string
	// DSN when non-empty is used as-is (proxy/cloud ready-made DSN); Address/Username may still be set for logging.
	DSN string

	// MaxIdleConns is the maximum idle connections in the pool.
	MaxIdleConns int
	// MaxOpenConns is the maximum open connections in the pool.
	MaxOpenConns int
	// TimeoutSec is the DSN timeout parameter in seconds.
	TimeoutSec int
	// SlowThreshold logs queries slower than this duration via GORM.
	SlowThreshold time.Duration

	// ConnMaxLifetime is the maximum lifetime of a pooled connection (default 5m).
	ConnMaxLifetime time.Duration
	// ConnMaxIdleTime closes idle connections after this duration (0 = driver default).
	ConnMaxIdleTime time.Duration
	// ConnectTimeout bounds startup connect+retry when > 0.
	ConnectTimeout time.Duration

	// Retry controls exponential-backoff during connect.
	Retry RetryConfig
	// Logger is an optional GORM logger; nil falls back to default warn logger.
	Logger glog.Interface
}

const (
	defaultTimeoutSec      = 20
	defaultSlowThreshold   = time.Second
	defaultMaxOpenConns    = 50
	defaultMaxIdleConns    = 10
	defaultConnMaxLifetime = 5 * time.Minute
)

// RetryConfig controls exponential-backoff retry behavior during connect.
type RetryConfig struct {
	// MaxAttempts is the total number of attempts (including the first).
	MaxAttempts int
	// InitialDelay is the first backoff pause.
	InitialDelay time.Duration
	// MaxDelay caps exponential backoff.
	MaxDelay time.Duration
	// Multiplier scales delay after each failed attempt.
	Multiplier float64
}

// ConfigDefault returns pool and timeout defaults; Address/Database/DSN must still be set before use.
func ConfigDefault() Config {
	return Config{
		TimeoutSec:      defaultTimeoutSec,
		SlowThreshold:   defaultSlowThreshold,
		MaxOpenConns:    defaultMaxOpenConns,
		MaxIdleConns:    defaultMaxIdleConns,
		ConnMaxLifetime: defaultConnMaxLifetime,
		Retry:           DefaultRetryConfig(),
	}
}

// DefaultRetryConfig returns production-safe retry defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Multiplier:   2.0,
	}
}

// Validate validates config and fills compatible defaults in place.
func (c *Config) Validate() error {
	if c.DSN == "" {
		if c.Address == "" {
			return fmt.Errorf("%w: address is required when dsn is empty", ErrInvalidConfig)
		}
		if c.Database == "" {
			return fmt.Errorf("%w: database is required when dsn is empty", ErrInvalidConfig)
		}
	}
	if c.TimeoutSec <= 0 {
		c.TimeoutSec = defaultTimeoutSec
	}
	if c.SlowThreshold <= 0 {
		c.SlowThreshold = defaultSlowThreshold
	}
	if c.MaxOpenConns <= 0 {
		c.MaxOpenConns = defaultMaxOpenConns
	}
	if c.MaxIdleConns <= 0 {
		c.MaxIdleConns = defaultMaxIdleConns
	}
	if c.ConnMaxLifetime <= 0 {
		c.ConnMaxLifetime = defaultConnMaxLifetime
	}
	if c.Retry.MaxAttempts <= 0 {
		c.Retry = DefaultRetryConfig()
	}
	if c.Retry.InitialDelay <= 0 {
		c.Retry.InitialDelay = DefaultRetryConfig().InitialDelay
	}
	if c.Retry.MaxDelay <= 0 {
		c.Retry.MaxDelay = DefaultRetryConfig().MaxDelay
	}
	if c.Retry.Multiplier <= 0 {
		c.Retry.Multiplier = DefaultRetryConfig().Multiplier
	}
	if c.Retry.MaxDelay < c.Retry.InitialDelay {
		c.Retry.MaxDelay = c.Retry.InitialDelay
	}
	return nil
}
