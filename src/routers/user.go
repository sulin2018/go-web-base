package routers

import (
	"github.com/sulin2018/go-web-base/src/controllers"
	"github.com/sulin2018/go-web-base/src/middleware"
)

func AddUserV1Router() {
	userApi := g.Group("/api/v1")
	userApi.Use(middleware.PermissionMiddleware("manage_user"))

	// user
	userApi.GET("/user/:id", controllers.UserGet)
	userApi.POST("/user", controllers.UserPost)
	userApi.PATCH("/user/:id", controllers.UserPatch)
	userApi.DELETE("/user/:id", controllers.UserDelete)
	userApi.GET("/users", controllers.UsersGet)

	// group
	userApi.GET("/group/:id", controllers.GroupGet)
	userApi.PATCH("/group/:id", controllers.GroupPatch)
	userApi.DELETE("/group/:id", controllers.GroupDelete)
	userApi.PUT("/group/:id", controllers.GroupPut)
	userApi.POST("/group", controllers.GroupPost)
	userApi.GET("/groups", controllers.GroupsGet)

	// permission
	userApi.GET("/permission/:id", controllers.PermissionGet)
	userApi.PATCH("/permission/:id", controllers.PermissionPatch)
	userApi.DELETE("/permission/:id", controllers.PermissionDelete)
	userApi.POST("/permission", controllers.PermissionPost)
	userApi.GET("/permissions", controllers.PermissionsGet)
}
