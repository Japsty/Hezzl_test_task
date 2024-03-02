package handlers

import (
	"Hezzl_test_task/internal/natsclient"
	"Hezzl_test_task/internal/storage/repos"
	"Hezzl_test_task/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
	"os"
	"strconv"
)

type goodsHandler struct {
	goodsRepository repos.Repository
	redisRepository repos.RedisRepository
	requestParams   requestParams
	validator       *validator.Validate
	natsConn        *natsclient.NATSClient
}

// Я понимаю, что данное решение совсем плохое, но я не смог понять каким образом
// мне реализовать инвалидацию не имея доступа к полям limit и offset при этом не
// совершая сохранение построчно, что накладно т.к. у нас будет огромное количество запросов
// в реляционную бд
//
// Буду рад получить ответ в фидбеке на то, как вернее было бы реализовать
type requestParams struct {
	Limit  int
	Offset int
}

// NewGoodsHandler функция для настройки эндпоинтов и создания экземпляра структуры goodsHandler
func NewGoodsHandler(repo repos.Repository, redis repos.RedisRepository, natsConn *natsclient.NATSClient) *gin.Engine {

	router := gin.Default()
	router.Use(logger.LogMiddleware())

	gh := goodsHandler{
		goodsRepository: repo,
		redisRepository: redis,
		validator:       validator.New(),
		natsConn:        natsConn,
	}

	router.POST("/good/create", gh.AddGood)
	router.PATCH("/good/update", gh.PatchGoodUpdate)
	router.DELETE("/good/remove", gh.DeleteGood)
	router.GET("/goods/list", gh.GetGoods)
	router.PATCH("/good/reprioritiize", gh.PatchGoodReprioritiize)

	return router
}

// AddGood example
//
//	@Summary Add a new good
//	@Description Add a new good to the database
//	@Tags good
//	@Accept json
//	@Produce json
//	@Param projectId query int true "Project ID"
//	@Param request body addGoodRequest true "Good details"
//	@Success 		201 	{string}	string "ok"
//	@Failure		400		{string}	string "bad input"
//	@Failure		500		{string}	string "Internal Server Error"
//	@Router /good/create [post]
func (gh *goodsHandler) AddGood(c *gin.Context) {
	var AddGoodRequest addGoodRequest

	if err := c.BindJSON(&AddGoodRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad input"})
		log.Println("AddGood BindJSON Error: ", err)
		return
	}

	if err := gh.validator.Struct(AddGoodRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad input"})
		log.Println("AddGood Validation Error: ", err)
		return
	}

	projectId, err := strconv.Atoi(c.Query(""))
	if err != nil || projectId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad input"})
		log.Println("Invalid 'projectId' parameter: ", projectId)
		return
	}

	good, err := gh.goodsRepository.CreateGood(c.Request.Context(), projectId, AddGoodRequest.Name, AddGoodRequest.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		log.Println("AddGood CreateGood Error: ", err)
		return
	}
	log.Println("Good added successfully")

	c.JSON(http.StatusCreated, good)
}

// PatchGoodUpdate example
//
//	@Summary Update an existing good
//	@Description Update details of an existing good
//	@Tags good
//	@Accept json
//	@Produce json
//	@Param id path int true "ID"
//	@Param projectId path int true "ProjectID"
//	@Param request body updateGoodRequest true "Updated good details"
//	@Success 		200 	{string} 	string "ok"
//	@Failure		400		{string}	string "bad input"
//	@Failure		404		{string}	string "not found"
//	@Failure		500		{string}	string "Internal Server Error"
//	@Router /good/update [patch]
func (gh *goodsHandler) PatchGoodUpdate(c *gin.Context) {
	var patchGoodRequest updateGoodRequest

	if err := c.BindJSON(&patchGoodRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad input"})
		log.Println("AddGood BindJSON Error: ", err)
		return
	}

	name := patchGoodRequest.Name
	description := patchGoodRequest.Description

	var ids idsRequest
	if err := c.ShouldBindQuery(&ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad input"})
		log.Println("PatchGoodUpdate ShouldBindQuery Error: ", err)
		return
	}

	if err := gh.validator.Struct(ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad input"})
		log.Println("PatchGoodUpdate Validation Error: ", err)
		return
	}
	if err := gh.validator.Struct(patchGoodRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad input"})
		log.Println("PatchGoodUpdate Validation Error: ", err)
		return
	}

	id := ids.Id
	projectId := ids.ProjectId

	existing, err := gh.goodsRepository.ExistionCheck(c, id, projectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		log.Println("PatchGoodUpdate UpdateGood Error: ", err)
		return
	}

	err = gh.redisRepository.InvalidateCache(c.Request.Context(), gh.requestParams.Limit, gh.requestParams.Offset)
	if err != nil {
		log.Println("PatchGoodUpdate CacheInvalidation Error:", err)
	}

	log.Println("Good updated successfully")

	c.JSON(http.StatusOK, good)

	clickhouseLog := repos.GoodToClickhouseLog(good)

	log.Println("NATS Subject:", os.Getenv("NATS_SUBJECT"))
	err = gh.natsConn.PublishMessage(os.Getenv("NATS_SUBJECT"), clickhouseLog)
	if err != nil {
		log.Println("PatchGoodUpdate NATS PublishMessage Error: ", err)
		return
	}
}

