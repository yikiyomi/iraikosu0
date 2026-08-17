package main

import (
	"log"
	"allto/config"
	"allto/database"
	"allto/model"
	"allto/router"
)


func main() {
	if err:=config.Load();err!=nil{
		log.Fatal("配置加载失败:",err)
	}
	// 连接mysql
	if err := database.InitDB(config.AppConfig.Database.DSN); err != nil {
    log.Fatal("mysql连接失败:", err)
    }
    if err := database.InitRedis(
		config.AppConfig.Redis.Addr,
		config.AppConfig.Redis.Password);
		err != nil {
    log.Printf("警告: Redis 连接失败，缓存功能降级: %v", err)
    }

	// 自动建表迁移
	database.GetDB().AutoMigrate(
		&model.User{}, 
		&model.Post{}, 
		&model.Comment{}, 
		&model.Like{},
		&model.Follow{},
	)

	//初始化gin
	r := router.SetupRouter()
	log.Println("服务器启动:"+config.AppConfig.Server.Port)
	r.Run(config.AppConfig.Server.Port)
}
