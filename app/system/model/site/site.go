package site

type BindForm struct {
	WebName         string `json:"web_name" form:"web_name"`
	WebUrl          string `json:"web_url" form:"web_url"`
	LogoImg         string `json:"logo_img" form:"logo_img"`
	KeyWord         string `json:"key_word"  form:"key_word"`
	SiteDescription string `json:"site_description" form:"site_description"`
	Copyright       string `json:"copyright"  form:"copyright"`
	RecordInfo      string `json:"record_info"  form:"record_info"`
	SiteStatus      int    `json:"site_status"`
}
