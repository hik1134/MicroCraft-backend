package main

import (
	"errors"
	"fmt"
	"log"
	"MicroCraft/internal/config"
	"MicroCraft/internal/dao/mysql"
	red "MicroCraft/internal/dao/redis"
	"MicroCraft/internal/model"
	"MicroCraft/internal/router"
	perr "MicroCraft/pkg/errors"
)

func fatalInit(step string, err error) {
	log.Printf("%s失败: code=%s, cause=%v", step, perr.GetCode(err), errors.Unwrap(err))
	log.Fatalf("退出: %v", err) 
}

func main() {
	//初始化配置
	if err := config.InitConfig(); err != nil {
		fatalInit("配置初始化", err)
	}

	//初始化数据库
	if err := mysql.InitDB(); err != nil {
		fatalInit("数据库初始化", err)
	}

	//自动建表
	if err := mysql.DB.AutoMigrate(&model.User{}, &model.Work{}, &model.Post{}, &model.PostLike{}, &model.Carrier{}); err != nil {
		fatalInit("自动迁移数据库表", perr.Wrap(perr.DB_CREATE_FAIL, err))
	}

	//初始化 Redis
	if err := red.InitRedis(); err != nil {
		fatalInit("Redis 初始化", err)
	}

	//初始化路由
	r := router.InitRouter()

	//启动服务器
	port := config.Conf.Server.Port
	fmt.Printf("服务已启动，正在监听端口: %d...\n", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		fatalInit("服务器启动", err)
	}
}