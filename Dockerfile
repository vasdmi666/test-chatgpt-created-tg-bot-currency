# Этап 1: Установка зависимостей и создание go.mod
FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Этап 2: Сборка бота
COPY . .

RUN go build -o bot main.go

# Этап 3: Создание образа для запуска
FROM golang:latest

WORKDIR /app

COPY --from=builder /app/bot .

CMD ["./bot", "-token", "${TELEGRAM_API_TOKEN}"]
