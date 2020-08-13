package e

const (
	DayTimeFormatDir  = "20060102"
	HourTimeFormatDir = "2006010215"
	TimeFormatDir     = "20060102150405"
	TimeFormat        = "2006-01-02 15:04:05"

	ImgListTask    = "img_list"
	ImgListErrTask = "img_list_retry"

	StroeDefaultDir = "mybed"

	UserId          = "uid"
	UserSessionInfo = "user_info"

	// 存储类型, cs 前缀表示 CloudStore
	StoreOss   = "cs-oss"   //oss存储
	StoreMinio = "cs-minio" //minio存储
	StoreCos   = "cs-cos"   //腾讯云存储
	StoreObs   = "cs-obs"   //华为云存储
	StoreBos   = "cs-bos"   //百度云存储
	StoreQiniu = "cs-qiniu" //七牛云储存
	StoreUpyun = "cs-upyun" //又拍云存储

	// 临时文件夹redis-list名称
	TempDirList = "pending-folder"
)
