package repository

import (
	"github.com/markdeesoft/golang-api/database"
	"github.com/markdeesoft/golang-api/model"
)

// GetByUsername ค้นหาผู้ใช้จาก Username/email/phone เพื่อนำข้อมูลไปเช็ครหัสผ่านต่อ
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	// ต้องดึงคอลัมน์ password ออกมาด้วยเพื่อใช้ตรวจสอบ และต้องเช็คว่ายังไม่ถูกลบ (Soft Delete)
	query := `
		SELECT id, name, email, phone, password, role 
		FROM public.users 
		WHERE email = $1 AND deleted_at IS NULL;
	`
	var user model.User
	err := database.DB.QueryRow(query, username).
		Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Password, &user.Role)
	if err != nil {
		return nil, err // จะส่งกลับเป็น sql.ErrNoRows หากไม่พบชื่อผู้ใช้งานนี้
	}
	return &user, nil
}
