package user

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"mybedv2/app/helper/e"
	"mybedv2/app/helper/util/calculate"
	"mybedv2/app/helper/util/pathdir"
	userModel "mybedv2/app/system/model/user"
	"strconv"
	"time"
)

//获取用户当前空间使用率
func GetUserMemPer(userId, mem int64) float64 {
	totalSize, err := userModel.GetImgdataSize(userId)
	if err != nil {
		return 0
	}
	return calculate.Decimal((float64(totalSize) / float64(mem*1024*1024)) * 100) //默认单位mb 这里要转换成byte
}

func GetUserTmpImgName(c *gin.Context) string {
	userId := GetUserId(c)
	tmpName := strconv.FormatInt(userId, 10) + "_" + time.Now().Format(e.TimeFormatDir) + ".jpg"
	return tmpName
}

func GetUserTmpImgDir(c *gin.Context) string {
	userId := GetUserId(c)
	tmpDir := "static/upload/" + strconv.FormatInt(userId, 10) + "/" + time.Now().Format(e.DayTimeFormatDir)
	pathdir.PathExists(tmpDir)
	return tmpDir
}

// 保存用户信息到session
func SaveUserToSession(user userModel.Entity, c *gin.Context) error {
	session := sessions.Default(c)
	session.Set(e.UserId, user.ID)
	if IsAdmin(int64(user.ID)) {
		user.Level = 99
	}
	tmp, _ := json.Marshal(user)
	session.Set(e.UserSessionInfo, string(tmp))
	return session.Save()
}

// 判断是否是系统管理员
func IsAdmin(userId int64) bool {
	var userEntity userModel.Entity
	u, err := userEntity.FindById(userId)
	if err != nil {
		return false
	}
	if u.Level == 99 {
		return true
	} else {
		return false
	}
}

func GetUserId(c *gin.Context) (authId int64) {
	s := sessions.Default(c)
	authStr := s.Get(e.UserId)
	if authStr == nil {
		authStr = 0
	}
	if _, ok := authStr.(int); ok == true {
		authId = int64(authStr.(int))
	} else {
		authId = 0
	}
	return authId
}

//判断是否登陆
func IsLogin(c *gin.Context) bool {
	session := sessions.Default(c)
	uid := session.Get(e.UserId)
	if uid != nil && uid != 0 {
		return true
	}
	return false
}

func Logout(c *gin.Context) error {
	session := sessions.Default(c)
	session.Delete(e.UserId)
	session.Delete(e.UserSessionInfo)
	return session.Save()
}
