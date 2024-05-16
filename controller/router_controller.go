package controller

import "github.com/gin-gonic/gin"

func InitRoutes() *gin.Engine {
	router := gin.Default()
	router.POST("/upload", UpLoadFileController)
	router.POST("/connection/show", TcpConnectionShow)
	return router
}