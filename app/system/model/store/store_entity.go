package store

import (
	"github.com/jinzhu/gorm"
	db2 "mybedv2/app/helper/db"
	"mybedv2/app/helper/util/str"
	"time"
)

var db = db2.Instance().Table(TableName())

type Entity struct {
	Id                  int    `gorm:"primary_key"`
	AccessKey           string `json:"access_key" store:"access-key"`
	SecretKey           string `json:"secret_key" store:"secret-key"`
	Endpoint            string `json:"endpoint" store:"endpoint"`
	Region              string `json:"region" store:"region"`
	AppId               string `json:"app-id" store:"app-id"`
	PublicBucket        string `json:"public_bucket" store:"public-bucket"`
	PublicBucketDomain  string `json:"public_bucket_domain" store:"public-bucket-domain"`
	PrivateBucket       string `json:"private_bucket" store:"private-bucket"`
	PrivateBucketDomain string `json:"private_bucket_domain" store:"private-bucket-domain"`
	Expire              int64  `json:"expire" store:"expire"`
	CloudType           string `gorm:"unique;not null"`
	Status              int
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func TableName() string {
	return "cloud_store_config"
}

func (e *Entity) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedAt", time.Now())
	scope.SetColumn("UpdatedAt", time.Now())
	return nil
}

func (e *Entity) BeforeUpdate(scope gorm.Scope) error {
	scope.SetColumn("UpdatedAt", time.Now())
	return nil
}

func (e *Entity) CreateOrUpdate(cs BindForm) (err error) {
	var count int
	db.Where("cloud_type = ?", cs.CloudType).Count(&count)
	if count > 0 {
		err = db.Where("cloud_type = ?", cs.CloudType).Omit("id").Updates(str.Struct2Tag(cs)).Error
	} else {
		err = db.Create(&cs).Error
	}
	return
}

func (e *Entity) FindOne() (cs Entity) {
	db.Where("status = 1").First(&cs)
	return
}

func GetNotNullConfig() (csc []Entity) {
	db.Where("cloud_type <> ''").Find(&csc)
	return
}

func GetStoreConfig(k, v string) (csc Entity) {
	db.Where(k+" = ?", v).First(&csc)
	return
}

// 只配置一个激活的存储源 当一个激活时，另一个激活的就被取消
func ChangeOtherStoreStatus(ct string) {
	db.Where("status = 1 AND cloud_type <> ?", ct).Update("status", 2)
}
