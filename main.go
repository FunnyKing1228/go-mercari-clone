package main

import (
	"log"
	"net/http"
	//"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 1. 定義商品結構 (加上 GORM 的 tag)
type Item struct {
	gorm.Model
	Name string `json:"name"`
	Price int `json:"price"`	
}

// 全域變數用來存資料庫連線
var db *gorm.DB

func main() {

	// 2. 連線到資料庫 (這是新加入的)
	//我們從環境變數讀取資料庫連線資訊(這些變數在docker-compose.yml裡設定好了)
	dsn := "host=db user=mercari password=secret dbname=mercari_db port=5432 sslmode=disable TimeZone=Asia/Taipei"

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		//如果連線線不到資料庫，程式直接掛掉(Panic)，因為沒有DB什麼都做不了
		log.Fatal("Failed to connect to database:", err)
	}


	// 3. 自動遷移(Auto Migration)
	// 這行是 GORM 的魔法! 它會看 Item 結構長怎樣，自動去資料庫建立對應的 Table
	// 如果 Table 已經存在，它會檢查有沒有欄位變更
	db.AutoMigrate(&Item{})

	r := gin.Default()
	

	// 4. API實作 (改成操作 DB)

	// GET: 取得所有商品
	r.GET("/items", func(c *gin.Context){
		var items []Item
		
		// SELECT * FROM items;
		// 找出所有商品，把結果填入 items slice
		db.Find(&items)
		c.JSON(http.StatusOK, gin.H{"items": items})
	})

	// POST: 新增商品
	r.POST("/items", func(c *gin.Context){
		var newItem Item
		if err := c.ShouldBindJSON(&newItem) ; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// INSERT INTO items (name, price) VALUES (...);
		// GORM 會幫我們把 newItem 存進去，並自動生成 ID
		result := db.Create(&newItem)
		if result.Error != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
		c.JSON(http.StatusCreated, newItem)
	})


	// 為了方便測試，我們再加一個 "取得單一商品" 的API
	// GET /items/:id (例如 /items/1)
	r.GET("/items/:id", func(c *gin.Context){
		id := c.Param("id")
		var item Item
		
		// SELECT * FROM items WHERE id = ? LIMIT 1;
		// 如果找不到，FIRST 會回傳錯誤
		if err := db.First(&item, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.JSON(http.StatusOK, item)
	})

	r.Run() //預設
}