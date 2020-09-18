package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var statusMessage = map[int]string{
	http.StatusBadRequest:       "参数有误",
	http.StatusUnauthorized:     "缺少认证信息",
	http.StatusForbidden:        "无权限",
	http.StatusMethodNotAllowed: "服务器未实现的请求方法",

	http.StatusInternalServerError: "服务器出错",
	http.StatusNotImplemented:      "服务器未实现",
}

func ResponseJson(c *gin.Context, httpCode int, data interface{}) {
	var result = map[string]interface{}{
		"code":    httpCode,
		"message": GetStatusMsg(httpCode),
		"data":    data,
	}
	c.JSON(httpCode, result)
}

func ResponseJsonMore(c *gin.Context, httpCode int, data interface{}, moreInfo map[string]interface{}) {
	var result = map[string]interface{}{
		"code":    httpCode,
		"message": GetStatusMsg(httpCode),
		"data":    data,
	}
	for k, v := range moreInfo {
		result[k] = v
	}
	c.JSON(httpCode, result)
}

func GetStatusMsg(code int) string {
	msg, ok := statusMessage[code]
	if ok {
		return msg
	}
	return http.StatusText(code)
}

func Ping(c *gin.Context) {
	ResponseJson(c, http.StatusOK, nil)
}
