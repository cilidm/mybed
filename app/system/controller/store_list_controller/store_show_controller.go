package store_list_controller

import (
	"github.com/gin-gonic/gin"
	"mybedv2/app/helper/e"
	"mybedv2/app/helper/util/str"
	"mybedv2/app/system/model"
	"mybedv2/app/system/model/img"
	"mybedv2/app/system/model/site"
	"mybedv2/app/system/model/store"
	"mybedv2/app/system/model/user"
	storeService "mybedv2/app/system/service/store"
	"net/http"
	"strconv"
	"strings"
)

func StoreSystem(c *gin.Context) {
	var siteEntity site.Entity
	storeType := c.Query("type")
	conf := siteEntity.FindOne()
	c.HTML(http.StatusOK, "store_system_list.html", gin.H{"type": storeType, "conf": conf})
}

//读取数据库存储的数据
func StoreSystemStorage(c *gin.Context) {
	storeType := c.Query("type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if storeType == "" {
		c.Redirect(http.StatusFound, "/system/page_not_found")
		return
	}
	storeNum := store.GetStoreNum(storeType)
	datas, count := img.GetImgdataBySource(storeNum, page, limit)
	var (
		storeMap   []map[string]interface{}
		userEntity user.Entity
	)
	for _, v := range datas {
		user, _ := userEntity.FindById(v.UserId)
		storeMap = append(storeMap, str.Struct2Map(store.ListShow{
			Id:        v.Id,
			ImgName:   v.ImgName,
			ImgUrl:    v.ImgUrl,
			User:      user.Username,
			Size:      str.SizeFormat(float64(v.Sizes)),
			Ip:        v.Abnormal,
			CreatedAt: v.CreatedAt.Format(e.TimeFormat),
		}))
	}
	c.JSON(http.StatusOK, model.LayuiResp{
		Code: 0, Count: count, Data: storeMap, Msg: "",
	})
}

func ImgdataDelMore(c *gin.Context) {
	ids := c.PostFormArray("ids")
	storeType := c.PostForm("type")
	for _, v := range ids {
		id, _ := strconv.Atoi(v)
		err := ImgdataDelHandler(id, storeType)
		if err != nil {
			c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ERROR, ErrorMsg: e.GetMsg(e.ERROR) + err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.SUCCESS, Msg: e.GetMsg(e.SUCCESS)})
}

func ImgdataDel(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("id"))
	storeType := c.PostForm("type")
	err := ImgdataDelHandler(id, storeType)
	if err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ERROR, ErrorMsg: e.GetMsg(e.ERROR) + err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.SUCCESS, Msg: e.GetMsg(e.SUCCESS)})
}

func ImgdataDelHandler(id int, storeType string) error {
	var (
		imgEntity img.Entity
		storePath string
	)
	img := img.GetImgdataById(id)
	storeConf := store.GetStoreConfig("cloud_type", storeType)
	cs, _ := storeService.NewCloudStoreByConf(storeConf, false)
	if storeType == "cs-minio" {
		storePath = strings.ReplaceAll(img.ImgUrl, storeConf.PublicBucketDomain+storeConf.PublicBucket, "")
	} else {
		storePath = strings.ReplaceAll(img.ImgUrl, storeConf.PublicBucketDomain, "")
	}
	if err := cs.Delete(storePath); err != nil {
		return err
	}
	if img.ImgThumb != "" {
		cs.Delete(strings.ReplaceAll(img.ImgThumb, storeConf.PublicBucketDomain, ""))
	}
	if err := imgEntity.Delete(id); err != nil {
		return err
	}
	return nil
}
