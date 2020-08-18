package sftp_controller

import (
	"github.com/gin-gonic/gin"
	"mybedv2/app/helper/e"
	"mybedv2/app/helper/pkg/sftp"
	"mybedv2/app/system/model"
	"mybedv2/conf"
	"net/http"
)

func SftpPage(c *gin.Context) {
	c.HTML(http.StatusOK, "sftp.html", gin.H{})
}

func SftpFormPage(c *gin.Context) {
	c.HTML(http.StatusOK, "sftp_form.html", gin.H{})
}

func SftpFormHandler(c *gin.Context) {
	var b sftp.SftpBindForm
	if err := c.Bind(&b); err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ERROR, Msg: e.GetMsg(e.ERROR)})
		return
	}
	err := sftp.NewSftp(sftp.SftpConf{
		Config: sftp.Config{
			Host: b.Host,
			User: b.User,
			Pwd:  b.Pwd,
			Port: b.Port,
		},
		Source: b.Source,
		Target: conf.Setting.UploadTmpDir,
	})
	if err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ERROR, Msg: e.GetMsg(e.ERROR)})
		return
	}
	c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.SUCCESS, Msg: e.GetMsg(e.SUCCESS)})
}
