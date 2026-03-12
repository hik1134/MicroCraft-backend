package mysql

import (
	"fmt"

	"MicroCraft/internal/config"
	perr "MicroCraft/pkg/errors"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//全局数据库连接对象
var DB *gorm.DB

func InitDB() error {
	if config.Conf == nil {
		return perr.New(perr.CONFIG_NOT_INIT)
	}
	database := config.Conf.Mysql
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		database.User,
		database.Password,
		database.Host,
		database.Port,
		database.DBName,
	)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return perr.Wrap(perr.DB_CONNECT_FAIL, err)
	}
	fmt.Println("MySQL数据库连接成功")
	return nil
}