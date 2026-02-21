package controllers

import (
	"net/http"

	"github.com/FunnyKing1228/go-mercari-clone/utils"
	"github.com/gin-gonic/gin"
)

// 1. 定義客人登入時必須填寫的表單
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 2. 登入櫃台的處理邏輯
func Login(c *gin.Context) {
	var req LoginRequest

	// 檢查客人有沒有乖乖填寫表單
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "請輸入帳號密碼"})
		return
	}

	// 3. 模擬資料庫驗證 (這裡先寫死一組通關密碼)
	if req.Username == "admin" && req.Password == "mercari" {
		// 帳密正確！呼叫剛剛寫好的工具，發行 VIP 手環 (我們假設這個管理員的 UserID 是 1)
		tokenString, err := utils.GenerateToken(1)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "手環發行失敗，機器故障"})
			return
		}

		// 雙手把手環 (Token) 奉上給客人
		c.JSON(http.StatusOK, gin.H{
			"message": "登入成功！這是你的 VIP 手環",
			"token":   tokenString,
		})
		return
	}

	// 4. 帳密錯誤，安管直接把人踢走 (401 Unauthorized)
	c.JSON(http.StatusUnauthorized, gin.H{"error": "帳號或密碼錯誤"})
}