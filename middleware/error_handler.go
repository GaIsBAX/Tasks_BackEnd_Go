package middleware

import (
	"log"
	"net/http"
	"tasks_backend/utils"
)

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("⚠ Ошибка сервера: %v", err)
				utils.SendErrorResponse(w, http.StatusInternalServerError, "Внутренняя ошибка сервера", "")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
