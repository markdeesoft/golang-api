package repository

import (
	"log"

	"github.com/markdeesoft/golang-api/model"
	"gorm.io/gorm"
)

type ProductHashtagRepository interface {
	Count() (int64, error)
	GetByID(id uint) (*model.ProductHashtag, error)
	GetAll(limit int, offset int) ([]model.ProductHashtag, error)
	Store(product_hashtag *model.ProductHashtag) error
	Update(id uint, updatedData *model.ProductHashtag) (*model.ProductHashtag, error)
}

type productHashtagRepository struct {
	db *gorm.DB
}

func NewProductHashtagRepository(db *gorm.DB) ProductHashtagRepository {
	return &productHashtagRepository{db: db}
}

func (r *productHashtagRepository) Count() (int64, error) {

	var total int64
	result := r.db.Model(&model.ProductHashtag{}).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (r *productHashtagRepository) GetAll(limit, offset int) ([]model.ProductHashtag, error) {

	var product_hashtags []model.ProductHashtag

	// สั่งดึงข้อมูลตามหน้า LIMIT / OFFSET
	err := r.db.Limit(limit).Offset(offset).Order("id asc").Find(&product_hashtags).Error

	return product_hashtags, err

}

func (r *productHashtagRepository) GetByID(id uint) (*model.ProductHashtag, error) {

	var product_hashtag model.ProductHashtag

	result := r.db.First(&product_hashtag, id)
	if result.Error != nil {
		// log.Printf("Error creating query : %v", result.Error)
		return nil, result.Error
	}

	return &product_hashtag, nil
}

func (r *productHashtagRepository) Store(product_hashtag *model.ProductHashtag) error {

	result := r.db.Create(product_hashtag)
	if result.Error != nil {
		log.Printf("Failed to execute insert query: %v", result.Error)
		return result.Error
	}

	return nil
}

func (r *productHashtagRepository) Update(id uint, updatedData *model.ProductHashtag) (*model.ProductHashtag, error) {

	var product_hashtag model.ProductHashtag

	result := r.db.First(&product_hashtag, id)
	if result.Error != nil {
		log.Printf("Error creating query : %v", result.Error)
		return nil, result.Error
	}

	err := r.db.Model(&product_hashtag).Updates(model.ProductHashtag{
		Name: updatedData.Name,
	}).Error

	if err != nil {
		return nil, err
	}

	return &product_hashtag, nil
}
