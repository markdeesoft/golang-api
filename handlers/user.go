package handler

import (
	"database/sql"
	"log"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/markdeesoft/golang-api/model"
	"github.com/markdeesoft/golang-api/repository"
	"github.com/markdeesoft/golang-api/utils"
	"golang.org/x/crypto/bcrypt"
)

// UserHandler โครงสร้างสำหรับผูกพึ่งพา (Dependency Injection) กับ Repository
type UserHandler struct {
	repo *repository.UserRepository
}

// NewUserHandler ฟังก์ชันเริ่มต้นใช้งาน Handler โดยดึง Repo เข้ามาผูกไว้
func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// Handler functions
func (h *UserHandler) GetUsers(c fiber.Ctx) error {

	//รับค่า page และ limit จาก Query Parameters (หากไม่ได้ส่งมา ให้ใส่ค่าเริ่มต้นไว้)
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "100")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	// 2. คำนวณหาค่า OFFSET (จุดเริ่มต้นในการดึงข้อมูลของหน้านั้นๆ)
	offset := (page - 1) * limit

	total, err := h.repo.Count()
	if err != nil {
		log.Printf("Failed to count users: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to count total records",
		})
	}

	users, err := h.repo.List(limit, offset)
	if err != nil {
		log.Printf("Handler Error - List failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve records",
		})
	}

	// คำนวณจำนวนหน้าทั้งหมด (Total Pages)
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	// ประกอบข้อมูลทั้งหมดส่งกลับให้ Client
	result := model.PaginationResult[model.User]{
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		Data:       users,
	}

	// ส่ง Array ของรายชื่อผู้ใช้กลับไปให้ Client ด้วยสถานะ 200 OK
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *UserHandler) GetUser(c fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	user, err := h.repo.GetByID(id)
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.JSON(user)
}

func (h *UserHandler) StoreUser(c fiber.Ctx) error {

	user := new(model.User)

	if err := c.Bind().Body(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	//validation
	if user.Name == "" || user.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name and Email cannot be empty",
		})
	}

	hashedPassword, err := getPassword("")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}
	user.Password = hashedPassword

	// สั่งงานผ่านเลเยอร์ Repository
	if err := h.repo.Store(user); err != nil {
		log.Printf("Failed to execute insert query: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user in database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) UpdateUser(c fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	user := new(model.User)
	if err := c.Bind().Body(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	//validation
	if user.Name == "" || user.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name and Email cannot be empty",
		})
	}

	if err := h.repo.Update(user, id); err != nil {
		// หากไม่พบ ID ดังกล่าวในฐานข้อมูล
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		log.Printf("Failed to execute update query: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user in database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *UserHandler) DeleteUser(c fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if h.repo.Delete(id) != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found or already deleted",
			})
		}
		log.Printf("Failed to execute soft delete: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully (Soft Delete)",
	})
}

func UploadPhotoUser(c fiber.Ctx) error {

	file, err := c.FormFile("photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// Save the file to the server
	err = c.SaveFile(file, "./uploads/"+file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendString("File uploaded successfully: " + file.Filename)

}

func (h *UserHandler) ResetPasswordUser(c fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	hashedPassword, err := getPassword("")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	user, err := h.repo.ResetPassword(id, hashedPassword)
	if err != nil {
		// หากไม่พบ ID ดังกล่าวในฐานข้อมูล
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		log.Printf("Failed to execute update query: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user in database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func getPassword(password string) (string, error) {

	if password == "" {
		password = utils.GetEnvStr("USER_PASSWORD_DEFAULT", "12345678")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
