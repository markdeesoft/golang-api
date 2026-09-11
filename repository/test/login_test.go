package repository_test

import (
	"testing"

	"github.com/markdeesoft/golang-api/model"
	"github.com/markdeesoft/golang-api/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserRepository_GetByUsername(t *testing.T) {

	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)

	// เตรียมข้อมูล User สำหรับทดสอบ
	user := &model.User{
		Name:     "Alice Smith",
		Email:    "alice@example.com",
		Phone:    "0898765432",
		Password: "password123",
	}
	err := repo.Store(user)
	assert.NoError(t, err)

	// 1. ค้นหาด้วย Email
	foundByEmail, err := repo.GetByUsername("alice@example.com")
	assert.NoError(t, err)
	assert.Equal(t, user.ID, foundByEmail.ID)

	// 2. ค้นหาด้วย Phone
	foundByPhone, err := repo.GetByUsername("0898765432")
	assert.NoError(t, err)
	assert.Equal(t, user.ID, foundByPhone.ID)

	// 3. ค้นหาด้วยค่าที่ไม่ไม่มีอยู่จริง
	notFound, err := repo.GetByUsername("unknown@example.com")
	assert.Error(t, err, "category data deleted")
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.Nil(t, notFound)
}
