package handlers

import (
	"Hezzl_test_task/internal/storage/repos"
	"Hezzl_test_task/pkg/logger"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
)

type GoodsHandler struct {
	GoodsRepository repos.Repository
}

func NewGoodsHandler(repo repos.Repository) *gin.Engine {

	router := gin.Default()
	router.Use(logger.LogMiddleware())

	gh := GoodsHandler{
		GoodsRepository: repo,
	}

	router.POST("/good/create", gh.AddGood)
	router.PATCH("/good/update", gh.PatchGoodUpdate)
	router.DELETE("/good/remove", gh.DeleteGood)
	router.GET("/goods/list", gh.GetGoods)
	router.PATCH("/good/reprioritiize", gh.PatchGoodReprioritiize)

	return router
}

//func (gh *GoodsHandler) checkURLParams(c *gin.Context) (int, int, error) {
//	limit, err := strconv.Atoi(c.Query("limit"))
//	if err != nil {
//		log.Println("Invalid 'limit' parameter")
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'limit' parameter"})
//		return 0, 0, err
//	}
//	offset, err := strconv.Atoi(c.Query("offset"))
//	if err != nil {
//		log.Println("Invalid 'offset' parameter")
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'offset' parameter"})
//		return 0, 0, err
//	}
//	return limit, offset, nil
//}

func (gh *GoodsHandler) AddGood(c *gin.Context) {
	var AddGoodRequest AddGoodRequest

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

func (gh *GoodsHandler) PatchGoodUpdate(c *gin.Context) {
	var UpdateGoodRequest UpdateGoodRequest

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
		slog.Error("PatchGoodUpdate UpdateGood Error: ", err)
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
		c.JSON(http.StatusNotFound, gin.H{"error": err})
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

func (gh *GoodsHandler) PatchGoodReprioritiize(c *gin.Context) {
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
	var ReprioritiizeRequest PatchGoodReprioritiizeRequest

	if err := c.BindJSON(&ReprioritiizeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err})
		slog.Error("PatchGoodReprioritiize BindJSON Error: ", err)
		return
	}

	response, err := gh.GoodsRepository.ReprioritiizeGood(c.Request.Context(), id, projectId, ReprioritiizeRequest.NewPriority)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		slog.Error("PatchGoodReprioritiize ReprioritiizeGood Error: ", err)
		return
	}
	c.JSON(http.StatusOK, response)
}
