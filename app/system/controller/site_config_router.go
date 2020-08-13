package controller

import (
	"mybedv2/app/helper/router"
	"mybedv2/app/system/controller/site_config_controller"
	"mybedv2/app/system/middleware"
)

func init() {
	r := router.New("system", "/system", middleware.AuthMiddleware, middleware.SidebarMiddleware)
	r.GET("site_config", site_config_controller.SiteConfig)
	r.POST("site_config", site_config_controller.SiteConfigHandler)
	r.GET("site_upload", site_config_controller.SiteUpload)
	r.POST("site_upload", site_config_controller.SiteUploadHandler)
	r.GET("site_mail", site_config_controller.SiteMail)
	r.POST("site_mail", site_config_controller.SiteMailHandler)
	r.POST("site_test_mail", site_config_controller.SiteMailTestMailHandler)
	r.POST("upload_base64", site_config_controller.UploadImgBase64)
	r.GET("site_examine", site_config_controller.SiteExamine)
	r.POST("site_examine", site_config_controller.SiteExamineHandler)
}
