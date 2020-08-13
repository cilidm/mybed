package e

const (
	SUCCESS                      = 200
	ERROR                        = 500
	ErrorLoginCheckRequired      = 10001
	ErrorLoginCheckPwd           = 10002
	ErrorLoginCheckStatus        = 10003
	ErrorSetSession              = 10004
	ErrorUploadForm              = 10005
	ErrorUploadSave              = 10006
	ErrorSaveImgdata             = 10007
	ErrorUploadStore             = 10008
	ErrorUploadImgBase64NullByte = 10009
	ErrorUploadImgBase64Save     = 10010

	// pic_bed
	UserStatusErr          = 1000
	FileTypeNotAllow       = 4000
	NotEnoughFreeSpace     = 4005
	ExceedsUploadSizeLimit = 4006
)
