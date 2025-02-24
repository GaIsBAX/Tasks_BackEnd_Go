package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	_ "tasks_backend/docs"
	"tasks_backend/handlers"
	"tasks_backend/repository"
	"tasks_backend/utils"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/lib/pq"
)

var port = ":8080"

func main() {

	host := os.Getenv("DB_HOST")
	portdb := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, portdb, user, password, dbname)

	if host == "" || portdb == "" || user == "" || password == "" || dbname == "" {
		log.Fatalf("❌ Ошибка: не все переменные окружения установлены! DB_HOST=%s, DB_PORT=%s, DB_USER=%s, DB_NAME=%s", host, port, user, dbname)
	}

	// connStr := "host=localhost port=5433 user=admin password=secret dbname=tasks_db sslmode=disable"

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
			utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Метод не поддерживается", "Ошибка с Get или Post запросом")
		}
	})

	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			taskHandler.UpdateTask(w, r)
		case http.MethodDelete:
			taskHandler.DeleteTask(w, r)
		default:
			utils.SendErrorResponse(w, http.StatusMethodNotAllowed, "Метод не поддерживается", "Ошибка с Put или Delete запросом")
		}
	})

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	//http://localhost:8080/swagger/index.html

	fmt.Println("Сервер запущен на http://localhost", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}

}
