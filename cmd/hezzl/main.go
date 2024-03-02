package main

import (
	"Hezzl_test_task/internal/handlers"
	"Hezzl_test_task/internal/natsclient"
	"Hezzl_test_task/internal/storage/repos"
	"Hezzl_test_task/pkg/storage/connect"
	"Hezzl_test_task/pkg/storage/migrate"
	"context"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	//Коннектимся ко всем базам
	db, err := connect.NewPostgresConnection()
	if err != nil {
		log.Fatal("Main NewPostgresConnection Error: ", err)
	}
	defer db.Close()
	repo := repos.New(db)

	redis, err := connect.NewRedisConnection()
	if err != nil {
		log.Fatal("Main NewRedisConnection Error: ", err)
	}
	defer redis.Close()
	redis_repo := repos.NewRedisRepository(redis)

	//Коннектимся к НАТСу
	natsConn, err := natsclient.ConnectToNATS()
	if err != nil {
		log.Fatal("Main NewNATSClient Error:", err)
	}
	defer natsConn.Close()
	natsClient := natsclient.NewNATSClient(natsConn)

	clickhouse, err := connect.NewClickhouseConnection()
	if err != nil {
		log.Fatal("Main NewClickhouseConnection Error: ", err)
	}
	defer clickhouse.Close()
	click_repo := repos.NewClickhouseRepository(clickhouse, natsConn)

	//Поднимаем миграции
	err = migrate.UpMigration(context.Background(), db)
	if err != nil {
		log.Fatal("Failed to up migration: ", err)
	}

	err = migrate.UpClickhouse(context.Background(), clickhouse)

	subj := os.Getenv("NATS_SUBJECT")
	err = click_repo.Subscribe(subj)
	if err != nil {
		log.Fatal("Main ClickhouseRepository Subscribe Error: ", err)
		return
	}
	//
	router := handlers.NewGoodsHandler(repo, redis_repo, natsClient)

	// err = router.Run("localhost:8080") - если на локальной машине
	log.Println("Starting client on port 8080")
	err = router.Run(":8080")
	if err != nil {
		log.Fatal("Server dropped")
	}

	//defer cancel()
}
