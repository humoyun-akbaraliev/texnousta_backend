# 🚀 Деплой на Render - Пошаговая инструкция

## ✅ Что уже настроено:

- ✅ PostgreSQL драйвер подключен
- ✅ База данных Neon настроена
- ✅ Код поддерживает SQLite (локально) + PostgreSQL (продакшен)
- ✅ CORS настроен для продакшена
- ✅ Render конфигурация готова

## 📋 Шаг 1: Подготовка Git репозитория

```bash
# Убедитесь что все файлы добавлены в git
git add .
git commit -m "Готов к деплою на Render с PostgreSQL"

# Если репозиторий не на GitHub, создайте:
# 1. Создайте новый репозиторий на GitHub
# 2. Добавьте remote:
git remote add origin https://github.com/your-username/texnousta-backend.git
git branch -M main
git push -u origin main
```

## 🌐 Шаг 2: Создание Web Service на Render

### 1. Перейдите на [render.com](https://render.com) и войдите

### 2. Создайте новый Web Service:
- Нажмите **"New +"** → **"Web Service"**
- Подключите ваш GitHub репозиторий
- Выберите репозиторий **TexnoUsta_Backend**

### 3. Настройте параметры деплоя:

#### Основные настройки:
```
Name: texnousta-api
Region: Frankfurt (EU Central) - ближе к вашей базе Neon
Branch: main
Root Directory: (пусто)
Runtime: Go
```

#### Build & Deploy:
```
Build Command: go build -o main cmd/main.go
Start Command: ./main
```

#### Advanced Settings:
```
Auto-Deploy: Yes (включить)
Health Check Path: /api/v1/products
```

## 🔐 Шаг 3: Переменные окружения

В разделе **Environment Variables** добавьте:

```bash
# Обязательные переменные
GIN_MODE=release
PORT=8080
DATABASE_URL=postgresql://neondb_owner:npg_UfiA4BeTNC3P@ep-flat-bonus-agmh7719-pooler.c-2.eu-central-1.aws.neon.tech/neondb?sslmode=require
JWT_SECRET=your-super-secret-jwt-key-here-minimum-32-characters-long

# Опциональные (для CORS)
CORS_ORIGINS=https://your-frontend-domain.com
```

### ⚠️ ВАЖНО: 
**Смените JWT_SECRET на свой секретный ключ!** 
Используйте генератор паролей для создания случайной строки длиной минимум 32 символа.

## ⚡ Шаг 4: Деплой

1. Нажмите **"Create Web Service"**
2. Render начнет автоматический деплой
3. Процесс займет 2-3 минуты

## 📊 Шаг 5: Проверка деплоя

После успешного деплоя ваш API будет доступен по ссылке:
```
https://texnousta-api.onrender.com
```

### Проверьте эндпоинты:

```bash
# Замените на ваш URL
API_URL="https://texnousta-api.onrender.com"

# Проверка API
curl "$API_URL/api/v1/products"

# Проверка Swagger
curl "$API_URL/swagger/index.html"

# Тест контактной формы
curl -X POST "$API_URL/api/v1/phone-contact" \
  -H "Content-Type: application/json" \
  -d '{"phone": "+998901234567"}'
```

## 🔧 Шаг 6: Мониторинг и логи

### Просмотр логов:
- В Render dashboard перейдите в **Logs**
- Смотрите в реальном времени что происходит

### Проверка базы данных:
- Подключитесь к Neon через psql
- Проверьте что таблицы созданы автоматически

```bash
# Подключение к базе
psql 'postgresql://neondb_owner:npg_UfiA4BeTNC3P@ep-flat-bonus-agmh7719-pooler.c-2.eu-central-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require'

# Проверка таблиц
\dt

# Проверка данных
SELECT * FROM contact_forms;
```

## 🎯 Шаг 7: Интеграция с фронтендом

Обновите базовый URL в React приложении:

```javascript
// src/services/api.js
const API_BASE_URL = 'https://texnousta-api.onrender.com/api/v1';
```

Обновите CORS_ORIGINS после деплоя фронтенда:
```bash
CORS_ORIGINS=https://your-frontend-domain.vercel.app
```

## 🔄 Автоматические обновления

После настройки каждый push в main ветку будет автоматически обновлять API:

```bash
# Внесите изменения
git add .
git commit -m "Обновление API"
git push origin main

# Render автоматически задеплоит изменения
```

## 🛠️ Полезные команды для мониторинга

### Проверка статуса API:
```bash
curl -s -o /dev/null -w "%{http_code}" https://texnousta-api.onrender.com/api/v1/products
```

### Проверка времени ответа:
```bash
curl -w "@-" -o /dev/null -s "https://texnousta-api.onrender.com/api/v1/products" <<'EOF'
     time_namelookup:  %{time_namelookup}\n
        time_connect:  %{time_connect}\n
     time_appconnect:  %{time_appconnect}\n
    time_pretransfer:  %{time_pretransfer}\n
       time_redirect:  %{time_redirect}\n
  time_starttransfer:  %{time_starttransfer}\n
                     ----------\n
          time_total:  %{time_total}\n
EOF
```

## ❗ Возможные проблемы и решения

### 1. Build Error: "Go not found"
- Убедитесь что выбрали Runtime: **Go**
- Проверьте что go.mod в корне репозитория

### 2. Database Connection Error
- Проверьте DATABASE_URL переменную
- Убедитесь что Neon база доступна
- Проверьте что PostgreSQL драйвер подключен

### 3. CORS Errors
- Обновите CORS_ORIGINS с правильным доменом фронтенда
- Для тестирования временно используйте "*" (небезопасно для продакшена)

### 4. Health Check Failed
- Проверьте что /api/v1/products отвечает 200
- Смотрите логи для диагностики

## 🎉 Готово!

После успешного деплоя у вас будет:

- ✅ **API**: https://texnousta-api.onrender.com
- ✅ **Swagger**: https://texnousta-api.onrender.com/swagger/index.html  
- ✅ **PostgreSQL база** на Neon
- ✅ **Автоматические обновления** при push в Git
- ✅ **HTTPS** сертификат
- ✅ **Мониторинг** и логи

**Ваш Go API теперь в продакшене! 🚀**