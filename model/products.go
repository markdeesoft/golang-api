package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID                uint             `gorm:"primaryKey" json:"id"`
	Name              string           `gorm:"unique;not null" json:"name"`
	NameEn            string           `json:"name_en"`
	Price             float64          `json:"price"`
	Description       string           `json:"description"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	DeletedAt         gorm.DeletedAt   `gorm:"index" json:"-"`
	UserID            uint             `json:"user_id"`
	ProductCategoryID uint             `json:"product_category_id"`
	ProductCategory   ProductCategory  `gorm:"constraint:OnDelete:RESTRICT;" json:"-"`
	ProductHashtags   []ProductHashtag `gorm:"many2many:product_hashtag_mappings;" json:"product_hashtags,omitempty"`
}

type ProductCategory struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"unique;not null" json:"name"`
	NameEn      string `json:"name_en"`
	Description string `json:"description"`
}

type ProductHashtag struct {
	ID      uint   `json:"id"`
	Name    string `gorm:"unique;not null" json:"name"`
	NameEn  string `json:"name_en"`
	Actived bool   `json:"actived"`
}
