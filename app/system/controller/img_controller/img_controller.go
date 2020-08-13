package img_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mybedv2/app/helper/util/page"
	"mybedv2/app/helper/util/str"
	"mybedv2/app/system/model"
	"mybedv2/app/system/model/img"
	"mybedv2/app/system/model/site"
	"mybedv2/app/system/model/store"
	storeService "mybedv2/app/system/service/store"
	"mybedv2/app/system/service/user"
	"mybedv2/conf"
	"net/http"
	"strconv"
	"strings"
)

func ImgListPage(c *gin.Context) {
	var siteEntity site.Entity
	conf := siteEntity.FindOne()
	c.HTML(http.StatusOK, "imgdata_list.html", gin.H{"conf": conf})
}

//图片列表
func ImgListJson(c *gin.Context) {
	var imgShow []map[string]interface{}
	page, _ := strconv.Atoi(strings.TrimSpace(c.DefaultQuery("page", "1")))
	limit := strings.TrimSpace(c.DefaultQuery("limit", "10"))
	imgData, count := img.GetImgDataByLimit(page, limit, user.GetUserId(c))
	for _, v := range imgData {
		imgShow = append(imgShow, str.Struct2Map(img.List{
			Id:        v.Id,
			ImgUrl:    v.ImgUrl,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
			Sizes:     str.SizeFormat(float64(v.Sizes)),
			Abnormal:  v.Abnormal,
		}))
	}
	c.JSON(http.StatusOK, model.LayuiResp{Code: 0, Msg: "", Count: count, Data: imgShow})
}

func ImgShowPage(c *gin.Context) {
	page := page.GetQueryPage(c)
	img, count := img.GetImgData(page, user.GetUserId(c))
	pageCount := (count / conf.Setting.PageSize) + 1
	var imgStr []string
	for _, v := range img {
		imgStr = append(imgStr, v.ImgUrl)
	}
	c.HTML(http.StatusOK, "img_gallery.html", gin.H{"img": imgStr, "pageCount": pageCount, "page": page})
}

func DelImg(c *gin.Context) {
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ErrorMsg: err.Error(), ResultCode: -1})
		return
	}
	img := img.GetImgdataById(id)
	if err := DelImgDomain(img); err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ErrorMsg: err.Error(), ResultCode: -1})
		return
	}
	c.JSON(http.StatusOK, model.AjaxResp{ResultCode: 200})
}

func DelImgList(c *gin.Context) {
	ids := c.PostFormArray("ids")
	imgs := img.GetImgdataByArr(ids)
	for _, v := range imgs {
		if err := DelImgDomain(v); err != nil {
			c.JSON(http.StatusOK, model.AjaxResp{ErrorMsg: err.Error(), ResultCode: -1})
			return
		}
	}
	c.JSON(http.StatusOK, model.AjaxResp{ResultCode: 200})
}

func DelImgDomain(v img.Entity) error {
	var imgEntity img.Entity
	cloudType := store.GetStoreStr(v.Source)
	csc := store.GetStoreConfig("cloud_type", cloudType) //todo 每次都要查数据库里store配置 修改为查询缓存？
	config := storeService.GetConfigType(csc)
	private := false
	cs, err := storeService.NewCloudStoreWithConfig(config, csc.CloudType, private) //todo 修改为创建相应的连接池，直接根据不同store调用
	if err != nil {
		fmt.Println(err)
		return err
	}
	var file string
	if cloudType == "cs-minio" {
		file = strings.ReplaceAll(v.ImgUrl, csc.PublicBucketDomain+csc.PublicBucket, "")
	} else {
		file = strings.ReplaceAll(v.ImgUrl, csc.PublicBucketDomain, "")
	}
	cs.Delete(file)              //删除store原图
	err = imgEntity.Delete(v.Id) //删除数据库
	if err != nil {
		return err
	}
	return nil
}
