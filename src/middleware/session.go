package middleware

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"github.com/sulin2018/go-web-base/src/app/config"
	"github.com/sulin2018/go-web-base/src/models"
	"github.com/wader/gormstore"
)

var store *gormstore.Store

func InitSessionStore() {
	if store == nil {
		store = gormstore.New(models.GetDB(), []byte(config.AppConfig.AppSecret))
		// db cleanup every hour
		// close quit channel to stop cleanup
		quit := make(chan struct{})
		go store.PeriodicCleanup(1*time.Hour, quit)
	}
}

func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cSession, err := store.Get(c.Request, "session")
		if err != nil {
			logrus.Error(err)
		}
		c.Set("session", cSession)

		// before request

		c.Next()

		// after request
	}
}

func GetSession(c *gin.Context, key string) (interface{}, error) {
	// session, err := store.Get(c.Request, "session")
	// if err != nil {
	// 	return nil, err
	// }
	tempValue, exists := c.Get("session")
	if !exists {
		logrus.Error("获取sessions有误")
		return nil, errors.New("session not exist")
	}
	cSession := tempValue.(*sessions.Session)

	value, ok := cSession.Values[key]
	if ok {
		return value, nil
	}
	return nil, errors.New("not exist")
}

func SetSession(c *gin.Context, key string, value interface{}) error {
	// session, err := store.Get(c.Request, "session")
	// if err != nil {
	// 	return err
	// }
	tempValue, exists := c.Get("session")
	if !exists {
		logrus.Error("获取sessions有误")
		return errors.New("session not exist")
	}
	cSession := tempValue.(*sessions.Session)

	cSession.Values[key] = value
	err := cSession.Save(c.Request, c.Writer)
	if err != nil {
		logrus.Error("设置session有误")
		logrus.Error(err)
		return err
	}
	return nil
}
