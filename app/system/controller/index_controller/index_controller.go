package index_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mybedv2/app/helper/e"
	"mybedv2/app/helper/util/pathdir"
	"mybedv2/app/helper/util/str"
	"mybedv2/app/system/model"
	"mybedv2/app/system/model/blacklist/img_blacklist"
	imgModel "mybedv2/app/system/model/img"
	"mybedv2/app/system/model/site"
	"mybedv2/app/system/model/store"
	user2 "mybedv2/app/system/model/user"
	"mybedv2/app/system/service/img"
	"mybedv2/app/system/service/user"
	"net/http"
	"strconv"
)

func Index(c *gin.Context) {
	var (
		userEntity user2.Entity
		siteEntity site.Entity
	)
	authID := user.GetUserId(c)
	isAdmin := user.IsAdmin(authID)
	authInfo, err := userEntity.FindById(authID) //用户信息
	if err != nil {
		fmt.Println(err)
	}
	config := siteEntity.FindOne()                                        //站点配置
	imgPer := int(img.GetUserMemPer(int64(authInfo.ID), authInfo.Memory)) //用户空间使用率
	c.HTML(http.StatusOK, "index.html", gin.H{
		"comInfo":   user2.Info{ID: authInfo.ID, Nickname: authInfo.Nickname, Level: authInfo.Level, Avatar: authInfo.Avatar},
		"isAdmin":   isAdmin,
		"config":    config,
		"imgper":    imgPer,
		"storeMenu": store.GetStoreMenu()})
}

func GetIndexLine(c *gin.Context) {
	ls := new(model.Lines)
	ig, _ := imgModel.GetImgNumByDay()
	for _, v := range ig {
		ls.Days = append(ls.Days, v.Day)
		ls.Nums = append(ls.Nums, v.Num)
	}
	ug, _ := user2.GetUserIndexLine()
	for _, v := range ug {
		ls.UserDays = append(ls.UserDays, v.Day)
		ls.UserNums = append(ls.UserNums, v.Num)
	}
	c.JSON(http.StatusOK, ls)
}

func IndexFrame(c *gin.Context) {
	userId := user.GetUserId(c)
	totalSize, err := imgModel.GetImgdataSize(userId)
	if err != nil {
		totalSize = 0
	}
	imgBlack := img_blacklist.GetAllNum()
	c.HTML(http.StatusOK, "frame_index.html", gin.H{
		"imgNum":    imgModel.GetImgdataNum(),
		"userNum":   user2.GetUserNum(),
		"totalSize": str.SizeFormat(float64(totalSize)),
		"blackNum":  imgBlack,
	})
}

func PageNotFound(c *gin.Context) {
	c.HTML(http.StatusOK, "page_not_found.html", gin.H{})
}

func EditPwd(c *gin.Context) {
	c.HTML(http.StatusOK, "edit_pwd.html", gin.H{})
}

func EditPwdHandler(c *gin.Context) {
	var (
		ep         user2.EditPwdForm
		userEntity user2.Entity
	)
	if err := c.Bind(&ep); err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ERROR, ErrorMsg: "提交不全，请重新填写"})
		return
	}
	if ep.Newpwd != ep.Confirmpwd {
		c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ERROR, ErrorMsg: "两次密码不一致，请重新提交"})
		return
	}
	uid := user.GetUserId(c)
	u, err := userEntity.FindById(uid)
	if err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ERROR, ErrorMsg: err.Error()})
		return
	}
	if u.Password != str.Md5([]byte(ep.Oldpwd+u.Salt)) {
		c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ERROR, ErrorMsg: "原密码错误，请重新提交"})
		return
	}
	u.Password, u.Salt = str.Password(10, ep.Newpwd)
	user2.UpdateUserPwd(&u)
	user.Logout(c)
	c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.SUCCESS, Msg: e.GetMsg(e.SUCCESS)})
}

// 修改资料页面
func ProfilePage(c *gin.Context) {
	var userEntity user2.Entity
	uid := user.GetUserId(c)
	user, err := userEntity.FindById(uid)
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
	}
	var userShow user2.Info
	str.CopyFields(&userShow, user)
	used, err := imgModel.GetImgdataSize(uid)
	if err != nil {
		userShow.UsedMem = "0KB"
		userShow.UsedPer = 0
	} else {
		userShow.UsedMem = str.SizeFormat(float64(used))
		userShow.UsedPer = (int64(used) / (userShow.Memory * 1024 * 1024)) * 100
	}
	c.HTML(http.StatusOK, "profile.html", userShow)
}

func ProfilePageHandler(c *gin.Context) {
	email := c.PostForm("email")
	nickName := c.PostForm("nickname")
	uid := user.GetUserId(c)
	if err := user2.UpdateUserProfile(uid, email, nickName); err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ERROR, ErrorMsg: e.GetMsg(e.ERROR) + err.Error()})
		return
	}
	c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.SUCCESS, Msg: e.GetMsg(e.SUCCESS)})
}

func ProfileUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ErrorMsg: err.Error(), ResultCode: e.ERROR})
		return
	}
	uid := strconv.FormatInt(user.GetUserId(c), 10)
	saveDir := "static/upload/" + uid + "/avatar/"
	pathdir.PathExists(saveDir)
	c.SaveUploadedFile(file, saveDir+file.Filename)
	var siteEntity site.Entity
	conf := siteEntity.FindOne()
	user2.UpdateUserAvatar(uid, conf.WebUrl+saveDir+file.Filename)
	c.JSON(http.StatusOK, model.AjaxResp{Msg: conf.WebUrl + saveDir + file.Filename, ResultCode: e.SUCCESS})
}
