package database

import (
	"fmt"
	"os"
	
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/FunnyKing1228/go-mercari-clone/models"
)

// Connect函式: 負責連線並回傳 DB 物件
func Connect() (*gorm.DB, error) {
	//1. 透過 os.Getenv 取得環境變數
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	//2. 使用 fmt.Sprintf 將變數組裝成連線字串

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Taipei",
	host, user, password, dbname, port)

	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	//自動遷移(Auto Migration)搬來這裡做
	//因為這裡引用了 models.Item，所以上面要 import models
	db.AutoMigrate(&models.Item{})

	return db, nil
}