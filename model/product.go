package model

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var validate = validator.New()

type Product struct {
	ID                uint             `gorm:"primaryKey" json:"id"`
	Name              string           `gorm:"unique;not null" json:"name" validate:"required"`
	NameEn            string           `json:"name_en"`
	Price             float64          `gorm:"type:decimal(10,2);not null;default:0.00" json:"price" validate:"gt=0"`
	Description       string           `json:"description"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	DeletedAt         gorm.DeletedAt   `gorm:"index" json:"-"`
	UserID            uint             `json:"user_id"`
	ProductCategoryID uint             `json:"product_category_id" validate:"required"`
	ProductCategory   ProductCategory  `gorm:"constraint:OnDelete:RESTRICT;" json:"-" validate:"-"`
	ProductHashtags   []ProductHashtag `gorm:"many2many:product_hashtag_mappings;" json:"product_hashtags,omitempty"`
}

func (p *Product) Validate() error {
	return validate.Struct(p)
}

type ProductCategory struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"unique;not null" json:"name" validate:"required"`
	NameEn      string `json:"name_en"`
	Description string `json:"description"`
}

func (p *ProductCategory) Validate() error {
	return validate.Struct(p)
}

type ProductHashtag struct {
	ID       uint   `json:"id"`
	Name     string `gorm:"unique;not null" json:"name" validate:"required"`
	NameEn   string `json:"name_en"`
	IsActive bool   `gorm:"default:false" json:"is_active"`
}

func (p *ProductHashtag) Validate() error {
	return validate.Struct(p)
}
