package models

import "gorm.io/gorm"

// Collect 收藏
type Collect struct {
	gorm.Model
	Identity     string `gorm:"type:varchar(36) not null unique; comment:'唯一标识'" json:"identity" `
	UserIdentity string `gorm:"type:varchar(36) not null unique; comment:'用户唯一标识'"  json:"userIdentity"`
	CollectId    string `gorm:"type:varchar(36) not null unique; comment:'被关注用户唯一标识'"  json:"collect_id"`
	BlogId       string `gorm:"type:varchar(36) not null unique; comment:'blog唯一标识'"  json:"blogID"`
}

func (c *Collect) TableName() string {
	return "collect_basic"
}
