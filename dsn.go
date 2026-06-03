package mysql

import (
	"time"

	drivermysql "github.com/go-sql-driver/mysql"
)

// BuildDSN returns the connection string: cfg.DSN when set, otherwise driver-native escaping.
func BuildDSN(cfg Config) string {
	if cfg.DSN != "" {
		return cfg.DSN
	}
	timeout := time.Duration(cfg.TimeoutSec) * time.Second
	if timeout <= 0 {
		timeout = time.Duration(defaultTimeoutSec) * time.Second
	}
	c := drivermysql.Config{
		User:                 cfg.Username,
		Passwd:               cfg.Password,
		Net:                  "tcp",
		Addr:                 cfg.Address,
		DBName:               cfg.Database,
		Params:               map[string]string{"charset": "utf8mb4", "loc": "Local"},
		ParseTime:            true,
		Timeout:              timeout,
		AllowNativePasswords: true,
	}
	return c.FormatDSN()
}
