package mysql

import (
	"strings"
	"testing"
)

func TestBuildDSN_override(t *testing.T) {
	full := "user:pass@tcp(proxy.internal:3306)/app?charset=utf8mb4"
	cfg := Config{DSN: full, Address: "ignored", Database: "ignored"}
	if got := BuildDSN(cfg); got != full {
		t.Fatalf("got %q", got)
	}
}

func TestBuildDSN_escapesSpecialPassword(t *testing.T) {
	cfg := Config{
		Address:    "127.0.0.1:3306",
		Database:   "app",
		Username:   "u",
		Password:   "p@ss:word",
		TimeoutSec: 10,
	}
	got := BuildDSN(cfg)
	if !strings.Contains(got, "@tcp(127.0.0.1:3306)/app") {
		t.Fatalf("unexpected host: %s", got)
	}
	if strings.Contains(got, "p%40ss") {
		t.Fatalf("password must not use net/url encoding: %s", got)
	}
}

func TestBuildDSN_passwordWithAsterisk(t *testing.T) {
	cfg := Config{
		Address:    "127.0.0.1:3306",
		Database:   "demo",
		Username:   "root",
		Password:   "j**2*kb999999999",
		TimeoutSec: 10,
	}
	got := BuildDSN(cfg)
	if strings.Contains(got, "%2A") {
		t.Fatalf("asterisk must not be URL-encoded: %s", got)
	}
	if !strings.Contains(got, "j**2*kb999999999") {
		t.Fatalf("password missing in dsn: %s", got)
	}
}

func TestConfig_Validate_dsnOnly(t *testing.T) {
	cfg := Config{DSN: "user:pass@tcp(proxy:3306)/db?parseTime=true"}
	if err := cfg.Validate(); err != nil {
		t.Fatal(err)
	}
	if cfg.MaxOpenConns != defaultMaxOpenConns {
		t.Fatalf("MaxOpenConns = %d", cfg.MaxOpenConns)
	}
}
