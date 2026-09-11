package repository

import (
	"github.com/markdeesoft/golang-api/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Count() (int64, error)
	Delete(id uint) error
	GetByID(id uint) (*model.User, error)
	GetAll(limit int, offset int) ([]model.User, error)
	ResetPassword(id uint, password string) error
	Store(user *model.User) error
	Update(id uint, user *model.User) error
	GetByUsername(username string) (*model.User, error)
}

// userRepository โครงสร้างรองรับการเรียกใช้งานคิวรี่
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository ฟังก์ชันสร้างอินสแตนซ์ของ Repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Count ดึงจำนวนผู้ใช้ทั้งหมดที่ยังไม่ถูกลบ
func (r *userRepository) Count() (int64, error) {
	var total int64
	err := r.db.Model(&model.User{}).Count(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

// GetAll ดึงรายการผู้ใช้ตามช่วงที่กำหนดด้วย limit และ offset
func (r *userRepository) GetAll(limit, offset int) ([]model.User, error) {
	var users []model.User
	err := r.db.Limit(limit).Offset(offset).Order("id asc").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetByID ดึงข้อมูลผู้ใช้รายบุคคลตาม ID
func (r *userRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err // จะ คืนค่า gorm.ErrRecordNotFound กรณีหาไม่เจอ
	}
	return &user, nil
}

// Store บันทึกข้อมูลผู้ใช้ใหม่ลงฐานข้อมูล
func (r *userRepository) Store(user *model.User) error {
	// GORM จะทำ Auto-increment ID และเติม CreatedAt/UpdatedAt ให้อัตโนมัติ
	return r.db.Create(user).Error
}

// Update แก้ไขข้อมูลผู้ใช้ตาม ID
func (r *userRepository) Update(id uint, user *model.User) error {
	// ใช้ Model(&model.User{}).Where("id = ?", id) เพื่อระบุ Record ที่ต้องการอัปเดต
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":  user.Name,
		"email": user.Email,
		"phone": user.Phone,
	}).Error
}

// Delete ลบข้อมูลแบบ Soft Delete
func (r *userRepository) Delete(id uint) error {
	// หาก Struct model.User มี field gorm.DeletedAt อยู่แล้ว
	// คำสั่ง Delete() ของ GORM จะทำการ UPDATE deleted_at = NOW() ให้อัตโนมัติ
	return r.db.Delete(&model.User{}, id).Error
}

// ResetPassword อัปเดตรหัสผ่านใหม่
func (r *userRepository) ResetPassword(id uint, password string) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("password", password).Error
}
