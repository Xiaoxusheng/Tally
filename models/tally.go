package models

import (
	"gorm.io/gorm"
)

type Tally struct {
	gorm.Model
	Identity     string  `gorm:"type:varchar(36) not null unique; comment:'唯一标识'"  `
	UserIdentity string  `gorm:"type:varchar(36) not null ; comment:'用户唯一标识'"  `
	Kind         int     `gorm:"type int not null  default=0  ;comment:'收入支出种类'" `
	Money        float64 `gorm:"type float not null  default=0 ; comment:'金额'" `
	Remark       string  `gorm:"type:varchar(255) not null; comment:'备注'"`
	Category     int     `gorm:"type:int not null; comment:'类别'"`
}

func (t *Tally) TableName() string {
	return "tally_basic"
}

//func (t *Tally) MarshalBinary() ([]byte, error) {
//	return json.Marshal(t)
//}
