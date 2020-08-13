package util

import (
	"encoding/json"
	"io/ioutil"
	"mybedv2/app/helper/util/str"
	"net/http"
	"strings"
)

//获取IP所在城市
func GetCityByIp(ip string) string {
	if ip == "" {
		return ""
	}

	if ip == "::1" || ip == "127.0.0.1" {
		return "内网IP"
	}

	url := "http://whois.pconline.com.cn/ipJson.jsp?json=true&ip=" + ip
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Accept-Charset", "GBK,utf-8;q=0.7,*;q=0.3")
	response, _ := client.Do(request)
	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		bodystr := string(body)
		tmp := str.ConvertToString(bodystr, "gbk", "utf-8")
		p := make(map[string]interface{}, 0)
		if err := json.Unmarshal([]byte(tmp), &p); err == nil {
			return p["city"].(string)
		}
	}
	return ""
}

func GetCurrentIP() string {
	r := http.Request{}
	// 这里也可以通过X-Forwarded-For请求头的第一个值作为用户的ip
	// 但是要注意的是这两个请求头代表的ip都有可能是伪造的
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		// 当请求头不存在即不存在代理时直接获取ip
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}
	return ip
}
