package models

import (
	"gorm.io/gorm"
	"time"
)

type Share struct {
	gorm.Model
	Identity     string    `gorm:"type:varchar(36) not null unique; comment:'唯一标识'" json:"identity" `
	UserIdentity string    `gorm:"type:varchar(36) not null; comment:'用户唯一标识'"  json:"userIdentity"`
	ImgUrl       string    `gorm:"type:varchar(1000) ; comment:'图片url,为空表示没有图片'" json:"imgUrl"`
	Text         string    `gorm:"type:varchar(2000); comment:'文本内容' " json:"text"`
	ExpiredTime  time.Time `gorm:"type:int; comment:'过期时间' " json:"expired_time"`
}
