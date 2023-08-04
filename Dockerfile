# Этап 1: Установка зависимостей и создание go.mod
FROM golang:latest AS builder

WORKDIR /app

# Копируем исходные файлы в контейнер
COPY . .

# Создаем go.mod внутри контейнера
RUN go mod init my_telegram_bot && go mod tidy

# Проверяем наличие main.go
COPY main.go ./
RUN test -f main.go

# Этап 2: Сборка бота
RUN go build -o bot main.go

# Этап 3: Создание образа для запуска
FROM golang:latest

WORKDIR /app

# Копируем исполняемый файл бота из предыдущего этапа
COPY --from=builder /app/bot .

# Устанавливаем аргумент сборки TELEGRAM_API_TOKEN в качестве переменной окружения
ARG TELEGRAM_API_TOKEN
ENV TELEGRAM_API_TOKEN=$TELEGRAM_API_TOKEN

# Запускаем бота с передачей токена через параметр -token
CMD ["./bot", "-token", "$TELEGRAM_API_TOKEN"]
