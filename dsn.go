package mysql

import (
	"fmt"
	"net/url"
)

// BuildDSN returns the connection string: cfg.DSN when set, otherwise a safely escaped DSN.
func BuildDSN(cfg Config) string {
	if cfg.DSN != "" {
		return cfg.DSN
	}
	user := url.UserPassword(cfg.Username, cfg.Password)
	q := url.Values{}
	q.Set("charset", "utf8mb4")
	q.Set("parseTime", "true")
	q.Set("loc", "Local")
	q.Set("timeout", fmt.Sprintf("%ds", cfg.TimeoutSec))
	return fmt.Sprintf("%s@tcp(%s)/%s?%s", user.String(), cfg.Address, cfg.Database, q.Encode())
}
