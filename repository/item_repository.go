package repository

import (
	"github.com/FunnyKing1228/go-mercari-clone/models"
	"gorm.io/gorm"
)

//1. 定義介面 (Interface): 這就是合約!
//不管是真資料還是假資料庫，都必須提供這兩個功能
type ItemRepository interface {
	FindAll() ([]models.Item, error)
	Create(item *models.Item) error
}

//2. 實作真實的 Repository (負責跟 GORM 溝通)
type itemRepository struct {
	db *gorm.DB
}

// 建構子
func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepository{db : db}
}

func (r *itemRepository) FindAll() ([]models.Item, error){
	var items []models.Item
	err := r.db.Find(&items).Error
	return items, err
}

func (r *itemRepository) Create(item *models.Item) error {
	return r.db.Create(item).Error
}