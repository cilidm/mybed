package middleware

import (
	"github.com/gin-gonic/gin"
	"mybedv2/app/helper/e"
	"mybedv2/app/system/model/upload"
	"mybedv2/app/system/service/user"
	"net/http"
)

func SidebarMiddleware(c *gin.Context) {
	authId := user.GetUserId(c)
	if (authId == 0) || (user.IsAdmin(authId) == false) {
		c.Redirect(http.StatusFound, "/page_not_found")
		c.Abort()
	}
	c.Next()
}

func AuthMiddleware(c *gin.Context) {
	if user.IsLogin(c) == false {
		c.Redirect(http.StatusFound, "/")
		c.Abort()
	} else {
		c.Next()
	}
}

func CheckLoginPage(c *gin.Context) {
	if user.IsLogin(c) == true {
		c.Redirect(http.StatusFound, "/")
	}
}

func CheckConfigUpload(c *gin.Context) {
	var uploadEntity upload.Entity
	config := uploadEntity.FindOne()
	if config.AllowVisitor == 2 && (user.IsLogin(c) == false) {
		c.JSON(http.StatusBadRequest, gin.H{"code": e.ERROR, "msg": "管理员已禁止游客上传"})
		c.Abort()
	}
}
