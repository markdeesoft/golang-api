package model

// PaginationResult โครงสร้างข้อมูลสำหรับส่งผลลัพธ์กลับแบบแบ่งหน้า
type PaginationResult[T any] struct {
	Total      int64 `json:"total"`       // จำนวนข้อมูลทั้งหมดที่ยังไม่ลบ
	Page       int   `json:"page"`        // หน้าปัจจุบัน
	Limit      int   `json:"limit"`       // จำนวนข้อมูลต่อหน้า
	TotalPages int   `json:"total_pages"` // จำนวนหน้าทั้งหมดที่มี
	Data       []T   `json:"data"`        // รายชื่อข้อมูลในหน้านั้นๆ
}
