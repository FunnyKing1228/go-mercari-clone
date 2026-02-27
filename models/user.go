package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string `gorm:"uniqueIndex;not null" json:"username"` //帳號必須唯一
	Email        string `gorm:"uniqueIndex;not null" json:"email"`    //信箱必須唯一
	PasswordHash string `gorm:"not null" json:"-"`                    //json:"-" 代表在轉成 JSON 傳給前端時，絕對不要把密碼交出去!
}
