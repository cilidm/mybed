package mail

type BindForm struct {
	EmailName     string `json:"email_name" form:"email_name"`
	EmailHost     string `json:"email_host" form:"email_host"`
	EmailPort     string `json:"email_port" form:"email_port"`
	EmailUser     string `json:"email_user" form:"email_user"`
	EmailPwd      string `json:"email_pwd" form:"email_pwd"`
	EmailTest     string `json:"email_test" form:"email_test"`
	EmailTemplate string `json:"email_template" form:"email_template"`
	EmailStatus   int    `json:"email_status"`
}
