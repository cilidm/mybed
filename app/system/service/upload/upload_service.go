package upload

import (
	"mybedv2/app/system/model/upload"
	"strings"
)

func GetVisConfig(uc upload.Entity, upType int) (upc upload.PageConfig) {
	var magn int
	if upType == 1 {
		magn = 1024 * 1024 //byte
	}
	if upType == 2 {
		magn = 1024 //kb
	}
	upc.FileSize = uc.VisitorImgSize * magn
	upc.Explains = uc.VisitorExplains
	upc.ImgCount = uc.VisitorImgNum
	upc.AllowImgUploadExt = strings.ReplaceAll(uc.AllowImgUploadExt, ",", "|")
	if uc.AllowVisitor == 1 {
		upc.AllowUpload = 1
	} else {
		upc.AllowUpload = 2
	}
	return upc
}

func GetMemConfig(uc upload.Entity, upType int) (upc upload.PageConfig) {
	var magn int
	if upType == 1 {
		magn = 1024 * 1024 //byte
	}
	if upType == 2 {
		magn = 1024 //kb
	}
	upc.AllowUpload = 1
	upc.FileSize = uc.MemberImgSize * magn
	upc.Explains = uc.MemberExplains
	upc.ImgCount = uc.MemberImgNum
	upc.AllowImgUploadExt = strings.ReplaceAll(uc.AllowImgUploadExt, ",", "|")
	return
}
