package main

import (
	"Hezzl_test_task/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.Use(logger.LogMiddleware())

	router.POST("/good/create")
	router.PATCH("/good/update")
	router.DELETE("/good/remove")
	router.GET("/goods/list")
	router.PATCH("/good/reprioritiize")
}
