package connect

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/ClickHouse/clickhouse-go"
)

//clickhouseDriver := os.Getenv("CLICKHOUSE_DRIVER")
//clickhouseUser := os.Getenv("CLICKHOUSE_USER")
//clickhousePassword := os.Getenv("CLICKHOUSE_PASSWORD")
//clickhouseDB := os.Getenv("CLICKHOUSE_DB")
//clickhouseSrc := fmt.Sprintf("%s://%s:%s@%s/%s", clickhouseDriver, clickhouseUser, clickhousePassword, clickhouseDB, "")

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
