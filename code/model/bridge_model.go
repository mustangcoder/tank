package model

import (
	"github.com/eyebluecn/tank/code/constant"
)

/**
 * the link table for Share and Matter.
 */
type Bridge struct {
	Base
	ShareUuid  string `json:"shareUuid" gorm:"type:char(36)"`
	MatterUuid string `json:"matterUuid" gorm:"type:char(36)"`
}

func (this *Bridge) TableName() string {
	return constant.TABLE_PREFIX + "bridge"
}
