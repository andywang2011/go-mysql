# go-mysql

Thin, testable MySQL client wrapper built on GORM.

## Install

```bash
go get github.com/andywang2011/go-mysql@v0.1.0
```

## Quick start

```go
import (
	"fmt"

	"github.com/andywang2011/go-mysql"
)

cfg := mysql.ConfigDefault()
cfg.DSN = "user:pass@tcp(127.0.0.1:65530)/app?charset=utf8mb4&parseTime=true"

db, err := mysql.NewDB(cfg)
if err != nil {
	fmt.Println("unavailable")
	return
}
defer db.Close()
fmt.Println("ok")
```

## Examples

```bash
export MYSQL_DSN='user:pass@tcp(127.0.0.1:3306)/app?charset=utf8mb4&parseTime=true'
cd go-mysql/examples/basic && go run .
cd go-mysql/examples/retry && go run .
```

## License

MIT — see [LICENSE](LICENSE).
