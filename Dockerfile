# Этап 1: Установка зависимостей и создание go.mod
FROM golang:latest AS builder

WORKDIR /app

# Копируем исходные файлы в контейнер
COPY . .

# Создаем go.mod внутри контейнера
RUN go mod init my_telegram_bot && go mod tidy

# Этап 2: Сборка бота
RUN go build -o bot main.go

# Этап 3: Создание образа для запуска
FROM golang:latest

WORKDIR /app

COPY --from=builder /app/bot .

CMD ["./bot", "-token", ${TELEGRAM_API_TOKEN}]
