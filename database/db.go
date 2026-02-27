package database

import (
	"fmt"
	"os"

	"github.com/FunnyKing1228/go-mercari-clone/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// defaultEnv 若環境變數為空則回傳預設值（本機開發時用）
func defaultEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// 暫時的(?) 全域且公開的 DB
var DB *gorm.DB

// Connect函式: 負責連線並回傳 DB 物件
func Connect() (*gorm.DB, error) {
	//1. 透過 os.Getenv 取得環境變數；本機執行時若沒設則用預設值（連到 localhost）
	//   Docker 內會用 .env 的 DB_HOST=db，本機 go run 時用預設 localhost
	host := defaultEnv("DB_HOST", "localhost")
	user := defaultEnv("DB_USER", "mercari")
	password := defaultEnv("DB_PASSWORD", "secret")
	dbname := defaultEnv("DB_NAME", "mercari_db")
	port := defaultEnv("DB_PORT", "5432")

	//2. 使用 fmt.Sprintf 將變數組裝成連線字串
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Taipei",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 把連線成功的區域變數 db，指派給剛剛宣告的公開變數 DB
	DB = db

	//自動遷移(Auto Migration)搬來這裡做
	//因為這裡引用了 models.Item，所以上面要 import models
	db.AutoMigrate(&models.Item{}, &models.User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
