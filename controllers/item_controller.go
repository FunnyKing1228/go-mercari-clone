package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/FunnyKing1228/go-mercari-clone/repository"
	"github.com/FunnyKing1228/go-mercari-clone/models"
)

// 改動1: 這裡不再綁死 gorm.DB ， 而是依賴 Interface
type ItemController struct {
	Repo repository.ItemRepository
}

// 改動2: 建構子也改成接收 Interface
func NewItemController(repo repository.ItemRepository) *ItemController {
	return &ItemController{Repo: repo}
}

//FindAll 對應原本的 GET /items
func (ctrl *ItemController) FindAll(c *gin.Context) {
	//改動3: 呼叫 Repo 的方法，Controller 現在根本不知道背後是 Postgres 還是 MySQL
	items, err := ctrl.Repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

//Create 也要修改 把ctrl.DB.Create 改成 ctrl.Repo.Create
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