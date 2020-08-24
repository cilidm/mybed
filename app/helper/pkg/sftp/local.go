package sftp

import (
	"errors"
	"fmt"
	"io/ioutil"
	"mime"
	"mybedv2/app/helper/util/pathdir"
	img2 "mybedv2/app/system/model/img"
	"mybedv2/app/system/model/store"
	storeService "mybedv2/app/system/service/store"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Local2StoreForm struct {
	Base      string `json:"base" form:"base"`
	Paths     string `json:"paths" form:"paths"`
	KeepLocal int    `json:"keep_local"`
}

type Local2StoreSftp struct {
	base     string
	paths    string
	countNum int
	allNum   int
}

var local Local2StoreSftp

func Local2Store(f Local2StoreForm) error {
	local.countNum = 0
	local.base = f.Base
	if strings.HasSuffix(f.Paths, "/") {
		f.Paths = strings.TrimRight(f.Paths, "/")
	}
	local.paths = f.Paths
	var err error
	cs, err = storeService.NewCloudStore(false)
	if err != nil {
		return errors.New("未检测到可用的store")
	}
	var (
		dirName []string
		dirs    []string
	)
	files, dir, _ := ShowDirWithFile(local.paths, dirName, dirs)
	local.allNum = len(files)
	for _, v := range files {
		sema.Add(1)
		go Handler(v, f.KeepLocal)
	}
	sema.Wait()
	if f.KeepLocal == 2 {
		for _, v := range dir {
			stat, _ := os.Stat(v)
			if stat.Size() == 0 && stat.IsDir() {
				os.RemoveAll(v)
			}
		}
	}
	return nil
}

func Handler(v string, keepLocal int) {
	local.countNum++
	fmt.Println(local.countNum, "/", local.allNum)
	defer sema.Done()
	md5, _ := pathdir.GetMd5V1(v)
	hasImg := img2.GetImgdataByMd5(md5)
	if hasImg.Id > 0 {
		return
	}
	stat, _ := os.Stat(v)
	rel, _ := filepath.Rel(local.base, v)
	var savePath string
	sysType := runtime.GOOS
	if sysType == "windows" {
		savePath = "local/" + strings.ReplaceAll(rel, "\\", "/") //服务器保存的路径，win/文件夹名/图片名称
	} else {
		savePath = "local/" + rel
	}
	ext := mime.TypeByExtension(path.Ext(v))
	if err := cs.Upload(v, savePath, map[string]string{"Content-Type": ext}); err != nil {
		fmt.Println(err)
		return
	}
	var storeEntity store.Entity
	csc := storeEntity.FindOne()
	var mcPath string
	mcPath = csc.PublicBucketDomain + csc.PublicBucket + "/" + savePath
	var img img2.Entity
	_, err := img.Insert(img2.Entity{
		ImgName:   stat.Name(),
		ImgUrl:    mcPath,
		UserId:    1,
		Sizes:     stat.Size(),
		Source:    store.GetStoreNum(csc.CloudType),
		Md5:       md5,
		Abnormal:  "127.0.0.1",
		ImgType:   8,
		CreatedAt: time.Now(),
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	if keepLocal == 2 {
		os.Remove(v)
	}
}

func ShowDirWithFile(path string, dirName []string, dirs []string) ([]string, []string, error) {
	dir, err := os.Stat(path)
	if err != nil {
		return dirName, dirs, err
	}
	if dir.IsDir() {
		rd, err := ioutil.ReadDir(path)
		if err != nil {
			return dirName, dirs, err
		}
		for _, v := range rd {
			if v.IsDir() {
				dirs = append(dirs, path+"/"+v.Name())
				dirName, dirs, err = ShowDirWithFile(path+"/"+v.Name(), dirName, dirs)
				if err != nil {
					return dirName, dirs, err
				}
			} else {
				if strings.HasPrefix(v.Name(), ".") {
					continue
				}
				dirName = append(dirName, path+"/"+v.Name())
			}
		}
	}
	return dirName, dirs, nil
}
