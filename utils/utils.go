package utils

import (
	"fmt"
	"strconv"
	"strings"
)

var validStatuses = map[string]bool{
	"pending":     true,
	"in progress": true,
	"completed":   true,
}

// IsValidStatus проверяет, является ли статус допустимым
func IsValidStatus(status string) bool {
	return validStatuses[status]
}

// extractID извлекает ID из URL-пути, например: "/tasks/5" → 5
func ExtractID(urlPath string) (int, error) {
	pathParts := strings.Split(urlPath, "/")
	if len(pathParts) < 3 || pathParts[2] == "" {
		return 0, fmt.Errorf("ID не указан")
	}
	return strconv.Atoi(pathParts[2])
}
