package models

import "gorm.io/gorm"

type Kind struct {
	gorm.Model
	Name         string `gorm:"type:varchar(255);not null;unique;common:分类名称"`
	SerialNumber int    `gorm:"type:int;not null;unique;common:类别编号"`
}

func (Kind) TableName() string {
	return "kind_basic"
}
