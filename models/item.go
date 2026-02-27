package models

import "gorm.io/gorm"

type Item struct {
	// gorm.Model 只在資料庫層使用，Swagger 不需要知道這些欄位
	gorm.Model `swaggerignore:"true"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Status   string `json:"status" gorm:"default:'available'"` // available, sold, deleted (刪除狀態)
	ImageURL string `json:"image_url"`
}
