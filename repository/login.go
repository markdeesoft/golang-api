package repository

import (
	"github.com/markdeesoft/golang-api/model"
)

// GetByUsername ค้นหาผู้ใช้จาก Username/email/phone เพื่อนำข้อมูลไปเช็ครหัสผ่านต่อ
// GetByUsername ค้นหาผู้ใช้จาก Username (หรือ Email/Phone ตามโครงสร้างของคุณ)
func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ? OR phone = ?", username, username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
