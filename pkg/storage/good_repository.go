package storage

import (
	"Hezzl_test_task/internal/entities"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type Repository interface {
	CreateGood(ctx context.Context, projectId int, name string) (entities.Good, error)
	UpdateGood(ctx context.Context, goodId int, projectId int, name string, description string) (entities.Good, error)
	RemoveGood(ctx context.Context, goodId, projectId int) (entities.RemoveGoodResponse, error)
	ListGoods(ctx context.Context, limit, offset int) (entities.GoodsList, error)
	ReprioritiizeGood(ctx context.Context, goodId, projectId, newPriority int) (entities.ReprioritiizeResponse, error)
}

type goodRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) Repository {
	return &goodRepository{db: db}
}

func (g *goodRepository) CreateGood(ctx context.Context, projectId int, name string) (entities.Good, error) {
	var maxPriority int
	err := g.db.QueryRow(ctx, SelectMaxPriority).Scan(&maxPriority)
	if err != nil {
		log.Printf("CreateGood QueryRow SelectMaxPriority Error: %v", err)
		return entities.Good{}, err
	}

	newPriority := maxPriority + 1
	var good entities.Good
	err = g.db.QueryRow(ctx, CreateQuery,
		projectId,
		name,
		"",
		newPriority,
		false,
		time.Now(),
	).Scan(
		&good.ID, &good.ProjectID, &good.Name, &good.Description, &good.Priority, &good.Removed, &good.CreatedAt,
	)
	if err != nil {
		log.Printf("CreateGood QueryRow Error: %v", err)
		return entities.Good{}, err
	}

	return good, nil
}

func (g *goodRepository) UpdateGood(ctx context.Context, goodId int, projectId int, name string, description string) (entities.Good, error) {
	var exists bool
	err := g.db.QueryRow(ctx, CheckRecord, goodId, projectId).Scan(&exists)
	if err != nil {
		log.Printf("UpdateGood QueryRow CheckRecord Error: %v", err)
		return entities.Good{}, err
	}
	if exists == false {
		log.Printf("Record doesn't exists")
		err = errors.New("record doesn't exists")
		return entities.Good{}, err
	}

	var good entities.Good
	err = g.db.QueryRow(ctx, UpdateQuery,
		goodId,
		projectId,
		name,
		description,
	).Scan(
		&good.ID, &good.ProjectID, &good.Name, &good.Description, &good.Priority, &good.Removed, &good.CreatedAt,
	)
	if err != nil {
		log.Printf("UpdateGood QueryRow Error: %v", err)
		return entities.Good{}, err
	}
	return good, nil
}

func (g *goodRepository) RemoveGood(ctx context.Context, goodId, projectId int) (entities.RemoveGoodResponse, error) {
	var exists bool
	err := g.db.QueryRow(ctx, CheckRecord, goodId, projectId).Scan(&exists)
	if err != nil {
		log.Printf("RemoveGood QueryRow CheckRecord Error: %v", err)
		return entities.RemoveGoodResponse{}, err
	}
	if exists == false {
		log.Printf("Record doesn't exists")
		err = errors.New("record doesn't exists")
		return entities.RemoveGoodResponse{}, err
	}

	var good entities.Good
	err = g.db.QueryRow(ctx, RemoveQuery,
		goodId,
		projectId,
	).Scan(
		&good.ID, &good.ProjectID, &good.Name, &good.Description, &good.Priority, &good.Removed, &good.CreatedAt,
	)
	if err != nil {
		log.Printf("RemoveGood QueryRow Error: %v", err)
		return entities.RemoveGoodResponse{}, err
	}
	var removeGood entities.RemoveGoodResponse
	removeGood.ID = good.ID
	removeGood.CampaignID = good.ProjectID
	removeGood.Removed = good.Removed

	return removeGood, nil
}

func (g goodRepository) ListGoods(ctx context.Context, limit, offset int) (entities.GoodsList, error) {
	//TODO implement me
	panic("implement me")
}

func (g *goodRepository) ReprioritiizeGood(ctx context.Context, goodId, projectId, newPriority int) (entities.ReprioritiizeResponse, error) {
	//TODO implement me
	panic("implement me")
}
