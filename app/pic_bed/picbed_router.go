package pic_bed

import (
	"mybedv2/app/helper/router"
	"mybedv2/app/system/middleware"
	"mybedv2/app/system/service/upload_handler"
)

func init() {
	r := router.New("pic_bed", "/")
	r.GET("/", WebUpload)
	r.GET("/list", IndexList)
	r.GET("/up_page", SingleFile)
	r.POST("uploads", middleware.CheckConfigUpload, upload_handler.UploadHandler)
	r.POST("islogin", IsLogin)
	r.POST("login", LoginHandler)
	r.POST("register",middleware.CheckSiteReg,RegisterHandler)
	r.GET("logout", Logout)
	r.POST("logout", Logout)
}
