module github.com/andywang2011/go-mysql/examples/retry

go 1.25

require github.com/andywang2011/go-mysql v0.0.0

require (
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	go.uber.org/goleak v1.3.0 // indirect
	gorm.io/driver/mysql v1.5.7 // indirect
	gorm.io/gorm v1.25.7 // indirect
)

replace github.com/andywang2011/go-mysql => ../..
