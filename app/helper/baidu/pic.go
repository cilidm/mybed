package baidu

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mybedv2/app/helper/redis"
	_ "mybedv2/app/helper/redis"
	"mybedv2/app/helper/util/str"
	"mybedv2/app/system/model/bd/bd_img"
	"mybedv2/app/system/model/blacklist/img_blacklist"
	"mybedv2/app/system/service/user"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	TokenUrl = "https://aip.baidubce.com/oauth/2.0/token"
	BaiduApi = "https://aip.baidubce.com/rest/2.0/solution/v1/img_censor/v2/user_defined"
)

type TokenResp struct {
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	SessionKey    string `json:"session_key"`
	SessionSecret string `json:"session_secret"`
	Scope         string `json:"scope"`
	ExpiresIn     int    `json:"expires_in"`
}

type Resp struct {
	Conclusion     string     `json:"conclusion"`
	LogId          int        `json:"log_id"`
	Data           []RespData `json:"data"`
	ConclusionType int        `json:"conclusionType"`
	ErrorMsg       string     `json:"error_msg"`
	ErrorCode      int        `json:"error_code"`
}

type RespData struct {
	Msg            string  `json:"msg"`
	Conclusion     string  `json:"conclusion"`
	Probability    float64 `json:"probability"`
	SubType        int     `json:"subType"`
	ConclusionType int     `json:"conclusionType"`
	Type           int     `json:"type"`
}

func accessToken(id string, secret string) (token string, err error) {
	fmt.Println("获取token")
	apiURL := fmt.Sprintf(TokenUrl+"?grant_type=client_credentials&client_id=%s&client_secret=%s", id, secret)
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var tokenResp TokenResp
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return "", err
	}
	token = tokenResp.AccessToken
	redis.Client.Set("baidu_token", token, time.Second*time.Duration(tokenResp.ExpiresIn-1800))
	return token, nil
}

func GetToken() (string, error) {
	token, err := redis.Client.Get("baidu_token").Result()
	if token == "" {
		bd := bd_img.FindOne()
		token, err = accessToken(bd.ApiKey, bd.SecretKey)
	}
	if err != nil {
		fmt.Println("accessToken()", err)
		return "", err
	}
	return token, nil
}

func GetIdentificationResult(filePath string, sourceType int) (Resp, error) {
	token, err := GetToken()
	if err != nil {
		fmt.Println("获取token失败")
		return Resp{}, err
	}
	var (
		urlPath = BaiduApi + "?access_token=" + token
		param   string
	)
	if sourceType == 1 {
		sourcestring, err := str.GetBase64ByFile(filePath)
		if err != nil {
			fmt.Println("图片转base64失败", err.Error())
			return Resp{}, err
		}
		param = "image=" + url.QueryEscape(sourcestring)
	}
	resp, err := http.Post(urlPath, "application/x-www-form-urlencoded", strings.NewReader(param))
	if err != nil {
		fmt.Println(err)
		return Resp{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return Resp{}, err
	}
	fmt.Println(string(body))
	var res Resp
	if err := json.Unmarshal(body, &res); err != nil {
		fmt.Println(err)
		return Resp{}, err
	}
	return res, nil
}

// 通用图片审核，传图片本地路径+md5
func BdPicExamine(fileTmpPath, md5 string, c *gin.Context) error {
	bdResp, err := GetIdentificationResult(fileTmpPath, 1)
	if err != nil {
		return err
	}
	if bdResp.ConclusionType != 1 {
		fs, _ := os.Stat(fileTmpPath)
		var blist img_blacklist.Entity
		err := blist.Insert(img_blacklist.BindForm{
			Info:     bdResp.Conclusion + ":" + bdResp.Data[0].Msg,
			UserIp:   c.ClientIP(),
			UserId:   int(user.GetUserId(c)),
			FileName: fs.Name(),
			FileSize: fs.Size(),
			FileMd5:  md5,
		})
		if err != nil {
			return err
		}
		return errors.New(bdResp.Conclusion + ":" + bdResp.Data[0].Msg)
	}
	return nil
}
