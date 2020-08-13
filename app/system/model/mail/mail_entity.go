package mail

import (
	"github.com/jinzhu/gorm"
	db2 "mybedv2/app/helper/db"
	"mybedv2/app/helper/util/str"
	"time"
)

var db = db2.Instance().Table(TableName())

type Entity struct {
	Id            int
	EmailName     string `json:"email_name"`
	EmailHost     string `json:"email_host"`
	EmailPort     string `json:"email_port"`
	EmailUser     string `json:"email_user"`
	EmailPwd      string `json:"email_pwd"`
	EmailTest     string `json:"email_test"`
	EmailTemplate string `json:"email_template"`
	EmailStatus   int    `json:"email_status"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func TableName() string {
	return "email_config"
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

func (e *Entity) CreateOrUpdate(s BindForm) (err error) {
	var oldConfig Entity
	db.Where("id = 1").First(&oldConfig)
	if oldConfig.Id > 0 {
		err = db.Where("id = ?", oldConfig.Id).Omit("id").Updates(str.Struct2Tag(s)).Error
	} else {
		err = db.Create(&s).Error
	}
	return
}

func (e *Entity) FindOne() (s Entity) {
	db.Where("id = 1").First(&s)
	return s
}
