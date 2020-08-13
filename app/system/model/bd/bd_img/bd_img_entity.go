package bd_img

import (
	db2 "mybedv2/app/helper/db"
	"mybedv2/app/helper/util/str"
	"time"
)

var db = db2.Instance().Table(TableName())

type Entity struct {
	Id        int
	AppID     string
	ApiKey    string
	SecretKey string
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func TableName() string {
	return "bd_examine"
}

func (e *Entity) Insert(b BindForm) error {
	return db.Create(&b).Error
}

func (e *Entity) CreateOrUpdate(s BindForm) error {
	var (
		oldConfig Entity
		err       error
	)
	db.Where("id = 1").First(&oldConfig)
	if oldConfig.Id > 0 {
		s.UpdatedAt = time.Now()
		err = db.Where("id = ?", oldConfig.Id).Omit("id").Updates(str.Struct2Tag(s)).Error
	} else {
		s.CreatedAt = time.Now()
		err = db.Create(&s).Error
	}
	return err
}

func FindOne() (s Entity) {
	db.Where("id = 1").First(&s)
	return s
}
