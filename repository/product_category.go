package repository

import (
	"log"

	"github.com/markdeesoft/golang-api/database"
	"github.com/markdeesoft/golang-api/model"
)

type ProductCategoryRepository interface {
	Count() (int64, error)
	Delete(id uint) error
	GetByID(id uint) (*model.ProductCategory, error)
	GetAll(limit int, offset int) ([]model.ProductCategory, error)
	Store(product_category *model.ProductCategory) error
	Update(id uint, updatedData *model.ProductCategory) (*model.ProductCategory, error)
}

type productCategoryRepository struct{}

// NewproductCategoryRepository ฟังก์ชันสร้างอินสแตนซ์ของ Repository
func NewProductCategoryRepository() ProductCategoryRepository {
	return &productCategoryRepository{}
}

func (r *productCategoryRepository) Count() (int64, error) {

	var total int64
	result := database.GDB.Model(&model.ProductCategory{}).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (r *productCategoryRepository) GetAll(limit, offset int) ([]model.ProductCategory, error) {

	var product_categories []model.ProductCategory

	// สั่งดึงข้อมูลตามหน้า LIMIT / OFFSET
	err := database.GDB.Limit(limit).Offset(offset).Order("id asc").Find(&product_categories).Error

	return product_categories, err

}

func (r *productCategoryRepository) GetByID(id uint) (*model.ProductCategory, error) {

	var product_category model.ProductCategory

	result := database.GDB.First(&product_category, id)
	if result.Error != nil {
		// log.Fatalf("Error creating query : %v", result.Error)
		return nil, result.Error
	}

	return &product_category, nil
}

func (r *productCategoryRepository) Store(product_category *model.ProductCategory) error {

	result := database.GDB.Create(product_category)
	if result.Error != nil {
		log.Fatalf("Failed to execute insert query: %v", result.Error)
		return result.Error
	}

	return nil
}

func (r *productCategoryRepository) Update(id uint, updatedData *model.ProductCategory) (*model.ProductCategory, error) {

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

func (r *productCategoryRepository) Delete(id uint) error {

	// var product_category model.ProductCategory

	result := database.GDB.Delete(&model.ProductCategory{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
