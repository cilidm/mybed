package conf

import (
	"github.com/go-ini/ini"
	"io/ioutil"
	"log"
	"mybedv2/app/helper/e"
	"mybedv2/app/helper/util/pathdir"
	"mybedv2/app/helper/util/str"
	"os"
	"sync"
	"time"
)

var (
	Cfg     *ini.File
	once    sync.Once
	Setting = &SettingConf{}
	server  = &Server{}
)

type Server struct {
	HTTPPort     int
	ReadTimeout  int
	WriteTimeout int
}

type SettingConf struct {
	RunMode      string
	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	PageSize     int
	JwtSecret    string
	UploadTmpDir string
	DBType       string
	DBUser       string
	DBPwd        string
	DBHost       string
	DBTableName  string
	DBPath       string
	RedisAddr    string
	RedisPWD     string
	RedisDB      int
}

func init() {
	iniPath := "conf/app.ini"
	once.Do(func() {
		var err error
		Cfg, err = ini.Load(iniPath)
		if err != nil {
			pathdir.PathExists("conf/")
			pathdir.CreateFile(iniPath)
			createIni(iniPath)
			os.Exit(-1)
		}
		err = Cfg.Section("runMode").MapTo(Setting)
		str.CheckErr(err, "映射配置文件出错，请检查runMode配置")
		err = Cfg.Section("app").MapTo(Setting)
		str.CheckErr(err, "映射配置文件出错，请检查app配置")
		err = Cfg.Section("database").MapTo(Setting)
		str.CheckErr(err, "映射配置文件出错，请检查database配置")
		err = Cfg.Section("redis").MapTo(Setting)
		str.CheckErr(err, "映射配置文件出错，请检查redis配置")

		err = Cfg.Section("server").MapTo(server)
		str.CheckErr(err, "映射配置文件出错，请检查server配置")
		Setting.HTTPPort = server.HTTPPort
		Setting.ReadTimeout = time.Duration(server.ReadTimeout) * time.Second
		Setting.WriteTimeout = time.Duration(server.WriteTimeout) * time.Second
	})
}

func createIni(path string) {
	ini := e.IniStr
	if err := ioutil.WriteFile(path, []byte(ini), os.ModePerm); err != nil {
		log.Println("写入默认配置信息失败，请手动创建")
	} else {
		log.Println("写入配置信息完毕，请重新启动")
	}
}
