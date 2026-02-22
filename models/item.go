package models

import "gorm.io/gorm"

// 2. Item 結構定義 (跟原來一樣)
type Item struct {
	gorm.Model
	Name   string `json:"name"`
	Price  int    `json:"price"`
	Status string `json:"status" gorm:"default:'available'"` // available, sold
}
