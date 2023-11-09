package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(10);not null "`
	Password string `gorm:"type:varchar(10);not null "`
	Phone    string `gorm:"type:varchar(11) not null unique"`
	Identity string `gorm:"type:varchar(36) not null unique"`
	IP       string `gorm:"type:varchar(64) not null"`
}

func (User) TableName() string {
	return "user_basic"
}
