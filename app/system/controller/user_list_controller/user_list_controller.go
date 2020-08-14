package user_list_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mybedv2/app/helper/e"
	"mybedv2/app/helper/pkg/xid"
	"mybedv2/app/helper/util/snowFlake"
	"mybedv2/app/helper/util/str"
	"mybedv2/app/system/model"
	"mybedv2/app/system/model/code"
	user2 "mybedv2/app/system/model/user"
	"mybedv2/app/system/service/img"
	"mybedv2/app/system/service/user"
	"mybedv2/conf"
	"net/http"
	"strconv"
	"time"
)

func CodeListPage(c *gin.Context) {
	c.HTML(http.StatusOK, "code_list.html", gin.H{})
}

func CodeFormPage(c *gin.Context) {
	c.HTML(http.StatusOK, "code_form.html", gin.H{})
}

func UseCode(c *gin.Context)  {
	codeForm := c.PostForm("code")
	codeInfo,err := code.FindCode(codeForm)
	if err != nil{
		c.JSON(http.StatusOK,model.AjaxResp{ResultCode: e.ERROR,Msg: e.GetMsg(e.ERROR) + ": 未找到此扩容码或已被使用"})
		return
	}
	uid := user.GetUserId(c)
	u,err := user2.FindUserById(uid)
	if err != nil{
		c.JSON(http.StatusOK,model.AjaxResp{ResultCode: e.ERROR,Msg: e.GetMsg(e.ERROR) + err.Error()})
		return
	}
	u.Memory = u.Memory + int64(codeInfo.Value)
	err = user2.UpdateMem(u)
	if err != nil{
		c.JSON(http.StatusOK,model.AjaxResp{ResultCode: e.ERROR,Msg: e.GetMsg(e.ERROR) + err.Error()})
		return
	}
	codeInfo.UserId = uid
	codeInfo.Status = 2
	codeInfo.UpdatedAt = time.Now()
	if err := code.UpdateStatus(codeInfo);err != nil{
		c.JSON(http.StatusOK,model.AjaxResp{ResultCode: e.ERROR,Msg: e.GetMsg(e.ERROR) + err.Error()})
		return
	}
	c.JSON(http.StatusOK,model.AjaxResp{ResultCode: e.SUCCESS,Msg: e.GetMsg(e.SUCCESS)})
}

func DeleteCodes(c *gin.Context)  {
	ids := c.PostFormArray("ids")
	if err := code.DeleteByIds(ids);err != nil{
		c.JSON(http.StatusOK,model.AjaxResp{ResultCode: e.ERROR,Msg: e.GetMsg(e.ERROR) + err.Error()})
		return
	}
	c.JSON(http.StatusOK,model.AjaxResp{ResultCode: e.SUCCESS,Msg: e.GetMsg(e.SUCCESS)})
}

func DeleteCode(c *gin.Context)  {
	id := c.PostForm("id")
	if err := code.DeleteById(id);err != nil{
		c.JSON(http.StatusOK,model.AjaxResp{ResultCode: e.ERROR,Msg: e.GetMsg(e.ERROR) + err.Error()})
		return
	}
	c.JSON(http.StatusOK,model.AjaxResp{ResultCode: e.SUCCESS,Msg: e.GetMsg(e.SUCCESS)})
}

func CreateCode(c *gin.Context)  {
	codeForm := new(code.BindForm)
	if err := c.Bind(&codeForm);err != nil{
		c.JSON(http.StatusOK,model.AjaxResp{ResultCode: e.ERROR,Msg: "请填写数量及容量"})
		return
	}
	for i:=0;i<codeForm.CodeNum;i++ {
		snowId,_ := snowFlake.Snow.GetSnowflakeId()
		err := code.InsertCode(code.InsertForm{Code: xid.GetXid() + strconv.FormatInt(snowId,10) ,Value: codeForm.CodeValue})
		if err != nil{
			c.JSON(http.StatusOK,model.AjaxResp{ResultCode: e.ERROR,Msg: e.GetMsg(e.ERROR) + err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK,model.AjaxResp{ResultCode: e.SUCCESS,Msg: e.GetMsg(e.SUCCESS)})
}

func CodeListJson(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(conf.Setting.PageSize)))
	list, count := code.GetCodeList(page, limit)
	var codeMap []map[string]interface{}
	for _, v := range list {
		u, err := user2.FindUserById(int64(v.UserId))
		if err != nil {
			if v.UserId == 0{
				u.Username = "未使用"
			}else{
				u.Username = "未找到此用户"
			}
		}
		useTime := ""
		if v.CreatedAt != v.UpdatedAt{
			useTime = v.UpdatedAt.Format(e.TimeFormat)
		}
		codeMap = append(codeMap, str.Struct2Map(code.ListJson{
			Id:        v.Id,
			Code:      v.Code,
			Value:     v.Value,
			Status:    v.Status,
			UserName:  u.Username,
			CreatedAt: v.CreatedAt.Format(e.TimeFormat),
			UpdatedAt: useTime,
		}))
	}
	c.JSON(http.StatusOK, model.LayuiResp{Code: 0, Count: count, Data: codeMap, Msg: ""})
}

func UserListPage(c *gin.Context) {
	c.HTML(http.StatusOK, "user_list.html", gin.H{})
}

func UserListJson(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(conf.Setting.PageSize)))
	user, count := user2.GetUserList(page, limit)
	var userMap []map[string]interface{}
	for _, v := range user {
		memPer := img.GetUserMemPer(int64(v.ID), v.Memory)
		userMap = append(userMap, str.Struct2Map(user2.ListShow{
			Id:            int(v.ID),
			Username:      v.Username,
			Status:        v.Status,
			Email:         v.Email,
			Avatar:        v.Avatar,
			Nickname:      v.Nickname,
			MemoryUsed:    int(memPer),
			Level:         v.Level,
			CreatedAt:     v.CreatedAt.Format(e.TimeFormat),
			LastLoginTime: v.LastLoginTime.Format(e.TimeFormat),
			LastLoginIp:   v.LastLoginIp,
		}))
	}
	c.JSON(http.StatusOK, model.LayuiResp{Code: 0, Count: count, Data: userMap, Msg: ""})
}

func UpdateUserStatus(c *gin.Context) {
	uid := c.PostForm("uid")
	intId, _ := strconv.ParseInt(uid, 10, 64)
	var userEntity user2.Entity
	user, err := userEntity.FindById(intId)
	fmt.Println(user)
	if user.ID == 0 || err != nil {
		c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.ERROR, ErrorMsg: "出错！请重试" + err.Error()})
		return
	}
	status := c.PostForm("status")
	user2.UpdateUserStatus(uid, status)
	c.JSON(http.StatusOK, model.AjaxResp{ResultCode: e.SUCCESS, Msg: e.GetMsg(e.SUCCESS)})
}
