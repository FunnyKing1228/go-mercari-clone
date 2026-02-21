package main

import (
	"github.com/gin-gonic/gin"
	"github.com/FunnyKing1228/go-mercari-clone/database"
	"github.com/FunnyKing1228/go-mercari-clone/controllers"
	"github.com/FunnyKing1228/go-mercari-clone/repository"
	"github.com/FunnyKing1228/go-mercari-clone/middlewares"
)

func main() {
	// 1.連線到資料庫 (呼叫 database package)
	db, err := database.Connect()
	if err != nil {
		panic("Failed to connect to database!")
	}

	// 2.初始化 Controller (把 db 丟進去)
	itemController := controllers.NewItemController(repository.NewItemRepository(db)) 

	r := gin.Default()

	//註冊登入櫃檯
	r.POST("/login", controllers.Login)
	r.GET("/items", itemController.FindAll)
	
	//建立一個「VIP專區」，把安管 (AuthMiddleware) 派到這個區塊的門口
	protected := r.Group("/")
	protected.Use(middlewares.AuthMiddleware())//啟用安管
	{
		//只有通過安管的人，才能進入這個 VIP 專區執行新增商品
		protected.POST("/items", itemController.Create)
	}

	r.Run() // listen and serve on 0.0.0.0:8080
}