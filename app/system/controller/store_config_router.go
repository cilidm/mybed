package controller

import (
	"mybedv2/app/helper/router"
	"mybedv2/app/system/controller/store_config_controller"
	"mybedv2/app/system/middleware"
)

func init() {
	r := router.New("system", "/system", middleware.AuthMiddleware, middleware.SidebarMiddleware)
	r.GET("store_config_minio", store_config_controller.MinioConfigPage)
	r.GET("store_config_oss", store_config_controller.OssConfigPage)
	r.GET("store_config_cos", store_config_controller.CosConfigPage)
	r.GET("store_config_bos", store_config_controller.BosConfigPage)
	r.GET("store_config_obs", store_config_controller.ObsConfigPage)
	r.GET("store_config_qiniu", store_config_controller.QiniuConfigPage)
	r.GET("store_config_upyun", store_config_controller.UpYunConfigPage)
	r.POST("store_config", store_config_controller.StoreConfigSave)
}
