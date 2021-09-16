package model

import (
	"github.com/eyebluecn/tank/code/constant"
	"time"
)

type User struct {
	Base
	Role           string    `json:"role" gorm:"type:varchar(45)"`
	Username       string    `json:"username" gorm:"type:varchar(45) not null;unique"`
	Password       string    `json:"-" gorm:"type:varchar(255)"`
	AvatarUrl      string    `json:"avatarUrl" gorm:"type:varchar(255)"`
	LastIp         string    `json:"lastIp" gorm:"type:varchar(128)"`
	LastTime       time.Time `json:"lastTime" gorm:"type:timestamp not null;default:'2018-01-01 00:00:00'"`
	SizeLimit      int64     `json:"sizeLimit" gorm:"type:bigint(20) not null;default:-1"`
	TotalSizeLimit int64     `json:"totalSizeLimit" gorm:"type:bigint(20) not null;default:-1"`
	TotalSize      int64     `json:"totalSize" gorm:"type:bigint(20) not null;default:0"`
	Status         string    `json:"status" gorm:"type:varchar(45)"`
}

// set User's table name to be `profiles`
func (this *User) TableName() string {
	return constant.TABLE_PREFIX + "user"
}
