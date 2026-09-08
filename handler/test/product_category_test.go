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

func TestGetAllProductCategory_Success(t *testing.T) {
	mockRepo := new(MockProductCategoryRepositoryTest)
	h := handler.NewProductCategoryHandler(mockRepo)

	app := fiber.New()
	app.Get("/productcategory", h.List)

	mockRepo.On("Count").Return(int64(2), nil)

	expectedData := []model.ProductCategory{
		{ID: 1, Name: "สุขภาพและความงาม"},
		{ID: 2, Name: "ของใช้ในบ้านและไลฟ์สไตล์"},
	}

	mockRepo.On("GetAll", mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/productcategory?limit=10&offset=0", nil)
	resp, err := app.Test(req, fiber.TestConfig{})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

// --- Test GetByID (Found) ---
func TestGetProductCategoryByID_Success(t *testing.T) {
	mockRepo := new(MockProductCategoryRepositoryTest)
	h := handler.NewProductCategoryHandler(mockRepo)

	app := fiber.New()
	app.Get("/productcategory/:id", h.View)

	mockProduct := &model.ProductCategory{
		ID:     1,
		Name:   "สุขภาพและความงาม",
		NameEn: "Health & Beauty",
	}

	// กำหนดว่าถ้าดึง ID=1 ให้คืน mockProduct
	mockRepo.On("GetByID", uint(1)).Return(mockProduct, nil)

	req := httptest.NewRequest(http.MethodGet, "/productcategory/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "สุขภาพและความงาม")
	mockRepo.AssertExpectations(t)
}

// --- Test GetByID (Not Found) ---
func TestGetProductCategoryByID_NotFound(t *testing.T) {
	mockRepo := new(MockProductCategoryRepositoryTest)
	h := handler.NewProductCategoryHandler(mockRepo)

	app := fiber.New()
	app.Get("/productcategory/:id", h.View)

	// กำหนดว่าถ้าดึง ID=99 ให้คืน error
	mockRepo.On("GetByID", uint(99)).Return(nil, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/productcategory/99", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

// --- Test Create ---
func TestCreateProductCategory_Success(t *testing.T) {

	mockRepo := new(MockProductCategoryRepositoryTest)
	h := handler.NewProductCategoryHandler(mockRepo)

	app := fiber.New()
	app.Post("/productcategory", h.Store)

	// กำหนด Expectation ว่า Create จะถูกเรียกด้วยอะไรก็ได้ที่สอดคล้อง แล้วไม่คืน error
	mockRepo.On("Store", mock.AnythingOfType("*model.ProductCategory")).Return(nil)

	jsonBody := `{
		"name":   "สุขภาพและความงาม",
		"name_en": "Health & Beauty"
	}`
	// "product_category_id": 1,
	// "product_category": {"name": "IT"}
	req := httptest.NewRequest(http.MethodPost, "/productcategory", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	// body, _ := io.ReadAll(resp.Body)
	// t.Logf("Response Status: %d, Body: %s", resp.StatusCode, string(body))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	mockRepo.AssertExpectations(t) // ตรวจสอบว่าถูกเรียกใช้จริงไหม
}

// --- Test Update ---
func TestUpdateProductCategory_Success(t *testing.T) {

	mockRepo := new(MockProductCategoryRepositoryTest)
	h := handler.NewProductCategoryHandler(mockRepo)

	app := fiber.New()
	app.Put("/productcategory/:id", h.Update)

	updatedCategory := &model.ProductCategory{ID: 1, Name: "สุขภาพและความงาม"}
	mockRepo.On("Update", uint(1), mock.AnythingOfType("*model.ProductCategory")).Return(updatedCategory, nil)

	jsonBody := `{
		"name":   "สุขภาพและความงาม",
		"name_en": "Health & Beauty"
	}`

	req := httptest.NewRequest(http.MethodPut, "/productcategory/1", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	body, _ := io.ReadAll(resp.Body)
	t.Logf("Response Status: %d, Body: %s", resp.StatusCode, string(body))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

// --- Test Delete Success ---
func TestDeleteProductCategory_Success(t *testing.T) {
	mockRepo := new(MockProductCategoryRepositoryTest)
	h := handler.NewProductCategoryHandler(mockRepo)

	app := fiber.New()
	app.Delete("/productcategory/:id", h.Delete)

	// ตั้ง Expectation: เมื่อรับ uint(1) ให้คืนค่า error = nil
	mockRepo.On("Delete", uint(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/productcategory/1", nil)
	resp, err := app.Test(req, fiber.TestConfig{})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode) // หรือ http.StatusNoContent (204) ตามที่ออกแบบไว้
	mockRepo.AssertExpectations(t)
}

// --- Test Delete Not Found ---
func TestDeleteProductCategory_NotFound(t *testing.T) {
	mockRepo := new(MockProductCategoryRepositoryTest)
	h := handler.NewProductCategoryHandler(mockRepo)

	app := fiber.New()
	app.Delete("/productcategory/:id", h.Delete)

	// ตั้ง Expectation: เมื่อดึง id ไม่มีจริง ให้คืนค่า error
	mockRepo.On("Delete", uint(99)).Return(gorm.ErrRecordNotFound)

	req := httptest.NewRequest(http.MethodDelete, "/productcategory/99", nil)
	resp, err := app.Test(req, fiber.TestConfig{})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}
