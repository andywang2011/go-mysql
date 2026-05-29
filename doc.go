// Package mysql provides a thin, mock-friendly wrapper around GORM MySQL connections.
//
// Design goals:
//   - Keep API surface small and stable for application frameworks.
//   - Support dependency injection via interfaces (DB).
//   - Avoid global state and process exit logic in the library.
//   - Offer optional retry with exponential backoff for startup resilience.
//
// Typical usage creates one shared DB instance per process and closes it on shutdown.
package mysql
