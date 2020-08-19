package img

type List struct {
	Id        int    `json:"id"`
	ImgName   string `json:"img_name"`
	ImgUrl    string `json:"img_url"`
	ImgThumb  string `json:"img_thumb"`
	UserId    int64  `json:"user_id"`
	Sizes     string `json:"sizes"`
	Abnormal  string `json:"abnormal"`
	Source    int    `json:"source"`
	ImgType   int    `json:"img_type"`
	Explains  string `json:"explains"`
	Md5       string `json:"md_5"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Line struct {
	Day string `json:"day"`
	Num int    `json:"num"`
}

type PageJson struct {
	Id        int    `json:"id"`
	ImgName   string `json:"img_name"`
	ImgUrl    string `json:"img_url"`
	ImgThumb  string `json:"img_thumb"`
	UserId    int64  `json:"user_id"`
	Sizes     string `json:"sizes"`
	Abnormal  string `json:"abnormal"`
	Source    string `json:"source"`
	ImgType   int    `json:"img_type"`
	Explains  string `json:"explains"`
	Md5       string `json:"md_5"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

//获取img_type归属
func GetImgType(imgType int) string {
	img := map[int]string{
		0: "页面上传",
		9: "sftp迁移",
	}
	return img[imgType]
}
