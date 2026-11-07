package mysql

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
)

func validConfig() Config {
	return Config{
		Address:    "127.0.0.1:3306",
		Database:   "app",
		Username:   "u",
		Password:   "p",
		TimeoutSec: 1,
	}
}

func TestNewDB_invalid(t *testing.T) {
	cfg := Config{Address: "127.0.0.1:59999", Database: "x", Username: "u", Password: "p", TimeoutSec: 1}
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
	_, err := NewDB(context.Background(), cfg)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewDB_openError(t *testing.T) {
	oldOpen := openGorm
	defer func() { openGorm = oldOpen }()
	openGorm = func(Config, glog.Interface) (*gorm.DB, error) {
		return nil, errors.New("open failed")
	}
	_, err := NewDB(context.Background(), validConfig())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewDB_ormDBError(t *testing.T) {
	oldOpen := openGorm
	oldGet := getSQLDB
	defer func() { openGorm = oldOpen }()
	defer func() { getSQLDB = oldGet }()
	openGorm = func(Config, glog.Interface) (*gorm.DB, error) {
		return &gorm.DB{}, nil
	}
	getSQLDB = func(*gorm.DB) (*sql.DB, error) {
		return nil, errors.New("db handle failed")
	}
	_, err := NewDB(context.Background(), validConfig())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewDB_success(t *testing.T) {
	oldOpen := openGorm
	defer func() { openGorm = oldOpen }()

	orm, sqlDB, mock := mustTestGorm(t)
	mock.ExpectClose()
	openGorm = func(Config, glog.Interface) (*gorm.DB, error) {
		return orm, nil
	}
	cfg := validConfig()
	cfg.Logger = glog.Discard
	db, err := NewDB(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	if db == nil || db.GORM() == nil {
		t.Fatal("expected db")
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestNewDBWithRetry_eventuallySuccess(t *testing.T) {
	oldNewDB := newDBWithConfig
	oldPing := pingSQLDB
	defer func() {
		newDBWithConfig = oldNewDB
		pingSQLDB = oldPing
	}()

	newDBWithConfig = func(context.Context, Config) (DB, error) { return &fakeDB{}, nil }
	attempts := 0
	pingSQLDB = func(DB, context.Context) error {
		attempts++
		if attempts < 2 {
			return errors.New("temporary ping error")
		}
		return nil
	}

	cfg := validConfig()
	cfg.Retry = RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 1 * time.Millisecond,
		MaxDelay:     2 * time.Millisecond,
		Multiplier:   1.1,
	}
	db, err := NewDBWithRetry(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d", attempts)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestNewDBWithRetry_exhausted(t *testing.T) {
	oldNewDB := newDBWithConfig
	oldPing := pingSQLDB
	defer func() {
		newDBWithConfig = oldNewDB
		pingSQLDB = oldPing
	}()

	newDBWithConfig = func(context.Context, Config) (DB, error) { return &fakeDB{}, nil }
	pingSQLDB = func(DB, context.Context) error {
		return errors.New("ping failed")
	}
	cfg := validConfig()
	cfg.Retry = RetryConfig{
		MaxAttempts:  2,
		InitialDelay: 1 * time.Millisecond,
		MaxDelay:     1 * time.Millisecond,
		Multiplier:   1,
	}
	_, err := NewDBWithRetry(context.Background(), cfg)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewDBWithRetry_openError(t *testing.T) {
	oldNewDB := newDBWithConfig
	defer func() { newDBWithConfig = oldNewDB }()
	newDBWithConfig = func(context.Context, Config) (DB, error) { return nil, errors.New("open failed") }

	cfg := validConfig()
	cfg.Retry = RetryConfig{
		MaxAttempts:  2,
		InitialDelay: 1 * time.Millisecond,
		MaxDelay:     1 * time.Millisecond,
		Multiplier:   1,
	}
	_, err := NewDBWithRetry(context.Background(), cfg)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestMysqlDB_Close_nil(t *testing.T) {
	var d *MysqlDB
	if err := d.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestMysqlDB_Close_dbError(t *testing.T) {
	d := &MysqlDB{gormDB: &gorm.DB{}}
	if err := d.Close(); err == nil {
		t.Fatal("expected error")
	}
}

func TestMysqlDB_Ping_success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectPing()
	mock.ExpectPing()
	orm, err := gorm.Open(gmysql.New(gmysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	d := &MysqlDB{gormDB: orm}
	if err := d.Ping(context.Background()); err != nil {
		t.Fatal(err)
	}
	_ = sqlDB.Close()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestNewDBWithRetry_connectTimeout(t *testing.T) {
	oldNewDB := newDBWithConfig
	defer func() { newDBWithConfig = oldNewDB }()
	newDBWithConfig = func(context.Context, Config) (DB, error) {
		time.Sleep(50 * time.Millisecond)
		return nil, errors.New("slow open")
	}
	cfg := validConfig()
	cfg.ConnectTimeout = 10 * time.Millisecond
	cfg.Retry = RetryConfig{
		MaxAttempts:  10,
		InitialDelay: 20 * time.Millisecond,
		MaxDelay:     20 * time.Millisecond,
		Multiplier:   1,
	}
	_, err := NewDBWithRetry(context.Background(), cfg)
	if err == nil {
		t.Fatal("expected timeout/cancel error")
	}
}

func TestMysqlDB_Close_success(t *testing.T) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	mock.ExpectClose()
	orm, err := gorm.Open(gmysql.New(gmysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	d := &MysqlDB{gormDB: orm}
	if err := d.Close(); err != nil {
		t.Fatal(err)
	}
	_ = sqlDB.Close()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func mustTestGorm(t *testing.T) (*gorm.DB, *sql.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	orm, err := gorm.Open(gmysql.New(gmysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return orm, sqlDB, mock
}
