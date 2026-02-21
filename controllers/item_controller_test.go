package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/FunnyKing1228/go-mercari-clone/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 1. 建立一個假的 Repository (Mock DB)
type MockItemRepository struct {
	mock.Mock
}

// 實作合約1: 假的 FindAll (這裡先留空實作，滿足介面即可)
func (m *MockItemRepository) FindAll() ([]models.Item, error){
	// m.Called() 會去抓你在測試裡設定好的假資料
	args := m.Called()
	return args.Get(0).([]models.Item), args.Error(1)
}

// 實作合約2: 假的 Create (這裡先留空實作，滿足介面即可)
func (m *MockItemRepository) Create(item *models.Item) error{
	args := m.Called(item)
	return args.Error(0)
}

// 2. 開始寫測試案例!
func TestItemController_FindAll_Success(t *testing.T) {
	//準備 Gin 的測試環境
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	//準備假資料
	mockItems := []models.Item{
		{Name: "PS5", Price: 15000},
		{Name: "Switch", Price: 8000},
	}

	//初始化 Mock Repo，並「教」它如何回話
	mockRepo := new(MockItemRepository)
	//白話文: 當有人呼叫 FindAll 時，請回傳 mockItems 和 nil(沒錯誤)
	mockRepo.On("FindAll").Return(mockItems, nil)

	//把假 Repo 塞給真 Controller
	controller := NewItemController(mockRepo)

	//執行測試動作
	controller.FindAll(c)

	// 3. 驗證結果 (Assert)
	assert.Equal(t, http.StatusOK, w.Code) //檢查 HTTP 狀態碼是不是 200
	assert.Contains(t, w.Body.String(), "PS5") //檢查回傳的 JSON 裡有沒有 PS5
	mockRepo.AssertExpectations(t) // 確認 Controller 真的有去呼叫FindAll
}

func TestItemController_FindAll_Error(t *testing.T) {
	//1. 準備環境
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	//2. 製造一個假的錯誤 (模擬資料庫斷線)
	mockError := errors.New("database connection lost")

	//3. 訓練 Mock Repo 扮演壞人
	mockRepo := new(MockItemRepository)
	//白話文: 當有人呼叫 FindAll 時，回傳空的資料，並丟出一個大錯誤！
	mockRepo.On("FindAll").Return([]models.Item{}, mockError)

	controller := NewItemController(mockRepo)

	// 4. 執行測試
	controller.FindAll(c)

	// 5. 驗證 Controller 有沒有好好把錯誤轉成 HTTP 500
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "database connection lost") // 檢查客人的畫面上有沒有出現錯誤訊息
	mockRepo.AssertExpectations(t)
}

func TestItemController_Create_Success(t *testing.T) {
	// ==========================================
	// 📍 第一部曲：Arrange (佈置片場)
	// ==========================================
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 1. 模擬客人傳來的 JSON 資料 (這段我幫你扛了)
	jsonBody := `{"name": "PS5", "price": 15000}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/items", bytes.NewBufferString(jsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// 2. 請演員 (Mock) 並發劇本
	mockRepo := new(MockItemRepository)
	
	// TODO: 告訴 mockRepo，當有人呼叫 "Create" 且參數不管傳什麼 (mock.Anything) 時，請回傳 nil (代表成功沒錯誤)
	// 你的程式碼：
	mockRepo.On("Create", mock.Anything).Return(nil)

	// 3. 依賴注入：把演員塞給主角
	// TODO: 初始化 NewItemController 並把 mockRepo 塞進去
	// 你的程式碼：
	controller := NewItemController(mockRepo)


	// ==========================================
	// 📍 第二部曲：Act (Action！正式開拍)
	// ==========================================
	// TODO: 讓 controller 實際去執行 Create 任務 (記得把 c 傳進去)
	// 你的程式碼：
	controller.Create(c)


	// ==========================================
	// 📍 第三部曲：Assert (導演看回放驗收)
	// ==========================================
	// TODO: 1. 檢查對客人的反應：HTTP 狀態碼是不是 201 (http.StatusCreated) 或 200 (http.StatusOK) 看你當初怎麼寫的
	// 你的程式碼：
	assert.Equal(t, http.StatusCreated, c.Writer.Status())
	// TODO: 2. 檢查演員有沒有摸魚：確保 Controller 真的有去呼叫 DB 的 Create
	// 你的程式碼：
	mockRepo.AssertExpectations(t)
}

func TestItemController_Create_BadInput(t *testing.T) {
	// ==========================================
	// 📍 第一部曲：Arrange (佈置片場)
	// ==========================================
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 1. 模擬奧客傳來的「壞掉的 JSON」 (price 故意給字串)
	badJsonBody := `{"name": "PS5", "price": "一萬五"}`
	c.Request, _ = http.NewRequest(http.MethodPost, "/items", bytes.NewBufferString(badJsonBody))
	c.Request.Header.Set("Content-Type", "application/json")

	// 2. 請演員 (Mock) 
	mockRepo := new(MockItemRepository)
	
	// 💡 注意：因為我們預期 Controller 會在第一關就擋下爛資料，
	// 所以這裡我們「不需要」設定 mockRepo.On("Create")... 的劇本！
	// (因為根本不該走到那一步)

	// 3. 依賴注入
	controller := NewItemController(mockRepo)


	// ==========================================
	// 📍 第二部曲：Act (Action！正式開拍)
	// ==========================================
	// TODO: 讓 controller 實際去執行 Create 任務
	// 你的程式碼：
    controller.Create(c)

	// ==========================================
	// 📍 第三部曲：Assert (導演看回放驗收)
	// ==========================================
	// TODO: 1. 檢查對客人的反應：因為客人亂傳資料，HTTP 狀態碼應該要是 400 (http.StatusBadRequest)
	// 你的程式碼：
	assert.Equal(t, http.StatusBadRequest, c.Writer.Status())


	// TODO: 2. 終極抓包：確保 Controller「沒有」笨笨地去呼叫 DB 的 Create！(使用 mockRepo.AssertNotCalled)
	// 你的程式碼：
	mockRepo.AssertNotCalled(t, "Create")
}