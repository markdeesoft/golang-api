package repository_test

import (
	"testing"

	"github.com/markdeesoft/golang-api/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB ทำหน้าที่เตรียมการเชื่อมต่อ DB สำหรับทดสอบ
func setupTestDB(t *testing.T) *gorm.DB {

	var err error

	// // ดึง DSN จาก Environment Variable หรือใช้ Default DSN สำหรับการ Test (postgres)
	// db_host := utils.GetEnvStr("DB_HOST", "localhost")
	// db_port := utils.GetEnvInt("DB_PORT", 5432)                 // default PostgreSQL port
	// db_user := utils.GetEnvStr("DB_USER", "myuser")             // as defined in docker-compose.yml
	// db_password := utils.GetEnvStr("DB_PASSWORD", "mypassword") // as defined in docker-compose.yml
	// db_name := "mydatabase_test"         // as defined in docker-compose.yml

	// // สร้าง Connection String
	// dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	// 	db_host, db_port, db_user, db_password, db_name)

	// database.GDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Skipf("ข้ามการทดสอบ: ไม่สามารถเชื่อมต่อกับ Postgres Test DB ได้: %v", err)
	}

	// ล็อกให้ใช้ Connection เดียวตลอดเวลา (ป้องกัน DB หลุด/หายระหว่าง Test)
	// sqlDB, err := db.DB()
	// assert.NoError(t, err)
	// sqlDB.SetMaxOpenConns(1)

	// ทำ AutoMigrate เพื่อเตรียม Table สำหรับ Test
	err = db.AutoMigrate(&model.User{}, &model.ProductCategory{}, &model.ProductHashtag{}, &model.Product{})
	assert.NoError(t, err)

	// คืนค่า Cleanup Function เพื่อล้างข้อมูลและลบ Table หลังเทสเสร็จ
	t.Cleanup(func() {
		db.Migrator().DropTable(&model.User{}, &model.ProductCategory{}, &model.ProductHashtag{}, &model.Product{})
	})

	return db
}
