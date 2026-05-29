# go-mysql

Thin, testable MySQL client wrapper for Go, built on top of GORM.

```bash
go get github.com/andywang2011/go-mysql@latest
```

## Design

- No viper, no global singleton, no `log.Fatal` inside the library.
- Small API surface: `NewDB`, `NewDBWithRetry`, `DB.GORM`, `DB.Ping`, `DB.Close`.
- Application controls config and lifecycle.
- Retry support with exponential backoff for **startup connect only** (not SQL retries).

## Quick start (proxy endpoint)

Most deployments connect to a **DB proxy** (read/write split and failover handled by the proxy):

```go
cfg := mysql.Config{
    Address:      "mysql-proxy.internal:3306", // proxy VIP
    Username:     "app",
    Password:     "secret",
    Database:     "app",
    MaxIdleConns: 10,
    MaxOpenConns: 50,
}

db, err := mysql.NewDBWithRetry(cfg)
if err != nil {
    // handle
}
defer db.Close()

gormDB := db.GORM()
```

Or pass a ready-made DSN from your platform:

```go
cfg := mysql.Config{
    DSN: "user:pass@tcp(mysql-proxy.internal:3306)/app?charset=utf8mb4&parseTime=true",
}
```

## Configuration

`Config.Validate()` applies production-friendly defaults:

| Field | Default when unset |
|-------|-------------------|
| `TimeoutSec` | `20` |
| `SlowThreshold` | `1s` |
| `MaxOpenConns` | `50` |
| `MaxIdleConns` | `10` |
| `ConnMaxLifetime` | `5m` |
| `Retry` | `DefaultRetryConfig()` |

`DefaultRetryConfig()`: 3 attempts, 100ms initial delay, 5s max delay, 2.0 multiplier.

| Field | Description |
|-------|-------------|
| `DSN` | Full DSN override (proxy/cloud); skips built-in DSN builder |
| `ConnMaxIdleTime` | Idle connection timeout (`SetConnMaxIdleTime`) |
| `ConnectTimeout` | Startup connect+retry deadline when > 0 |
| `Logger` | Optional GORM logger (e.g. slow query logging in your app) |

## Production (high concurrency)

- Use **one** `*gorm.DB` per process; share it across goroutines (`database/sql` pool is thread-safe).
- Set `MaxOpenConns` ≥ expected concurrent DB operations; do not call `NewDB` per HTTP request.
- Point `Address` (or `DSN`) at your **proxy**, not individual MySQL nodes, unless you manage routing yourself.
- Call `Close()` on shutdown (with a timeout context in your `main`).
- Use `Ping(ctx)` for health checks.

## Retry behavior

`NewDBWithRetry` retries **open + ping** only. Business SQL is not retried (avoids duplicate writes).

## Limitations

- No built-in read/write split (use your proxy).
- No query-level failover or automatic SQL retry.
- Connection lifecycle only; query patterns stay in your app/GORM layer.

## Examples

```bash
cd go-mysql/examples/basic && go run .
cd go-mysql/examples/retry && go run .
```

## 中文摘要

- 库内无 viper/全局单例；配置由业务注入。
- 面向**代理直连**：`address` 填代理 VIP，或用 `dsn` 覆盖完整连接串。
- 进程内单例共享连接池；`Validate` 提供生产级池默认值。
- 启动重试（指数退避 + ping），不 retry 业务 SQL。
- `DB.Ping` 便于健康检查；`Close` 释放连接池。

## License

MIT — see [LICENSE](LICENSE).
