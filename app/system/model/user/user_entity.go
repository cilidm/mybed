package user

import (
	db2 "mybedv2/app/helper/db"
	"mybedv2/app/system/model/upload"
	"time"
)

var db = db2.Instance().Table(TableName())

type Entity struct {
	ID            int
	Username      string    `json:"username"`
	Password      string    `json:"-"`
	Salt          string    `json:"-"`
	Status        int       `json:"status"`
	Email         string    `json:"email"`
	Avatar        string    `json:"avatar"`
	Nickname      string    `json:"nickname"`
	Birthday      time.Time `json:"birthday"`
	Level         int64     `json:"level"`
	Uid           string    `json:"-"`
	Isok          uint      `json:"-"`
	Memory        int64     `json:"memory"`
	Groupid       int64     `json:"-"`
	LastLoginTime time.Time `json:"last_login_time"`
	LastLoginIp   string    `json:"last_login_ip"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"-"`
}

func TableName() string {
	return "user"
}

func (e *Entity) Insert(user Entity) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.LastLoginTime = time.Now()
	user.Status = 1
	user.Avatar = "/static/admin/images/users/1_0.jpg"
	var up upload.Entity
	mem := up.FindOne()
	if mem.Id > 0 && mem.MemberImgTotalSize > 0 {
		user.Memory = int64(mem.MemberImgTotalSize) //如果已经设置过会员最大上传空间就使用设置的值
	} else {
		user.Memory = 100 //没设置过 默认100mb
	}
	return db.Create(&user).Error
}

func (e *Entity) Update() {
	db.Omit("id").Updates(e)
}

func (e *Entity) Delete(user Entity) error {
	return db.Where("id = ?", user.ID).Delete(&user).Error
}

func (e *Entity) FindById(id int64) (user Entity, err error) {
	err = db.Where("id = ?", id).First(&user).Error
	return
}

func UpdateMem(user Entity) error {
	return db.Where("id = ?",user.ID).UpdateColumns(Entity{Memory: user.Memory,UpdatedAt: time.Now()}).Error
}

func FindUserById(id int64) (user Entity, err error) {
	err = db.Where("id = ?", id).First(&user).Error
	return
}

func (e *Entity) FindByName(username string) (Entity, error) {
	var user Entity
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

// 获取用户分页列表
func GetUserList(page, limit int) (user []Entity, count int) {
	pageSize := limit
	offset := (page - 1) * pageSize
	query := db.Model(Entity{})
	query.Count(&count)
	query.Order("id desc").Limit(pageSize).Offset(offset).Find(&user)
	return
}

// 获取用户当前空间使用量
func GetImgdataSize(userId int64) (size int, err error) {
	rows, err := db.Select("sum(sizes) as total").Where("user_id = ?", userId).Rows()
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&size)
	}
	return
}

// 获取用户数量
func GetUserNum() (count int) {
	db.Where("level = 1").Count(&count)
	return
}

// 更新用户密码
func UpdateUserPwd(user *Entity) {
	db.Where("id = ?", user.ID).UpdateColumns(Entity{Salt: user.Salt, Password: user.Password})
}

// 更新用户状态
func UpdateUserStatus(id, status string) {
	db.Where("id = ?", id).Update("status", status)
}

// 更新用户头像
func UpdateUserAvatar(id, avatar string) {
	db.Where("id = ?", id).Update("avatar", avatar)
}

// 更新用户资料
func UpdateUserProfile(id int64, email, nickname string) error {
	return db.Where("id = ?", id).Updates(Entity{Nickname: nickname, Email: email}).Error
}

func GetUserIndexLine() (l []IndexLine,err error) {
	err = db.Raw("select DATE_FORMAT(created_at,'%Y-%m-%d') AS day, count(*) AS num from user where (level <> 99) group by day").Scan(&l).Error
	return
}