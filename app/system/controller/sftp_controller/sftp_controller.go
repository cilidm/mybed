package sftp_controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SftpPage(c *gin.Context) {
	c.HTML(http.StatusOK, "sftp.html", gin.H{})
}
