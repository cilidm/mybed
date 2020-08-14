package code

import "time"

type BindForm struct {
	CodeNum   int `json:"code" form:"code_num"`
	CodeValue int `json:"value" form:"code_value"`
}

type InsertForm struct {
	Code      string `json:"code"`
	Value     int    `json:"value"`
	Status    int    `json:"status"`
	UserId    int    `json:"user_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ListJson struct {
	Id        int    `json:"id"`
	Code      string `json:"code"`
	Value     int    `json:"value"`
	Status    int    `json:"status"`
	UserName  string `json:"user_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
