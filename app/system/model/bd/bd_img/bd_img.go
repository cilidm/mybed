package bd_img

import "time"

type BindForm struct {
	AppID     string    `json:"app_id" form:"app_id"`
	ApiKey    string    `json:"api_key" form:"api_key" binding:"required"`
	SecretKey string    `json:"secret_key" form:"secret_key" binding:"required"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
