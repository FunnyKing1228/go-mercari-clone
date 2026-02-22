package controllers

import (
	"net/http"
	"strconv"

	"github.com/FunnyKing1228/go-mercari-clone/models"
	"github.com/FunnyKing1228/go-mercari-clone/repository"
	"github.com/gin-gonic/gin"
)

// 改動1: 這裡不再綁死 gorm.DB ， 而是依賴 Interface
type ItemController struct {
	Repo repository.ItemRepository
}

// 改動2: 建構子也改成接收 Interface
func NewItemController(repo repository.ItemRepository) *ItemController {
	return &ItemController{Repo: repo}
}

// FindAll 對應原本的 GET /items
func (ctrl *ItemController) FindAll(c *gin.Context) {
	//改動3: 呼叫 Repo 的方法，Controller 現在根本不知道背後是 Postgres 還是 MySQL
	items, err := ctrl.Repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"items": items})
}

// Create 也要修改 把ctrl.DB.Create 改成 ctrl.Repo.Create
func (ctrl *ItemController) Create(c *gin.Context) {
	var newItem models.Item
	if err := c.ShouldBindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.Repo.Create(&newItem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newItem)
}

func (ctrl *ItemController) BuyItem(c *gin.Context) {
	//1. 從網址取得商品 ID (例如 /items/3/buy)
	idStr := c.Param("id")
	itemID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的商品 ID"})
		return
	}

	//實戰提醒: 真實情況下，這裡的 userID 應該從 JWT Token 裡解碼拿出來
	//為了專注在 Transaction 測試，這裡先假設購買者的 ID 是 1
	buyerID := uint(1)

	//2. 呼叫 Repository 執行帶有悲觀鎖的購買交易
	err = ctrl.Repo.BuyItem(uint(itemID), buyerID)
	if err != nil {
		//如果錯誤訊息是 "商品已被買走"，回傳 409 Conflict (狀態衝突)
		if err.Error() == "商品已被買走" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		//其他資料庫錯誤
		c.JSON(http.StatusInternalServerError, gin.H{"error": "購買失敗，請稍後再試"})
		return
	}

	//3. 交易成功
	c.JSON(http.StatusOK, gin.H{"message": "購買成功!"})

}
