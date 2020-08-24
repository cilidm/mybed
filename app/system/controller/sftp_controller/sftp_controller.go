package sftp_controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"mybedv2/app/helper/e"
	"mybedv2/app/helper/pkg/sftp"
	"mybedv2/app/helper/redis"
	"mybedv2/app/helper/util/str"
	"mybedv2/app/system/model"
	"mybedv2/app/system/model/img"
	"mybedv2/app/system/model/store"
	"mybedv2/app/system/service/user"
	"mybedv2/conf"
	"net/http"
	"strconv"
)

func LocalPage(c *gin.Context) {
	c.HTML(http.StatusOK, "sftp.html", gin.H{})
}

func LocalJson(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	datas, count := img.GetImgdataSftp(page, limit)
	var dataMap []map[string]interface{}
	for _, v := range datas {
		dataMap = append(dataMap, str.Struct2Map(img.PageJson{
			Id:        v.Id,
			ImgName:   v.ImgName,
			ImgUrl:    v.ImgUrl,
			Sizes:     str.SizeFormat(float64(v.Sizes)),
			Source:    store.GetStoreNameByCloudType(store.GetStoreStr(v.Source)),
			CreatedAt: v.CreatedAt.Format(e.TimeFormat),
		}))
	}
	c.JSON(http.StatusOK, model.LayuiResp{
		Code: 0, Count: count, Data: dataMap, Msg: "",
	})
}

func LocalFormPage(c *gin.Context) {
	var b sftp.SftpBindForm
	confJson, _ := redis.Client.Get("local:form:" + strconv.Itoa(int(user.GetUserId(c)))).Result()
	if confJson != "" {
		json.Unmarshal([]byte(confJson), &b)
	}
	c.HTML(http.StatusOK, "local_form.html", gin.H{"conf": b})
}

func LocalFormHandler(c *gin.Context) {
	var b sftp.SftpBindForm
	if err := c.Bind(&b); err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ERROR, Msg: e.GetMsg(e.ERROR)})
		return
	}
	switchBtn := c.PostForm("keep_local")
	if switchBtn == "on" {
		b.KeepLocal = 1
	} else {
		b.KeepLocal = 2
	}
	formJson, _ := json.Marshal(b)
	uid := strconv.Itoa(int(user.GetUserId(c)))
	redis.Client.Set("sftp:form:"+uid, formJson, 0)
	err := sftp.NewSftp(sftp.SftpConf{
		Config: sftp.Config{
			Host: b.Host,
			User: b.User,
			Pwd:  b.Pwd,
			Port: b.Port,
		},
		Source:    b.Source,
		Target:    conf.Setting.UploadTmpDir + "/sftp/" + uid,
		KeepLocal: b.KeepLocal,
		Uid:       strconv.Itoa(int(user.GetUserId(c))),
		ClientIP:  c.ClientIP(),
	})
	if err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ERROR, Msg: e.GetMsg(e.ERROR)})
		return
	}
	c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.SUCCESS, Msg: e.GetMsg(e.SUCCESS)})
}

func SftpPage(c *gin.Context) {
	c.HTML(http.StatusOK, "sftp.html", gin.H{})
}

func SftpJson(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	datas, count := img.GetImgdataSftp(page, limit)
	var dataMap []map[string]interface{}
	for _, v := range datas {
		dataMap = append(dataMap, str.Struct2Map(img.PageJson{
			Id:        v.Id,
			ImgName:   v.ImgName,
			ImgUrl:    v.ImgUrl,
			Sizes:     str.SizeFormat(float64(v.Sizes)),
			Source:    store.GetStoreNameByCloudType(store.GetStoreStr(v.Source)),
			CreatedAt: v.CreatedAt.Format(e.TimeFormat),
		}))
	}
	c.JSON(http.StatusOK, model.LayuiResp{
		Code: 0, Count: count, Data: dataMap, Msg: "",
	})
}

func SftpFormPage(c *gin.Context) {
	var b sftp.SftpBindForm
	confJson, _ := redis.Client.Get("sftp:form:" + strconv.Itoa(int(user.GetUserId(c)))).Result()
	if confJson != "" {
		json.Unmarshal([]byte(confJson), &b)
	}
	c.HTML(http.StatusOK, "sftp_form.html", gin.H{"conf": b})
}

func SftpFormHandler(c *gin.Context) {
	var b sftp.SftpBindForm
	if err := c.Bind(&b); err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ERROR, Msg: e.GetMsg(e.ERROR)})
		return
	}
	switchBtn := c.PostForm("keep_local")
	if switchBtn == "on" {
		b.KeepLocal = 1
	} else {
		b.KeepLocal = 2
	}
	formJson, _ := json.Marshal(b)
	uid := strconv.Itoa(int(user.GetUserId(c)))
	redis.Client.Set("sftp:form:"+uid, formJson, 0)
	err := sftp.NewSftp(sftp.SftpConf{
		Config: sftp.Config{
			Host: b.Host,
			User: b.User,
			Pwd:  b.Pwd,
			Port: b.Port,
		},
		Source:    b.Source,
		Target:    conf.Setting.UploadTmpDir + "/sftp/" + uid,
		KeepLocal: b.KeepLocal,
		Uid:       strconv.Itoa(int(user.GetUserId(c))),
		ClientIP:  c.ClientIP(),
	})
	if err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ERROR, Msg: e.GetMsg(e.ERROR)})
		return
	}
	c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.SUCCESS, Msg: e.GetMsg(e.SUCCESS)})
}
