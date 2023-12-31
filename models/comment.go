package models

import "gorm.io/gorm"

// Comment 评论
type Comment struct {
	gorm.Model
	Identity     string `gorm:"type:varchar(36) not null unique; comment:'唯一标识'" json:"identity" `
	UserIdentity string `gorm:"type:varchar(36) not null unique; comment:'用户唯一标识'"  json:"userIdentity"`
	BlogID       string `gorm:"type:varchar(36) not null ; comment:'blog唯一标识'"  json:"blogID"`
	ParentID     string `gorm:"type:varchar(36) not null ; comment:'父评论ID'"  json:"parentID"`
	Text         string `gorm:"type:varchar(36) not null ; comment:'评论内容'"  json:"text"`
	Ip           string `gorm:"type:varchar(36) not null ; comment:'ip地址'"  json:"ip"`
	IsTop        bool   `gorm:"type bool not null default = false ;  comment:'是否为父评论'"  json:"is_top"`
	ViolateRule  bool   `gorm:"type bool not null ; comment:'评论内容是否违规'"  json:"violateRule"`
}

func (c *Comment) TableName() string {
	return "comment_basic"
}
