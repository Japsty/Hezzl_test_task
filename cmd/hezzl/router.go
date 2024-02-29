package main

import (
	"Hezzl_test_task/internal/handlers"
	"Hezzl_test_task/pkg/logger"
	"Hezzl_test_task/pkg/storage"
	"github.com/gin-gonic/gin"
)

func SetupRouter(repo storage.Repository) *gin.Engine {
	router := gin.Default()
	router.Use(logger.LogMiddleware())

	gh := handlers.GoodsHandler{
		GoodsRepository: repo,
	}

	router.POST("/good/create", gh.AddGood)
	router.PATCH("/good/update", gh.PatchGood)
	router.DELETE("/good/remove", gh.DeleteGood)
	router.GET("/goods/list", gh.GetGoods)
	router.PATCH("/good/reprioritiize")

	return router
}
