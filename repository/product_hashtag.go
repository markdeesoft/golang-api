package repository

import (
	"log"

	"github.com/markdeesoft/golang-api/database"
	"github.com/markdeesoft/golang-api/model"
)

type ProductHashtagRepository interface {
	Count() (int64, error)
	GetByID(id uint) (*model.ProductHashtag, error)
	GetAll(limit int, offset int) ([]model.ProductHashtag, error)
	Store(product_hashtag *model.ProductHashtag) error
	Update(id uint, updatedData *model.ProductHashtag) (*model.ProductHashtag, error)
}

type productHashtagRepository struct{}

func NewProductHashtagRepository() ProductHashtagRepository {
	return &productHashtagRepository{}
}

func (r *productHashtagRepository) Count() (int64, error) {

	var total int64
	result := database.GDB.Model(&model.ProductHashtag{}).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (r *productHashtagRepository) GetAll(limit, offset int) ([]model.ProductHashtag, error) {

	var product_hashtags []model.ProductHashtag

	// สั่งดึงข้อมูลตามหน้า LIMIT / OFFSET
	err := database.GDB.Limit(limit).Offset(offset).Order("id asc").Find(&product_hashtags).Error

	return product_hashtags, err

}

func (r *productHashtagRepository) GetByID(id uint) (*model.ProductHashtag, error) {

	var product_hashtag model.ProductHashtag

	result := database.GDB.First(&product_hashtag, id)
	if result.Error != nil {
		// log.Fatalf("Error creating query : %v", result.Error)
		return nil, result.Error
	}

	return &product_hashtag, nil
}

func (r *productHashtagRepository) Store(product_hashtag *model.ProductHashtag) error {

	result := database.GDB.Create(product_hashtag)
	if result.Error != nil {
		log.Fatalf("Failed to execute insert query: %v", result.Error)
		return result.Error
	}

	return nil
}

func (r *productHashtagRepository) Update(id uint, updatedData *model.ProductHashtag) (*model.ProductHashtag, error) {

	var product_hashtag model.ProductHashtag

	result := database.GDB.First(&product_hashtag, id)
	if result.Error != nil {
		log.Fatalf("Error creating query : %v", result.Error)
		return nil, result.Error
	}

	err := database.GDB.Model(&product_hashtag).Updates(model.ProductHashtag{
		Name: updatedData.Name,
	}).Error

	if err != nil {
		return nil, err
	}

	return &product_hashtag, nil
}
