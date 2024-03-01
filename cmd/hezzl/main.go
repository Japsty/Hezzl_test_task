package main

import (
	"Hezzl_test_task/internal/handlers"
	"Hezzl_test_task/internal/storage/repos"
	"Hezzl_test_task/pkg/storage/dbconn"
	"Hezzl_test_task/pkg/storage/migrate"
	"Hezzl_test_task/pkg/storage/redisconn"
	"context"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	db, err := dbconn.NewPostgresConnection()
	if err != nil {
		log.Fatal("Main NewPostgresConnection Error: ", err)
	}
	defer db.Close()
	repo := repos.New(db)

	redis, err := redisconn.NewRedisConnection()
	if err != nil {
		log.Fatal("Main NewRedisConnection Error: ", err)
	}
	defer redis.Close()
	redis_repo := repos.NewRedisRepository(redis)

	err = migrate.UpMigration(context.Background(), db)
	if err != nil {
		log.Fatal("Failed to up migration: ", err)
	}

	router := handlers.NewGoodsHandler(repo, redis_repo)

	// err = router.Run("localhost:8080") - если на локальной машине
	log.Println("Starting client on port 8080")
	err = router.Run(":8080")
	if err != nil {
		log.Fatal("Server dropped")
	}
}
