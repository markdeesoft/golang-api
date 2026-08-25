package repository

import (
	"log"

	"github.com/markdeesoft/golang-api/database"
	"github.com/markdeesoft/golang-api/model"
)

type ProductCategoryRepository struct{}

// NewProductCategoryRepository ฟังก์ชันสร้างอินสแตนซ์ของ Repository
func NewProductCategoryRepository() *ProductCategoryRepository {
	return &ProductCategoryRepository{}
}

func (r *ProductCategoryRepository) Count() (int64, error) {

	var total int64
	result := database.GDB.Model(&model.ProductCategory{}).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (r *ProductCategoryRepository) List(limit, offset int) ([]model.ProductCategory, error) {

	var product_categories []model.ProductCategory

	// สั่งดึงข้อมูลตามหน้า LIMIT / OFFSET
	err := database.GDB.Limit(limit).Offset(offset).Order("id asc").Find(&product_categories).Error

	return product_categories, err

}

func (r *ProductCategoryRepository) GetByID(id int) (*model.ProductCategory, error) {

	var product_category model.ProductCategory

	result := database.GDB.First(&product_category, id)
	if result.Error != nil {
		// log.Fatalf("Error creating query : %v", result.Error)
		return nil, result.Error
	}

	return &product_category, nil
}

func (r *ProductCategoryRepository) Store(product_category *model.ProductCategory) error {

	result := database.GDB.Create(product_category)
	if result.Error != nil {
		log.Fatalf("Failed to execute insert query: %v", result.Error)
		return result.Error
	}

	return nil
}

func (r *ProductCategoryRepository) Update(id string, updatedData *model.ProductCategory) (*model.ProductCategory, error) {

	var product_category model.ProductCategory

	result := database.GDB.First(&product_category, id)
	if result.Error != nil {
		log.Fatalf("Error creating query : %v", result.Error)
		return nil, result.Error
	}

	err := database.GDB.Model(&product_category).Updates(model.ProductCategory{
		Name: updatedData.Name,
	}).Error

	if err != nil {
		return nil, err
	}

	// ส่งข้อมูลเวอร์ชันอัปเดตล่าสุดกลับไป
	return &product_category, nil
}

func (r *ProductCategoryRepository) Delete(id int) error {

	// var product_category model.ProductCategory

	result := database.GDB.Delete(&model.ProductCategory{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
