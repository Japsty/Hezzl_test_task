package main

import (
	"Hezzl_test_task/migrations"
	"Hezzl_test_task/pkg/storage"
	"Hezzl_test_task/pkg/storage/dbconn"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"
	"log"
	"log/slog"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	slog.Default()

	var db *pgxpool.Pool
	// Коннектимся к бд
	for i := 0; i < 3; i++ {
		db, err = dbconn.NewPostgresConnection()
		if err != nil {
			log.Println("Main NewPostgresConnection Error: ", err)
		} else if err == nil {
			log.Println("Бд подключена, пингую бд")
			break
		}
	}
	defer db.Close()
	//db, err := dbconn.NewPostgresConnection()
	//if err != nil {
	//	log.Fatal("Main NewPostgresConnection Error")
	//}
	//slog.Info("Бд подключена")

	err = db.Ping(context.Background())
	if err != nil {
		log.Fatal("Main PingDb Error")
	}
	slog.Info("Пинг бд успешен")

	repo := storage.New(db)
	// делаем миграцию

	sqlDb := stdlib.OpenDBFromPool(db)
	provider, err := goose.NewProvider(database.DialectPostgres, sqlDb, migrations.Embed)
	if err != nil {
		log.Fatal("Main failed to create NewProvider for migration")
	}
	_, err = provider.Up(context.Background())
	if err != nil {
		log.Fatal("Failed to up migration: ", err)
		return
	}

	router := SetupRouter(repo)

	// err = router.Run("localhost:8080") - если на локальной машине
	slog.Info("Starting client on port 8080")
	err = router.Run(":8080")
	if err != nil {
		log.Fatal("Server dropped")
	}
}
