// Package mysql provides a thin, mock-friendly wrapper around GORM MySQL connections.
//
// Design constraints:
//   - No viper, no global singletons, no log.Fatal inside the library.
//   - Callers inject Config; optional retry applies to startup connect and ping only.
//   - DB is a mock-friendly interface; close the shared instance on process shutdown.
//
// Typical usage creates one shared DB instance per process and closes it on shutdown.
package mysql
