package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/markdeesoft/golang-api/model"
	"github.com/markdeesoft/golang-api/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB ตัวแปร Global (ตัวใหญ่) เพื่อให้แพ็กเกจอื่นเรียกใช้ได้
var DB *sql.DB
var GDB *gorm.DB

// InitDB ทำหน้าที่เชื่อมต่อฐานข้อมูลตอนเริ่มโปรแกรม
func InitDB() {
	var err error

	db_host := utils.GetEnvStr("DB_HOST", "localhost")
	db_port := utils.GetEnvInt("DB_PORT", 5432)                 // default PostgreSQL port
	db_user := utils.GetEnvStr("DB_USER", "myuser")             // as defined in docker-compose.yml
	db_password := utils.GetEnvStr("DB_PASSWORD", "mypassword") // as defined in docker-compose.yml
	db_name := utils.GetEnvStr("DB_NAME", "mydatabase")         // as defined in docker-compose.yml

	// สร้าง Connection String
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		db_host, db_port, db_user, db_password, db_name)

	// เชื่อมต่อฐานข้อมูล (ใช้ = ไม่ใช่ := เพราะเราต้องการกำหนดค่าลงตัวแปร Global ด้านบน)
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// ทดสอบการเชื่อมต่อจริง
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// New logger for detailed SQL logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	GDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger, // add Logger
	})
	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the schema
	GDB.AutoMigrate(&model.ProductCategory{}, &model.ProductHashtag{}, &model.Product{})

	SeedData()

	fmt.Println("Database connection successfully established!")
}

func CheckConnected() bool {

	// ทดสอบการเชื่อมต่อจริง
	if DB == nil {
		log.Fatalf("Error connecting to the database")
		return false
	}

	return true
}
