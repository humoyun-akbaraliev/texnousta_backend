# API Endpoints для Frontend

## Base URL
```
http://localhost:8080/api/v1
```

## 📋 Список всех доступных эндпоинтов

### 🔐 Авторизация (Публичные)

#### 1. Регистрация пользователя
```bash
POST /api/v1/register

# Пример запроса:
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Иван Петров",
    "email": "ivan@example.com",
    "password": "password123",
    "phone": "+998901234567"
  }'

# Ответ:
{
  "message": "Пользователь успешно зарегистрирован",
  "user": {
    "id": 2,
    "name": "Иван Петров",
    "email": "ivan@example.com",
    "phone": "+998901234567",
    "role": "user"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 2. Вход пользователя
```bash
POST /api/v1/login

# Пример запроса:
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@texnousta.com",
    "password": "password"
  }'

# Ответ:
{
  "message": "Успешная авторизация",
  "user": {
    "id": 1,
    "name": "Администратор",
    "email": "admin@texnousta.com",
    "phone": "+998901234567",
    "role": "admin"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 📱 Товары (Публичные)

#### 3. Получить список товаров
```bash
GET /api/v1/products

# С параметрами:
GET /api/v1/products?page=1&limit=12&category=1&search=iPhone&featured=true

# Пример запроса:
curl "http://localhost:8080/api/v1/products?page=1&limit=2"

# Ответ:
{
  "products": [
    {
      "id": 1,
      "name": "iPhone 15 Pro",
      "description": "Новейший флагманский смартфон от Apple",
      "price": 1200,
      "old_price": 1300,
      "image": "",
      "category_id": 1,
      "brand": "Apple",
      "model": "iPhone 15 Pro",
      "stock": 50,
      "is_active": true,
      "is_featured": true,
      "created_at": "2025-12-08T03:41:33.216996+05:00",
      "updated_at": "2025-12-08T03:41:33.216996+05:00",
      "category": {
        "id": 1,
        "name": "Смартфоны",
        "description": "Мобильные телефоны и смартфоны",
        "image": "",
        "is_active": true,
        "created_at": "2025-12-08T03:41:33.215475+05:00",
        "updated_at": "2025-12-08T03:41:33.215475+05:00"
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 2,
    "total": 4,
    "total_pages": 2
  }
}
```

#### 4. Получить товар по ID
```bash
GET /api/v1/products/{id}

# Пример запроса:
curl http://localhost:8080/api/v1/products/1

# Ответ:
{
  "product": {
    "id": 1,
    "name": "iPhone 15 Pro",
    "description": "Новейший флагманский смартфон от Apple",
    "price": 1200,
    "old_price": 1300,
    "category": {
      "id": 1,
      "name": "Смартфоны"
    }
  }
}
```

### 🏷️ Категории (Публичные)

#### 5. Получить список категорий
```bash
GET /api/v1/categories

# Пример запроса:
curl http://localhost:8080/api/v1/categories

# Ответ:
{
  "categories": [
    {
      "id": 1,
      "name": "Смартфоны",
      "description": "Мобильные телефоны и смартфоны",
      "image": "",
      "is_active": true,
      "created_at": "2025-12-08T03:41:33.215475+05:00",
      "updated_at": "2025-12-08T03:41:33.215475+05:00"
    },
    {
      "id": 2,
      "name": "Ноутбуки",
      "description": "Портативные компьютеры",
      "image": "",
      "is_active": true
    }
  ]
}
```

### 📞 Контактная форма (Публичные)

#### 6. Отправить полную контактную форму
```bash
POST /api/v1/contact

# Пример запроса:
curl -X POST http://localhost:8080/api/v1/contact \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Мария Иванова",
    "email": "maria@example.com",
    "phone": "+998901234568",
    "subject": "Консультация",
    "message": "Хочу узнать о наличии iPhone 15 Pro"
  }'

# Ответ:
{
  "message": "Ваше обращение успешно отправлено. Мы свяжемся с вами в ближайшее время.",
  "id": 3
}
```

#### 7. Отправить быструю заявку (только имя + телефон)
```bash
POST /api/v1/quick-contact

# Пример запроса:
curl -X POST http://localhost:8080/api/v1/quick-contact \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Иван Петров",
    "phone": "+998901234567"
  }'

# Ответ:
{
  "message": "Заявка принята! Мы перезвоним вам в течение 15 минут.",
  "id": 2
}
```

### 👤 Профиль (Требует авторизации)

#### 8. Получить профиль пользователя
```bash
GET /api/v1/profile
Authorization: Bearer {token}

# Пример запроса:
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/v1/profile

# Ответ:
{
  "user": {
    "id": 1,
    "name": "Администратор",
    "email": "admin@texnousta.com",
    "phone": "+998901234567",
    "role": "admin",
    "created_at": "2025-12-08T03:41:33.215998+05:00"
  }
}
```

#### 9. Обновить профиль пользователя
```bash
PUT /api/v1/profile
Authorization: Bearer {token}

# Пример запроса:
curl -X PUT http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Новое Имя",
    "phone": "+998901234569"
  }'

# Ответ:
{
  "message": "Профиль успешно обновлен",
  "user": {
    "id": 1,
    "name": "Новое Имя",
    "email": "admin@texnousta.com",
    "phone": "+998901234569",
    "role": "admin"
  }
}
```

## 🔒 Админские эндпоинты (Требует роль admin)

### Управление товарами

#### 10. Создать товар
```bash
POST /api/v1/admin/products
Authorization: Bearer {admin_token}

# Пример запроса:
curl -X POST http://localhost:8080/api/v1/admin/products \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 16",
    "description": "Новый iPhone 16",
    "price": 1500,
    "old_price": 1600,
    "category_id": 1,
    "brand": "Apple",
    "model": "iPhone 16",
    "stock": 25,
    "is_active": true,
    "is_featured": true
  }'
```

#### 11. Обновить товар
```bash
PUT /api/v1/admin/products/{id}
Authorization: Bearer {admin_token}
```

#### 12. Удалить товар
```bash
DELETE /api/v1/admin/products/{id}
Authorization: Bearer {admin_token}
```

### Управление категориями

#### 13. Создать категорию
```bash
POST /api/v1/admin/categories
Authorization: Bearer {admin_token}

# Пример запроса:
curl -X POST http://localhost:8080/api/v1/admin/categories \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Планшеты",
    "description": "iPad и Android планшеты",
    "is_active": true
  }'
```

#### 14. Обновить категорию
```bash
PUT /api/v1/admin/categories/{id}
Authorization: Bearer {admin_token}
```

#### 15. Удалить категорию
```bash
DELETE /api/v1/admin/categories/{id}
Authorization: Bearer {admin_token}
```

### Управление пользователями

#### 16. Получить список пользователей
```bash
GET /api/v1/admin/users
Authorization: Bearer {admin_token}

# С параметрами:
GET /api/v1/admin/users?page=1&limit=20
```

#### 17. Обновить пользователя
```bash
PUT /api/v1/admin/users/{id}
Authorization: Bearer {admin_token}
```

#### 18. Удалить пользователя
```bash
DELETE /api/v1/admin/users/{id}
Authorization: Bearer {admin_token}
```

### Управление контактными обращениями

#### 19. Получить список контактных обращений
```bash
GET /api/v1/admin/contacts
Authorization: Bearer {admin_token}

# С фильтрами:
GET /api/v1/admin/contacts?unread=true&page=1&limit=20

# Пример запроса:
curl -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  "http://localhost:8080/api/v1/admin/contacts"

# Ответ:
{
  "contacts": [
    {
      "id": 3,
      "name": "Мария Иванова",
      "email": "maria@example.com",
      "phone": "+998901234568",
      "subject": "Консультация",
      "message": "Хочу узнать о наличии iPhone 15 Pro и возможности рассрочки",
      "is_read": false,
      "created_at": "2025-12-08T03:44:46.603018+05:00"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 3,
    "total_pages": 1
  }
}
```

#### 20. Получить конкретное обращение
```bash
GET /api/v1/admin/contacts/{id}
Authorization: Bearer {admin_token}
```

#### 21. Пометить обращение как прочитанное
```bash
PUT /api/v1/admin/contacts/{id}/read
Authorization: Bearer {admin_token}

# Пример запроса:
curl -X PUT -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  http://localhost:8080/api/v1/admin/contacts/1/read

# Ответ:
{
  "message": "Обращение помечено как прочитанное"
}
```

#### 22. Удалить обращение
```bash
DELETE /api/v1/admin/contacts/{id}
Authorization: Bearer {admin_token}
```

## 🎯 Параметры запросов

### Параметры товаров (GET /api/v1/products):
- `page` - номер страницы (по умолчанию: 1)
- `limit` - количество на странице (по умолчанию: 12)
- `category` - ID категории для фильтрации
- `search` - текст для поиска в названии и описании
- `featured` - показать только рекомендуемые товары (true/false)
- `sort` - поле для сортировки (created_at, price, name)
- `order` - порядок сортировки (asc, desc)

### Параметры контактов (GET /api/v1/admin/contacts):
- `page` - номер страницы (по умолчанию: 1)
- `limit` - количество на странице (по умолчанию: 20)
- `unread` - показать только непрочитанные (true/false)

## 🔐 Авторизация

Для защищенных эндпоинтов добавляйте заголовок:
```
Authorization: Bearer YOUR_JWT_TOKEN
```

## 📊 Коды ответов

- `200` - Успешно
- `201` - Создано
- `400` - Ошибка в запросе
- `401` - Не авторизован
- `403` - Доступ запрещен
- `404` - Не найдено
- `500` - Ошибка сервера

## 🧪 Тестовые данные

### Администратор:
- **Email**: admin@texnousta.com
- **Password**: password

### Товары:
- iPhone 15 Pro (ID: 1)
- Samsung Galaxy S24 (ID: 2)
- MacBook Pro 16" (ID: 3)
- LG OLED TV 55" (ID: 4)

### Категории:
- Смартфоны (ID: 1)
- Ноутбуки (ID: 2)
- Телевизоры (ID: 3)
- Бытовая техника (ID: 4)
- Аксессуары (ID: 5)

## 🚨 Примечания для Frontend разработки

1. **CORS настроен** для `http://localhost:3000`
2. **JWT токены** действительны **7 дней**
3. **База данных**: SQLite файл `texnousta.db`
4. **Загрузка файлов**: `/uploads` (статические файлы)
5. **Все даты** в формате ISO 8601 с таймзоной