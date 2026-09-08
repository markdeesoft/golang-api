package handler_test

import (
	"database/sql"
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

func TestGetAllUser_Success(t *testing.T) {

	mockRepo := new(MockUserRepositoryTest)
	h := handler.NewUserHandler(mockRepo)

	app := fiber.New()
	app.Get("/user", h.List)

	mockRepo.On("Count").Return(int64(2), nil)

	expectedData := []model.User{
		{ID: 1, Name: "Emily Johnson", Phone: "+81 965-431-3024", Email: "emily.johnson@deesoft.com", Role: "admin"},
		{ID: 1, Name: "Michael Williams", Phone: "258-627-6644", Email: "michael.williams@deesoft.com"},
	}

	mockRepo.On("GetAll", mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(expectedData, nil)

	req := httptest.NewRequest(http.MethodGet, "/user?limit=10&offset=0", nil)
	resp, err := app.Test(req, fiber.TestConfig{})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

// --- Test GetByID (Found) ---
func TestGetUserByID_Success(t *testing.T) {
	mockRepo := new(MockUserRepositoryTest)
	h := handler.NewUserHandler(mockRepo)

	app := fiber.New()
	app.Get("/user/:id", h.View)

	mockUser := &model.User{
		ID:    1,
		Name:  "Emily Johnson",
		Phone: "+81 965-431-3024",
		Email: "emily.johnson@deesoft.com",
		Role:  "admin",
	}

	// กำหนดว่าถ้าดึง ID=1 ให้คืน mockUser
	mockRepo.On("GetByID", uint(1)).Return(mockUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/user/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Emily")
	mockRepo.AssertExpectations(t)
}

// --- Test GetByID (Not Found) ---
func TestGetUserByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepositoryTest)
	h := handler.NewUserHandler(mockRepo)

	app := fiber.New()
	app.Get("/user/:id", h.View)

	// กำหนดว่าถ้าดึง ID=99 ให้คืน error
	mockRepo.On("GetByID", uint(99)).Return(nil, errors.New("not found"))

	req := httptest.NewRequest(http.MethodGet, "/user/99", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

// --- Test Create ---
func TestCreateUser_Success(t *testing.T) {

	mockRepo := new(MockUserRepositoryTest)
	h := handler.NewUserHandler(mockRepo)

	app := fiber.New()
	app.Post("/user", h.Store)

	// กำหนด Expectation ว่า Create จะถูกเรียกด้วยอะไรก็ได้ที่สอดคล้อง แล้วไม่คืน error
	mockRepo.On("Store", mock.AnythingOfType("*model.User")).Return(nil)

	jsonBody := `{
		"name":  "Emily Johnson",
		"phone": "+81 965-431-3024",
		"email": "emily.johnson@deesoft.com",
		"role":  "admin"
	}`

	req := httptest.NewRequest(http.MethodPost, "/user", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	// body, _ := io.ReadAll(resp.Body)
	// t.Logf("Response Status: %d, Body: %s", resp.StatusCode, string(body))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	mockRepo.AssertExpectations(t) // ตรวจสอบว่าถูกเรียกใช้จริงไหม
}

// --- Test Update ---
func TestUpdateUser_Success(t *testing.T) {

	mockRepo := new(MockUserRepositoryTest)
	h := handler.NewUserHandler(mockRepo)

	app := fiber.New()
	app.Put("/user/:id", h.Update)

	// updated := &model.User{
	// 	ID:    1,
	// 	Name:  "Emily Johnson",
	// 	Phone: "+81 965-431-3024",
	// 	Email: "emily.johnson@deesoft.com",
	// 	Role:  "admin",
	// }
	mockRepo.On("Update", uint(1), mock.AnythingOfType("*model.User")).Return(nil)

	jsonBody := `{
		"name":  "Emily Johnson",
		"phone": "+81 965-431-3024",
		"email": "emily.johnson@deesoft.com",
		"role":  "admin"
	}`

	req := httptest.NewRequest(http.MethodPut, "/user/1", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	// body, _ := io.ReadAll(resp.Body)
	// t.Logf("Response Status: %d, Body: %s", resp.StatusCode, string(body))

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

// --- Test Delete Success ---
func TestDeleteUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepositoryTest)
	h := handler.NewUserHandler(mockRepo)

	app := fiber.New()
	app.Delete("/user/:id", h.Delete)

	// ตั้ง Expectation: เมื่อรับ uint(1) ให้คืนค่า error = nil
	mockRepo.On("Delete", uint(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/user/1", nil)
	resp, err := app.Test(req, fiber.TestConfig{})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode) // หรือ http.StatusNoContent (204) ตามที่ออกแบบไว้
	mockRepo.AssertExpectations(t)
}

// --- Test Delete Not Found ---
func TestDeleteUser_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepositoryTest)
	h := handler.NewUserHandler(mockRepo)

	app := fiber.New()
	app.Delete("/user/:id", h.Delete)

	// ตั้ง Expectation: เมื่อดึง id ไม่มีจริง ให้คืนค่า error
	mockRepo.On("Delete", uint(99)).Return(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodDelete, "/user/99", nil)
	resp, err := app.Test(req, fiber.TestConfig{})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

// --- Test Delete Test 500 Internal Server Error (Database Error) ---
func TestDeleteUser_InternalServerError(t *testing.T) {
	mockRepo := new(MockUserRepositoryTest)
	h := handler.NewUserHandler(mockRepo)

	app := fiber.New()
	app.Delete("/user/:id", h.Delete)

	// ตั้ง Expectation: เมื่อดึง id ไม่มีจริง ให้คืนค่า error
	mockRepo.On("Delete", uint(99)).Return(errors.New("db query execution failed"))

	req := httptest.NewRequest(http.MethodDelete, "/user/99", nil)
	resp, err := app.Test(req, fiber.TestConfig{})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

// --- 1. Test Case: ResetPassword สำเร็จ (200 OK) ---
func TestResetPassword_Success(t *testing.T) {
	mockRepo := new(MockUserRepositoryTest)
	h := handler.NewUserHandler(mockRepo)

	app := fiber.New()
	app.Post("/resetpass/:id", h.ResetPasswordUser)

	// Mock คืนค่า error = nil เมื่อรับ Token และ Password ถูกต้อง
	mockRepo.On("ResetPassword", uint(1)).Return(nil)

	jsonBody := `{}`

	req := httptest.NewRequest(http.MethodPost, "/resetpass/1", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, fiber.TestConfig{})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockRepo.AssertExpectations(t)
}

// --- Test Case: ResetPassword ไม่สำเร็จ - Validation Error (400 Bad Request) ---
// func TestResetPassword_ValidationError(t *testing.T) {

// 	mockRepo := new(MockUserRepositoryTest)
// 	h := handler.NewUserHandler(mockRepo)

// 	app := fiber.New()
// 	app.Post("/resetpass/:id", h.ResetPasswordUser)

// 	// Password สั้นเกินไป (ขัดกับ validate:"min=6")
// 	jsonBody := `{
// 		"token": "valid-token-123",
// 	}`

// 	req := httptest.NewRequest(http.MethodPost, "/resetpass/1", strings.NewReader(jsonBody))
// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := app.Test(req, fiber.TestConfig{})

// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
// 	// mockRepo ต้องไม่ถูกเรียกเลยเพราะติด validation ตั้งแต่แรก
// 	mockRepo.AssertNotCalled(t, "ResetPassword", mock.Anything, mock.Anything)
// }
