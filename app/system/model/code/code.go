package code

import "time"

type BindForm struct {
	Code      string    `json:"code" form:"code"`
	Value     int       `json:"value" form:"value"`
	Status    int       `json:"status" form:"status"`
	UserId    int       `json:"user_id" form:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
