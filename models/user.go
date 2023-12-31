package models

import (
	"gorm.io/gorm"
)

// User 用户
type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(10);not null ; comment:'用户名'" json:"username,omitempty"`
	Password string `gorm:"type:varchar(36);not null ; comment:'密码'" json:"password,omitempty"`
	Phone    string `gorm:"type:varchar(11) not null unique ; comment:'手机号'" json:"phone,omitempty"`
	Identity string `gorm:"type:varchar(36) not null unique ; comment:'唯一标识'" json:"identity,omitempty"`
	GithubId string `gorm:"type:varchar(36) not null unique ; comment:'Github账号'" json:"githubId,omitempty"`
	Status   int    `gorm:"type:int ; comment:'0表示正常, 1表示封禁'" json:"status"`
	IsHide   bool   `gorm:"type:int ; comment:'是否隐私账号'" json:"isHide"`
	IP       string `gorm:"type:varchar(64) not null ; comment:'IP地址'" json:"IP,omitempty"`
}

func (u *User) TableName() string {
	return "user_basic"
}
