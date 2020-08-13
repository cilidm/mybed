package controller

import (
	"mybedv2/app/helper/router"
	"mybedv2/app/system/controller/user_list_controller"
	"mybedv2/app/system/middleware"
)

func init() {
	r := router.New("system", "/system", middleware.AuthMiddleware, middleware.SidebarMiddleware)
	r.GET("user_list", user_list_controller.UserListPage)
	r.GET("user_list_json", user_list_controller.UserListJson)
	r.POST("user_change_status", user_list_controller.UpdateUserStatus)
}
