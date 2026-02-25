package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// 全域變數，讓 Controller 之後可以拿來用
var RedisClient *redis.Client
var Ctx = context.Background() // Redis 操作必備的 Context

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // 剛剛 Docker 啟動的 port
		Password: "",               // 預設沒有密碼
		DB:       0,                // 使用預設的 0 號資料庫
	})

	// 測試連線 (Ping 一下看它有沒有活著)
	pong, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("無法連線到 Redis: %v", err))
	}
	fmt.Printf("成功連線到 Redis! (回應: %s)\n", pong)
}
