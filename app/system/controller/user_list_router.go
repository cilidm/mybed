package controller

import (
	"mybedv2/app/helper/router"
	"mybedv2/app/system/controller/user_list_controller"
	"mybedv2/app/system/middleware"
)

func init() {
	r := router.New("system", "/system", middleware.AuthMiddleware, middleware.SidebarMiddleware)
	r.GET("code", user_list_controller.CodeListPage)
	r.GET("code_json", user_list_controller.CodeListJson)
	r.GET("code_form", user_list_controller.CodeFormPage)
	r.POST("code_form", user_list_controller.CreateCode)
	r.POST("code_dels", user_list_controller.DeleteCodes)
	r.POST("code_del", user_list_controller.DeleteCode)
	r.POST("use_code", user_list_controller.UseCode)
	r.GET("user_list", user_list_controller.UserListPage)
	r.GET("user_list_json", user_list_controller.UserListJson)
	r.POST("user_change_status", user_list_controller.UpdateUserStatus)
}
