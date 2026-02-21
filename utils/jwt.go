package utils

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret 就像是夜店老闆的「防偽印章」
// 注意：真實世界裡，這個絕對不能寫死在程式碼裡，要放在 .env 變數，不然駭客拿到印章就能自己印手環了！
var jwtSecret = []byte("super-secret-mercari-key")

// GenerateToken：發行 VIP 手環 (傳入使用者的 ID，回傳一段加密字串)
func GenerateToken(userID uint) (string, error) {
	// 1. 填寫手環上的資訊 (在 JWT 裡稱為 Claims)
	claims := jwt.MapClaims{
		"user_id": userID,                                  // 記錄這是哪個客人的 ID
		"exp":     time.Now().Add(time.Hour * 24).Unix(),   // 手環有效期限：24 小時後過期
	}

	// 2. 選擇防偽簽名的演算法 (HS256 是一種對稱加密算法) 並建立手環本體
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 3. 蓋上老闆的私密印章 (jwtSecret)，把它壓縮成最終那一長串的亂碼字串
	return token.SignedString(jwtSecret)
}

// VerifyToken: 安管的驗鈔機，用來檢查手環的防偽簽名和是否過期
func VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token)(interface{}, error){
		//拿出老闆的私密印章來核對簽名
		return jwtSecret, nil
	})
	return token, err
}