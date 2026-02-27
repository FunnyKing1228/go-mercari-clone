package controllers

import (
	"github.com/FunnyKing1228/go-mercari-clone/database"
	"github.com/FunnyKing1228/go-mercari-clone/models"
	"github.com/FunnyKing1228/go-mercari-clone/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterInput 定義了使用者註冊時必須提供的資料格式
type RegisterInput struct {
	Username string `json:"username" binding:"required"`       //必填
	Email    string `json:"email" binding:"required,email"`    //必填，且必須符合 Email 格式!
	Password string `json:"password" binding:"required,min=6"` //必填，且密碼至少要 6 個字!
}

// LoginInput 定義了使用者登入時的資料格式
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 1. 定義客人登入時必須填寫的表單
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 註冊新使用者
// @Summary 註冊新帳號
// @Description 建立一組新的帳號密碼，密碼將會被 bcrypt 演算法單向加密
// @Tags auth
// @Accept json
// @Produce json
// @Param input body RegisterInput true "註冊資料 (需包含 email 與至少 6 碼的密碼)"
// @Success 201 {object} map[string]string
// @Router /register [post]
func Register(c *gin.Context) {
	var input RegisterInput

	// 1. 驗證前端傳來的資料是否符合我們剛才定義的規則
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "資料格式錯誤: " + err.Error()})
		return
	}

	// 2. 把明文密碼丟進 bcrypt 絞肉機 (Cost 設為 DefaultCost 10，數字越大越安全但也越慢)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "密碼加密過程發生錯誤"})
		return
	}

	// 3. 準備好要存入資料庫的 User 資料
	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashedPassword), // 👈 注意：我們存的是攪碎後的亂碼！
	}

	// 4. 存入 PostgreSQL！(如果帳號或信箱重複，這裡會報錯)
	// 註：這裡為了快速展示先用 database.DB，之後有空可以把它搬到 UserRepository 裡！
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": "帳號或信箱可能已被使用"})
		return
	}

	// 5. 成功回傳！
	c.JSON(201, gin.H{"message": "註冊成功！歡迎加入 Mercari Clone！"})
}

// Login 登入拿取 Token
// @Summary 登入系統取得手環
// @Description 輸入註冊過的帳號密碼來換取 JWT Token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body LoginInput true "請輸入帳號與密碼"
// @Success 200 {object} map[string]string
// @Router /login [post]
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "請完整輸入帳號與密碼"})
		return
	}

	// 1. 去資料庫尋找這個帳號
	var user models.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		// 💡 資安細節：就算查無此人，也要統一回傳「帳號或密碼錯誤」，避免駭客用這支 API 測試有哪些帳號存在
		c.JSON(401, gin.H{"error": "帳號或密碼錯誤"})
		return
	}

	// 2. 拿出 bcrypt 驗證器，比對「資料庫裡的亂碼」與「使用者剛輸入的明文」
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		c.JSON(401, gin.H{"error": "帳號或密碼錯誤"})
		return
	}

	// 3. 密碼正確！核發專屬的 JWT 手環 (這裡呼叫你 Phase 4 寫好的 Token 產生器)
	// 註：請確認你原本產生 Token 的函式名稱，如果叫 GenerateToken 就維持不變
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "無法生成 Token"})
		return
	}

	c.JSON(200, gin.H{"token": token})
}
