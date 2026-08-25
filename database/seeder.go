package database

import (
	"log"
	"math/rand/v2"
	"strconv"

	"github.com/markdeesoft/golang-api/model"
	"github.com/markdeesoft/golang-api/utils"
	"golang.org/x/crypto/bcrypt"
)

func SeedData() {

	// เช็คเบื้องต้นว่ามีข้อมูลอยู่แล้วหรือยัง เพื่อไม่ให้ข้อมูลซ้ำซ้อนตอนสั่งรันใหม่
	var count int64
	result := GDB.Model(&model.User{}).Count(&count)
	if result.Error != nil {

		log.Println("User Query Error", result.Error)

	} else if count == 0 {

		log.Println("Seeding User database...")

		// เตรียมข้อมูลที่ต้องการจำลอง
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(utils.GetEnvStr("USER_PASSWORD_DEFAULT", "12345678")), bcrypt.DefaultCost)

		datas := []model.User{
			{Name: "admin", Phone: "123456", Email: "admin@example.com", Password: string(hashedPassword)},
			{Name: "john_doe", Phone: "123456", Email: "john@example.com", Password: string(hashedPassword)},
			{Name: "jane_smith", Phone: "123456", Email: "jane@example.com", Password: string(hashedPassword)},
		}

		// 3. ใช้ GORM บันทึกข้อมูลลงตารางแบบ Bulk Insert (รวดเดียวทั้งหมด)
		if err := GDB.Create(&datas).Error; err != nil {
			log.Fatalf("Failed to seed user: %v", err)
		}

		log.Println("User Database seeded successfully!")
	}

	result = GDB.Model(&model.ProductCategory{}).Count(&count)
	if result.Error != nil {

		log.Println("ProductCategory Query Error", result.Error)

	} else if count == 0 {

		log.Println("Seeding ProductCategory database...")

		datas := []model.ProductCategory{
			{Name: "สุขภาพและความงาม", NameEn: "Health & Beauty"},
			{Name: "ของใช้ในบ้านและไลฟ์สไตล์", NameEn: "Home & Living"},
			{Name: "อาหารและเครื่องดื่ม", NameEn: "Food & Beverage"},
			{Name: "เสื้อผ้าและเครื่องแต่งกาย", NameEn: "Apparel & Clothing"},
			{Name: "ของเล่นและเกม", NameEn: "Toys & Games"},
			{Name: "กีฬาและอุปกรณ์กลางแจ้ง", NameEn: "Sports & Outdoors"},
			{Name: "สินค้าอิเล็กทรอนิกส์และแกดเจ็ต", NameEn: "Electronics & Gadgets"},
		}

		// 3. ใช้ GORM บันทึกข้อมูลลงตารางแบบ Bulk Insert (รวดเดียวทั้งหมด)
		if err := GDB.Create(&datas).Error; err != nil {
			log.Fatalf("Failed to seed product category: %v", err)
		}

		log.Println("ProductCategory Database seeded successfully!")
	}

	result = GDB.Model(&model.ProductHashtag{}).Count(&count)
	if result.Error != nil {

		log.Println("ProductCategory Query Error", result.Error)

	} else if count == 0 {

		log.Println("Seeding ProductHashtag database...")

		datas := []model.ProductHashtag{
			{Name: "สินค้าเด่นประจำวัน", NameEn: "ProductOfTheDay"},
			{Name: "สินค้าขายดี", NameEn: "BestSeller"},
			{Name: "สินค้าเข้าใหม่", NameEn: "NewArrival"},
			{Name: "ของมันต้องมี", NameEn: "MustHave"},
			{Name: "โปรโมชันลดราคา", NameEn: "Sale"},
			{Name: "สินค้าจำนวนจำกัด", NameEn: "LimitedEdition"},
			{Name: "ส่งฟรี", NameEn: "FreeShipping"},
			{Name: "ช้อปเลย", NameEn: "ShopNow"},
		}

		// 3. ใช้ GORM บันทึกข้อมูลลงตารางแบบ Bulk Insert (รวดเดียวทั้งหมด)
		if err := GDB.Create(&datas).Error; err != nil {
			log.Fatalf("Failed to seed product category: %v", err)
		}

		log.Println("ProductHashtag Database seeded successfully!")
	}

	result = GDB.Model(&model.Product{}).Count(&count)
	if result.Error != nil {

		log.Println("ProductCategory Query Error", result.Error)

	} else if count == 0 {

		log.Println("Seeding Product database...")

		datas := []model.Product{}

		count := 5
		for i := 0; i < count; i++ {

			var product_hashtags = []model.ProductHashtag{}

			for i := 0; i < randomIntBetween(1, 3); i++ {

				product_hashtag := model.ProductHashtag{ID: uint(randomIntBetween(1, 8))}
				product_hashtags = append(product_hashtags, product_hashtag)
			}

			var product = model.Product{
				Name:              "Product" + strconv.Itoa(i),
				Price:             randomFloatBetween(100, 1000),
				ProductCategoryID: uint(randomIntBetween(1, 7)),
				ProductHashtags:   product_hashtags,
			}

			datas = append(datas, product)
		}

		// 3. ใช้ GORM บันทึกข้อมูลลงตารางแบบ Bulk Insert (รวดเดียวทั้งหมด)
		if err := GDB.Create(&datas).Error; err != nil {
			log.Fatalf("Failed to seed product category: %v", err)
		}

		log.Println("Product Database seeded successfully!")
	}

}

func randomIntBetween(min, max int) int {
	return rand.IntN(max-min+1) + min
}

func randomFloatBetween(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
