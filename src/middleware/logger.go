package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sulin2018/go-web-base/src/utils"
)

func AppLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s %s %s %s %s %d %s %s %s \n",
			param.TimeStamp.Format(utils.TIMEFORMAT),
			param.ClientIP,
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}
