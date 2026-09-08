package handler_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/markdeesoft/golang-api/handler"
	"github.com/markdeesoft/golang-api/model"
	"github.com/stretchr/testify/assert"
)

// --- Test GetByID (Found) ---
func TestGetUserByUsername_Success(t *testing.T) {
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
func TestGetUserByUsername_NotFound(t *testing.T) {
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
