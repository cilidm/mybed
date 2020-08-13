package img

import (
	"errors"
	"fmt"
	"mybedv2/app/helper/util/calculate"
	imgModel "mybedv2/app/system/model/img"
	"mybedv2/app/system/model/store"
	u2 "mybedv2/app/system/model/user"
	"os"
)

//获取用户当前空间使用率
func GetUserMemPer(userId, mem int64) float64 {
	totalSize, err := imgModel.GetImgdataSize(userId)
	if err != nil {
		return 0
	}
	return calculate.Decimal((float64(totalSize) / float64(mem*1024*1024)) * 100) //默认单位mb 这里要转换成byte
}

// 判断用户图片空间是否超限
func IsMaxSize(userId, imgSize int64) bool {
	totalSize, err := imgModel.GetImgdataSize(userId)
	if err != nil {
		totalSize = 0
	}
	u := u2.Entity{}
	user, err := u.FindById(userId)
	if err != nil {
		fmt.Println("获取用户出错:", err)
		return true
	}
	return (int64(totalSize) + imgSize) > user.Memory*1024*1024
}

func CreateImgdata(file, imgUrl, userIp, md5 string, uid int64, cloudType string) (int, error) {
	var (
		img    imgModel.Entity
		entity imgModel.Entity
	)
	fileinfo, err := os.Stat(file)
	if err != nil {
		return 0, errors.New("os stat失败," + err.Error())
	}
	img.ImgUrl = imgUrl
	//img.ImgThumb = getTmpPath(fileinfo.Name())
	img.ImgName = fileinfo.Name()
	img.Sizes = fileinfo.Size()
	img.Abnormal = userIp
	img.UserId = uid
	img.Md5 = md5
	img.Source = store.GetStoreNum(cloudType)
	insertId, err := entity.Insert(img)
	if err != nil {
		return 0, errors.New("数据库保存失败," + err.Error())
	}
	return insertId, nil
}
