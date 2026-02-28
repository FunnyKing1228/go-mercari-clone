package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/FunnyKing1228/go-mercari-clone/database"
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

// FindAll 取得所有商品
// @Summary 取得商品列表
// @Description 支援 limit, offset 分頁與 search 關鍵字搜尋
// @Tags items
// @Accept json
// @Produce json
// @Param limit query int false "限制回傳筆數 (預設 10)"
// @Param offset query int false "跳過幾筆資料 (預設 0)"
// @Param search query string false "關鍵字搜尋 (例如: PS5)"
// @Success 200 {object} map[string]interface{}
// @Router /items [get]
func (ctrl *ItemController) FindAll(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 10
	}

	offsetStr := c.DefaultQuery("offset", "0")
	offset, _ := strconv.Atoi(offsetStr)
	if offset < 0 {
		offset = 0
	}

	// 👇 新增：取得網址列上的 search 參數 (例如 /items?search=ps5)
	search := c.Query("search")

	var items []models.Item
	if database.RedisClient != nil {
		// 👇 重要改動：把 search 字串也加入 Redis 的 Key 裡面！
		cacheKey := fmt.Sprintf("items:limit:%d:offset:%d:search:%s", limit, offset, search)

		if cachedData, err := database.RedisClient.Get(database.Ctx, cacheKey).Result(); err == nil {
			if err := json.Unmarshal([]byte(cachedData), &items); err == nil {
				fmt.Println("[極速] 從 Redis 快取取得資料!")
				c.IndentedJSON(http.StatusOK, gin.H{
					"source": "redis",
					"limit":  limit,
					"offset": offset,
					"search": search, // 回傳告訴前端現在搜了什麼
					"items":  items,
				})
				return
			}
		}
	}

	fmt.Println("[緩慢] 從 PostgreSQL 撈取資料...")
	// 👇 修改：把 search 傳進 Repo 裡
	items, err := ctrl.Repo.FindAll(limit, offset, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if database.RedisClient != nil {
		cacheKey := fmt.Sprintf("items:limit:%d:offset:%d:search:%s", limit, offset, search)
		if itemsJSON, err := json.Marshal(items); err == nil {
			database.RedisClient.Set(database.Ctx, cacheKey, itemsJSON, 60*time.Second)
		}
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"source": "database",
		"limit":  limit,
		"offset": offset,
		"search": search,
		"items":  items,
	})
}

// Create 新增商品
// @Summary 新增拍賣商品
// @Description 建立一個新的商品並存入資料庫，會自動綁定發送請求的賣家身分
// @Tags items
// @Accept json
// @Produce json
// @Param item body models.Item true "商品資料 (填寫 name 和 price 即可)"
// @Success 201 {object} models.Item
// @Security BearerAuth
// @Router /items [post]
func (ctrl *ItemController) Create(c *gin.Context) {
	var newItem models.Item
	if err := c.ShouldBindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// --- 👇 核心改動：從 JWT 驗證通過的 Context 中取出 user_id ---
	// 註：這裡的 "userID" 字串，必須與你 AuthMiddleware 裡 c.Set("xxx", id) 的 key 一致！
	rawUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "找不到使用者身分，請確認 Token 是否有效"})
		return
	}

	// 安全的型別轉換 (因為 JWT 解析完的數字很常變成 float64)
	switch v := rawUserID.(type) {
	case float64:
		newItem.UserID = uint(v)
	case uint:
		newItem.UserID = v
	case int:
		newItem.UserID = uint(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "身分驗證資料型別異常"})
		return
	}
	// -----------------------------------------------------------

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
