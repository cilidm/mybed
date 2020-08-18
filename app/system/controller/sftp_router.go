package controller

import (
	"mybedv2/app/helper/router"
	"mybedv2/app/system/controller/store_list_controller"
	"mybedv2/app/system/middleware"
)

func init() {
	r := router.New("system", "/system", middleware.AuthMiddleware, middleware.SidebarMiddleware)
	r.GET("sftp", store_list_controller.SftpPage)
	r.GET("storemenu_json", store_list_controller.StoreSystemStorage)
	r.POST("imgdata_del", store_list_controller.ImgdataDel)
	r.POST("imgdata_del_more", store_list_controller.ImgdataDelMore)
}
