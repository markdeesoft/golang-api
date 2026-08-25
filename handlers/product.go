package handler

import (
	"errors"
	"log"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/markdeesoft/golang-api/model"
	"github.com/markdeesoft/golang-api/repository"
	"gorm.io/gorm"
)

type ProductHandler struct {
	repo *repository.ProductRepository
}

// NewProductHandler ฟังก์ชันเริ่มต้นใช้งาน Handler โดยดึง Repo เข้ามาผูกไว้
func NewProductHandler(repo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

func (h *ProductHandler) List(c fiber.Ctx) error {

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
		log.Printf("Failed to count product: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to count total records",
		})
	}

	products, err := h.repo.List(limit, offset)
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
	result := model.PaginationResult[model.Product]{
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		Data:       products,
	}

	// ส่ง Array ของรายชื่อผู้ใช้กลับไปให้ Client ด้วยสถานะ 200 OK
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *ProductHandler) View(c fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	product, err := h.repo.GetByID(id)
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.JSON(product)
}

func (h *ProductHandler) Store(c fiber.Ctx) error {

	product := new(model.Product)

	if err := c.Bind().Body(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	//validation
	if product.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name cannot be empty",
		})
	}

	// สั่งงานผ่านเลเยอร์ Repository
	if err := h.repo.Store(product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create  product  in database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

func (h *ProductHandler) Update(c fiber.Ctx) error {

	id := c.Params("id")
	request := new(model.Product)

	if err := c.Bind().Body(&request); err != nil {
		return err
	}

	//validation
	if request.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name cannot be empty",
		})
	}

	updatedUser, err := h.repo.Update(id, request)
	if err != nil {

		// ดักตรวจสอบกรณีไม่พบข้อมูล ID นี้ในระบบ (ซอฟต์ดีลีทไปแล้ว หรือไม่มีอยู่จริง)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}

		log.Printf("Handler Error - Update failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update product  in database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(updatedUser)
}

func (h *ProductHandler) Delete(c fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.repo.Delete(id); err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "product not found or already deleted",
			})
		}

		log.Printf("Failed to execute soft delete: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete product ",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "product  deleted successfully (Soft Delete)",
	})
}
