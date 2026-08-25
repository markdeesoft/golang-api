package repository

import (
	"errors"
	"log"

	"github.com/markdeesoft/golang-api/database"
	"github.com/markdeesoft/golang-api/model"
)

type ProductRepository struct{}

// NewProductRepository ฟังก์ชันสร้างอินสแตนซ์ของ Repository
func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

func (r *ProductRepository) Count() (int64, error) {

	var total int64
	result := database.GDB.Model(&model.Product{}).Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return total, nil
}

func (r *ProductRepository) List(limit, offset int) ([]model.Product, error) {

	var products []model.Product

	// สั่งดึงข้อมูลตามหน้า LIMIT / OFFSET
	err := database.GDB.Limit(limit).Offset(offset).
		Preload("ProductCategory").
		Preload("ProductHashtags").
		Order("id asc").
		Find(&products).Error

	return products, err

}

func (r *ProductRepository) GetByID(id int) (*model.Product, error) {

	var product model.Product

	result := database.GDB.Preload("ProductCategory").
		Preload("ProductHashtags").
		First(&product, id)
	if result.Error != nil {
		// log.Fatalf("Error creating query : %v", result.Error)
		return nil, result.Error
	}

	return &product, nil
}

func (r *ProductRepository) Store(product *model.Product) error {

	// เช็คว่ามี CategoryID นี้อยู่จริงไหม ป้องกันคีย์กำพร้า
	var count int64
	database.GDB.Model(&model.ProductCategory{}).Where("id = ?", product.ProductCategoryID).Count(&count)
	if count == 0 {
		return errors.New("product_category_id not found")
	}

	result := database.GDB.Create(product)
	if result.Error != nil {
		log.Fatalf("Failed to execute insert query: %v", result.Error)
		return result.Error
	}

	return nil
}

func (r *ProductRepository) Update(id string, updatedData *model.Product) (*model.Product, error) {

	var product model.Product

	result := database.GDB.First(&product, id)
	if result.Error != nil {
		log.Fatalf("Error creating query : %v", result.Error)
		return nil, result.Error
	}

	// เช็คความถูกต้องของหมวดหมู่ใหม่หากมีการแก้ไข
	var count int64
	database.GDB.Model(&model.ProductCategory{}).Where("id = ?", updatedData.ProductCategoryID).Count(&count)
	if count == 0 {
		return nil, errors.New("product_category_id not found")
	}

	// 💡 เทคนิคสําหรับ Many-to-Many: สั่งล้างข้อมูลแฮชแท็กเก่าในตารางกลางก่อน แล้วจึงอัปเดตชุดใหม่เข้าไปแทน
	database.GDB.Model(&product).Association("ProductHashtags").Replace(updatedData.ProductHashtags)

	result = database.GDB.Model(&product).Updates(model.Product{
		Name:              updatedData.Name,
		Price:             updatedData.Price,
		ProductCategoryID: updatedData.ProductCategoryID,
	})
	if result.Error != nil {
		return nil, result.Error
	}

	// ส่งข้อมูลเวอร์ชันอัปเดตล่าสุดกลับไป
	return &product, nil
}

func (r *ProductRepository) Delete(id int) error {

	// var product model.Product

	result := database.GDB.Delete(&model.Product{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
