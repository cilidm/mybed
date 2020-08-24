package main

import (
	"fmt"
	_ "mybedv2/app/helper/redis"
	"mybedv2/app/helper/router"
	_ "mybedv2/app/pic_bed"
	_ "mybedv2/app/system/controller"
	"mybedv2/conf"
	"net/http"
)

func main() {
	r := router.InitRouter()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", conf.Setting.HTTPPort),
		Handler:        r,
		ReadTimeout:    conf.Setting.ReadTimeout,
		WriteTimeout:   conf.Setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
