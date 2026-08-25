package repository

import (
	"log"

	"github.com/markdeesoft/golang-api/database"
	"github.com/markdeesoft/golang-api/model"
)

type ProductHashtagRepository struct{}

// NewProductHashtagRepository ฟังก์ชันสร้างอินสแตนซ์ของ Repository
func NewProductHashtagRepository() *ProductHashtagRepository {
	return &ProductHashtagRepository{}
}

func (r *ProductHashtagRepository) Count() (int64, error) {

	var total int64
	result := database.GDB.Model(&model.ProductHashtag{}).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (r *ProductHashtagRepository) List(limit, offset int) ([]model.ProductHashtag, error) {

	var product_hashtags []model.ProductHashtag

	// สั่งดึงข้อมูลตามหน้า LIMIT / OFFSET
	err := database.GDB.Limit(limit).Offset(offset).Order("id asc").Find(&product_hashtags).Error

	return product_hashtags, err

}

func (r *ProductHashtagRepository) GetByID(id int) (*model.ProductHashtag, error) {

	var product_hashtag model.ProductHashtag

	result := database.GDB.First(&product_hashtag, id)
	if result.Error != nil {
		// log.Fatalf("Error creating query : %v", result.Error)
		return nil, result.Error
	}

	return &product_hashtag, nil
}

func (r *ProductHashtagRepository) Store(product_hashtag *model.ProductHashtag) error {

	result := database.GDB.Create(product_hashtag)
	if result.Error != nil {
		log.Fatalf("Failed to execute insert query: %v", result.Error)
		return result.Error
	}

	return nil
}

func (r *ProductHashtagRepository) Update(id string, updatedData *model.ProductHashtag) (*model.ProductHashtag, error) {

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

	// ส่งข้อมูลเวอร์ชันอัปเดตล่าสุดกลับไป
	return &product_hashtag, nil
}

// func (r *ProductHashtagRepository) Active(status bool, id int) error {

// 	var product_hashtag model.ProductHashtag

// 	result := database.GDB.First(&product_hashtag, id)
// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	return nil
// }
