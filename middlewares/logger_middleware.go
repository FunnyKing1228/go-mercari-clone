package middlewares

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// StructuredLogger 是一個 Gin Middleware，用來記錄每一筆 HTTP 請求
func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now() // 記錄客人進門的時間
		path := c.Request.URL.Path

		// c.Next() 讓客人進去消費 (執行對應的 Controller)
		c.Next()

		// 客人離開了，計算他待了多久
		latency := time.Since(start)

		// 使用 slog 寫下 JSON 格式的紀錄
		slog.Info("HTTP Request",
			slog.Int("status", c.Writer.Status()),
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("ip", c.ClientIP()),
			slog.String("latency", latency.String()),
		)
	}
}
