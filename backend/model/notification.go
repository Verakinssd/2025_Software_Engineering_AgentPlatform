package model

import (
	"backend/db"
	"github.com/jinzhu/gorm"
)

type Notification struct {
	gorm.Model
	Title      string `gorm:"type:varchar(255);not null"`
	Content    string `gorm:"type:text;not null"`
	Type       string `gorm:"type:varchar(255);not null"`
	FromUserId uint
	FromUser   *User `gorm:"foreignKey:FromUserId"`
	ToUserId   uint
	ToUser     *User `gorm:"foreignKey:ToUserId"`
	IsRead     bool  `gorm:"type:boolean;default:false"`
}

func init() {
	db.DB.AutoMigrate(&Notification{})
}
