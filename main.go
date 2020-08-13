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
	"syscall"
)

func main() {
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
