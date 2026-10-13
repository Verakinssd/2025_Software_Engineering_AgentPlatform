package model

import (
	"backend/db"
	"github.com/jinzhu/gorm"
	"time"
)

type UserPenaltyRecord struct {
	gorm.Model
	AdminId  uint      `gorm:"not null"`
	Admin    *User     `gorm:"foreignKey:AdminId"`
	UserId   uint      `gorm:"not null"`
	User     *User     `gorm:"foreignKey:UserId"`
	Type     string    `gorm:"type:enum('warning', 'restrict_comment', 'restrict_agent_publish', 'temporary_ban', 'permanent_ban');not null"`
	Reason   string    `gorm:"type:varchar(255);not null"`
	ExpireAt time.Time `gorm:"not null"`
}

func init() {
	db.DB.AutoMigrate(&UserPenaltyRecord{})
}
