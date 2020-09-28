package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sulin2018/go-web-base/src/app/config"
	"github.com/sulin2018/go-web-base/src/app/log"
	"github.com/sulin2018/go-web-base/src/middleware"
	"github.com/sulin2018/go-web-base/src/models"
	"github.com/sulin2018/go-web-base/src/routers"
)

func init() {
	config.InitConfig(*flag.String("config", "config.yaml", "config file path"))
	log.InitLogrus()
	models.DBInit()
	middleware.InitSessionStore()
}

func main() {
	if config.AppConfig.AppRunMode == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	ginEngine := routers.InitGinEngine()
	//readTimeout := config.AppConfig.AppReadTimeout
	//writeTimeout := config.AppConfig.AppWriteTimeout
	endPoint := fmt.Sprintf("%s:%d", config.AppConfig.AppAddr, config.AppConfig.AppPort)
	//maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:    endPoint,
		Handler: ginEngine,
		//ReadTimeout:    readTimeout,
		//WriteTimeout:   writeTimeout,
		//MaxHeaderBytes: maxHeaderBytes,
	}

	logrus.Info("Start http server, listening at ", endPoint)

	err := server.ListenAndServe()
	if err != nil {
		logrus.Error(err)
	}

	// ginEngine.Run(endPoint)
}
