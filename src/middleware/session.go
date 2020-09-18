package middleware

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
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

func GetSession(c *gin.Context, key string) (interface{}, error) {
	session, err := store.Get(c.Request, "session")
	if err != nil {
		return nil, err
	}
	value, ok := session.Values[key]
	if ok {
		return value, nil
	}
	return nil, errors.New("not exist")
}

func SetSession(c *gin.Context, key string, value interface{}) error {
	session, err := store.Get(c.Request, "session")
	if err != nil {
		return err
	}
	session.Values[key] = value
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		return err
	}
	return nil
}
