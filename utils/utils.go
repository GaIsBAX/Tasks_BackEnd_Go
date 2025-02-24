package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	model "tasks_backend/models"
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

func SendErrorResponse(w http.ResponseWriter, statusCode int, errMsg string, details string) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := model.ErrorResponse{
		Error:   errMsg,
		Details: details,
	}

	json.NewEncoder(w).Encode(response)

}
