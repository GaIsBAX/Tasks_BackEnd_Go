# 1. Используем официальный образ Go для сборки
FROM golang:1.23 AS builder

# 2. Устанавливаем рабочую директорию
WORKDIR /app

# 3. Копируем модули и устанавливаем зависимости
COPY go.mod go.sum ./

RUN go mod download

# 4. Копируем исходный код
COPY . .

# 5. Собираем бинарный файл ДЛЯ LINUX
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/tasks_backend

# 6. Создаём финальный контейнер на базе Debian (лучше, чем Alpine)
FROM debian:latest

# 7. Устанавливаем рабочую директорию
WORKDIR /root/

# 8. Копируем скомпилированный бинарник
COPY --from=builder /app/tasks_backend .

# 9. Даём права на выполнение
RUN chmod +x /root/tasks_backend

# 10. Запускаем приложение
CMD ["./tasks_backend"]
