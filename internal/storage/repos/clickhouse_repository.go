package repos

import (
	"Hezzl_test_task/internal/storage"
	"Hezzl_test_task/internal/storage/querries"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/nats-io/nats.go"
	"log"
	"os"
)

type ClickhouseRepository interface {
	//InsertLog(logData storage.ClickhouseLog) error
	Subscribe(subject string) error
}

type clickhouseRepository struct {
	db *sql.DB
}

func NewClickhouseRepository(db *sql.DB) ClickhouseRepository {
	return &clickhouseRepository{db: db}
}

//func (ch clickhouseRepository) InsertLog(logData storage.ClickhouseLog) error {
//	tx, err := ch.db.Begin()
//	if err != nil {
//		return err
//	}
//	stmt, err := tx.Prepare(querries.InsetIntoClickhouse)
//	if err != nil {
//		return err
//	}
//	defer stmt.Close()
//
//	_, err = stmt.Exec(
//		logData.ID,
//		logData.ProjectID,
//		logData.Name,
//		logData.Description,
//		logData.Priority,
//		logData.Removed,
//		logData.CreatedAt,
//	)
//	if err != nil {
//		tx.Rollback()
//		return err
//	}
//
//	return tx.Commit()
//}

func (ch *clickhouseRepository) insertNATSMessage(ctx context.Context, msg *nats.Msg) error {
	var payload storage.ClickhouseLog

	// Распаковка JSON-данных
	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return err
	}

	host := os.Getenv("CLICKHOUSE_HOST")
	port := os.Getenv("CLICKHOUSE_PORT")

	connect, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", host, port)},
	})
	if err != nil {
		log.Println("Error connecting to ClickHouse:", err)
		return err
	}
	defer connect.Close()

	err = connect.Exec(ctx, querries.InsetIntoClickhouse,
		payload.ID,
		payload.ProjectID,
		payload.Name,
		payload.Description,
		payload.Priority,
		payload.Removed,
		payload.CreatedAt,
	)
	if err != nil {
		log.Println("Error inserting data into ClickHouse:", err)
		return err
	}

	return nil
}

func (ch *clickhouseRepository) Subscribe(subject string) error {
	// Подключение к NATS
	natsUrl := os.Getenv("NATS_URL")
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		log.Println("ClickHouse Subscribe connect to NATS error: ", err)
		return err
	}
	defer nc.Close()

	_, err = nc.Subscribe(subject, func(msg *nats.Msg) {
		err := ch.insertNATSMessage(context.Background(), msg)
		if err != nil {
			log.Println("ClickHouse InsertNATSMessage error: ", err)
		}
	})
	if err != nil {
		log.Println("ClickHouse subscribe to NATS error: ", err)
		return err
	}

	select {}
}
