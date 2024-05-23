package controller

import (
	"github.com/gin-gonic/gin"
	"scripts-manage/common"
)

func InitRoutes() *gin.Engine {

	router := gin.Default()

	// http认证中间件
	if common.Conf.Server.Auth.Enable {
		router.Use(BasicAuthMiddleware())
		common.Log.Infof("http basic auth enable")
	}

	// 源IP限流中间件
	if common.Conf.Server.Limit.Enable {
		router.Use(RateLimiterMiddleware())
		common.Log.Infof("request rate limit enable")
	}

	// 白名单中间件
	router.Use(WhiteListMiddleware())

	router.POST("/upload", UpLoadFileController)
	router.GET("/connection/show", TcpConnectionShow)

	return router
}
