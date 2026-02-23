package main

import (
	"log/slog"
	"os"

	"github.com/FunnyKing1228/go-mercari-clone/controllers"
	"github.com/FunnyKing1228/go-mercari-clone/database"
	"github.com/FunnyKing1228/go-mercari-clone/middlewares"
	"github.com/FunnyKing1228/go-mercari-clone/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	//在程式一啟動時，設定 slog 全域使用 JSON 格式輸出
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 1.連線到資料庫 (呼叫 database package)
	db, err := database.Connect()
	if err != nil {
		panic("Failed to connect to database!")
	}

	// 2.初始化 Controller (把 db 丟進去)
	itemController := controllers.NewItemController(repository.NewItemRepository(db))

	//開除預設的紀錄員，把原來的 r := gin.Default() 改成下面這些:
	r := gin.New()
	r.Use(gin.Recovery())                 // Gin 內建的防崩潰保鏢
	r.Use(middlewares.StructuredLogger()) // 我們剛剛自己寫的 JSON 紀錄員

	//註冊登入櫃檯
	r.POST("/login", controllers.Login)
	r.GET("/items", itemController.FindAll)

	//1. 【展示櫃】這行很重要：告訴 Gin ，只要網址是 /uploads 開頭，就去硬碟的 ./uploads 資料夾找檔案
	r.Static("/uploads", "./uploads")

	//先用 Go 內建套件確保 uploads 資料夾存在，免得存檔時報錯 (記得 import "os")
	os.MkdirAll("uploads", os.ModePerm)

	//建立一個「VIP專區」，把安管 (AuthMiddleware) 派到這個區塊的門口
	protected := r.Group("/")
	protected.Use(middlewares.AuthMiddleware()) //啟用安管
	{
		//只有通過安管的人，才能進入這個 VIP 專區執行新增商品
		protected.POST("/items", itemController.Create)
		//處理購買動作的路由
		protected.POST("/items/:id/buy", itemController.BuyItem)

		//2. 【收發室】加上上傳圖片的路由 (需要帶手環才能上傳)
		protected.POST("/items/:id/image", itemController.UploadImage)
	}

	r.Run() // listen and serve on 0.0.0.0:8080
}
