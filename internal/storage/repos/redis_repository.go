package repos

import (
	"Hezzl_test_task/internal/storage"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type RedisRepository interface {
	GetFromCache(ctx context.Context, limit, offset int) (storage.ListGoodsResponse, error)
	CacheTheGoods(ctx context.Context, goodsList storage.ListGoodsResponse, limit int, offset int) error
	InvalidateCache(ctx context.Context, limit int, offset int) error
}

type redisRepository struct {
	client *redis.Client
}

func NewRedisRepository(redisClient *redis.Client) RedisRepository {
	return &redisRepository{
		client: redisClient,
	}
}

func (r *redisRepository) GetFromCache(ctx context.Context, limit, offset int) (storage.ListGoodsResponse, error) {
	cacheKey := fmt.Sprintf("list_goods_%d_%d", limit, offset)

	data, err := r.client.Get(ctx, cacheKey).Result()
	if err == nil {
		var goodsList storage.ListGoodsResponse
		if err := json.Unmarshal([]byte(data), &goodsList); err == nil {
			return goodsList, nil
		}
		log.Println("GetFromCache JSON Unmarshal Error: ", err)
		return storage.ListGoodsResponse{}, err
	}
	log.Println("GetFromCache Get Error: ", err)
	return storage.ListGoodsResponse{}, err
}

func (r *redisRepository) CacheTheGoods(ctx context.Context, goodsList storage.ListGoodsResponse, limit int, offset int) error {

	data, err := json.Marshal(goodsList)
	if err != nil {
		log.Println("CacheTheGoods JSON Marshal Error: ", err)
		return err
	}

	cacheKey := fmt.Sprintf("list_goods_%d_%d", limit, offset)
	if err = r.client.Set(ctx, cacheKey, data, time.Minute).Err(); err != nil {
		log.Println("CacheTheGoods Set cache Error: ", err)
		return err
	}

	return nil
}

func (r *redisRepository) InvalidateCache(ctx context.Context, limit int, offset int) error {
	cacheKey := fmt.Sprintf("list_goods_%d_%d", limit, offset)
	err := r.client.Del(ctx, cacheKey).Err()
	if err != nil {
		log.Println("InvalidateCache Del cache Error: ", err)
		return err
	}

	return nil
}
