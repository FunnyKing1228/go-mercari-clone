package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/FunnyKing1228/go-mercari-clone/models"
)

// ItemController 是一個結構，它持有 DB 連線
type ItemController struct {
	DB *gorm.DB
}

//NewItemController 是建構子，用來建立一個新的Controller
func NewItemController(db *gorm.DB) *ItemController {
	return &ItemController{DB: db}
}

//FindAll 對應原本的 GET /items
func (ctrl *ItemController) FindAll(c *gin.Context) {
	var items []models.Item //引用 models.Item
	ctrl.DB.Find(&items)
	c.JSON(http.StatusOK, gin.H{"items": items})
}

//Create 對應原本的POST /items
func (ctrl *ItemController) Create(c *gin.Context) {
	var newItem models.Item
	if err := c.ShouldBindJSON(&newItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.DB.Create(&newItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, newItem)
}