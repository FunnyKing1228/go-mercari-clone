package repository

import (
	"errors"

	"github.com/FunnyKing1228/go-mercari-clone/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 1. 定義介面 (Interface): 這就是合約!
// 不管是真資料還是假資料庫，都必須提供這兩個功能
type ItemRepository interface {
	FindAll(limit int, offset int, search string) ([]models.Item, error)
	Create(item *models.Item) error
	BuyItem(itemID uint, userID uint) error
	UpdateImage(itemID uint, imageURL string) error
}

// 2. 實作真實的 Repository (負責跟 GORM 溝通)
type itemRepository struct {
	db *gorm.DB
}

// 建構子
func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepository{db: db}
}

func (r *itemRepository) FindAll(limit int, offset int, search string) ([]models.Item, error) {
	var items []models.Item

	// 1. 先建立一個包含賣家資訊的基礎查詢
	query := r.db.Preload("User")

	// 2. 如果有傳入搜尋字串，就疊加模糊搜尋條件 (% 代表前後可以有任何字元)
	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	// 3. 最後加上分頁限制，並執行查詢
	err := query.Limit(limit).Offset(offset).Find(&items).Error

	return items, err
}

func (r *itemRepository) Create(item *models.Item) error {
	return r.db.Create(item).Error
}

func (r *itemRepository) BuyItem(itemID uint, userID uint) error {
	//開啟 GORM 的 Transaction
	return r.db.Transaction(func(tx *gorm.DB) error {
		var item models.Item

		//1. 悲觀鎖 (Pessimistic Lock): 使用 FOR UPDATE 鎖定這筆資料
		//在這個 Transaction 結束前，其他任何想讀取或修改這筆 Item 的請求都會被阻擋等待
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&item, itemID).Error; err != nil {
			return err //找不到商品
		}

		//2. 檢查商品狀態
		if item.Status == "sold" {
			return errors.New("商品已被買走")
		}

		//3. 執行購買邏輯 (更新狀態，真實世界還會在這裡建立 Order 訂單邏輯)
		item.Status = "sold"
		if err := tx.Save(&item).Error; err != nil {
			return err //更新失敗，將自動 Rollback
		}

		//回傳 nil 代表沒有錯誤， GORM 會自動幫我們 Commit 這次的交易
		return nil
	})
}

func (r *itemRepository) UpdateImage(itemID uint, imageURL string) error {
	// 告訴 GORM: 去 items 表格，找到對應的 ID，把 image_url 欄位更新
	return r.db.Model(&models.Item{}).Where("id = ?", itemID).Update("image_url", imageURL).Error
}
