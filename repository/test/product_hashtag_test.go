package repository_test

import (
	"fmt"
	"testing"

	"github.com/markdeesoft/golang-api/model"
	"github.com/markdeesoft/golang-api/repository"
	"github.com/stretchr/testify/assert"
)

func TestProductHashtagRepository_Store_And_GetByID(t *testing.T) {

	db := setupTestDB(t)

	repo := repository.NewProductHashtagRepository(db)

	// 1. ทดสอบการ Store (Insert)
	newHashtag := &model.ProductHashtag{
		Name: "NameTest",
	}

	err := repo.Store(newHashtag)
	assert.NoError(t, err)
	assert.NotZero(t, newHashtag.ID) // ควรสร้าง ID ให้อัตโนมัติ

	// 2. ทดสอบการ GetByID
	fetchedHashtag, err := repo.GetByID(newHashtag.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedHashtag)
	assert.Equal(t, newHashtag.Name, fetchedHashtag.Name)
}

func TestProductHashtagRepository_GetAll_And_Count(t *testing.T) {

	db := setupTestDB(t)

	repo := repository.NewProductHashtagRepository(db)

	// เตรียมข้อมูล 5 รายการ
	for i := 1; i <= 5; i++ {
		hashtag := &model.ProductHashtag{
			Name: fmt.Sprintf("Hashtag %d", i),
		}
		_ = repo.Store(hashtag)
	}

	// 1. ทดสอบ Count
	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)

	// 2. ทดสอบ GetAll พร้อม Limit & Offset (Pagination)
	limit := 2
	offset := 0
	product_hashtags, err := repo.GetAll(limit, offset)

	assert.NoError(t, err)
	assert.Len(t, product_hashtags, 2)
	assert.Equal(t, "Hashtag 1", product_hashtags[0].Name)
}

func TestProductHashtagRepository_Update(t *testing.T) {

	db := setupTestDB(t)

	repo := repository.NewProductHashtagRepository(db)

	// สร้างข้อมูลเริ่มต้น
	hashtag := &model.ProductHashtag{
		Name: "Old Name",
	}
	_ = repo.Store(hashtag)

	// สั่ง Update
	updatedData := &model.ProductHashtag{
		Name: "New Name",
	}

	result, err := repo.Update(hashtag.ID, updatedData)
	assert.NoError(t, err)
	assert.Equal(t, "New Name", result.Name)

	// ตรวจสอบข้อมูลใน DB อีกครั้งเพื่อความมั่นใจ
	fetched, _ := repo.GetByID(hashtag.ID)
	assert.Equal(t, "New Name", fetched.Name)
}
