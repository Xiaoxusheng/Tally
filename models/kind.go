package models

import (
	"encoding/json"
	"gorm.io/gorm"
)

type Kind struct {
	gorm.Model
	Name         string `gorm:"type:varchar(255);not null;unique;common:分类名称"`
	SerialNumber int    `gorm:"type:int;not null;unique;common:类别编号"`
}

func (k *Kind) TableName() string {
	return "kind_basic"
}

func (k *Kind) MarshalBinary() ([]byte, error) {
	// 在这里编写将 User 类型转换为字节切片的逻辑
	return json.Marshal(k)
}
