package code

import (
	"github.com/jinzhu/gorm"
	db2 "mybedv2/app/helper/db"
	"time"
)

var db = db2.Instance().Table(TableName())

type Entity struct {
	Id        int
	Code      string
	Value     int
	Status    int
	UserId    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func TableName() string {
	return "code"
}

func (e *Entity) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("Status", 1)
	scope.SetColumn("CreatedAt", time.Now())
	scope.SetColumn("UpdatedAt", time.Now())
	return nil
}

func (e *Entity) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedAt", time.Now())
	return nil
}

func (e *Entity) Insert(b BindForm) error {
	return db.Create(&b).Error
}
