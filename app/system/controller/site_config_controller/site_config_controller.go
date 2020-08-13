package site_config_controller

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mybedv2/app/helper/e"
	"mybedv2/app/helper/pkg/gomail"
	"mybedv2/app/system/model"
	"mybedv2/app/system/model/bd/bd_img"
	"mybedv2/app/system/model/mail"
	"mybedv2/app/system/model/site"
	store2 "mybedv2/app/system/model/store"
	"mybedv2/app/system/model/upload"
	"mybedv2/app/system/service/store"
	"mybedv2/app/system/service/user"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func SiteConfig(c *gin.Context) {
	var siteEntity site.Entity
	config := siteEntity.FindOne()
	c.HTML(http.StatusOK, "site_config.html", config)
}

func SiteConfigHandler(c *gin.Context) {
	var (
		sc         site.BindForm
		siteEntity site.Entity
	)
	if err := c.ShouldBind(&sc); err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ERROR, ErrorMsg: e.GetMsg(e.ERROR) + err.Error()})
		return
	}
	emailStatus := c.PostForm("site_status")
	if emailStatus == "on" {
		sc.SiteStatus = 1
	} else {
		sc.SiteStatus = 2
	}
	if err := siteEntity.CreateOrUpdate(sc); err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ERROR, ErrorMsg: e.GetMsg(e.ERROR) + err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.SUCCESS, Msg: e.GetMsg(e.SUCCESS)})
}

func SiteUpload(c *gin.Context) {
	var uploadEntity upload.Entity
	config := uploadEntity.FindOne()
	c.HTML(http.StatusOK, "site_upload.html", config)
}

func SiteUploadHandler(c *gin.Context) {
	var (
		uc           upload.BindForm
		uploadEntity upload.Entity
	)
	if err := c.ShouldBind(&uc); err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ERROR, ErrorMsg: e.GetMsg(e.ERROR) + err.Error()})
		return
	}
	if c.PostForm("allow_visitor_button") == "on" {
		uc.AllowVisitor = 1
	} else {
		uc.AllowVisitor = 2
	}
	if err := uploadEntity.CreateOrUpdate(uc); err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ERROR, ErrorMsg: e.GetMsg(e.ERROR) + err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.SUCCESS, Msg: e.GetMsg(e.SUCCESS)})
}

func SiteExamine(c *gin.Context) {
	config := bd_img.FindOne()
	c.HTML(http.StatusOK, "site_examine.html", config)
}

func SiteExamineHandler(c *gin.Context) {
	var (
		bdForm      bd_img.BindForm
		bdImgEntity bd_img.Entity
	)
	if err := c.ShouldBind(&bdForm); err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ERROR, ErrorMsg: e.GetMsg(e.ERROR) + err.Error()})
		return
	}
	if c.PostForm("status_form") == "on" {
		bdForm.Status = 1
	} else {
		bdForm.Status = 2
	}
	if err := bdImgEntity.CreateOrUpdate(bdForm); err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ERROR, ErrorMsg: e.GetMsg(e.ERROR) + err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.SUCCESS, Msg: e.GetMsg(e.SUCCESS)})
}

func SiteMail(c *gin.Context) {
	var mailEntity mail.Entity
	config := mailEntity.FindOne()
	c.HTML(http.StatusOK, "site_email.html", config)
}

// 邮箱
func SiteMailHandler(c *gin.Context) {
	var (
		ec         mail.BindForm
		mailEntity mail.Entity
	)
	if err := c.ShouldBind(&ec); err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ERROR, ErrorMsg: e.GetMsg(e.ERROR) + err.Error()})
		return
	}
	emailStatus := c.PostForm("email_status")
	if emailStatus == "on" {
		ec.EmailStatus = 1
	} else {
		ec.EmailStatus = 2
	}
	if err := mailEntity.CreateOrUpdate(ec); err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ERROR, ErrorMsg: e.GetMsg(e.ERROR) + err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.SUCCESS, Msg: e.GetMsg(e.SUCCESS)})
}

func SiteMailTestMailHandler(c *gin.Context) {
	var ec mail.BindForm
	if err := c.ShouldBind(&ec); err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ERROR, ErrorMsg: e.GetMsg(e.ERROR) + err.Error()})
		return
	}
	var testMail gomail.Config
	testMail.Config = ec
	testMail.MailTo = append(testMail.MailTo, ec.EmailTest)
	testMail.Subject = "测试邮件"
	if err := gomail.SendMailV1(testMail); err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ERROR, ErrorMsg: e.GetMsg(e.ERROR) + err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.SUCCESS, Msg: e.GetMsg(e.SUCCESS)})
}

func UploadImgBase64(c *gin.Context) {
	imgData := c.PostForm("img_data")
	if imgData == "" {
		c.JSON(http.StatusOK, model.AjaxResp{ErrorMsg: e.GetMsg(e.ErrorUploadImgBase64NullByte), ResultCode: e.ErrorUploadImgBase64NullByte})
		return
	}
	imgData = strings.ReplaceAll(imgData, "data:image/jpeg;base64,", "")
	imgBuf, _ := base64.StdEncoding.DecodeString(imgData) //拼接成图片文件并把文件写入到buffer
	userTmpFile := path.Join(user.GetUserTmpImgDir(c), user.GetUserTmpImgName(c))
	err := ioutil.WriteFile(userTmpFile, imgBuf, 0666) //buffer输出到jpg文件中（不做处理，直接写到文件）
	if err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ErrorMsg: e.GetMsg(e.ErrorUploadImgBase64Save) + "," + err.Error(), ResultCode: e.ErrorUploadImgBase64Save})
		return
	}
	storePath, err := store.DefaultUploadStore(userTmpFile, c)
	if err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ErrorMsg: e.GetMsg(e.ErrorUploadStore) + "," + err.Error(), ResultCode: e.ErrorUploadStore})
		return
	}
	defer os.Remove(userTmpFile)
	var (
		storeEntity store2.Entity
		url         string
	)
	csc := storeEntity.FindOne()
	if csc.CloudType == "cs-minio" {
		url = csc.PublicBucketDomain + filepath.Join(csc.PublicBucket+storePath)
	} else {
		url = csc.PublicBucketDomain + storePath
	}
	c.JSON(http.StatusOK, model.AjaxResp{
		ResultCode: e.SUCCESS,
		Url:        url,
	})
}
