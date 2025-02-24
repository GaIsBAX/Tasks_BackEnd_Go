package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"tasks/repository"

	"tasks/utils"

	model "tasks/models"
)

type TaskHandler struct {
	Repo *repository.TaskRepository
}

// GetTasks handles HTTP GET requests for retrieving a list of tasks.
// @Summary Получить список задач
// @Description Возвращает массив задач
// @Tags tasks
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Task "Массив задач"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /tasks [get]

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.Repo.GetTasks()
	if err != nil {
		http.Error(w, "Ошибка получения задач", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

// func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var newTask model.Task
// 	err := json.NewDecoder(r.Body).Decode(&newTask)
// 	if err != nil {
// 		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
// 		return
// 	}

// 	// if strings.TrimSpace(newTask.Title) == "" || strings.TrimSpace(newTask.Description) == "" || strings.TrimSpace(newTask.Status) == "" {
// 	// 	http.Error(w, "Заполните все поля", http.StatusBadRequest)
// 	// 	return
// 	// }

// 	if !utils.IsValidStatus(newTask.Status) {
// 		http.Error(w, "Недопустимый статус. Используйте: pending, in progress, completed", http.StatusBadRequest)
// 		return
// 	}

// 	// Добавляем в БД
// 	err = h.Repo.CreateTask(&newTask)
// 	if err != nil {
// 		http.Error(w, "Ошибка создания задачи", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(newTask)
// }

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var newTask model.Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		log.Printf("Ошибка декодирования JSON: %v", err)
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(newTask.Title) == "" || strings.TrimSpace(newTask.Description) == "" || strings.TrimSpace(newTask.Status) == "" {
		http.Error(w, "Заполните все поля", http.StatusBadRequest)
		return
	}

	err = h.Repo.CreateTask(&newTask)
	if err != nil {
		log.Printf("Ошибка добавления задачи в БД: %v", err)
		http.Error(w, "Ошибка создания задачи", http.StatusInternalServerError)
		return
	}

	log.Printf("Задача создана: %+v", newTask)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	id, err := utils.ExtractID(r.URL.Path)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	var updatedTask model.Task
	err = json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(updatedTask.Title) == "" && strings.TrimSpace(updatedTask.Description) == "" && strings.TrimSpace(updatedTask.Status) == "" {
		http.Error(w, "Заполните хотя бы одно поле", http.StatusBadRequest)
		return
	} // потом сделай по человечески

	if !utils.IsValidStatus(updatedTask.Status) {
		http.Error(w, "Недопустимый статус. Используйте: pending, in progress, completed", http.StatusBadRequest)
		return
	}

	// Обновляем задачу в БД через репозиторий
	err = h.Repo.UpdateTask(id, &updatedTask)
	if err != nil {
		http.Error(w, "Ошибка обновления задачи", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	id, err := utils.ExtractID(r.URL.Path)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	err = h.Repo.DeleteTask(id)
	if err != nil {
		http.Error(w, "Ошибка удаления задачи", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
