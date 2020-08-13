package db

import "time"

type User struct {
	ID            int
	Username      string    `json:"username"`
	Password      string    `json:"-"`
	Salt          string    `json:"-"`
	Status        int       `json:"status"`
	Email         string    `json:"email"`
	Avatar        string    `json:"avatar"`
	Nickname      string    `json:"nickname"`
	Birthday      time.Time `json:"birthday"`
	Level         int64     `json:"level"`
	Uid           string    `json:"-"`
	Isok          uint      `json:"-"`
	Memory        int64     `json:"memory"`
	Groupid       int64     `json:"-"`
	LastLoginTime time.Time `json:"last_login_time"`
	LastLoginIp   string    `json:"last_login_ip"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"-"`
}

type CloudStoreConfig struct {
	Id                  int
	AccessKey           string
	SecretKey           string
	Endpoint            string
	Region              string
	AppId               string
	PublicBucket        string
	PublicBucketDomain  string
	PrivateBucket       string
	PrivateBucketDomain string
	Expire              int64
	CloudType           string
	Status              uint
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type EmailConfig struct {
	Id            int
	EmailName     string
	EmailHost     string
	EmailPort     string
	EmailUser     string
	EmailPwd      string
	EmailTest     string
	EmailTemplate string
	EmailStatus   int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ImgData struct {
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

type SystemOperLog struct {
	Id           int
	Title        string
	Method       string
	OperName     string
	OperUrl      string
	OperIp       string
	OperLocation string
	OperParam    string
	JsonResult   string
	Status       int
	ErrorMsg     string
	OperTime     time.Time
}

type SiteConfig struct {
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

type UploadConfig struct {
	Id                int
	AllowUploadExt    string
	MaxUploadSize     int
	AllowImgUploadExt string

	MemberImgTotalSize  int
	MemberImgSize       int
	MemberImgNum        int
	VisitorImgTotalSize int
	VisitorImgSize      int
	VisitorImgNum       int
	AllowVisitor        int    //是否允许游客上传
	IpBlacklist         string //IP黑名单
	VisitorExplains     int    //游客图片保存期限
	MemberExplains      int    //会员图片保存期限
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type Code struct {
	Id        int
	Code      string
	Value     int
	Status    int
	UserId    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BdExamine struct {
	Id        int
	AppID     string
	ApiKey    string
	SecretKey string
	Status    int
	UseMd5    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ImgBlacklist struct {
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

func initModel() {
	checkMigrate(&User{})
	checkMigrate(&ImgData{})
	checkMigrate(&CloudStoreConfig{})
	checkMigrate(&SiteConfig{})
	checkMigrate(&EmailConfig{})
	checkMigrate(&UploadConfig{})
	checkMigrate(&SystemOperLog{})
	checkMigrate(&Code{})
	checkMigrate(&BdExamine{})
	checkMigrate(&ImgBlacklist{})
}

func checkMigrate(tb interface{}) {
	if db.HasTable(tb) == false {
		db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(tb)
	}
}

func initAdmin() {
	var userNum, uploadNum, siteNum int
	db.Model(User{}).Where("id = 1").Count(&userNum)
	db.Model(UploadConfig{}).Where("id = 1").Count(&uploadNum)
	db.Model(SiteConfig{}).Where("id = 1").Count(&siteNum)
	if userNum == 0 {
		db.Create(User{
			ID: 1, Username: "admin", Password: "4b4b7e95aa5ad1fc86bb1f2b3bd6d7ee", Salt: "BtBFNwXINC",
			Nickname: "admin", Status: 1, Level: 99, Avatar: "/static/upload/1/avatar/2354117_0.jpg",
			Memory: 1000, CreatedAt: time.Now(),
		})
	}

	if uploadNum == 0 {
		db.Create(UploadConfig{
			Id: 1, AllowUploadExt: "pdf,log,txt", MaxUploadSize: 0, AllowImgUploadExt: "png,jpg,jpeg,gif,bmp",
			MemberImgTotalSize: 1000, MemberImgSize: 10, MemberImgNum: 100,
			VisitorImgTotalSize: 10, VisitorImgSize: 1, VisitorImgNum: 2, AllowVisitor: 2,
			VisitorExplains: 0, MemberExplains: 0, CreatedAt: time.Now(),
		})
	}

	if siteNum == 0 {
		db.Create(SiteConfig{
			Id: 1, WebName: "Mybed", WebUrl: "http://localhost:8000/", SiteStatus: 1, CreatedAt: time.Now(),
		})
	}
}
