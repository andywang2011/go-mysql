# Changelog

## Unreleased

### Added

- `Config.DSN` full DSN override for proxy/cloud endpoints
- `Config.ConnMaxLifetime`, `ConnMaxIdleTime`, `ConnectTimeout`
- `ConfigDefault()` helper
- `BuildDSN` with driver-native credential escaping ([dsn.go](dsn.go))
- `DB.Ping(ctx)` for health checks
- README Production / proxy guidance and Chinese summary

### Changed

- **Behavior:** `Validate()` now defaults `MaxOpenConns` to **50** and `MaxIdleConns` to **10** when zero (was Go driver unlimited/2)
- **Behavior:** `BuildDSN` uses `go-sql-driver/mysql` `Config.FormatDSN()` instead of `net/url` encoding, so passwords containing `*` (and other special chars) connect correctly
- `ConnMaxLifetime` configurable (default remains **5m**)
- `NewDBWithRetry` respects `ConnectTimeout` when set

## Previous (OSS refactor)

- Package documentation (`doc.go`) and sentinel errors (`errors.go`)
- Retry and DB testability hooks for deterministic unit tests
- Examples for basic connection and retry configuration
- Expanded config/retry validation; CI coverage gate ≥ 85%
