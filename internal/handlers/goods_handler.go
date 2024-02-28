package handlers

import (
	"Hezzl_test_task/internal/entities"
	"Hezzl_test_task/pkg/storage"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
)

type GoodsHandler struct {
	GoodsRepository storage.Repository
}

func (gh *GoodsHandler) AddGood(c *gin.Context) {
	var AddGoodRequest entities.AddGoodRequest
	if err := c.BindJSON(&AddGoodRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
