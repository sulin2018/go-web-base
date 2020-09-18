package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	logs "github.com/sirupsen/logrus"
	"github.com/sulin2018/go-web-base/src/app/config"
	"github.com/sulin2018/go-web-base/src/models"
	"github.com/sulin2018/go-web-base/src/utils"
)

func UserGet(c *gin.Context) {
	var user models.User
	err := c.BindUri(&user)

	if err != nil {
		logs.Error(err)
		ResponseJson(c, http.StatusBadRequest, "ID错误")
		return
	}

	err = models.Detail(&user)
	if err != nil {
		logs.Error(err)
		ResponseJson(c, http.StatusBadRequest, "ID错误")
		return
	}

	user.GetAllAssociationIds()
	ResponseJson(c, http.StatusOK, struct {
		*models.User
		// 忽略password 不传递给客户端
		Password bool `json:"password,omitempty"`
	}{
		User: &user,
	})
}

func UserPost(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		logs.Error(err)
		ResponseJson(c, http.StatusBadRequest, "参数错误")
		return
	}

	if user.Password == "" {
		user.Password = config.AppConfig.UserBasePassword
	}

	user.EncryptPassword()

	user.Create()
	ResponseJson(c, http.StatusCreated, user)
}

func UserPatch(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		logs.Error(err)
		ResponseJson(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := c.ShouldBindUri(&user); err != nil || user.ID == 0 {
		logs.Error(err)
		ResponseJson(c, http.StatusBadRequest, "参数错误: 用户id不能为空")
		return
	}

	user.Update()
	ResponseJson(c, http.StatusOK, user)
}

func UserDelete(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindUri(&user); err != nil || user.ID == 0 {
		logs.Error(err)
		ResponseJson(c, http.StatusBadRequest, "参数错误: 用户id不能为空")
		return
	}

	user.Delete()
	ResponseJson(c, http.StatusNoContent, nil)
}

func UsersGet(c *gin.Context) {
	var users []*models.User
	var page uint
	var count uint

	page = utils.StrTo(c.Query("page")).Uint()
	models.PageColumn(&users, &count, page, "id, username, chinese_name, active, superuser")

	ResponseJsonMore(c, http.StatusOK, users, map[string]interface{}{"count": count})
}

func UserLogin(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		logs.Error(err)
		ResponseJson(c, http.StatusBadRequest, "参数错误")
		return
	}

	if user.CheckPassword() {
		user.GetAllAssociations()
		ResponseJson(c, http.StatusOK, user)
	} else {
		ResponseJson(c, http.StatusBadRequest, "密码错误")
	}
}

func GroupGet(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindUri(&group); err != nil {
		logs.Error(err)
		ResponseJson(c, http.StatusBadRequest, "ID错误")
		return
	}

	err := models.Detail(&group)
	if err != nil {
		logs.Error(err)
		ResponseJson(c, http.StatusBadRequest, "ID错误")
		return
	}

	group.GetAllAssociationIds()
	ResponseJson(c, http.StatusOK, group)
}

func GroupPut(c *gin.Context) {
	// todo: 未完全实现, put实践要求更新所有字段, 此处只更新了传递的字段
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		logs.Error(err)
		ResponseJson(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := c.ShouldBindUri(&group); err != nil {
		logs.Error(err)
		ResponseJson(c, http.StatusBadRequest, "ID错误")
		return
	}

	group.CreateOrUpdate()
	ResponseJson(c, http.StatusOK, group)
}

func GroupPost(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		logs.Error(err)
		ResponseJson(c, http.StatusBadRequest, "参数错误")
		return
	}

	group.Create()
	ResponseJson(c, http.StatusCreated, group)
}

func GroupPatch(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		logs.Error(err)
		ResponseJson(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := c.ShouldBindUri(&group); err != nil || group.ID == 0 {
		logs.Error(err)
		ResponseJson(c, http.StatusBadRequest, "ID错误")
		return
	}

	group.Update()
	ResponseJson(c, http.StatusOK, group)
}

func GroupDelete(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindUri(&group); err != nil || group.ID == 0 {
		logs.Error(err)
		ResponseJson(c, http.StatusBadRequest, "ID错误")
		return
	}

	group.Delete()
	ResponseJson(c, http.StatusNoContent, nil)
}

func GroupsGet(c *gin.Context) {
	var group []*models.Group
	var page uint
	var count uint

	page = utils.StrTo(c.Query("page")).Uint()
	models.Page(&group, &count, page)

	ResponseJsonMore(c, http.StatusOK, group, map[string]interface{}{"count": count})
}

func PermissionGet(c *gin.Context) {
	var permission models.Permission
	if err := c.ShouldBindUri(&permission); err != nil {
		logs.Error(err)
		ResponseJson(c, http.StatusBadRequest, "ID错误")
		return
	}

	err := models.Detail(&permission)
	if err != nil {
		logs.Error(err)
		ResponseJson(c, http.StatusBadRequest, "ID错误")
		return
	}

	permission.GetAllAssociationIds()
	ResponseJson(c, http.StatusOK, permission)
}

func PermissionPost(c *gin.Context) {
	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		logs.Error(err)
		ResponseJson(c, http.StatusBadRequest, "参数错误")
		return
	}

	permission.Create()
	ResponseJson(c, http.StatusCreated, permission)
}

func PermissionPatch(c *gin.Context) {
	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		logs.Error(err)
		ResponseJson(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := c.ShouldBindUri(&permission); err != nil || permission.ID == 0 {
		logs.Error(err)
		ResponseJson(c, http.StatusBadRequest, "ID错误")
		return
	}

	permission.Update()
	ResponseJson(c, http.StatusOK, permission)
}

func PermissionDelete(c *gin.Context) {
	var permission models.Permission
	if err := c.ShouldBindUri(&permission); err != nil || permission.ID == 0 {
		logs.Error(err)
		ResponseJson(c, http.StatusBadRequest, "ID错误")
		return
	}

	permission.Delete()
	ResponseJson(c, http.StatusNoContent, nil)
}

func PermissionsGet(c *gin.Context) {
	var permissions []*models.Permission
	var page uint
	var count uint

	page = utils.StrTo(c.Query("page")).Uint()
	models.Page(&permissions, &count, page)

	ResponseJsonMore(c, http.StatusOK, permissions, map[string]interface{}{"count": count})
}
