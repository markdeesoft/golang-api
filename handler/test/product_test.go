package handler_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/markdeesoft/golang-api/handler"
	"github.com/markdeesoft/golang-api/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestGetAllProduct_Success(t *testing.T) {

	mockRepo := new(MockProductRepositoryTest)
	h := handler.NewProductHandler(mockRepo)

	app := fiber.New()
	app.Get("/product", h.List)

	mockRepo.On("Count").Return(int64(2), nil)

	expectedData := []model.Product{
		{ID: 1, Name: "เอร่า", Price: 22000},
		{ID: 2, Name: "deesoft", Price: 400000},
	}

	mockRepo.On("GetAll", mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/product?limit=10&offset=0", nil)
	resp, err := app.Test(req, fiber.TestConfig{})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

// --- Test GetByID (Found) ---
func TestGetProductByID_Success(t *testing.T) {
	mockRepo := new(MockProductRepositoryTest)
	h := handler.NewProductHandler(mockRepo)

	app := fiber.New()
	app.Get("/product/:id", h.View)

	mockProduct := &model.Product{
		ID:     1,
		Name:   "เอร่า",
		NameEn: "Aera",
		Price:  2200,
	}

	// กำหนดว่าถ้าดึง ID=1 ให้คืน mockProduct
	mockRepo.On("GetByID", uint(1)).Return(mockProduct, nil)

	req := httptest.NewRequest(http.MethodGet, "/product/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Aera")
	mockRepo.AssertExpectations(t)
}

// --- Test GetByID (Not Found) ---
func TestGetProductByID_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepositoryTest)
	h := handler.NewProductHandler(mockRepo)

	app := fiber.New()
	app.Get("/product/:id", h.View)

	// กำหนดว่าถ้าดึง ID=99 ให้คืน error
	mockRepo.On("GetByID", uint(99)).Return(nil, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/product/99", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

// --- Test Create ---
func TestCreateProduct_Success(t *testing.T) {

	mockRepo := new(MockProductRepositoryTest)
	h := handler.NewProductHandler(mockRepo)

	app := fiber.New()
	app.Post("/product", h.Store)

	// กำหนด Expectation ว่า Create จะถูกเรียกด้วยอะไรก็ได้ที่สอดคล้อง แล้วไม่คืน error
	mockRepo.On("Store", mock.AnythingOfType("*model.Product")).Return(nil)

	jsonBody := `{
		"name":   "เอร่า",
		"name_en": "Aera",
		"price": 44000,
		"product_category_id": 1,
		"product_category": {"name": "สุขภาพและความงาม"}
	}`

	req := httptest.NewRequest(http.MethodPost, "/product", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	// body, _ := io.ReadAll(resp.Body)
	// t.Logf("Response Status: %d, Body: %s", resp.StatusCode, string(body))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	mockRepo.AssertExpectations(t) // ตรวจสอบว่าถูกเรียกใช้จริงไหม
}

// --- Test Update ---
func TestUpdateProduct_Success(t *testing.T) {

	mockRepo := new(MockProductRepositoryTest)
	h := handler.NewProductHandler(mockRepo)

	app := fiber.New()
	app.Put("/product/:id", h.Update)

	updated := &model.Product{
		ID:    1,
		Name:  "เอร่า",
		Price: 22000,
	}
	mockRepo.On("Update", uint(1), mock.AnythingOfType("*model.Product")).Return(updated, nil)

	jsonBody := `{
		"name":   "เอร่า",
		"name_en": "Aera",
		"price": 44000,
		"product_category_id": 1,
		"product_category": {"name": "สุขภาพและความงาม"}
	}`

	req := httptest.NewRequest(http.MethodPut, "/product/1", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	// body, _ := io.ReadAll(resp.Body)
	// t.Logf("Response Status: %d, Body: %s", resp.StatusCode, string(body))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

// --- Test Delete Success ---
func TestDeleteProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepositoryTest)
	h := handler.NewProductHandler(mockRepo)

	app := fiber.New()
	app.Delete("/product/:id", h.Delete)

	// ตั้ง Expectation: เมื่อรับ uint(1) ให้คืนค่า error = nil
	mockRepo.On("Delete", uint(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/product/1", nil)
	resp, err := app.Test(req, fiber.TestConfig{})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode) // หรือ http.StatusNoContent (204) ตามที่ออกแบบไว้
	mockRepo.AssertExpectations(t)
}

// --- Test Delete Not Found ---
func TestDeleteProduct_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepositoryTest)
	h := handler.NewProductHandler(mockRepo)

	app := fiber.New()
	app.Delete("/product/:id", h.Delete)

	// ตั้ง Expectation: เมื่อดึง id ไม่มีจริง ให้คืนค่า error
	mockRepo.On("Delete", uint(99)).Return(gorm.ErrRecordNotFound)

	req := httptest.NewRequest(http.MethodDelete, "/product/99", nil)
	resp, err := app.Test(req, fiber.TestConfig{})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}
