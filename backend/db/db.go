package db

import (
	"backend/config"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func init() {
	c := config.Config

	var err error
	DB, err = gorm.Open(c.Database.Driver, fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=true&loc=Local",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
		c.Database.Charset,
	))
	DB.LogMode(c.Server.Mode == "debug")
	if err != nil {
		panic(err)
	}
}
