package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sulin2018/go-web-base/src/models"
)

func GetLoginUser(c *gin.Context) *models.User {
	username, err := GetSession(c, "username")
	if err != nil || username == "" {
		return nil
	}

	user := models.User{Username: username.(string)}
	err = models.Detail(&user)
	if err != nil {
		return nil
	}

	return &user
}

func LoginPermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if GetLoginUser(c) == nil {
			c.AbortWithStatus(http.StatusForbidden)
		}

		// before request
		c.Next()
		// after request
	}
}

func PermissionMiddleware(permName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !CheckPermission(c, permName) {
			c.AbortWithStatusJSON(http.StatusForbidden, map[string]interface{}{"code": http.StatusForbidden, "message": "无权限"})
		}

		// before request
		c.Next()
		// after request
	}
}

func CheckPermission(c *gin.Context, permName string) bool {
	user := GetLoginUser(c)
	if user == nil {
		return false
	}

	if user.Superuser {
		return true
	}

	tempPerm := &models.Permission{Name: permName}
	err := models.LoadColumns(tempPerm, []string{"id"})
	if err != nil {
		return false
	}

	// 先从用户关联权限寻找
	err = user.LoadPermAssociationIds()
	if err != nil {
		return false
	}
	for _, permId := range user.PermissionIds {
		if permId == tempPerm.ID {
			return true
		}
	}

	// 从关联组里面寻找是否有对应权限
	err = user.LoadGroupAssociationIds()
	if err != nil {
		return false
	}
	err = tempPerm.LoadGroupAssociationIds()
	if err != nil {
		return false
	}
	for _, userGroupId := range user.GroupIds {
		for _, permGroupId := range tempPerm.GroupIds {
			if userGroupId == permGroupId {
				return true
			}
		}
	}

	return false
}
