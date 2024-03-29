package models

import (
	"gorm.io/gorm"
)

// Blog 博客
type Blog struct {
	gorm.Model
	Identity     string `gorm:"type:varchar(36) not null unique; comment:'唯一标识'" json:"identity" `
	UserIdentity string `gorm:"type:varchar(36) not null; comment:'用户唯一标识'"  json:"userIdentity"`
	ImgUrl       string `gorm:"type:varchar(1000) ; comment:'图片url,为空表示没有图片'" json:"imgUrl"`
	Text         string `gorm:"type:varchar(2000); comment:'文本内容' " json:"text"`
	IsHide       bool   `gorm:"type:int ; comment:'文章是否私密'" json:"isHide"`
	Likes        int32  `gorm:"type:int ; comment:'点赞数量'" json:"likes"`
	IP           string `gorm:"type:varchar(64) not null ; comment:'IP地址'" json:"IP,omitempty"`
	ViolateRule  bool   `gorm:"type:bool not null ; comment:'文章内容是否违规'"  json:"violateRule"`
}

func (b *Blog) TableName() string {
	return "tally_blog_basic"
}
