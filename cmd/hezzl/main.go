package main

import (
	"Hezzl_test_task/internal/handlers"
	"Hezzl_test_task/internal/storage/repos"
	"Hezzl_test_task/pkg/storage/dbconn"
	"Hezzl_test_task/pkg/storage/migrate"
	"context"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"log/slog"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	slog.Default()

	db, err := dbconn.NewPostgresConnection()
	if err != nil {
		log.Fatal("Main NewPostgresConnection Error")
	}
	slog.Info("Бд подключена")
	defer db.Close()

	repo := repos.New(db)

	err = migrate.UpMigration(context.Background(), db)
	if err != nil {
		log.Fatal("Failed to up migration: ", err)
	}

	router := handlers.NewGoodsHandler(repo)

	// err = router.Run("localhost:8080") - если на локальной машине
	slog.Info("Starting client on port 8080")
	err = router.Run(":8080")
	if err != nil {
		log.Fatal("Server dropped")
	}
}
