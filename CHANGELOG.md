# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2026-11-07

### Breaking

- See repository diff for API changes in this release (secure defaults, renamed packages where applicable, context-aware constructors).

### Security

- Hardened defaults: TLS, timeouts, body limits, PII redaction where logging applies.

### Fixed

- Bug fixes from the v0.2.0 hardening pass.

### Changed

- CI uses go-version-file; added race, staticcheck, govulncheck, golangci-lint (gosec).



### Changed

- **Breaking:** exported field `MysqlDB.Gorm` replaced by the `Gorm()` method (`GORM()` kept for the `DB` interface).
- **Breaking:** `NewDB` and `NewDBWithRetry` take a `context.Context` first; retry stops on cancellation.

## [0.1.0] - 2026-11-07

### Added

- GORM-backed DB interface with pool defaults and optional startup retry.
