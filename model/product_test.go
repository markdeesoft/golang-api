package model_test

import (
	"testing"

	"github.com/markdeesoft/golang-api/model"
	"github.com/stretchr/testify/assert"
)

func TestProductCategory_Validate(t *testing.T) {

	tests := []struct {
		name            string
		productCategory model.ProductCategory
		wantErr         bool
	}{
		{
			name:            "Valid ProductCategory Name",
			productCategory: model.ProductCategory{Name: "Health & Beauty", NameEn: "Health & Beauty"},
			wantErr:         false,
		},
		{
			name:            "Invalid ProductCategory Name is not null",
			productCategory: model.ProductCategory{Name: "", NameEn: "Health & Beauty"},
			wantErr:         true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.productCategory.Validate()
			// if (err != nil) != tc.wantErr {
			// 	t.Errorf("ProductCategory Validate() error = %v, wantErr %v", err, tc.wantErr)
			// }
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProductHashtag_Validate(t *testing.T) {

	tests := []struct {
		name           string
		ProductHashtag model.ProductHashtag
		wantErr        bool
	}{
		{
			name:           "Valid ProductHashtag",
			ProductHashtag: model.ProductHashtag{Name: "สินค้าขายดี", NameEn: "BestSeller", IsActive: true},
			wantErr:        false,
		},
		{
			name:           "Invalid ProductHashtag Name is not null",
			ProductHashtag: model.ProductHashtag{Name: "", NameEn: "BestSeller", IsActive: true},
			wantErr:        true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.ProductHashtag.Validate()
			// if (err != nil) != tc.wantErr {
			// 	t.Errorf("ProductHashtag Validate() error = %v, wantErr %v", err, tc.wantErr)
			// }
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProduct_Validate(t *testing.T) {

	tests := []struct {
		name    string
		product model.Product
		wantErr bool
	}{
		{
			name: "Valid all data",
			product: model.Product{
				Name:              "เอร่า",
				NameEn:            "Aera",
				Price:             1000,
				Description:       "Description",
				ProductCategoryID: 1,
				ProductCategory: model.ProductCategory{
					Name: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid Product Name is not null",
			product: model.Product{
				Name:              "",
				NameEn:            "Aera",
				Price:             1000,
				Description:       "Description",
				ProductCategoryID: 1,
				ProductCategory: model.ProductCategory{
					Name: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "valid Product Price is gt 0",
			product: model.Product{
				Name:              "เอร่า",
				NameEn:            "Aera",
				Price:             1000,
				Description:       "Description",
				ProductCategoryID: 1,
				ProductCategory: model.ProductCategory{
					Name: "test",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid Product Price is lt 0",
			product: model.Product{
				Name:              "เอร่า",
				NameEn:            "Aera",
				Price:             -1,
				Description:       "Description",
				ProductCategoryID: 1,
				ProductCategory: model.ProductCategory{
					Name: "test",
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid Product ProductCategoryID is wrong",
			product: model.Product{
				Name:        "เอร่า",
				NameEn:      "Aera",
				Price:       1000,
				Description: "Description",
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.product.Validate()
			// if (err != nil) != tc.wantErr {
			// 	t.Errorf("Product Validate() error = %v, wantErr %v", err, tc.wantErr)
			// }
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
