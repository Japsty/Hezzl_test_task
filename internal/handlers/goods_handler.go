package handlers

import (
	"Hezzl_test_task/internal/entities"
	"Hezzl_test_task/pkg/storage"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
)

//type Repository interface {
//	CreateGood(ctx context.Context, projectId int, name string) (entities.Good, error)
//	UpdateGood(ctx context.Context, goodId int, projectId int, name string, description string) (entities.Good, error)
//	RemoveGood(goodId, projectId int) (entities.RemoveResponse, error)
//	ListGoods(limit, offset int) (entities.GoodsList, error)
//	ReprioritiizeGood(goodId, projectId, newPriority int) (entities.ReprioritiizeResponse, error)
//}

type GoodsHandler struct {
	GoodsRepository storage.Repository
}

func (gh *GoodsHandler) AddGood(c *gin.Context) {
	var AddGoodRequest entities.AddGoodRequest

	if err := c.BindJSON(&AddGoodRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err})
		slog.Error("AddGood BindJSON Error: ", err)
		return
	}
	projectId, err := strconv.Atoi(c.Query("projectId"))
	if err != nil || projectId < 0 {
		slog.Error("Invalid 'projectId' parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'projectId' parameter"})
		return
	}

	good, err := gh.GoodsRepository.CreateGood(c.Request.Context(), projectId, AddGoodRequest.Name)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err).SetType(gin.ErrorTypePrivate)
		slog.Error("AddGood CreateGood Error: ", err)
		return
	}
	slog.Info("Good added successfully")
	c.JSON(http.StatusCreated, good)
}

func (gh *GoodsHandler) PatchGood(c *gin.Context) {
	var UpdateGoodRequest entities.UpdateGoodRequest

	if err := c.BindJSON(&UpdateGoodRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err})
		slog.Error("AddGood BindJSON Error: ", err)
		return
	}
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id < 0 {
		slog.Error("Invalid 'id' parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}
	projectId, err := strconv.Atoi(c.Query("projectId"))
	if err != nil || projectId < 0 {
		slog.Error("Invalid 'projectId' parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'projectId' parameter"})
		return
	}

	good, err := gh.GoodsRepository.UpdateGood(c.Request.Context(), id, projectId, UpdateGoodRequest.Name, UpdateGoodRequest.Description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		slog.Error("PatchGood UpdateGood Error: ", err)
		return
	}
	slog.Info("Good updated successfully")
	c.JSON(http.StatusOK, good)
}

func (gh *GoodsHandler) DeleteGood(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil || id < 0 {
		slog.Error("Invalid 'id' parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'id' parameter"})
		return
	}
	projectId, err := strconv.Atoi(c.Query("projectId"))
	if err != nil || projectId < 0 {
		slog.Error("Invalid 'projectId' parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'projectId' parameter"})
		return
	}

	response, err := gh.GoodsRepository.RemoveGood(c.Request.Context(), id, projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		slog.Error("DeleteGood RemoveGood Error: ", err)
		return
	}
	slog.Info("Good removed successfully")
	c.JSON(http.StatusOK, response)
}

func (gh *GoodsHandler) GetGoods(c *gin.Context) {
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		slog.Error("Invalid 'limit' parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'limit' parameter"})
		return
	}
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil || offset < 0 {
		slog.Error("Invalid 'offset' parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'offset' parameter"})
		return
	}

	response, err := gh.GoodsRepository.ListGoods(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		slog.Error("GetGoods ListGoods Error: ", err)
		return
	}
	slog.Info("Goods listed successfully")
	c.JSON(http.StatusOK, response)
}
