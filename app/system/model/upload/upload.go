package upload

type BindForm struct {
	AllowUploadExt      string `json:"allow_upload_ext" form:"allow_upload_ext"`
	MaxUploadSize       int    `json:"max_upload_size" form:"max_upload_size"`
	AllowImgUploadExt   string `json:"allow_img_upload_ext" form:"allow_img_upload_ext"`
	MemberImgTotalSize  int    `json:"member_img_total_size" form:"member_img_total_size"`
	MemberImgSize       int    `json:"member_img_size" form:"member_img_size"`
	MemberImgNum        int    `json:"member_img_num" form:"member_img_num"`
	VisitorImgTotalSize int    `json:"visitor_img_total_size" form:"visitor_img_total_size"`
	VisitorImgSize      int    `json:"visitor_img_size" form:"visitor_img_size"`
	VisitorImgNum       int    `json:"visitor_img_num" form:"visitor_img_num"`
	AllowVisitor        int    `json:"allow_visitor" form:"allow_visitor"`       //是否允许游客上传
	IpBlacklist         string `json:"ip_blacklist" form:"ip_blacklist"`         //IP黑名单
	VisitorExplains     int    `json:"visitor_explains" form:"visitor_explains"` //游客图片保存期限
	MemberExplains      int    `json:"member_explains" form:"member_explains"`   //会员图片保存期限
}

type PageConfig struct {
	FileSize          int
	ImgCount          int
	AllowUpload       int
	AllowImgUploadExt string
	Explains          int
}
