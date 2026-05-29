package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
)

// MysqlDB wraps a GORM handle.
type MysqlDB struct {
	Gorm *gorm.DB
}

var (
	openGorm = func(cfg Config, logger glog.Interface) (*gorm.DB, error) {
		return gorm.Open(mysql.Open(BuildDSN(cfg)), &gorm.Config{Logger: logger})
	}
	getSQLDB = func(orm *gorm.DB) (*sql.DB, error) {
		return safeGetSQLDB(orm)
	}
	newDBWithConfig = func(cfg Config) (DB, error) {
		return NewDB(cfg)
	}
	pingSQLDB = func(db DB, ctx context.Context) error {
		return db.Ping(ctx)
	}
)

func safeGetSQLDB(orm *gorm.DB) (sqlDB *sql.DB, err error) {
	if orm == nil {
		return nil, errors.New("mysql: gorm db is nil")
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("mysql: gorm db handle panic: %v", r)
		}
	}()
	return orm.DB()
}

func applyPoolSettings(sqlDB *sql.DB, cfg Config) {
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	if cfg.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}
}

// NewDB creates a MySQL connection without retry.
func NewDB(cfg Config) (DB, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	logger := cfg.Logger
	if logger == nil {
		logger = glog.Default.LogMode(glog.Warn)
	}
	orm, err := openGorm(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("mysql open: %w", err)
	}
	sqlDB, err := getSQLDB(orm)
	if err != nil {
		return nil, err
	}
	applyPoolSettings(sqlDB, cfg)
	return &MysqlDB{Gorm: orm}, nil
}

// NewDBWithRetry creates a MySQL connection with exponential-backoff retry.
func NewDBWithRetry(cfg Config) (DB, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	ctx := context.Background()
	if cfg.ConnectTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.ConnectTimeout)
		defer cancel()
	}
	rm := newRetryManager(cfg.Retry)
	var db DB
	var lastErr error
	err := rm.execute(ctx, func(ctx context.Context) error {
		opened, openErr := newDBWithConfig(cfg)
		if openErr != nil {
			lastErr = openErr
			return openErr
		}
		if pingErr := pingSQLDB(opened, ctx); pingErr != nil {
			_ = opened.Close()
			lastErr = pingErr
			return pingErr
		}
		db = opened
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("mysql connect retry failed after %d attempts: %w", cfg.Retry.MaxAttempts, lastErr)
	}
	return db, nil
}

// Ping checks database reachability.
func (d *MysqlDB) Ping(ctx context.Context) error {
	if d == nil || d.Gorm == nil {
		return ErrInvalidConfig
	}
	sqlDB, err := getSQLDB(d.Gorm)
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// Close closes the underlying SQL connection pool.
func (d *MysqlDB) Close() error {
	if d == nil || d.Gorm == nil {
		return nil
	}
	sqlDB, err := getSQLDB(d.Gorm)
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
