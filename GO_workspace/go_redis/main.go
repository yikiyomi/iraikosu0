package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"context"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	})

	// 测试连接
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("连接失败:", err)
		return
	}
	fmt.Println("连接成功:", pong)

	// 设置值
	err = rdb.Set(context.Background(), "name", "李四", 0).Err()
	if err != nil {
		fmt.Println("设置失败:", err)
	}

	// 取值
	val, err := rdb.Get(context.Background(), "name").Result()
	if err != nil {
		fmt.Println("获取失败:", err)
	}
	fmt.Println("name =", val)
}