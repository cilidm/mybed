package store_config_controller

import (
	"github.com/gin-gonic/gin"
	"mybedv2/app/helper/e"
	store2 "mybedv2/app/system/model/store"
	"net/http"
	"strings"
)

func MinioConfigPage(c *gin.Context) {
	store := store2.GetStoreConfig("cloud_type", "cs-minio")
	c.HTML(http.StatusOK, "store_config_minio.html", gin.H{"store": store, "activeMenu": "minio"})
}

func OssConfigPage(c *gin.Context) {
	store := store2.GetStoreConfig("cloud_type", "cs-oss")
	c.HTML(http.StatusOK, "store_config_oss.html", gin.H{"store": store, "activeMenu": "oss"})
}

func CosConfigPage(c *gin.Context) {
	store := store2.GetStoreConfig("cloud_type", "cs-cos")
	c.HTML(http.StatusOK, "store_config_cos.html", gin.H{"store": store, "activeMenu": "cos"})
}

func BosConfigPage(c *gin.Context) {
	store := store2.GetStoreConfig("cloud_type", "cs-bos")
	c.HTML(http.StatusOK, "store_config_bos.html", gin.H{"store": store, "activeMenu": "bos"})
}

func ObsConfigPage(c *gin.Context) {
	store := store2.GetStoreConfig("cloud_type", "cs-obs")
	c.HTML(http.StatusOK, "store_config_obs.html", gin.H{"store": store, "activeMenu": "obs"})
}

func QiniuConfigPage(c *gin.Context) {
	store := store2.GetStoreConfig("cloud_type", "cs-qiniu")
	c.HTML(http.StatusOK, "store_config_qiniu.html", gin.H{"store": store, "activeMenu": "qiniu"})
}

func UpYunConfigPage(c *gin.Context) {
	store := store2.GetStoreConfig("cloud_type", "cs-upyun")
	c.HTML(http.StatusOK, "store_config_upyun.html", gin.H{"store": store, "activeMenu": "upyun"})
}

func StoreConfigSave(c *gin.Context) {
	var (
		sc          store2.BindForm
		storeEntity store2.Entity
	)
	if err := c.ShouldBind(&sc); err != nil {
		c.JSON(http.StatusOK, gin.H{"resultCode": e.ERROR, "errorMsg": "参数不足，请重新填写" + err.Error()})
		return
	}
	//规定PublicBucketDomain http开头 /结尾
	if strings.HasPrefix(sc.PublicBucketDomain, "http") == false {
		sc.PublicBucketDomain = "http://" + sc.PublicBucketDomain
	}
	if strings.HasSuffix(sc.PublicBucketDomain, "/") == false {
		sc.PublicBucketDomain = sc.PublicBucketDomain + "/"
	}
	status := c.PostForm("status")
	if status == "on" {
		store2.ChangeOtherStoreStatus(sc.CloudType)
		sc.Status = 1
	} else {
		sc.Status = 2
	}
	if err := storeEntity.CreateOrUpdate(sc); err != nil {
		c.JSON(http.StatusOK, gin.H{"resultCode": e.ERROR, "errorMsg": "配置保存失败"})
	}
	c.JSON(http.StatusOK, gin.H{"resultCode": 200})
}
