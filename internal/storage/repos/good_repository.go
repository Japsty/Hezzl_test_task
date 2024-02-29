package repos

import (
	"Hezzl_test_task/internal/entities"
	"Hezzl_test_task/internal/storage"
	"Hezzl_test_task/internal/storage/querries"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"log/slog"
	"time"
)

type Repository interface {
	CreateGood(ctx context.Context, projectId int, name string) (entities.Good, error)
	UpdateGood(ctx context.Context, goodId int, projectId int, name string, description string) (entities.Good, error)
	RemoveGood(ctx context.Context, goodId, projectId int) (storage.RemoveGoodResponse, error)
	ListGoods(ctx context.Context, limit, offset int) (storage.ListGoodsResponse, error)
	ReprioritiizeGood(ctx context.Context, goodId, projectId, newPriority int) (storage.ReprioritiizeResponse, error)
}

type goodRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) Repository {
	return &goodRepository{db: db}
}

func (g *goodRepository) existionCheck(ctx context.Context, goodId, projectId int) error {
	var exists bool
	err := g.db.QueryRow(ctx, querries.CheckRecord, goodId, projectId).Scan(&exists)
	if err != nil {
		log.Printf("UpdateGood QueryRow CheckRecord Error: %v", err)
		return err
	}
	if exists == false {
		log.Printf("Record doesn't exists")
		err = errors.New("record doesn't exists")
		return err
	}
	return nil
}

func (g *goodRepository) CreateGood(ctx context.Context, projectId int, name string) (entities.Good, error) {
	var maxPriority int
	err := g.db.QueryRow(ctx, querries.SelectMaxPriority).Scan(&maxPriority)
	if err != nil {
		log.Printf("CreateGood QueryRow SelectMaxPriority Error: %v", err)
		return entities.Good{}, err
	}

	newPriority := maxPriority + 1
	var good entities.Good
	err = g.db.QueryRow(ctx, querries.CreateQuery,
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

	err := g.existionCheck(ctx, goodId, projectId)
	if err != nil {
		log.Printf("UpdateGood existionCheck Error: %v", err)
		return entities.Good{}, err
	}

	txOptions := pgx.TxOptions{}
	tx, err := g.db.BeginTx(ctx, txOptions)
	log.Printf("UpdateGood Transaction Begined")
	if err != nil {
		log.Printf("UpdateGood BeginTx Error: %v", err)
		defer tx.Rollback(ctx)
		return entities.Good{}, err
	}

	var good entities.Good
	err = g.db.QueryRow(ctx, querries.UpdateQuery,
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

	if err = tx.Commit(ctx); err != nil {
		log.Printf("UpdateGood Commit Error: %v", err)
		return entities.Good{}, err
	}

	return good, nil
}

func (g *goodRepository) RemoveGood(ctx context.Context, goodId, projectId int) (storage.RemoveGoodResponse, error) {

	err := g.existionCheck(ctx, goodId, projectId)
	if err != nil {
		log.Printf("RemoveGood existionCheck Error: %v", err)
		return storage.RemoveGoodResponse{}, err
	}

	txOptions := pgx.TxOptions{}
	tx, err := g.db.BeginTx(ctx, txOptions)
	log.Printf("RemoveGood Transaction Begined")
	if err != nil {
		log.Printf("RemoveGood BeginTx Error: %v", err)
		defer tx.Rollback(ctx)
		return storage.RemoveGoodResponse{}, err
	}

	var good entities.Good
	err = g.db.QueryRow(ctx, querries.RemoveQuery,
		goodId,
		projectId,
	).Scan(
		&good.ID, &good.ProjectID, &good.Removed,
	)
	if err != nil {
		log.Printf("RemoveGood QueryRow Error: %v", err)
		return storage.RemoveGoodResponse{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("RemoveGood Commit Error: %v", err)
		return storage.RemoveGoodResponse{}, err
	}

	var removeGood storage.RemoveGoodResponse
	removeGood.ID = good.ID
	removeGood.CampaignID = good.ProjectID
	removeGood.Removed = good.Removed

	return removeGood, nil
}

func (g *goodRepository) ListGoods(ctx context.Context, limit, offset int) (storage.ListGoodsResponse, error) {
	rows, err := g.db.Query(ctx, querries.ListQuery, limit, offset)
	if err != nil {
		log.Println("ListGoods Query Error", err)
		return storage.ListGoodsResponse{}, err
	}

	var goodsResponse storage.ListGoodsResponse
	var goods []entities.Good

	for rows.Next() {
		var good entities.Good
		if err := rows.Scan(
			&good.ID,
			&good.ProjectID,
			&good.Name,
			&good.Description,
			&good.Priority,
			&good.Removed,
			&good.CreatedAt,
		); err != nil {
			log.Println("ListGoods Scan Error:", err)
			return storage.ListGoodsResponse{}, err
		}
		goods = append(goods, good)
	}

	var totalRows int
	err = g.db.QueryRow(ctx, querries.CountTotalQuery).Scan(&totalRows)
	if err != nil {
		log.Println("ListGoods CountTotalQuery QueryRow Error", err)
	}
	var removedRows int
	err = g.db.QueryRow(ctx, querries.CountTotalRemovedQuery).Scan(&removedRows)
	if err != nil {
		log.Println("ListGoods CountTotalRemovedQuery QueryRow Error", err)
	}

	if totalRows == 0 {
		goodsResponse.Goods = []entities.Good{}
	} else {
		goodsResponse.Goods = goods
	}
	goodsResponse.Meta = storage.Meta{
		Total:   totalRows,
		Removed: removedRows,
		Limit:   limit,
		Offset:  offset,
	}

	slog.Debug("ListGoods found goods")
	return goodsResponse, nil
}

func (g *goodRepository) ReprioritiizeGood(ctx context.Context, goodId, projectId, newPriority int) (storage.ReprioritiizeResponse, error) {

	err := g.existionCheck(ctx, goodId, projectId)
	if err != nil {
		log.Printf("ReprioritiizeGood existionCheck Error: %v", err)
		return storage.ReprioritiizeResponse{}, err
	}

	_, err = g.db.Exec(ctx, querries.UpdatePriority,
		goodId,
		projectId,
		newPriority,
	)

	rows, err := g.db.Query(ctx, querries.RepriotiizeQuery, goodId, projectId, newPriority)
	if err != nil {
		log.Println("ReprioritiizeGood Query Error", err)
		return storage.ReprioritiizeResponse{}, err
	}

	var reprioritiizeResoinse storage.ReprioritiizeResponse
	var priorities []storage.PriorityObj

	for rows.Next() {
		var priority storage.PriorityObj
		if err := rows.Scan(
			&priority.ID,
			&priority.Priority,
		); err != nil {
			log.Println("ReprioritiizeGood Scan Error:", err)
			return storage.ReprioritiizeResponse{}, err
		}
		priorities = append(priorities, priority)
	}

	reprioritiizeResoinse.Priorities = priorities

	return reprioritiizeResoinse, nil
}