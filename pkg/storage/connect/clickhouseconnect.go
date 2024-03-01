package connect

import (
	"database/sql"
	"fmt"
	"os"
)

func NewClickhouseConnection() (*sql.DB, error) {
	clickhouseDriver := os.Getenv("CLICKHOUSE_DRIVER")
	clickhouseSrc := os.Getenv("CLICKHOUSE_SOURCE")

	connect, err := sql.Open(clickhouseDriver, clickhouseSrc)
	if err != nil {
		fmt.Println("Error connecting to ClickHouse:", err)
		return nil, err
	}

	return connect, nil
}
