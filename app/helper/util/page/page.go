package page

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mybedv2/app/helper/util/str"
	"net/http"
	"strconv"
	"strings"
)

func GetQueryPage(c *gin.Context) int {
	pageQuery := c.DefaultQuery("page", "1")
	if strings.HasPrefix(pageQuery, "#") {
		pageQuery = str.Substr(pageQuery, 0, -1)
	}
	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		fmt.Println(http.StatusBadRequest, "分页数据有误")
		return 1
	}
	return page
}
