package img

import (
	"github.com/jinzhu/gorm"
	db2 "mybedv2/app/helper/db"
	"mybedv2/conf"
	"strconv"
	"time"
)

var db = db2.Instance().Table(TableName())

type Entity struct {
	Id        int
	ImgName   string
	ImgUrl    string
	ImgThumb  string
	UserId    int64
	Sizes     int64
	Abnormal  string
	Source    int
	ImgType   int
	Explains  string
	Md5       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func TableName() string {
	return "img_data"
}

func (e *Entity) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedAt", time.Now())
	scope.SetColumn("UpdatedAt", time.Now())
	return nil
}

func (e *Entity) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedAt", time.Now())
	return nil
}

func (e *Entity) Insert(img Entity) (int, error) {
	err := db.Create(&img).Error
	return img.Id, err
}

func (e *Entity) Delete(id int) error {
	return db.Where("id = ?", id).Delete(Entity{}).Error
}

//func (e *Entity) Find(ids []string) (imgs []Entity) {
//	db.Where("id in (?)", ids).Find(&imgs)
//	return
//}

func GetImgdataByArr(ids []string) (imgs []Entity) {
	db.Where("id in (?)", ids).Find(&imgs)
	return
}

// md5重复的图片信息创建
func CreateImgdataByMd5(imgdata Entity) (int, error) {
	err := db.Omit("id").Create(&imgdata).Error
	return imgdata.Id, err
}

func GetImgDataByLimit(page int, limit string, userId int64) (imgdata []Entity, count int) {
	var pageSize int
	if limit != "" {
		pageSize, _ = strconv.Atoi(limit)
	} else {
		pageSize = conf.Setting.PageSize
	}
	offset := (page - 1) * pageSize
	query := db.Where("user_id = ?", userId)
	query.Count(&count)
	query.Order("id desc").Limit(pageSize).Offset(offset).Find(&imgdata)
	return
}

func GetImgData(page int, userId int64) (imgdata []Entity, count int) {
	pageSize := conf.Setting.PageSize
	offset := (page - 1) * pageSize
	query := db.Where("user_id = ?", userId)
	query.Count(&count)
	query.Order("id desc").Limit(pageSize).Offset(offset).Find(&imgdata)
	return
}

func GetImgdataBySource(source, page, limit int) (imgdata []Entity, count int) {
	pageSize := conf.Setting.PageSize
	if limit > 0 {
		pageSize = limit
	}
	offset := (page - 1) * pageSize
	query := db.Where("source = ?", source)
	query.Count(&count)
	query.Order("id desc").Limit(pageSize).Offset(offset).Find(&imgdata)
	return
}

func GetImgdataSftp(page, limit int) (imgdata []Entity, count int) {
	pageSize := conf.Setting.PageSize
	if limit > 0 {
		pageSize = limit
	}
	offset := (page - 1) * pageSize
	query := db.Where("img_type = 9")
	query.Count(&count)
	query.Order("id desc").Limit(pageSize).Offset(offset).Find(&imgdata)
	return
}

func GetImgdataByMd5(md5 string) (img Entity) {
	db.Where("md5 = ?", md5).First(&img)
	return img
}

func GetImgdataNum() (count int) {
	db.Count(&count)
	return
}

func GetImgdataById(id int) (img Entity) {
	db.Where("id = ?", id).First(&img)
	return
}

// 获取用户当前空间使用量
func GetImgdataSize(userId int64) (size int, err error) {
	rows, err := db.Select("sum(sizes) as total").Where("user_id = ?", userId).Rows()
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&size)
	}
	return
}

func GetImgNumByDay() (l []Line, err error) {
	err = db.Raw("select DATE_FORMAT(created_at,'%m-%d') AS day, count(*) AS num from img_data group by day").Scan(&l).Error
	return
}

func UpdateImgdataThumb(id int, thumb string) error {
	return db.Where("id = ?", id).UpdateColumn("img_thumb", thumb).Error
}

func DelImgdataByIds(ids []string) error {
	return db.Where("id IN (?)", ids).Delete(Entity{}).Error
}
