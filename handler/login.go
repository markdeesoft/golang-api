package handler

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/markdeesoft/golang-api/model"
	"github.com/markdeesoft/golang-api/utils"
	"golang.org/x/crypto/bcrypt"
)

func (h *UserHandler) Login(c fiber.Ctx) error {

	request := new(model.LoginRequest)
	// แปลง JSON Body เข้าตัวแปร req พร้อม check error
	if err := c.Bind().Body(&request); err != nil {
		return err
	}

	// ตรวจสอบค่าว่างเบื้องต้น
	if request.Email == "" || request.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and Password are required",
		})
	}

	//เรียกใช้งาน Repository เพื่อค้นหาผู้ใช้จาก "username" ในระบบ
	user, err := h.repo.GetByUsername(request.Email)
	if err != nil {

		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}
		log.Printf("Login Error - Database query failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	// 🔒 ตรวจสอบความถูกต้องของรหัสผ่านที่ส่งมา เทียบกับตารางในฐานข้อมูลด้วย bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		// ถ้ารหัสผ่านไม่ตรงกัน (err != nil) ให้ปฏิเสธการเข้าสู่ระบบทันที
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}
	fmt.Println("user", user)
	// Create the Claims
	claims := jwt.MapClaims{
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	t, err := token.SignedString([]byte(utils.GetEnvStr("JWT_SECRET", "")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"user": user.Name, "role": user.Role, "token": t})
}
