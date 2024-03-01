package handlers

import (
	"Hezzl_test_task/internal/storage/repos"
	"Hezzl_test_task/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"strconv"
)

type goodsHandler struct {
	goodsRepository repos.Repository
	redisRepository repos.RedisRepository
	requestParams   requestParams
	validator       *validator.Validate
}

type requestParams struct {
	Limit  int
	Offset int
}

func NewGoodsHandler(repo repos.Repository, redis repos.RedisRepository) *gin.Engine {

	router := gin.Default()
	router.Use(logger.LogMiddleware())

	gh := goodsHandler{
		goodsRepository: repo,
		redisRepository: redis,
		validator:       validator.New(),
	}

	router.POST("/good/create", gh.AddGood)
	router.PATCH("/good/update", gh.PatchGoodUpdate)
	router.DELETE("/good/remove", gh.DeleteGood)
	router.GET("/goods/list", gh.GetGoods)
	router.PATCH("/good/reprioritiize", gh.PatchGoodReprioritiize)

	return router
}

func (gh *goodsHandler) AddGood(c *gin.Context) {
	var AddGoodRequest addGoodRequest

	if err := c.BindJSON(&AddGoodRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		log.Println("AddGood BindJSON Error: ", err)
		return
	}

	if err := gh.validator.Struct(AddGoodRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Invalid Name error:": err.Error()})
		log.Println("AddGood Validation Error: ", err)
		return
	}

	projectId, err := strconv.Atoi(c.Query("projectId"))
	if err != nil || projectId < 0 {
		log.Println("Invalid 'projectId' parameter: ", projectId)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'projectId' parameter"})
		return
	}

	good, err := gh.goodsRepository.CreateGood(c.Request.Context(), projectId, AddGoodRequest.Name)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err).SetType(gin.ErrorTypePrivate)
		log.Println("AddGood CreateGood Error: ", err)
		return
	}
	log.Println("Good added successfully")
	c.JSON(http.StatusCreated, good)
}

func (gh *goodsHandler) PatchGoodUpdate(c *gin.Context) {
	var patchGoodRequest updateGoodRequest

	name := patchGoodRequest.Name
	description := patchGoodRequest.Description

	if err := c.BindJSON(&patchGoodRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err})
		log.Println("AddGood BindJSON Error: ", err)
		return
	}
	var ids idsRequest
	if err := c.ShouldBindQuery(&ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err})
		log.Println("PatchGoodUpdate ShouldBindQuery Error: ", err)
		return
	}

	if err := gh.validator.Struct(patchGoodRequest); err == nil {
		if err = gh.validator.Struct(ids); err != nil {
			c.JSON(http.StatusBadRequest, "Invalid ids error")
			log.Println("PatchGoodUpdate Validation Error: ", err)
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, "Invalid request body error:")
		log.Println("PatchGoodUpdate Validation Error: ", err)
		return
	}

	id := ids.Id
	projectId := ids.ProjectId

	existing, err := gh.goodsRepository.ExistionCheck(c, id, projectId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		log.Println("PatchGoodUpdate ExistionCheck Error: ", err)
		return
	}
	if !existing {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    3,
			"message": "errors.good.notFound",
			"details": gin.H{},
		})
		log.Println("PatchGoodUpdate ExistionCheck: record not found")
		return
	}

	good, err := gh.goodsRepository.UpdateGood(c.Request.Context(), id, projectId, name, description)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		log.Println("PatchGoodUpdate UpdateGood Error: ", err)
		return
	}

	err = gh.redisRepository.InvalidateCache(c.Request.Context(), gh.requestParams.Limit, gh.requestParams.Offset)
	if err != nil {
		log.Println("PatchGoodUpdate CacheInvalidation Error")
	}
	log.Println("PatchGoodUpdate CacheInvalidation Error")

	log.Println("Good updated successfully")
	c.JSON(http.StatusOK, good)
}

