package repository_test

import (
	"fmt"
	"testing"

	"github.com/markdeesoft/golang-api/model"
	"github.com/markdeesoft/golang-api/repository"
	"github.com/stretchr/testify/assert"
)

func TestProductCategoryRepository_Store_And_GetByID(t *testing.T) {

	db := setupTestDB(t)

	repo := repository.NewProductCategoryRepository(db)

	// 1. ทดสอบการ Store (Insert)
	newCategory := &model.ProductCategory{
		Name: "NameTest",
	}

	err := repo.Store(newCategory)
	assert.NoError(t, err)
	assert.NotZero(t, newCategory.ID) // ควรสร้าง ID ให้อัตโนมัติ

	// 2. ทดสอบการ GetByID
	fetchedCategory, err := repo.GetByID(newCategory.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedCategory)
	assert.Equal(t, newCategory.Name, fetchedCategory.Name)
}

func TestProductCategoryRepository_GetAll_And_Count(t *testing.T) {

	db := setupTestDB(t)

	repo := repository.NewProductCategoryRepository(db)

	// เตรียมข้อมูล 5 รายการ
	for i := 1; i <= 5; i++ {
		category := &model.ProductCategory{
			Name: fmt.Sprintf("Category %d", i),
		}
		_ = repo.Store(category)
	}

	// 1. ทดสอบ Count
	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)

	// 2. ทดสอบ GetAll พร้อม Limit & Offset (Pagination)
	limit := 2
	offset := 0
	categories, err := repo.GetAll(limit, offset)

	assert.NoError(t, err)
	assert.Len(t, categories, 2)
	assert.Equal(t, "Category 1", categories[0].Name)
}

func TestProductCategoryRepository_Update(t *testing.T) {

	db := setupTestDB(t)

	repo := repository.NewProductCategoryRepository(db)

	// สร้างข้อมูลเริ่มต้น
	category := &model.ProductCategory{
		Name: "Old Name",
	}
	_ = repo.Store(category)

	// สั่ง Update
	updatedData := &model.ProductCategory{
		Name: "New Name",
	}

	result, err := repo.Update(category.ID, updatedData)
	assert.NoError(t, err)
	assert.Equal(t, "New Name", result.Name)

	// ตรวจสอบข้อมูลใน DB อีกครั้งเพื่อความมั่นใจ
	fetched, _ := repo.GetByID(category.ID)
	assert.Equal(t, "New Name", fetched.Name)
}

func TestProductCategoryRepository_Delete(t *testing.T) {

	db := setupTestDB(t)

	repo := repository.NewProductCategoryRepository(db)

	// สร้างข้อมูลเริ่มต้น
	category := &model.ProductCategory{
		Name: "To Be Deleted",
	}
	_ = repo.Store(category)

	// สั่ง Delete
	err := repo.Delete(category.ID)
	assert.NoError(t, err)

	// ลองดึงข้อมูลเพื่อดูว่าหายไปแล้วหรือยัง
	fetched, err := repo.GetByID(category.ID)
	assert.Error(t, err, "category data deleted") // ต้องคืน Error (Record Not Found)
	assert.Nil(t, fetched)
}
