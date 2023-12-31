package models

import (
	"gorm.io/gorm"
)

// Follow 关注
type Follow struct {
	gorm.Model
	Identity string `gorm:"type:varchar(36) not null unique ; comment:'关注记录唯一标识'" json:"identity,omitempty" json:"identity,omitempty"`
	UserId   string ` gorm:"type:varchar(36) not null unique ; comment:'用户唯一标识'"  json:"userId"`
	FollowId string `gorm:"type:varchar(36) not null unique ; comment:'关注用户唯一标识'" json:"followId,omitempty"`
}

func (c *Follow) TableName() string {
	return "follow_basic"
}
