package main

import (
	"github.com/gin-gonic/gin"
	"github.com/FunnyKing1228/go-mercari-clone/database"
	"github.com/FunnyKing1228/go-mercari-clone/controllers"
	"github.com/FunnyKing1228/go-mercari-clone/repository"
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

	// 3.設定路由 (呼叫 controller 的方法)
	r.GET("/items", itemController.FindAll)
	r.POST("/items", itemController.Create)

	r.Run() // listen and serve on 0.0.0.0:8080
}