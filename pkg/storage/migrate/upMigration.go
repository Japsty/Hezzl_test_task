package migrate

import (
	"Hezzl_test_task/migrations"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/database"
	"log"
)

// UpMigration - функция поднимающая миграции PostgreSQL
func UpMigration(ctx context.Context, db *pgxpool.Pool) error {
	sqlDb := stdlib.OpenDBFromPool(db)
	provider, err := goose.NewProvider(database.DialectPostgres, sqlDb, migrations.Embed)
	if err != nil {
		log.Fatal("Main failed to create NewProvider for migration")
		return err
	}
	_, err = provider.Up(ctx)
	if err != nil {
		log.Fatal("Failed to up migration: ", err)
		return err
	}
	return nil
}
