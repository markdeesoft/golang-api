package utils

import (
	"os"
	"strconv"
)

// GetEnvStr ดึงค่า Environment เป็น String (ขึ้นต้นด้วยตัวพิมพ์ใหญ่เพื่อให้เรียกใช้นอกแพ็กเกจได้)
func GetEnvStr(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt ดึงค่า Environment เป็น Int
func GetEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
