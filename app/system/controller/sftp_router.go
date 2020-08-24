package controller

import (
	"mybedv2/app/helper/router"
	"mybedv2/app/system/controller/sftp_controller"
	"mybedv2/app/system/middleware"
)

func init() {
	r := router.New("system", "/system", middleware.AuthMiddleware, middleware.SidebarMiddleware)
	r.GET("sftp", sftp_controller.SftpPage)
	r.GET("sftp_form", sftp_controller.SftpFormPage)
	r.POST("sftp_form", sftp_controller.SftpFormHandler)
	r.GET("sftp_json", sftp_controller.SftpJson)

	r.GET("local", sftp_controller.LocalPage)
	r.GET("local_form", sftp_controller.LocalFormPage)
	r.POST("local_form", sftp_controller.LocalFormHandler)
	r.GET("local_json", sftp_controller.LocalJson)
}
