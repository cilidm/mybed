package controller

import (
	"mybedv2/app/helper/router"
	"mybedv2/app/system/controller/index_controller"
	"mybedv2/app/system/middleware"
)

func init() {
	r := router.New("system", "/system", middleware.AuthMiddleware)
	r.GET("/", index_controller.Index)
	r.GET("frame_index", index_controller.IndexFrame)
	r.GET("page_not_found", index_controller.PageNotFound)
	r.GET("profile", index_controller.ProfilePage)
	r.GET("edit_pwd", index_controller.EditPwd)
	r.POST("profile", index_controller.ProfilePageHandler)
	r.POST("profile_upload", index_controller.ProfileUpload)
	r.POST("edit_pwd", index_controller.EditPwdHandler)
}
