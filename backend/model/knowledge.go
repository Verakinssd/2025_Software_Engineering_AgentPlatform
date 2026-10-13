package model

import (
	"backend/db"

	"github.com/jinzhu/gorm"
)

const (
	KnowledgeStatusPrivate = "private"
	KnowledgeStatusPublish = "publish"
)

type Knowledge struct {
	gorm.Model
	Title       string `json:"title" gorm:"type:varchar(255);"`
	Description string `json:"description" gorm:"type:text;"`
	Status      string `json:"status" gorm:"type:varchar(255);"`
	UserId      uint
	BelongsTo   User   `gorm:"foreignKey:UserId"`
	Path        string `json:"path" gorm:"type:varchar(255);"`
}

// KnowledgeFile 表示上传到知识库的文件
type KnowledgeFile struct {
	gorm.Model
	FileName    string `json:"fileName" gorm:"type:varchar(255);"`  // 文件名
	FilePath    string `json:"filePath" gorm:"type:varchar(255);"`  // 存储路径
	FileSize    int64  `json:"fileSize" gorm:"type:bigint;"`        // 文件大小
	FileType    string `json:"fileType" gorm:"type:varchar(50);"`   // 文件类型
	UserId      uint   `json:"userId" gorm:"not null;index"`        // 上传用户ID
	Content     string `json:"content" gorm:"type:longtext;"`       // 文件内容 (文本类型)
	Description string `json:"description" gorm:"type:text;"`       // 文件描述
}

func init() {
	db.DB.AutoMigrate(&Knowledge{}, &KnowledgeFile{})
}
