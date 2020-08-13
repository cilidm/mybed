package qrcode

import (
	"fmt"
	"github.com/gin-gonic/gin"
	qrcode "github.com/yeqown/go-qrcode"
	"net/http"
)

func CreateQrcode(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.String(http.StatusOK, "请传入网址")
		return
	}
	qrc, err := qrcode.New(url)
	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err)
	}
	header := c.Writer.Header()
	header.Add("Content-Type", "image/jpeg")
	qrc.SaveTo(c.Writer)
}
