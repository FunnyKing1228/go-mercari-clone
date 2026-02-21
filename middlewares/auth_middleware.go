package middlewares

import (
		"net/http"
		"strings"

		"github.com/FunnyKing1228/go-mercari-clone/utils"
		"github.com/gin-gonic/gin"
)

// AuthMiddleware 就是我們的夜店安管
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 客人走到門口，安管要求看 Header 裡的 Authorization 憑證 (通常長這樣: "Bearer eyJ...")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "抱歉，你沒有戴 VIP 手環，禁止進入！"})
			c.Abort() // 🚨 無情踢出，後面的 Controller 絕對不會被執行
			return
		}

		// 2. 檢查手環格式 (規定必須是 "Bearer " 開頭加上 Token)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "手環格式錯誤！"})
			c.Abort()
			return
		}

		// 3. 把手環放進驗鈔機 (呼叫剛剛寫的 VerifyToken)
		tokenString := parts[1]
		token, err := utils.VerifyToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "手環是假的或已過期！安管踢人！"})
			c.Abort()
			return
		}

		// 4. 驗證通過！手環是真的！安管幫你開門，請進～
		c.Next()
	}
}