package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"tasks_backend/repository"

	"tasks_backend/utils"

	model "tasks_backend/models"
)

type TaskHandler struct {
	Repo *repository.TaskRepository
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.Repo.GetTasks()
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Ошибка получения задач", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask model.Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		log.Printf("Ошибка декодирования JSON: %v", err)
		utils.SendErrorResponse(w, http.StatusBadRequest, "Некорректный JSON", err.Error())
		return
	}

	if strings.TrimSpace(newTask.Title) == "" || strings.TrimSpace(newTask.Description) == "" || strings.TrimSpace(newTask.Status) == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Ошибка заполения", "Поля title, description, status не должны быть пустыми")
		return
	}

	if !utils.IsValidStatus(newTask.Status) {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Ошибка статуса", "Используйте статусы: pending, in progress, completed")
		return
	}

	err = h.Repo.CreateTask(&newTask)
	if err != nil {
		log.Printf("Ошибка добавления задачи в БД: %v", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Ошибка создания задачи", err.Error())
		return
	}

	log.Printf("Задача создана: %+v", newTask)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractID(r.URL.Path)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Некорректный ID", err.Error())
		return
	}

	var updatedTask model.Task
	err = json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		log.Printf("Ошибка декодирования JSON: %v", err)
		utils.SendErrorResponse(w, http.StatusBadRequest, "Некорректный JSON", err.Error())
		return
	}

	if strings.TrimSpace(updatedTask.Title) == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Ошибка валидации", "Поле title не должно быть пустым")
		return
	}

	if strings.TrimSpace(updatedTask.Description) == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Ошибка валидации", "Поле description не должно быть пустым")
		return
	}

	if strings.TrimSpace(updatedTask.Status) == "" {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Ошибка валидации", "Поле status не должно быть пустым")
		return
	}

	if !utils.IsValidStatus(updatedTask.Status) {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Ошибка статуса", "Используйте статусы: pending, in progress, completed")
		return
	}

	// Обновляем задачу в БД через репозиторий
	err = h.Repo.UpdateTask(id, &updatedTask)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Ошибка обновления задачи", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractID(r.URL.Path)
	if err != nil {

		utils.SendErrorResponse(w, http.StatusBadRequest, "Некорректный ID", err.Error())
		return
	}

	err = h.Repo.DeleteTask(id)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Ошибка удаления задачи", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
