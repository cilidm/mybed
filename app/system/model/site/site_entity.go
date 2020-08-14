package site

import (
	"github.com/jinzhu/gorm"
	db2 "mybedv2/app/helper/db"
	"mybedv2/app/helper/util/str"
	"strings"
	"time"
)

var db = db2.Instance().Table(TableName())

type Entity struct {
	Id              int
	WebName         string
	WebUrl          string
	LogoImg         string
	KeyWord         string
	SiteDescription string
	Copyright       string
	RecordInfo      string
	SiteStatus      int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func TableName() string {
	return "site_config"
}

func (e *Entity) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedAt", time.Now())
	scope.SetColumn("UpdatedAt", time.Now())
	scope.SetColumn("SiteDescription", strings.TrimSpace(e.SiteDescription))
	return nil
}

func (e *Entity) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedAt", time.Now())
	scope.SetColumn("WebName", strings.TrimSpace(e.WebName))
	scope.SetColumn("SiteDescription", strings.TrimSpace(e.SiteDescription))
	return nil
}

func (e *Entity) CreateOrUpdate(s BindForm) error {
	var (
		oldConfig Entity
		err       error
	)
	db.Where("id = 1").First(&oldConfig)
	if oldConfig.Id > 0 {
		err = db.Where("id = ?", oldConfig.Id).Omit("id").Updates(str.Struct2Tag(s)).Error
	} else {
		err = db.Create(&s).Error
	}
	return err
}

func (e *Entity) FindOne() (s Entity) {
	db.Where("id = 1").First(&s)
	return s
}

func FindSiteConfig() (s Entity) {
	db.Where("id = 1").First(&s)
	return s
}