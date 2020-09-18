package routers

import (
	"github.com/sulin2018/go-web-base/src/controllers"
)

func AddUserV1Router() {
	// user
	apiv1.GET("/user/:id", controllers.UserGet)
	apiv1.POST("/user", controllers.UserPost)
	apiv1.PATCH("/user/:id", controllers.UserPatch)
	apiv1.DELETE("/user/:id", controllers.UserDelete)
	apiv1.GET("/users", controllers.UsersGet)

	// group
	apiv1.GET("/group/:id", controllers.GroupGet)
	apiv1.PATCH("/group/:id", controllers.GroupPatch)
	apiv1.DELETE("/group/:id", controllers.GroupDelete)
	apiv1.PUT("/group/:id", controllers.GroupPut)
	apiv1.POST("/group", controllers.GroupPost)
	apiv1.GET("/groups", controllers.GroupsGet)

	// permission
	apiv1.GET("/permission/:id", controllers.PermissionGet)
	apiv1.PATCH("/permission/:id", controllers.PermissionPatch)
	apiv1.DELETE("/permission/:id", controllers.PermissionDelete)
	apiv1.POST("/permission", controllers.PermissionPost)
	apiv1.GET("/permissions", controllers.PermissionsGet)
}
