package user_list_controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mybedv2/app/helper/e"
	"mybedv2/app/helper/util/str"
	"mybedv2/app/system/model"
	user2 "mybedv2/app/system/model/user"
	"mybedv2/app/system/service/img"
	"mybedv2/conf"
	"net/http"
	"strconv"
)

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
