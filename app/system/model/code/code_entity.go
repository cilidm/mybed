package code

import (
	db2 "mybedv2/app/helper/db"
	"time"
)

var db = db2.Instance().Table(TableName())

type Entity struct {
	Id        int
	Code      string
	Value     int
	Status    int		// 初始1 已用2
	UserId    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func TableName() string {
	return "code"
}

func (e *Entity) UpdateCode() error {
	return db.Omit("id").Updates(e).Error
}

func UpdateStatus(e Entity) error {
	return db.Where("id = ?",e.Id).UpdateColumns(Entity{
		UserId: e.UserId,
		Status: e.Status,
		UpdatedAt: time.Now(),
	}).Error
}

func InsertCode(b InsertForm) error {
	b.Status = 1
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
	return db.Create(&b).Error
}

func FindCode(code string) (e Entity,err error) {
	err = db.Where("status = 1 AND code = ?",code).First(&e).Error
	return
}

func FindByCode(code string) (e Entity,err error) {
	err = db.Where("code = ?",code).First(&e).Error
	return
}

func GetCodeList(page,limit int) (code []Entity,count int){
	pageSize := limit
	offset := (page - 1) * pageSize
	query := db.Model(Entity{})
	query.Count(&count)
	query.Order("id desc").Limit(pageSize).Offset(offset).Find(&code)
	return
}

func DeleteByIds(ids []string) error {
	return db.Where("id IN (?)", ids).Delete(Entity{}).Error
}

func DeleteById(id string) error {
	return db.Where("id = (?)", id).Delete(Entity{}).Error
}