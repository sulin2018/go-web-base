package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sulin2018/go-web-base/src/app/config"
	"github.com/sulin2018/go-web-base/src/middleware"
	"github.com/sulin2018/go-web-base/src/models"
	"github.com/sulin2018/go-web-base/src/utils"
)

func GetLoginUser(c *gin.Context) *models.User {
	username, err := middleware.GetSession(c, "username")
	if err != nil || username == "" {
		return nil
	}

	user := models.User{Username: username.(string)}
	err = models.Detail(user)
	if err != nil {
		return nil
	}

	return &user
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

// 详情
func UserGet(c *gin.Context) {
	var user models.User
	err := c.BindUri(&user)
	if err != nil {
		logrus.Error(err)
		ResponseJson(c, http.StatusBadRequest, "ID错误: "+err.Error())
		return
	}

	err = models.Detail(&user)
	if err != nil {
		ResponseJson(c, http.StatusBadRequest, err.Error())
		return
	}

	err = user.LoadAllAssociationIds()
	if err != nil {
		ResponseJson(c, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseJson(c, http.StatusOK, struct {
		*models.User
		// 忽略password 不传递给客户端
		Password bool `json:"password,omitempty"`
	}{
		User: &user,
	})
}

// 新增
func UserPost(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		logrus.Error(err)
		ResponseJson(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if user.Password == "" {
		user.Password = config.AppConfig.UserBasePassword
	}

	err := user.EncryptPassword()
	if err != nil {
		ResponseJson(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = user.Create()
	if err != nil {
		ResponseJson(c, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseJson(c, http.StatusCreated, user)
}

// 更新
func UserPatch(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindUri(&user); err != nil || user.ID == 0 {
		logrus.Error(err)
		ResponseJson(c, http.StatusBadRequest, "ID错误: "+err.Error())
		return
	}

	var newUserData models.User
	if err := c.ShouldBindJSON(&newUserData); err != nil {
		logrus.Error(err)
		ResponseJson(c, http.StatusBadRequest, err.Error())
		return
	}

	err := models.UpdateByMapOrStruct(&user, &newUserData)
	if err != nil {
		ResponseJson(c, http.StatusInternalServerError, err.Error())
		return
	}
	ResponseJson(c, http.StatusOK, user)
}

// 全量更新
func UserPut(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindUri(&user); err != nil || user.ID == 0 {
		logrus.Error(err)
		ResponseJson(c, http.StatusBadRequest, "ID错误: "+err.Error())
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		logrus.Error(err)
		ResponseJson(c, http.StatusBadRequest, err.Error())
		return
	}

	err := user.Update()
	if err != nil {
		ResponseJson(c, http.StatusInternalServerError, err.Error())
		return
	}
	ResponseJson(c, http.StatusOK, user)
}

// 删除
func UserDelete(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindUri(&user); err != nil || user.ID == 0 {
		logrus.Error(err)
		ResponseJson(c, http.StatusBadRequest, "ID错误: "+err.Error())
		return
	}

	err := user.Delete()
	if err != nil {
		ResponseJson(c, http.StatusInternalServerError, err.Error())
		return
	}
	ResponseJson(c, http.StatusNoContent, nil)
}

// 列表
func UsersGet(c *gin.Context) {
	var users []*models.User
	var count uint

	page := utils.StrTo(c.Query("page")).Uint()
	pageSize := utils.StrTo(c.Query("pagesize")).Uint()
	err := models.PageColumns(&users, &count, page, pageSize, "id, username, chinese_name, active, superuser, created_at, updated_at")
	if err != nil {
		ResponseJson(c, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseJsonMore(c, http.StatusOK, users, map[string]interface{}{"count": count})
}

// 登录
func UserLogin(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		logrus.Error(err)
		ResponseJson(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	if user.CheckPassword() {
		err := middleware.SetSession(c, "username", user.Username)
		if err != nil {
			ResponseJson(c, http.StatusInternalServerError, err.Error())
			return
		}

		err = user.LoadAllAssociations()
		if err != nil {
			ResponseJson(c, http.StatusInternalServerError, err.Error())
			return
		}

		user.Password = ""
		ResponseJson(c, http.StatusOK, user)
	} else {
		ResponseJson(c, http.StatusBadRequest, "账号或密码错误")
	}
}

func GroupGet(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindUri(&group); err != nil {
		logrus.Error(err)
		ResponseJson(c, http.StatusBadRequest, "ID错误: "+err.Error())
		return
	}

	err := models.Detail(&group)
	if err != nil {
		ResponseJson(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = group.LoadAllAssociationIds()
	if err != nil {
		ResponseJson(c, http.StatusInternalServerError, err.Error())
		return
	}
	ResponseJson(c, http.StatusOK, group)
}

func GroupPut(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindUri(&group); err != nil {
		logrus.Error(err)
		ResponseJson(c, http.StatusBadRequest, "ID错误: "+err.Error())
		return
	}

	if err := c.ShouldBindJSON(&group); err != nil {
		logrus.Error(err)
		ResponseJson(c, http.StatusBadRequest, err.Error())
		return
	}

	err := group.Update()
	if err != nil {
		ResponseJson(c, http.StatusInternalServerError, err.Error())
		return
	}
	ResponseJson(c, http.StatusOK, group)
}

func GroupPost(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		logrus.Error(err)
		ResponseJson(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	err := group.Create()
	if err != nil {
		ResponseJson(c, http.StatusInternalServerError, err.Error())
		return
	}
	ResponseJson(c, http.StatusCreated, group)
}

func GroupPatch(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindUri(&group); err != nil || group.ID == 0 {
		logrus.Error(err)
		ResponseJson(c, http.StatusBadRequest, "ID错误: "+err.Error())
		return
	}

	var newData models.Group
	if err := c.ShouldBindJSON(&newData); err != nil {
		logrus.Error(err)
		ResponseJson(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := models.UpdateByMapOrStruct(&group, &newData); err != nil {
		ResponseJson(c, http.StatusBadRequest, err.Error())
		return
	}
	ResponseJson(c, http.StatusOK, group)
}

func GroupDelete(c *gin.Context) {
	var group models.Group
	if err := c.ShouldBindUri(&group); err != nil || group.ID == 0 {
		logrus.Error(err)
		ResponseJson(c, http.StatusBadRequest, "ID错误: "+err.Error())
		return
	}

	err := group.Delete()
	if err != nil {
		ResponseJson(c, http.StatusInternalServerError, err.Error())
		return
	}
	ResponseJson(c, http.StatusNoContent, nil)
}

func GroupsGet(c *gin.Context) {
	var group []*models.Group
	var count uint

	page := utils.StrTo(c.Query("page")).Uint()
	pageSize := utils.StrTo(c.Query("pagesize")).Uint()
	err := models.Page(&group, &count, page, pageSize)
	if err != nil {
		ResponseJson(c, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseJsonMore(c, http.StatusOK, group, map[string]interface{}{"count": count})
}

func PermissionGet(c *gin.Context) {
	var permission models.Permission
	if err := c.ShouldBindUri(&permission); err != nil {
		logrus.Error(err)
		ResponseJson(c, http.StatusBadRequest, "ID错误: "+err.Error())
		return
	}

	err := models.Detail(&permission)
	if err != nil {
		logrus.Error(err)
		ResponseJson(c, http.StatusBadRequest, err.Error())
		return
	}

	err = permission.LoadAllAssociationIds()
	if err != nil {
		ResponseJson(c, http.StatusInternalServerError, err.Error())
		return
	}
	ResponseJson(c, http.StatusOK, permission)
}

func PermissionPost(c *gin.Context) {
	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		logrus.Error(err)
		ResponseJson(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	err := permission.Create()
	if err != nil {
		ResponseJson(c, http.StatusInternalServerError, err.Error())
		return
	}
	ResponseJson(c, http.StatusCreated, permission)
}

func PermissionPatch(c *gin.Context) {
	var permission models.Permission
	if err := c.ShouldBindUri(&permission); err != nil || permission.ID == 0 {
		logrus.Error(err)
		ResponseJson(c, http.StatusBadRequest, "ID错误: "+err.Error())
		return
	}

	if err := c.ShouldBindJSON(&permission); err != nil {
		logrus.Error(err)
		ResponseJson(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	err := permission.Update()
	if err != nil {
		ResponseJson(c, http.StatusInternalServerError, err.Error())
		return
	}
	ResponseJson(c, http.StatusOK, permission)
}

func PermissionDelete(c *gin.Context) {
	var permission models.Permission
	if err := c.ShouldBindUri(&permission); err != nil || permission.ID == 0 {
		logrus.Error(err)
		ResponseJson(c, http.StatusBadRequest, "ID错误: "+err.Error())
		return
	}

	err := permission.Delete()
	if err != nil {
		ResponseJson(c, http.StatusInternalServerError, err.Error())
		return
	}
	ResponseJson(c, http.StatusNoContent, nil)
}

func PermissionsGet(c *gin.Context) {
	var permissions []*models.Permission
	var count uint

	page := utils.StrTo(c.Query("page")).Uint()
	pageSize := utils.StrTo(c.Query("pagesize")).Uint()
	err := models.Page(&permissions, &count, page, pageSize)
	if err != nil {
		ResponseJson(c, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseJsonMore(c, http.StatusOK, permissions, map[string]interface{}{"count": count})
}
