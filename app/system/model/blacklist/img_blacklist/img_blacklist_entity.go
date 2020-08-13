package img_blacklist

import (
	db2 "mybedv2/app/helper/db"
	"time"
)

var db = db2.Instance().Table(TableName())

type Entity struct {
	Id        int
	FileName  string
	FileSize  string
	FileMd5   string
	UserId    int
	UserIp    string
	Info      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func TableName() string {
	return "img_blacklist"
}

func (e *Entity) Insert(b BindForm) error {
	b.CreatedAt = time.Now()
	return db.Create(&b).Error
}

func (e *Entity) Delete(id int) error {
	return db.Where("id = ?", id).Delete(Entity{}).Error
}

func FindByMd5(str string) (count int) {
	db.Where("file_md5 = ?", str).Count(&count)
	return
}

func GetAllNum() (count int) {
	db.Count(&count)
	return
}
