package models

import (
	"encoding/json"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(10);not null "`
	Password string `gorm:"type:varchar(36);not null "`
	Phone    string `gorm:"type:varchar(11) not null unique"`
	Identity string `gorm:"type:varchar(36) not null unique"`
	IP       string `gorm:"type:varchar(64) not null"`
}

func (u *User) TableName() string {
	return "user_basic"
}

func (u *User) MarshalBinary() ([]byte, error) {
	// 在这里编写将 User 类型转换为字节切片的逻辑
	return json.Marshal(u)
}
