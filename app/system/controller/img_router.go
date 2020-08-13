package controller

import (
	"mybedv2/app/helper/router"
	"mybedv2/app/helper/util/qrcode"
	"mybedv2/app/system/controller/img_controller"
	"mybedv2/app/system/middleware"
)

func init() {
	r := router.New("system", "/system", middleware.AuthMiddleware)
	r.GET("img_list", img_controller.ImgListPage)
	r.GET("imglist_json", img_controller.ImgListJson)
	r.GET("qrcode", qrcode.CreateQrcode)
	r.POST("delImgList", img_controller.DelImgList)
	r.POST("delImg", img_controller.DelImg)
	r.GET("img_show", img_controller.ImgShowPage)
}
