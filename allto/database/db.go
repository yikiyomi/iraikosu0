package database

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)
// 全局变量
var (
	db  *gorm.DB
	rdb *redis.Client
	ctx = context.Background()
)

func InitDB(dsn string) error {
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
}

func InitRedis(addr, password string) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	return rdb.Ping(ctx).Err()
}
func GetDB() *gorm.DB { return db }
func GetRedis() *redis.Client { return rdb }
func GetCtx() context.Context {return ctx}