// DeleteGood example
//
//	@Summary Delete a good
//	@Description Delete a good by ID
//	@Tags good
//	@Accept json
//	@Produce json
//	@Param id path int true "Good ID"
//	@Param projectId path int true "Project ID"
//	@Success 		200 	{string} 	string "ok"
//	@Failure		400		{string}	string "bad input"
//	@Failure		404		{string}	string "not found"
//	@Failure		500		{string}	string "Internal Server Error"
//	@Router /good/remove [delete]
func (gh *goodsHandler) DeleteGood(c *gin.Context) {
	var ids idsRequest
	if err := c.ShouldBindQuery(&ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad input"})
		log.Println("DeleteGood ShouldBindQuery Error: ", err)
		return
	}

	if err := gh.validator.Struct(ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad input"})
		log.Println("DeleteGood Validation Error: ", err)
		return
	}

	id := ids.Id
	projectId := ids.ProjectId

	existing, err := gh.goodsRepository.ExistionCheck(c, ids.Id, ids.ProjectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		log.Println("DeleteGood RemoveGood Error: ", err)
		return
	}

	err = gh.redisRepository.InvalidateCache(c.Request.Context(), gh.requestParams.Limit, gh.requestParams.Offset)
	if err != nil {
		log.Println("DeleteGood CacheInvalidation Error:", err)
	}
	log.Println("Good removed successfully")

	c.JSON(http.StatusOK, response)

	good, err := gh.goodsRepository.GetGoodById(c.Request.Context(), id, projectId)
	if err != nil {
		log.Println("DeleteGood GetGoodById Error: ", err)
	}
	clickhouseLog := repos.GoodToClickhouseLog(good)

	err = gh.natsConn.PublishMessage(os.Getenv("NATS_SUBJECT"), clickhouseLog)
	if err != nil {
		log.Println("DeleteGood NATS PublishMessage Error: ", err)
		return
	}
}

// GetGoods
//
//	@Summary Get list of goods
//	@Description Get list of goods with pagination parameters
//	@Tags good
//	@Accept json
//	@Produce json
//	@Param limit query int false "Limit"
//	@Param offset query int false "Offset"
//	@Success 		200 	{string} 	string "ok"
//	@Failure		400		{string}	string "bad input"
//	@Failure		500		{string}	string "Internal Server Error"
//	@Router /goods/list [get]
func (gh *goodsHandler) GetGoods(c *gin.Context) {
	var goodsRequestParams goodsRequest
	if err := c.ShouldBindQuery(&goodsRequestParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": "bad input"})
		log.Println("GetGoods ShouldBindQuery Error: ", err)
		return
	}
	if err := gh.validator.Struct(goodsRequestParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: ": "bad input"})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
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

// PatchGoodReprioritiize
//
//	@Summary Reprioritize a good
//	@Description Change the priority of a good
//	@Tags goods
//	@Accept json
//	@Produce json
//	@Param id path int true "Good ID"
//	@Param projectId path int true "Project ID"
//	@Param request body patchGoodReprioritiizeRequest true "New priority details"
//	@Success 		200 	{string} 	string "ok"
//	@Failure		400		{string}	string "bad input"
//	@Failure		404		{string}	string "not found"
//	@Failure		500		{string}	string "Internal Server Error"
//	@Router /good/reprioritiize [patch]
func (gh *goodsHandler) PatchGoodReprioritiize(c *gin.Context) {
	var ids idsRequest
	if err := c.ShouldBindQuery(&ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad input"})
		log.Println("PatchGoodReprioritiize ShouldBindQuery Error: ", err)
		return
	}

	if err := gh.validator.Struct(ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad input"})
		log.Println("PatchGoodReprioritiize Validation Error: ", err)
		return
	}

	existing, err := gh.goodsRepository.ExistionCheck(c, ids.Id, ids.ProjectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
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
		c.JSON(http.StatusBadRequest, gin.H{"error: ": "bad input"})
		log.Println("PatchGoodReprioritiize BindJSON Error: ", err)
		return
	}
	id := ids.Id
	projectId := ids.ProjectId

	response, err := gh.goodsRepository.ReprioritiizeGood(c.Request.Context(), id, projectId, ReprioritiizeRequest.NewPriority)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		log.Println("PatchGoodReprioritiize ReprioritiizeGood Error: ", err)
		return
	}

	err = gh.redisRepository.InvalidateCache(c.Request.Context(), gh.requestParams.Limit, gh.requestParams.Offset)
	if err != nil {
		log.Println("PatchGoodReprioritiize CacheInvalidation Error: ", err)
	}

	c.JSON(http.StatusOK, response)

	good, err := gh.goodsRepository.GetGoodById(c.Request.Context(), id, projectId)
	if err != nil {
		log.Println("PatchGoodReprioritiize GetGoodById Error: ", err)
	}
	clickhouseLog := repos.GoodToClickhouseLog(good)

	log.Println("NATS Subject:", os.Getenv("NATS_SUBJECT"))
	err = gh.natsConn.PublishMessage(os.Getenv("NATS_SUBJECT"), clickhouseLog)
	if err != nil {
		log.Println("PatchGoodReprioritiize NATS PublishMessage Error: ", err)
		return
	}
}
