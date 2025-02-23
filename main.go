package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	_ "tasks/docs"
	"tasks/handlers"
	"tasks/repository"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/lib/pq"
)

var port = ":8080"

func main() {
	connStr := "host=localhost port=5433 user=admin password=secret dbname=tasks_db sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	taskRepo := &repository.TaskRepository{DB: db}
	taskHandler := &handlers.TaskHandler{Repo: taskRepo}

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.GetTasks(w, r)
		case http.MethodPost:
			taskHandler.CreateTask(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			taskHandler.UpdateTask(w, r)
		case http.MethodDelete:
			taskHandler.DeleteTask(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	//http://localhost:8080/swagger/index.html

	fmt.Println("Сервер запущен на http://localhost", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}

}
