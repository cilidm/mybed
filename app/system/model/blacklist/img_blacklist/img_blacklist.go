package img_blacklist

import "time"

type BindForm struct {
	FileName  string `json:"file_name" form:"file_name"`
	FileSize  int64  `json:"file_size" form:"file_size"`
	FileMd5   string `json:"file_md5" form:"file_md5"`
	UserIp    string `json:"user_ip" form:"user_ip"`
	UserId    int    `json:"user_id" form:"user_id"`
	Info      string `json:"info" form:"info"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
