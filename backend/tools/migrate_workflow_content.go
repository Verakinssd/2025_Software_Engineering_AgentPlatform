package main

import (
	"backend/config"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	// 加载配置
	c := config.Config

	// 连接数据库
	db, err := gorm.Open(c.Database.Driver, fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=true&loc=Local",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
		c.Database.Charset,
	))
	if err != nil {
		log.Fatal("无法连接数据库:", err)
	}
	defer db.Close()

	// 执行迁移SQL
	sql := "ALTER TABLE workflows MODIFY COLUMN content LONGTEXT NOT NULL"
	
	fmt.Println("开始执行工作流内容字段迁移...")
	fmt.Printf("执行SQL: %s\n", sql)
	
	if err := db.Exec(sql).Error; err != nil {
		log.Fatal("迁移失败:", err)
	}
	
	fmt.Println("迁移成功完成！")
	fmt.Println("工作流内容字段已从TEXT类型扩展为LONGTEXT类型")
	fmt.Println("现在可以存储更大的工作流数据（最大4GB）")
} 