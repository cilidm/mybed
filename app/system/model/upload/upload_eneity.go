package upload

import (
	"github.com/jinzhu/gorm"
	db2 "mybedv2/app/helper/db"
	"mybedv2/app/helper/util/str"
	"time"
)

var db = db2.Instance().Table(TableName())

type Entity struct {
	Id                  int
	AllowUploadExt      string `json:"allow_upload_ext"`
	MaxUploadSize       int    `json:"max_upload_size"`
	AllowImgUploadExt   string `json:"allow_img_upload_ext"`
	MemberImgTotalSize  int    `json:"member_img_total_size"`
	MemberImgSize       int    `json:"member_img_size"`
	MemberImgNum        int    `json:"member_img_num"`
	VisitorImgTotalSize int    `json:"visitor_img_total_size"`
	VisitorImgSize      int    `json:"visitor_img_size"`
	VisitorImgNum       int    `json:"visitor_img_num"`
	AllowVisitor        int    `json:"allow_visitor"`    //是否允许游客上传
	IpBlacklist         string `json:"ip_blacklist"`     //IP黑名单
	VisitorExplains     int    `json:"visitor_explains"` //游客图片保存期限
	MemberExplains      int    `json:"member_explains"`  //会员图片保存期限
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func TableName() string {
	return "upload_config"
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
