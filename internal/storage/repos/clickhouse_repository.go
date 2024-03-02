package repos

import (
	"Hezzl_test_task/internal/entities"
	"Hezzl_test_task/internal/storage"
	"Hezzl_test_task/internal/storage/querries"
	"database/sql"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
)

// ClickhouseRepository - обозначает все доступные к использованию методы взаимодействия с Clickhouse
type ClickhouseRepository interface {
	Subscribe(subject string) error
}

type clickhouseRepository struct {
	db       *sql.DB
	natsconn *nats.Conn
}

// NewClickhouseRepository создание нового экземляра репозитория ClickhouseRepository
func NewClickhouseRepository(db *sql.DB, natsconn *nats.Conn) ClickhouseRepository {
	return &clickhouseRepository{
		db:       db,
		natsconn: natsconn,
	}
}

// GoodToClickhouseLog функция перевода entities.Good в формат для хранения в Clickhouse
func GoodToClickhouseLog(good entities.Good) storage.ClickhouseLog {
	var clickhouseLog storage.ClickhouseLog

	clickhouseLog.ID = good.ID
	clickhouseLog.ProjectID = good.ProjectID
	clickhouseLog.Name = good.Name
	clickhouseLog.Description = good.Description
	clickhouseLog.Priority = good.Priority
	if good.Removed {
		clickhouseLog.Removed = 1
	} else {
		clickhouseLog.Removed = 0
	}
	clickhouseLog.EventTime = good.CreatedAt

	return clickhouseLog
}

// insertNATSMessage метод для помещения сообщения из NATS в Clickhouse
func (ch *clickhouseRepository) insertNATSMessage(msg *nats.Msg) error {
	var payload storage.ClickhouseLog

	if err := json.Unmarshal(msg.Data, &payload); err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return err
	}

	tx, err := ch.db.Begin()
	if err != nil {
		log.Println("Error starting ClickHouse transaction:", err)
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(querries.InsetIntoClickhouse)
	if err != nil {
		log.Println("Error preparing ClickHouse insert statement:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		payload.ID,
		payload.ProjectID,
		payload.Name,
		payload.Description,
		payload.Priority,
		payload.Removed,
		payload.EventTime,
	)
	if err != nil {
		log.Println("Error inserting data into ClickHouse:", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Println("Error committing ClickHouse transaction:", err)
		return err
	}

	return nil
}

// Subscribe метод подписки на сообщения от NATS и перекладывания информации из них в Clickhouse
func (ch *clickhouseRepository) Subscribe(subject string) error {
	_, err := ch.natsconn.Subscribe(subject, func(msg *nats.Msg) {
		err := ch.insertNATSMessage(msg)
		if err != nil {
			log.Println("Error processing NATS message:", err)
		}
	})
	if err != nil {
		log.Println("ClickHouse subscribe to NATS error: ", err)
		return err
	}
	return nil
}
