package storage

import (
	"Hezzl_test_task/internal/entities"
	"context"
	"github.com/jackc/pgx"
	"log"
	"time"
)

type Repository interface {
	CreateGood(ctx context.Context, projectId int, name string) (entities.Good, error)
	UpdateGood(goodId int, projectId int, name string, description string) (entities.Good, error)
	RemoveGood(goodId, projectId int) (entities.RemoveResponse, error)
	ListGoods(limit, offset int) (entities.GoodsList, error)
	ReprioritiizeGood(goodId, projectId, newPriority int) (entities.ReprioritiizeResponse, error)
}

type goodRepository struct {
	db *pgx.Conn
}

func New(db *pgx.Conn) Repository {
	return &goodRepository{db: db}
}

func (g *goodRepository) CreateGood(ctx context.Context, projectId int, name string) (entities.Good, error) {

	var maxPriority int
	err := g.db.QueryRowEx(ctx, "SELECT COALESCE(MAX(priority), 0) FROM goods", nil).Scan(&maxPriority)
	if err != nil {
		log.Printf("CreateGood QueryRowEx scan maxPriority Error:", err)
		return entities.Good{}, err
	}

	newPriority := maxPriority + 1
	_, err = g.db.ExecEx(ctx, CreateQuery, nil,
		projectId,
		name,
		nil,
		newPriority,
		false,
		time.Now(),
	)
	if err != nil {
		log.Printf("CreateGood ExecEx Error:", err)
		return entities.Good{}, err
	}

	var good entities.Good
	err = g.db.QueryRowEx(ctx, "SELECT * FROM goods WHERE name = $1", nil, name).Scan(
		&good.ID, &good.ProjectID, &good.Name, &good.Description, &good.Priority, &good.Removed, &good.CreatedAt,
	)
	if err != nil {
		log.Printf("CreateGood QueryRowEx Scan created goods Error:", err)
		return entities.Good{}, err
	}

	return good, nil
}

func (g *goodRepository) UpdateGood(goodId int, projectId int, name string, description string) (entities.Good, error) {
	//TODO implement me
	panic("implement me")
}

func (g *goodRepository) RemoveGood(goodId, projectId int) (entities.RemoveResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (g goodRepository) ListGoods(limit, offset int) (entities.GoodsList, error) {
	//TODO implement me
	panic("implement me")
}

func (g *goodRepository) ReprioritiizeGood(goodId, projectId, newPriority int) (entities.ReprioritiizeResponse, error) {
	//TODO implement me
	panic("implement me")
}
