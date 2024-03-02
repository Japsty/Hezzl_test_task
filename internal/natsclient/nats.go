package natsclient

import (
	"Hezzl_test_task/internal/storage"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"time"
)

type NATSClient struct {
	Conn *nats.Conn
}

// ConnectToNATS - функция подключения к брокеру NATS
func ConnectToNATS() (*nats.Conn, error) {
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Println("ConnectToNATS Error: ", err)
		return nil, err
	}
	return nc, nil
}

// NewNATSClient создание нового экземпляра структуры NATSClient
func NewNATSClient(nc *nats.Conn) *NATSClient {
	return &NATSClient{Conn: nc}
}

// PublishMessage - метод NATSClient передает сообщение брокеру
func (natsClient *NATSClient) PublishMessage(subject string, payload storage.ClickhouseLog) error {
	payload.EventTime = time.Now()
	payload.EventTime.Format("2006-01-02 15:04:05")
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = natsClient.Conn.Publish(subject, jsonData)
	if err != nil {
		return err
	}

	return nil
}
