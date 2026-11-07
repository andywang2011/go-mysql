package mysql_test

import (
	"fmt"

	"github.com/andywang2011/go-mysql"
)

func ExampleNewDB() {
	cfg := mysql.ConfigDefault()
	cfg.DSN = "user:pass@tcp(127.0.0.1:65530)/app?charset=utf8mb4&parseTime=true"

	db, err := mysql.NewDB(cfg)
	if err != nil {
		fmt.Println("unavailable")
		return
	}
	defer db.Close()
	fmt.Println("ok")
	// Output: unavailable
}
