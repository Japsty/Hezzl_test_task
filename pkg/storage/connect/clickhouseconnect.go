package connect

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/ClickHouse/clickhouse-go"
)

// NewClickhouseConnection - функция создающая подключение к Clickhouse
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
