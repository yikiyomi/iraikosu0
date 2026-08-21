package main

import (
	"allto/config"
	"allto/database"
	"allto/model"
	"allto/router"
	"allto/util"
	"log"
	"net/http"
  _ "net/http/pprof"
	"go.uber.org/zap"
)


func main() {
	if err:=config.Load();err!=nil{
		log.Fatal("配置加载失败:",err)
	}

	//日志:为true时好看为false时json难看但性能好
	util.InitLogger(true)
    defer util.SyncLogger()

	// 连接mysql
	if err := database.InitDB(config.AppConfig.Database.DSN); err != nil {
    log.Fatal("mysql连接失败:", err)
    }
    if err := database.InitRedis(
		config.AppConfig.Redis.Addr,
		config.AppConfig.Redis.Password);
		err != nil {
	util.Logger.Error("Redis 连接失败，缓存功能降级",zap.Error(err))
    }

	// 自动建表迁移
	database.GetDB().AutoMigrate(
		&model.User{}, 
		&model.Post{}, 
		&model.Comment{}, 
		&model.Like{},
		&model.Follow{},
		&model.Notification{},
	)

	//pprof
	go func(){
		util.Logger.Info("pprof 启动",zap.String("addr",":6060"))
		if err :=http.ListenAndServe(":6060",nil);err!=nil{
			util.Logger.Error("pprof 启动失败",zap.Error(err))
		}
	}()
	//初始化gin
	r := router.SetupRouter()
	util.Logger.Info("服务启动",zap.String("port",config.AppConfig.Server.Port))
	r.Run(config.AppConfig.Server.Port)
}
