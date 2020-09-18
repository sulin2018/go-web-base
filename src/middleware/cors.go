package middleware

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/sulin2018/go-web-base/src/app/config"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")

		// 过滤，防止不合法的域名访问
		var isAccess = false
		for _, v := range config.AppConfig.AppCorsOrigin {
			match, _ := regexp.MatchString(v, origin)
			if match {
				isAccess = true
			}
		}

		if isAccess {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, Accept, Origin, Cache-Control")
			c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PUT, DELETE")
			c.Header("Access-Control-Max-Age", "172800")
			c.Set("content-type", "application/json")
		}

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, nil)
		}

		c.Next()
	}
}
