package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

//創建表以及跟數據庫交互時使用的struct

type Member struct {
	gorm.Model

	Name string `gorm:"type:varchar(36);not null;unique"`

	Pwd string `gorm:"type:varchar(200);not null;unique"`

	Email Email

	IsSeller *bool `gorm:"default:0"`

	Balanc float64 `gorm:"type:double"`
}

type Email struct {
	MemberId int64 `gorm:"primary_key;autoIncrement:false"`

	CreatedAt time.Time

	UpdatedAt time.Time

	Addrs string `gorm:"type:varchar(255);not null;unique"`
}
