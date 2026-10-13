package model

import (
	"backend/db"

	"github.com/jinzhu/gorm"
)

type CommentBase struct {
	Content string `gorm:"type:text;not null"`
	Likes   uint   `gorm:"default:0"` // cannot use 'type:unsigned' here

	CommenterID uint  `gorm:"not null;index"`
	Commenter   *User `gorm:"foreignKey:CommenterID"`
}

// WorkflowComment represents a comment on a workflow
type WorkflowComment struct {
	gorm.Model
	CommentBase

	// If this comment is a direct reply to a workflow, it will have a workflow ID
	WorkflowID *uint     `gorm:"index"`
	Workflow   *Workflow `gorm:"foreignKey:WorkflowID"`

	// If this comment is a secondary comment, it will have a parent comment ID
	// (i.e. reply to another workflow comment)
	ParentID *uint            `gorm:"index"`
	Parent   *WorkflowComment `gorm:"foreignKey:ParentID"`
	Children []*WorkflowComment

	ReplyTo     *uint
	ReplyToUser *User `gorm:"foreignKey:ReplyTo"`
}

type LikeRecordBase struct {
	CommentID uint `gorm:"not null;index"`
	UserID    uint `gorm:"not null;index"`
}

type WorkflowCommentLikeRecord struct {
	gorm.Model
	LikeRecordBase
}

type WorkflowCommentPublic struct {
	WorkflowComment
	LikeStatus bool
}

func init() {
	db.DB.AutoMigrate(&WorkflowComment{}, &WorkflowCommentLikeRecord{})
}
