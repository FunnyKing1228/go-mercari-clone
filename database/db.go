package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/FunnyKing1228/go-mercari-clone/models"
)

// Connect函式: 負責連線並回傳 DB 物件
func Connect() (*gorm.DB, error) {
	//這裡記得改成 host=db (Docker用的)
	dsn := "host=db user=mercari password=secret dbname=mercari_db port=5432 sslmode=disable TimeZone=Asia/Taipei"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	//自動遷移(Auto Migration)搬來這裡做
	//因為這裡引用了 models.Item，所以上面要 import models
	db.AutoMigrate(&models.Item{})

	return db, nil
}