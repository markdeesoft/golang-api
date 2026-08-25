package repository

import (
	"context"
	"log"

	"github.com/markdeesoft/golang-api/database"
	"github.com/markdeesoft/golang-api/model"
)

// UserRepository โครงสร้างรองรับการเรียกใช้งานคิวรี่
type UserRepository struct{}

// NewUserRepository ฟังก์ชันสร้างอินสแตนซ์ของ Repository
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Count ดึงจำนวนผู้ใช้ทั้งหมดที่ยังไม่ถูกลบ (สำหรับคำนวณหน้า)
func (r *UserRepository) Count() (int64, error) {
	var total int64
	query := "SELECT COUNT(*) FROM public.users WHERE deleted_at IS NULL"

	err := database.DB.QueryRow(query).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// List ดึงรายการผู้ใช้ตามช่วงที่กำหนดด้วย limit และ offset
func (r *UserRepository) List(limit, offset int) ([]model.User, error) {

	query := `
		SELECT id, name, email, phone, created_at 
		FROM public.users 
		WHERE deleted_at IS NULL 
		ORDER BY id ASC 
		LIMIT $1 OFFSET $2;
	`

	rows, err := database.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []model.User{}
	// วนลูปอ่านข้อมูลทีละแถวด้วย rows.Next()
	for rows.Next() {
		var user model.User
		// สแกนข้อมูลในแถวนั้นๆ เข้าตัวแปร user
		err := rows.Scan(&user.ID, &user.Name, &user.Phone, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		// เพิ่ม user เข้าไปในกลุ่ม Array
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Rows error: %v", err)
		return nil, err
	}

	return users, nil
}

// GetByID ดึงข้อมูลผู้ใช้รายบุคคลตาม ID
func (r *UserRepository) GetByID(id int) (*model.User, error) {
	query := `
		SELECT id, name, role FROM users 
		WHERE id = $1 AND deleted_at IS NULL;
	`
	var user model.User
	err := database.DB.QueryRow(query, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)

	if err != nil {
		return nil, err // ส่ง err กลับไปให้ Handler ไปตรวจสอบว่าเป็น sql.ErrNoRows หรือไม่
	}
	return &user, nil
}

// Create บันทึกข้อมูลผู้ใช้ใหม่ลงฐานข้อมูล
func (r *UserRepository) Store(user *model.User) error {

	tx, err := database.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	// ใช้ defer ร่วมกับ tx.Rollback() เพื่อป้องกันกรณีโปรแกรมแครชกลางทาง
	defer tx.Rollback()

	query := `
		INSERT INTO public.users (name, email, phone, password, updated_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, created_at;
	`

	// สั่งรันคิวรี่ผ่านตัวแปร Global DB และ Scan ค่ากลับมาอัปเดตลงใน struct
	err = database.DB.QueryRow(query, user.Name, user.Phone, user.Password).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err // หากคำสั่งนี้พัง ข้อมูลผู้ใช้ในคำสั่งก่อนหน้าก็จะถูก Rollback ไปด้วย
	}

	// บันทึกข้อมูลทั้งหมดลงฐานข้อมูลจริงเมื่อทุกคำสั่งทำงานผ่านฉลุย
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// แก้ไข
func (r *UserRepository) Update(user *model.User, id int) error {

	tx, err := database.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	// ใช้ defer ร่วมกับ tx.Rollback() เพื่อป้องกันกรณีโปรแกรมแครชกลางทาง
	defer tx.Rollback()

	query := `
		UPDATE public.users SET 
		name = $1, email = $2, phone = $3, updated_at = NOW()
		WHERE id = $4 
		RETURNING updated_at;
	`

	err = database.DB.QueryRow(query, user.Name, user.Email, user.Phone, id).Scan(&user.UpdatedAt)
	if err != nil {
		return err // หากคำสั่งนี้พัง ข้อมูลผู้ใช้ในคำสั่งก่อนหน้าก็จะถูก Rollback ไปด้วย
	}

	// บันทึกข้อมูลทั้งหมดลงฐานข้อมูลจริงเมื่อทุกคำสั่งทำงานผ่านฉลุย
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// ลบ softdelete
func (r *UserRepository) Delete(id int) error {

	query := `
		UPDATE public.users 
		SET deleted_at = NOW() 
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at;
	`
	_, err := database.DB.Exec(query, id)
	return err
}

// reset password to default
func (r *UserRepository) ResetPassword(id int, password string) (*model.User, error) {

	query := `
		UPDATE users SET 
		password = $1, updated_at = NOW()
		WHERE id = $2 
		RETURNING updated_at;
	`

	var user model.User
	row := database.DB.QueryRow(query, password, id)
	err := row.Scan(&user.UpdatedAt)

	if err != nil {
		return nil, err // ส่ง err กลับไปให้ Handler ไปตรวจสอบว่าเป็น sql.ErrNoRows หรือไม่
	}

	return &user, nil
}
