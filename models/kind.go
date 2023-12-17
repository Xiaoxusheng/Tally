package models

import (
	"gorm.io/gorm"
)

type Kind struct {
	gorm.Model
	Name         string `gorm:"type:varchar(255);not null;unique;common:分类名称" json:"name,omitempty"`
	SerialNumber int    `gorm:"type:int;not null;unique;common:类别编号" json:"serialNumber,omitempty"`
}

func (k *Kind) TableName() string {
	return "kind_basic"
}
