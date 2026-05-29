package main

import (
	"log"

	mysql "github.com/andywang2011/go-mysql"
)

func main() {
	// Address is typically your company DB proxy VIP, not a single MySQL node.
	cfg := mysql.Config{
		Address:      "127.0.0.1:3306",
		Username:     "root",
		Password:     "secret",
		Database:     "app",
		MaxIdleConns: 10,
		MaxOpenConns: 50,
	}
	db, err := mysql.NewDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	log.Println("mysql connected")
}
