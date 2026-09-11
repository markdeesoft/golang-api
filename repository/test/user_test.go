package repository_test

import (
	"fmt"
	"testing"

	"github.com/markdeesoft/golang-api/model"
	"github.com/markdeesoft/golang-api/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserRepository_Store_And_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)

	// 1. ทดสอบการ Store (Insert)
	newUser := &model.User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Phone:    "0812345678",
		Password: "hashedpassword123",
	}

	err := repo.Store(newUser)
	assert.NoError(t, err, "ควรจะสร้าง User ได้โดยไม่เกิด Error")
	assert.NotZero(t, newUser.ID, "ควรสร้าง ID ให้ User อัตโนมัติ")

	// 2. ทดสอบการ GetByID (ค้นหาเจอ)
	fetchedUser, err := repo.GetByID(newUser.ID)
	assert.NoError(t, err, "ควรค้นหา User ตาม ID เจอ")
	assert.NotNil(t, fetchedUser)
	assert.Equal(t, newUser.Name, fetchedUser.Name)
	assert.Equal(t, newUser.Email, fetchedUser.Email)

	// 3. ทดสอบการ GetByID (กรณีหาไม่เจอ)
	nonExistID := uint(999)
	notFoundUser, err := repo.GetByID(nonExistID)
	assert.Error(t, err, "ควรเกิด Error เมื่อหา ID ไม่พบ")
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound, " Error ควรเป็น gorm.ErrRecordNotFound")
	assert.Nil(t, notFoundUser)
}

func TestUserRepository_GetAll_And_Count(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)

	// สร้างข้อมูลจำลอง 5 รายการ
	for i := 1; i <= 5; i++ {
		user := &model.User{
			Name:     fmt.Sprintf("User %d", i),
			Email:    fmt.Sprintf("user%d@example.com", i),
			Phone:    fmt.Sprintf("080000000%d", i),
			Password: "password",
		}
		_ = repo.Store(user)
	}

	// 1. ทดสอบ Count
	total, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, int64(5), total)

	// 2. ทดสอบ GetAll พร้อม Pagination (Limit = 2, Offset = 0)
	users, err := repo.GetAll(2, 0)
	assert.NoError(t, err)
	assert.Len(t, users, 2, "ควรดึงข้อมูลมาแค่ 2 รายการตาม Limit")
	assert.Equal(t, "User 1", users[0].Name)

	// 3. ทดสอบ Pagination หน้าถัดไป (Limit = 2, Offset = 2)
	nextPageUsers, err := repo.GetAll(2, 2)
	assert.NoError(t, err)
	assert.Len(t, nextPageUsers, 2)
	assert.Equal(t, "User 3", nextPageUsers[0].Name)
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)

	user := &model.User{
		Name:     "Old Name",
		Email:    "old@example.com",
		Phone:    "0811111111",
		Password: "password",
	}
	_ = repo.Store(user)

	// ข้อมูลที่จะอัปเดต
	updatedInfo := &model.User{
		Name:  "New Name",
		Email: "new@example.com",
		Phone: "0822222222",
	}

	err := repo.Update(user.ID, updatedInfo)
	assert.NoError(t, err)

	// ดึงข้อมูลมาตรวจสอบอีกครั้ง
	fetched, err := repo.GetByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "New Name", fetched.Name)
	assert.Equal(t, "new@example.com", fetched.Email)
	assert.Equal(t, "0822222222", fetched.Phone)
}

func TestUserRepository_ResetPassword(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)

	user := &model.User{
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "oldpassword",
	}
	_ = repo.Store(user)

	// สั่ง Reset Password
	newPassword := "newsecretpassword"
	err := repo.ResetPassword(user.ID, newPassword)
	assert.NoError(t, err)

	// ตรวจสอบข้อมูลใน DB
	fetched, _ := repo.GetByID(user.ID)
	assert.Equal(t, newPassword, fetched.Password)
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewUserRepository(db)

	user := &model.User{
		Name:  "User To Delete",
		Email: "delete@example.com",
	}
	_ = repo.Store(user)

	// 1. สั่ง Delete (Soft Delete)
	err := repo.Delete(user.ID)
	assert.NoError(t, err)

	// 2. ลอง GetByID ดู ควรหาไม่เจอ
	fetched, err := repo.GetByID(user.ID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.Nil(t, fetched)

	// 3. ตรวจสอบ Count ควรเหลือ 0
	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}
