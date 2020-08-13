package user

import "time"

type ListShow struct {
	Id            int    `json:"id"`
	Username      string `json:"username"`
	Status        int    `json:"status"`
	Email         string `json:"email"`
	Avatar        string `json:"avatar"`
	Nickname      string `json:"nickname"`
	MemoryUsed    int    `json:"memory"`
	Level         int64  `json:"level"`
	LastLoginTime string `json:"last_login_time"`
	LastLoginIp   string `json:"last_login_ip"`
	CreatedAt     string `json:"created_at"`
}

// 前端显示字段
type Info struct {
	ID            int       `json:"id"`
	Username      string    `json:"username"`
	Status        int       `json:"status"`
	Email         string    `json:"email"`
	Avatar        string    `json:"avatar"`
	Nickname      string    `json:"nickname"`
	Level         int64     `json:"level"`
	Memory        int64     `json:"memory"`
	UsedMem       string    `json:"used_mem"`
	UsedPer       int64     `json:"used_per"`
	LastLoginTime time.Time `json:"last_login_time"`
	LastLoginIp   string    `json:"last_login_ip"`
	CreatedAt     time.Time `json:"created_at"`
}

// 数据接收绑定
type BindForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type EditPwdForm struct {
	Oldpwd     string `json:"oldpwd" form:"oldpwd" binding:"required"`
	Newpwd     string `json:"newpwd" form:"newpwd" binding:"required"`
	Confirmpwd string `json:"confirmpwd" form:"confirmpwd" binding:"required"`
}
