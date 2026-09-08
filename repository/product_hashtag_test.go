package repository

// "gorm.io/driver/postgres"
// "gorm.io/gorm"

// func setupTestDB() *gorm.DB {

// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
// 	}
// 	defer db.Close()

// 	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
// 	if err != nil {
// 		t.Fatalf("an error '%s' was not expected when opening a gorm database", err)
// 	}
// 	return db
// }

// func TestGetProductByID(t *testing.T) {

// 	// 1. สร้าง DB Mock
// 	db, mock, err := sqlmock.New()
// 	assert.NoError(t, err)
// 	defer db.Close()

// 	// 2. ตั้ง Expectation สำหรับ Query
// 	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Monitor")
// 	mock.ExpectQuery("SELECT (.+) FROM products WHERE id = ?").
// 		WithArgs(1).
// 		WillReturnRows(rows)

// 	// 3. รัน Code Repository และ Assert
// 	// repo := repository.NewProductRepository(db)
// 	// product, err := repo.GetByID(1)

// 	assert.NoError(t, err)
// 	assert.NoError(t, mock.ExpectationsWereMet()) // ตรวจสอบว่า Query ถูกเรียกตาม Expectation หรือไม่
// }
