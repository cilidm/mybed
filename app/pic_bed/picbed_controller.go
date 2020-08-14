package pic_bed

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mybedv2/app/helper/e"
	"mybedv2/app/helper/util/str"
	"mybedv2/app/system/model"
	"mybedv2/app/system/model/site"
	"mybedv2/app/system/model/upload"
	user2 "mybedv2/app/system/model/user"
	uploadservice "mybedv2/app/system/service/upload"
	"mybedv2/app/system/service/user"
	"net/http"
	"strings"
	"time"
)

// 前端页面返回
type PicbedPageResp struct {
	Config  site.Entity       `json:"config"`
	Upload  upload.PageConfig `json:"upload"`
	Islogin bool              `json:"islogin"`
}

func LoginPage(c *gin.Context) {
	var (
		site site.Entity
	)
	config := site.FindOne()
	c.HTML(http.StatusOK, "login.html", gin.H{"config": config})
}

func WebUpload(c *gin.Context) {
	var (
		up   upload.Entity
		upc  upload.PageConfig
		site site.Entity
	)
	uploadConfig := up.FindOne()
	if user.IsLogin(c) {
		upc = uploadservice.GetMemConfig(uploadConfig, 1) //获取登陆用户的图片空间配置
	} else {
		upc = uploadservice.GetVisConfig(uploadConfig, 1) //获取未登陆用户的图片空间配置
	}
	config := site.FindOne()
	c.HTML(http.StatusOK, "upload_page.html", PicbedPageResp{
		Islogin: user.IsLogin(c),
		Upload:  upc,
		Config:  config,
	})
}

func IndexList(c *gin.Context) {
	var (
		up   upload.Entity
		upc  upload.PageConfig
		site site.Entity
	)
	uploadConfig := up.FindOne()
	if user.IsLogin(c) {
		upc = uploadservice.GetMemConfig(uploadConfig, 2) //获取登陆用户的图片空间配置
	} else {
		upc = uploadservice.GetVisConfig(uploadConfig, 2) //获取未登陆用户的图片空间配置
	}
	fmt.Println(upc)
	config := site.FindOne()
	c.HTML(http.StatusOK, "index_list.html", PicbedPageResp{
		Islogin: user.IsLogin(c),
		Upload:  upc,
		Config:  config,
	})
}

func SingleFile(c *gin.Context) {
	var (
		up   upload.Entity
		upc  upload.PageConfig
		site site.Entity
	)
	uploadConfig := up.FindOne()
	if user.IsLogin(c) {
		upc = uploadservice.GetMemConfig(uploadConfig, 2) //获取登陆用户的图片空间配置
	} else {
		upc = uploadservice.GetVisConfig(uploadConfig, 2) //获取未登陆用户的图片空间配置
	}
	fmt.Println(upc)
	config := site.FindOne()
	//c.HTML(http.StatusOK, "picbed.html", PicbedPageResp{Upload: upc,Config: config})
	c.HTML(http.StatusOK, "index_new.html", PicbedPageResp{
		Islogin: user.IsLogin(c),
		Upload:  upc,
		Config:  config,
	})
}

func IsLogin(c *gin.Context) {
	if user.IsLogin(c) {
		c.JSON(http.StatusOK, model.IsLogin{Lgoinret: 1})
	} else {
		c.JSON(http.StatusOK, gin.H{})
	}
}

func LoginHandler(c *gin.Context) {
	userEntity := user2.Entity{}
	username := strings.TrimSpace(c.PostForm("username"))
	password := strings.TrimSpace(c.PostForm("password"))
	if username != "" && password != "" {
		u, err := userEntity.FindByName(username)
		if err != nil || u.Password != str.Md5([]byte(password+u.Salt)) {
			c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ErrorLoginCheckPwd, ErrorMsg: e.GetMsg(e.ErrorLoginCheckPwd)})
			return
		} else if u.Status == 2 {
			c.JSON(http.StatusOK, model.AjaxResp{ErrorMsg: e.GetMsg(e.ErrorLoginCheckStatus), ResultCode: e.ErrorLoginCheckStatus})
			return
		} else {
			u.LastLoginIp = c.ClientIP()
			u.LastLoginTime = time.Now()
			userEntity.Insert(u)
			//util.Cac.Set("uid"+strconv.Itoa(int(u.ID)), u, cache.DefaultExpiration)
			if err := user.SaveUserToSession(u, c); err != nil {
				c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ErrorSetSession, ErrorMsg: e.GetMsg(e.ErrorSetSession)})
			}
			c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.SUCCESS, Url: "/"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, model.AjaxResp{ResultCode: e.ErrorLoginCheckRequired, ErrorMsg: e.GetMsg(e.ErrorLoginCheckRequired)})
		return
	}
}

func RegisterHandler(c *gin.Context) {
	var (
		user       user2.Entity
		userEntity user2.Entity
	)
	user.Username = strings.TrimSpace(c.PostForm("username"))
	pwd := strings.TrimSpace(c.PostForm("password"))
	user.Nickname = user.Username
	user.Email = strings.TrimSpace(c.PostForm("email"))
	has, err := userEntity.FindByName(user.Username)
	if err == nil && has.ID > 0 {
		c.JSON(http.StatusOK, model.RegisterResp{Ret: -2, Zctype: 2})
		return
	}
	user.Password, user.Salt = str.Password(15, pwd)
	user.LastLoginIp = c.ClientIP()
	user.Level = 1
	if err := userEntity.Insert(user); err != nil {
		c.JSON(http.StatusOK, model.RegisterResp{Ret: -5, Zctype: 2, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.RegisterResp{Ret: 1, Zctype: 2})
}

func Logout(c *gin.Context) {
	if c.Request.Method == "POST" {
		if user.IsLogin(c) {
			if err := user.Logout(c); err != nil {
				c.JSON(http.StatusOK, gin.H{"code": e.ERROR, "msg": err.Error()})
			}
		}
		c.JSON(http.StatusOK, gin.H{"code": e.SUCCESS})
	} else {
		if user.IsLogin(c) {
			if err := user.Logout(c); err != nil {
				c.String(http.StatusOK, err.Error())
				return
			}
		}
		c.Redirect(http.StatusFound, "/")
	}
}
