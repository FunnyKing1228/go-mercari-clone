package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//1. 定義商品長什麼樣子 (struct)
type Item struct {
	ID int `json:"id"`			//把GO的ID轉換成JSON的"id"
	Name string `json:"name"`	//把GO的Name轉換成JSON的"name"
	Price int `json:"price"`	//把GO的Price轉換成JSON的"price"
}

//2. 把items搬到main外面 變成全域變數
//這樣大家才能一起用它(模擬暫時的資料庫)
var items = []Item{
	{ID: 1, Name: "AirPods Pro", Price: 300},
	{ID: 2, Name: "iPhone 15", Price: 1000},
	{ID: 3, Name: "Pokemon Card", Price: 50},
}

func main() {
	r := gin.Default()

	//取得所有商品(GET)
	r.GET("/items", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{"items": items})
	})

	r.POST("/items", func(c *gin.Context){
		var newItem Item

		//BindJSON會幫我們做三件事:
		//A. 讀取使用者傳來的JSON
		//B. 檢查個是正不正確
		//C. 把資料塞進去 newItem 變數裡
		if err := c.ShouldBindJSON(&newItem); err != nil{
			//如果格式不對(例如傳了字串給 Price)，回傳 400 錯誤
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		//把新商品加入Slice
		items = append(items, newItem)
		c.JSON(http.StatusCreated, newItem)
	})
	
	

	r.Run()
}