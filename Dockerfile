# Используем официальный образ Golang в качестве базового
FROM golang:latest

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем исходные файлы в контейнер
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем исполняемый файл с передачей параметра -token при сборке
RUN go build -o bot main.go

# Запускаем бота при старте контейнера с передачей параметра -token
CMD ["./bot", "-token", "${TELEGRAM_API_TOKEN}"]
