package e

var MsgFlags = map[int]string{
	SUCCESS:                      "操作成功",
	ERROR:                        "操作失败",
	ErrorLoginCheckRequired:      "用户名或密码不能为空",
	ErrorLoginCheckPwd:           "账号或密码错误",
	ErrorLoginCheckStatus:        "该账号已禁用",
	ErrorSetSession:              "服务器有误，请稍后重试",
	ErrorUploadForm:              "上传失败，请重试",
	ErrorUploadSave:              "图片保存失败",
	ErrorSaveImgdata:             "图片信息保存失败",
	ErrorUploadStore:             "上传到store失败，请检查配置源",
	ErrorUploadImgBase64NullByte: "图片base64数据接收失败，请检查后重试",
	ErrorUploadImgBase64Save:     "图片保存失败，请检查后重试",
	ExceedsUploadSizeLimit:       "图片超出上传大小限制",
	NotEnoughFreeSpace:           "可用空间不足",
	FileTypeNotAllow:             "文件类型不符合要求",
	UserStatusErr:                "上传失败，管理员已经关闭你的上传功能",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
