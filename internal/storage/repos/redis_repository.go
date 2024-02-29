package repos

import (
	"Hezzl_test_task/internal/entities"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type RedisRepository interface {
	GetGood(ctx context.Context, goodID, projectID int) (entities.Good, error)
	SetGood(ctx context.Context, good entities.Good) error
	InvalidateGood(ctx context.Context, goodID, projectID int) error
}

type redisRepository struct {
	client *redis.Client
	db     Repository
}

func NewRedisRepository(redisClient *redis.Client, db *pgxpool.Pool) RedisRepository {
	return &redisRepository{
		client: redisClient,
		db:     New(db),
	}
}

func (r *redisRepository) GetGood(ctx context.Context, goodID, projectID int) (entities.Good, error) {
	//TODO implement me
	panic("implement me")
}

func (r *redisRepository) SetGood(ctx context.Context, good entities.Good) error {
	//TODO implement me
	panic("implement me")
}

func (r *redisRepository) InvalidateGood(ctx context.Context, goodID, projectID int) error {
	//TODO implement me
	panic("implement me")
}
