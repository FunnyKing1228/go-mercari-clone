package models

import "gorm.io/gorm"

type Item struct {
	// gorm.Model 只在資料庫層使用，Swagger 不需要知道這些欄位
	gorm.Model `swaggerignore:"true"`
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Status     string `json:"status" gorm:"default:'available'"` // available, sold, deleted (刪除狀態)
	ImageURL   string `json:"image_url"`

	// 加兩行來建立賣家關聯
	UserID uint `json:"user_id"`                         // 外鍵: 紀錄是哪個 User (ID) 賣的
	User   User `json:"seller" gorm:"foreignKey:UserID"` // 讓 GORM 知道要用 UserID 去關聯 User 表格
}
