package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/sulin2018/go-web-base/src/controllers"
	"github.com/sulin2018/go-web-base/src/middleware"
)

var apiv1 *gin.RouterGroup
var g *gin.Engine

func InitGinEngine() *gin.Engine {
	g = gin.New()
	// g.Use(gin.Logger())
	g.Use(gin.Recovery())

	g.Use(middleware.AppLogger())
	g.Use(middleware.SessionMiddleware())
	g.Use(middleware.CorsMiddleware())

	apiv1 = g.Group("/api/v1")
	apiv1.GET("/ping", controllers.Ping)
	apiv1.POST("/user/login", controllers.UserLogin)

	AddUserV1Router()

	return g
}
