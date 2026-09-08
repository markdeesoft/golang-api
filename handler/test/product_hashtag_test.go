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
)

func TestGetAllProductHashtag_Success(t *testing.T) {
	mockRepo := new(MockProductHashtagRepositoryTest)
	h := handler.NewProductHashtagHandler(mockRepo)

	app := fiber.New()
	app.Get("/producthash", h.List)

	mockRepo.On("Count").Return(int64(2), nil)

	expectedData := []model.ProductHashtag{
		{ID: 1, Name: "sale2026"},
		{ID: 2, Name: "tech"},
	}

	mockRepo.On("GetAll", mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/producthash?limit=10&offset=0", nil)
	resp, err := app.Test(req, fiber.TestConfig{})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

// --- Test GetByID (Found) ---
func TestGetProductHashtagByID_Success(t *testing.T) {
	mockRepo := new(MockProductHashtagRepositoryTest)
	h := handler.NewProductHashtagHandler(mockRepo)

	app := fiber.New()
	app.Get("/producthash/:id", h.View)

	mockProduct := &model.ProductHashtag{
		ID:     1,
		Name:   "ช้อปเลย",
		NameEn: "ShopNow",
	}

	// กำหนดว่าถ้าดึง ID=1 ให้คืน mockProduct
	mockRepo.On("GetByID", uint(1)).Return(mockProduct, nil)

	req := httptest.NewRequest(http.MethodGet, "/producthash/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "ShopNow")
	mockRepo.AssertExpectations(t)
}

// --- Test GetByID (Not Found) ---
func TestGetProductHashtagByID_NotFound(t *testing.T) {
	mockRepo := new(MockProductHashtagRepositoryTest)
	h := handler.NewProductHashtagHandler(mockRepo)

	app := fiber.New()
	app.Get("/producthash/:id", h.View)

	// กำหนดว่าถ้าดึง ID=99 ให้คืน error
	mockRepo.On("GetByID", uint(99)).Return(nil, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/producthash/99", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

// --- Test Create ---
func TestCreateProductHashtag_Success(t *testing.T) {

	mockRepo := new(MockProductHashtagRepositoryTest)
	h := handler.NewProductHashtagHandler(mockRepo)

	app := fiber.New()
	app.Post("/producthash", h.Store)

	// กำหนด Expectation ว่า Create จะถูกเรียกด้วยอะไรก็ได้ที่สอดคล้อง แล้วไม่คืน error
	mockRepo.On("Store", mock.AnythingOfType("*model.ProductHashtag")).Return(nil)

	jsonBody := `{
		"name":   "ช้อปเลย",
		"name_en": "ShopNow"
	}`
	req := httptest.NewRequest(http.MethodPost, "/producthash", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	mockRepo.AssertExpectations(t) // ตรวจสอบว่าถูกเรียกใช้จริงไหม
}

// --- Test Update ---
func TestUpdateProductHashtag_Success(t *testing.T) {

	mockRepo := new(MockProductHashtagRepositoryTest)
	h := handler.NewProductHashtagHandler(mockRepo)

	app := fiber.New()
	app.Put("/producthashtag/:id", h.Update)

	updatedHashtag := &model.ProductHashtag{ID: 1, Name: "ช้อปเลย", IsActive: true}
	mockRepo.On("Update", uint(1), mock.AnythingOfType("*model.ProductHashtag")).Return(updatedHashtag, nil)

	jsonBody := `{
		"name":   "ช้อปเลย",
		"name_en": "ShopNow",
		"is_active": true
	}`

	req := httptest.NewRequest(http.MethodPut, "/producthashtag/1", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	body, _ := io.ReadAll(resp.Body)
	t.Logf("Response Status: %d, Body: %s", resp.StatusCode, string(body))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}