func (gh *goodsHandler) DeleteGood(c *gin.Context) {
	var ids idsRequest
	if err := c.ShouldBindQuery(&ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err})
		log.Println("DeleteGood ShouldBindQuery Error: ", err)
		return
	}

	if err := gh.validator.Struct(ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Invalid ids error:": err})
		log.Println("DeleteGood Validation Error: ", err)
		return
	}

	id := ids.Id
	projectId := ids.ProjectId

	existing, err := gh.goodsRepository.ExistionCheck(c, ids.Id, ids.ProjectId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		log.Println("DeleteGood ExistionCheck Error: ", err)
		return
	}
	if !existing {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    3,
			"message": "errors.good.notFound",
			"details": gin.H{},
		})
		log.Println("DeleteGood ExistionCheck: record not found")
		return
	}

	response, err := gh.goodsRepository.RemoveGood(c.Request.Context(), id, projectId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		log.Println("DeleteGood RemoveGood Error: ", err)
		return
	}

	err = gh.redisRepository.InvalidateCache(c.Request.Context(), gh.requestParams.Limit, gh.requestParams.Offset)
	if err != nil {
		log.Println("DeleteGood CacheInvalidation Error")
	}
	log.Println("DeleteGood CacheInvalidation Error")

	log.Println("Good removed successfully")
	c.JSON(http.StatusOK, response)
}

func (gh *goodsHandler) GetGoods(c *gin.Context) {
	var goodsRequestParams goodsRequest
	if err := c.ShouldBindQuery(&goodsRequestParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err})
		log.Println("GetGoods ShouldBindQuery Error: ", err)
		return
	}
	if err := gh.validator.Struct(goodsRequestParams); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid parameters")
		log.Println("GetGoods Validation Error: ", err)
		return
	}

	limit := goodsRequestParams.Limit
	offset := goodsRequestParams.Offset

	gh.requestParams.Limit = limit
	gh.requestParams.Offset = offset

	response, err := gh.redisRepository.GetFromCache(c.Request.Context(), limit, offset)
	if err != nil {
		log.Printf("GetGoods GetFromCache Error: %v Trying to Get from DB", err)
		response, err = gh.goodsRepository.ListGoods(c.Request.Context(), limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			log.Println("GetGoods ListGoods Error: ", err)
			return
		}
		err = gh.redisRepository.CacheTheGoods(c.Request.Context(), response, limit, offset)
		if err != nil {
			log.Println("GetGoods CacheTheGoods Failed to put into cache: ", err)
		}
	}
	log.Println("Goods listed successfully")
	c.JSON(http.StatusOK, response)
}

func (gh *goodsHandler) PatchGoodReprioritiize(c *gin.Context) {
	var ids idsRequest
	if err := c.ShouldBindQuery(&ids); err != nil {
		c.JSON(http.StatusBadRequest, "Query error")
		log.Println("PatchGoodReprioritiize ShouldBindQuery Error: ", err)
		return
	}

	if err := gh.validator.Struct(ids); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid ids error")
		log.Println("PatchGoodReprioritiize Validation Error: ", err)
		return
	}

	existing, err := gh.goodsRepository.ExistionCheck(c, ids.Id, ids.ProjectId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		log.Println("PatchGoodReprioritiize ExistionCheck Error: ", err)
		return
	}
	if !existing {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    3,
			"message": "errors.good.notFound",
			"details": gin.H{},
		})
		log.Println("PatchGoodReprioritiize ExistionCheck: record not found")
		return
	}

	var ReprioritiizeRequest patchGoodReprioritiizeRequest

	if err := c.BindJSON(&ReprioritiizeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": err})
		log.Println("PatchGoodReprioritiize BindJSON Error: ", err)
		return
	}
	id := ids.Id
	projectId := ids.ProjectId

	response, err := gh.goodsRepository.ReprioritiizeGood(c.Request.Context(), id, projectId, ReprioritiizeRequest.NewPriority)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		log.Println("PatchGoodReprioritiize ReprioritiizeGood Error: ", err)
		return
	}

	err = gh.redisRepository.InvalidateCache(c.Request.Context(), gh.requestParams.Limit, gh.requestParams.Offset)
	if err != nil {
		log.Println("PatchGoodReprioritiize CacheInvalidation Error")
	}
	log.Println("PatchGoodReprioritiize CacheInvalidated")

	c.JSON(http.StatusOK, response)
}
