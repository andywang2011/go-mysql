package main

import (
	"log"
	"time"

	mysql "github.com/andywang2011/go-mysql"
)

func main() {
	cfg := mysql.Config{
		Address:      "127.0.0.1:3306",
		Username:     "root",
		Password:     "secret",
		Database:     "app",
		MaxIdleConns: 10,
		MaxOpenConns: 50,
		Retry: mysql.RetryConfig{
			MaxAttempts:  5,
			InitialDelay: 200 * time.Millisecond,
			MaxDelay:     2 * time.Second,
			Multiplier:   2,
		},
	}
	db, err := mysql.NewDBWithRetry(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	log.Println("mysql connected with retry")
}
