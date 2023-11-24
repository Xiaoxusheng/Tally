package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(10);not null " json:"username,omitempty"`
	Password string `gorm:"type:varchar(36);not null " json:"password,omitempty"`
	Phone    string `gorm:"type:varchar(11) not null unique" json:"phone,omitempty"`
	Identity string `gorm:"type:varchar(36) not null unique" json:"identity,omitempty"`
	GithubId string `gorm:"type:varchar(36) not null unique" json:"githubId,omitempty"`
	IP       string `gorm:"type:varchar(64) not null" json:"IP,omitempty"`
}

func (u *User) TableName() string {
	return "user_basic"
}

//func (u *User) MarshalBinary() ([]byte, error) {
//	// 在这里编写将 User 类型转换为字节切片的逻辑
//	return json.Marshal(u)
//}
