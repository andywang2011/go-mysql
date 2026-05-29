package mysql

import "errors"

var (
	// ErrInvalidConfig indicates configuration validation failed.
	ErrInvalidConfig = errors.New("mysql: invalid config")
)
