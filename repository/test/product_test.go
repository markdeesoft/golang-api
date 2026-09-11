package repository_test

import (
	"fmt"
	"testing"

	"github.com/markdeesoft/golang-api/model"
	"github.com/markdeesoft/golang-api/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestProductRepository_Store_And_GetByID(t *testing.T) {

	db := setupTestDB(t)

	category_id := storeAndGetCategoryID(t, db)

	repo := repository.NewProductRepository(db)

	// 1. ทดสอบการ Store (Insert)
	new_product := &model.Product{
		Name:              "NameTest",
		Price:             44000,
		ProductCategoryID: category_id,
	}

	err := repo.Store(new_product)
	assert.NoError(t, err)
	assert.NotZero(t, new_product.ID) // ควรสร้าง ID ให้อัตโนมัติ

	// 2. ทดสอบการ GetByID
	fetched, err := repo.GetByID(new_product.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetched)
	assert.Equal(t, new_product.Name, fetched.Name)
}

func TestProductRepository_GetAll_And_Count(t *testing.T) {

	var err error

	db := setupTestDB(t)

	category_id := storeAndGetCategoryID(t, db)

	repo := repository.NewProductRepository(db)

	// เตรียมข้อมูล 5 รายการ
	for i := 1; i <= 5; i++ {
		product := &model.Product{
			Name:              fmt.Sprintf("test %d", i),
			Price:             44000,
			ProductCategoryID: category_id,
		}
		err = repo.Store(product)
	}

	// 1. ทดสอบ Count
	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)

	// 2. ทดสอบ GetAll พร้อม Limit & Offset (Pagination)
	limit := 2
	offset := 0
	products, err := repo.GetAll(limit, offset)

	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "test 1", products[0].Name)
}

func TestProductRepository_Update(t *testing.T) {

	db := setupTestDB(t)

	category_id := storeAndGetCategoryID(t, db)

	repo := repository.NewProductRepository(db)

	// สร้างข้อมูลเริ่มต้น
	product := &model.Product{
		Name:              "Old Name",
		Price:             44000,
		ProductCategoryID: category_id,
	}
	_ = repo.Store(product)

	// สั่ง Update
	updatedData := &model.Product{
		Name:              "New Name",
		Price:             33000,
		ProductCategoryID: category_id,
	}

	result, err := repo.Update(product.ID, updatedData)
	assert.NoError(t, err)
	assert.Equal(t, "New Name", result.Name)
	assert.Equal(t, float64(33000), result.Price)
	assert.Equal(t, category_id, result.ProductCategoryID)

	// ตรวจสอบข้อมูลใน DB อีกครั้งเพื่อความมั่นใจ
	fetched, _ := repo.GetByID(product.ID)
	assert.Equal(t, "New Name", fetched.Name)
}

func TestProductRepository_Delete(t *testing.T) {

	db := setupTestDB(t)

	// category_id := storeAndGetCategoryID(t, db)

	repo := repository.NewProductRepository(db)

	// สร้างข้อมูลเริ่มต้น
	product := &model.Product{
		Name: "To Be Deleted",
		// Price:             33000,
		// ProductCategoryID: category_id,
	}
	_ = repo.Store(product)

	// สั่ง Delete
	err := repo.Delete(product.ID)
	assert.NoError(t, err)

	// ลองดึงข้อมูลเพื่อดูว่าหายไปแล้วหรือยัง
	fetched, err := repo.GetByID(product.ID)
	assert.Error(t, err) // ต้องคืน Error (Record Not Found)
	assert.Nil(t, fetched)
}

func storeAndGetCategoryID(t *testing.T, db *gorm.DB) uint {

	repo_category := repository.NewProductCategoryRepository(db)

	// 1. ทดสอบการ Store (Insert)
	newCategory := &model.ProductCategory{
		Name: "NameTest",
	}

	err := repo_category.Store(newCategory)
	assert.NoError(t, err)
	assert.NotZero(t, newCategory.ID) // ควรสร้าง ID ให้อัตโนมัติ

	return newCategory.ID
}
