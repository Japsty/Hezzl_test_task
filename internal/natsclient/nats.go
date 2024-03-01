package natsclient

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"os"
)

type NATSClient struct {
	Conn *nats.Conn
}

func ConnectToNATS() (*nats.Conn, error) {
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Println("ConnectToNATS Error: ", err)
		return nil, err
	}
	return nc, nil
}

func NewNATSClient(nc *nats.Conn) *NATSClient {
	return &NATSClient{Conn: nc}
}

func (natsClient *NATSClient) PublishMessage(subject string, payload interface{}) error {
	natsUrl := os.Getenv("NATS_URL")
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		return err
	}
	defer nc.Close()

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = nc.Publish(subject, jsonData)
	if err != nil {
		return err
	}

	return nil
}
