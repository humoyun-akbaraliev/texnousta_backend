# Многоэтапная сборка для оптимизации размера образа
FROM golang:1.21-alpine AS builder

# Установка необходимых пакетов
RUN apk add --no-cache git

# Установка рабочей директории
WORKDIR /app

# Копирование go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main main.go

# Финальный образ
FROM alpine:latest

# Установка ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates tzdata

# Создание пользователя без привилегий
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /root/

# Копирование собранного приложения
COPY --from=builder /app/main .
COPY --from=builder /app/.env.example .env

# Создание директории для загрузок
RUN mkdir -p uploads && \
    chown -R appuser:appgroup uploads

# Изменение владельца файлов
RUN chown -R appuser:appgroup /root/

# Переключение на пользователя без привилегий
USER appuser

# Открытие порта
EXPOSE 8080

# Команда запуска
CMD ["./main"]