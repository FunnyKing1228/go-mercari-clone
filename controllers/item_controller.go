package controllers

import (
	"fmt"
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
	limit := 10
	offset := 0
	if l := c.Query("limit"); l != "" { if n, _ := strconv.Atoi(l); n > 0 { limit = n } }
	if o := c.Query("offset"); o != "" { if n, _ := strconv.Atoi(o); n >= 0 { offset = n } }
	items, err := ctrl.Repo.FindAll(limit, offset)
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

func (ctrl *ItemController) UploadImage(c *gin.Context) {
	//1. 取得商品ID
	idStr := c.Param("id")
	itemID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的商品 ID"})
		return
	}

	//2. 從請求中拿出名為 "image" 的檔案
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "找不到上傳的檔案"})
		return
	}

	//3. 準備存檔路徑 (我們先在專案目錄下建一個 uploads 資料夾)
	// 檔名加上 ID 前綴避免重複，例如: 3_my_ps5.jpg
	fileName := fmt.Sprintf("%d_%s", itemID, file.Filename)
	dst := "uploads/" + fileName

	//4. 請 Gin 幫我們把檔案存到硬碟裡!
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "圖片存檔失敗"})
		return
	}

	//5. 將圖片網址存入資料庫 (假設我們的圖片網址開頭是 /uploads/)
	imageURL := "/" + dst
	if err := ctrl.Repo.UpdateImage(uint(itemID), imageURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新資料失敗"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "圖片上傳成功!",
		"image_url": imageURL,
	})
}
