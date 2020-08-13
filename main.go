package main

import (
	"fmt"
	"github.com/fvbock/endless"
	"log"
	_ "mybedv2/app/helper/redis"
	"mybedv2/app/helper/router"
	_ "mybedv2/app/pic_bed"
	_ "mybedv2/app/system/controller"
	"mybedv2/conf"
	"net/http"
	"runtime"
	"syscall"
)

func main() {
	sysType := runtime.GOOS
	if sysType == "windows" {
		r := router.InitRouter()
		s := &http.Server{
			Addr:           fmt.Sprintf(":%d", conf.Setting.HTTPPort),
			Handler:        r,
			ReadTimeout:    conf.Setting.ReadTimeout,
			WriteTimeout:   conf.Setting.WriteTimeout,
			MaxHeaderBytes: 1 << 20,
		}
		s.ListenAndServe()
	} else {
		endless.DefaultReadTimeOut = conf.Setting.ReadTimeout
		endless.DefaultWriteTimeOut = conf.Setting.WriteTimeout
		endless.DefaultMaxHeaderBytes = 1 << 20
		endPoint := fmt.Sprintf(":%d", conf.Setting.HTTPPort)
		server := endless.NewServer(endPoint, router.InitRouter())
		server.BeforeBegin = func(add string) {
			log.Printf("Actual pid is %d And port is %d", syscall.Getpid(), conf.Setting.HTTPPort)
		}
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Server err: %v", err)
		}
	}
}
