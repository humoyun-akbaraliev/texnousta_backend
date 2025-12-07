# Makefile для TexnoUsta Backend

# Переменные
BINARY_NAME=texnousta-server
DOCKER_IMAGE=texnousta-backend
DOCKER_TAG=latest

# Основные команды
.PHONY: build run dev test clean docker-build docker-run help

# Сборка приложения
build:
	@echo "Сборка приложения..."
	go build -o $(BINARY_NAME) cmd/main.go

# Запуск в режиме разработки
dev:
	@echo "Запуск в режиме разработки..."
	go run cmd/main.go

# Запуск собранного приложения
run: build
	@echo "Запуск приложения..."
	./$(BINARY_NAME)

# Установка зависимостей
deps:
	@echo "Установка зависимостей..."
	go mod tidy
	go mod download

# Форматирование кода
fmt:
	@echo "Форматирование кода..."
	go fmt ./...

# Проверка кода
vet:
	@echo "Проверка кода..."
	go vet ./...

# Тестирование
test:
	@echo "Запуск тестов..."
	go test -v ./...

# Очистка
clean:
	@echo "Очистка..."
	rm -f $(BINARY_NAME)
	go clean

# Docker команды
docker-build:
	@echo "Сборка Docker образа..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run:
	@echo "Запуск Docker контейнера..."
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

# Миграции базы данных (если нужны отдельно)
migrate:
	@echo "Запуск миграций..."
	go run cmd/main.go -migrate

# Создание .env из примера
env:
	@if [ ! -f .env ]; then \
		echo "Создание .env файла..."; \
		cp .env.example .env; \
		echo "Не забудьте настроить переменные в .env файле"; \
	else \
		echo ".env файл уже существует"; \
	fi

# Проверка безопасности
security:
	@echo "Проверка безопасности..."
	@command -v gosec >/dev/null 2>&1 || { echo "gosec не установлен. Установите: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; exit 1; }
	gosec ./...

# Полная проверка перед коммитом
check: fmt vet test
	@echo "Все проверки пройдены успешно!"

# Справка
help:
	@echo "Доступные команды:"
	@echo "  build        - Сборка приложения"
	@echo "  dev          - Запуск в режиме разработки"
	@echo "  run          - Запуск собранного приложения"
	@echo "  deps         - Установка зависимостей"
	@echo "  fmt          - Форматирование кода"
	@echo "  vet          - Проверка кода"
	@echo "  test         - Запуск тестов"
	@echo "  clean        - Очистка"
	@echo "  docker-build - Сборка Docker образа"
	@echo "  docker-run   - Запуск Docker контейнера"
	@echo "  env          - Создание .env файла из примера"
	@echo "  check        - Полная проверка (fmt + vet + test)"
	@echo "  help         - Показать эту справку"