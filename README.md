# TexnoUsta Backend API

Backend API для интернет-магазина TexnoUsta, написанный на Go с использованием Gin фреймворка.

## Возможности

- 🔐 Аутентификация и авторизация с JWT
- 📱 CRUD операции для товаров
- 🏷️ Управление категориями товаров
- 👥 Управление пользователями (админ панель)
- 🔍 Поиск и фильтрация товаров
- 📄 Пагинация результатов
- 🌐 CORS поддержка для фронтенда
- 📸 Загрузка файлов (изображения)

## Технологии

- **Go 1.21+**
- **Gin** - веб-фреймворк
- **GORM** - ORM для работы с базой данных
- **PostgreSQL** - база данных
- **JWT** - токены авторизации
- **bcrypt** - хеширование паролей

## Установка и запуск

### 1. Клонирование репозитория
```bash
git clone <repository-url>
cd TexnoUsta_Backend
```

### 2. Установка зависимостей
```bash
go mod tidy
```

### 3. Настройка базы данных
Установите PostgreSQL и создайте базу данных:
```sql
CREATE DATABASE texnousta;
```

### 4. Настройка переменных окружения
Скопируйте файл `.env` и настройте параметры:
```bash
cp .env.example .env
```

### 5. Запуск сервера
```bash
go run cmd/main.go
```

## 🌐 Доступ к сервисам

### Локальная разработка:
- **🔗 API Server:** `http://localhost:8080`
- **📚 Swagger Documentation:** `http://localhost:8080/swagger/index.html`
- **🧪 HTML API Tester:** Откройте файл `api-tester.html` в браузере
- **📖 Полная документация:** см. файлы `API_ENDPOINTS.md` и `SWAGGER_GUIDE.md`

### 🚀 Деплой в продакшен:
- **Railway (рекомендуется):** См. подробные инструкции в `DEPLOY.md`
- **Render:** Альтернативный вариант с бесплатным планом
- **DigitalOcean App Platform:** Для масштабируемых проектов
- ⚠️ **Vercel НЕ поддерживает Go** - используйте альтернативы выше

Основные параметры в `.env`:
```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=texnousta
JWT_SECRET=your_jwt_secret_key
FRONTEND_URL=http://localhost:3000
```

### 5. Запуск приложения
```bash
# Режим разработки
go run cmd/main.go

# Или сборка и запуск
go build -o server cmd/main.go
./server
```

## API Эндпоинты

### Аутентификация
- `POST /api/v1/register` - Регистрация пользователя
- `POST /api/v1/login` - Вход пользователя
- `GET /api/v1/profile` - Получить профиль (требует авторизации)
- `PUT /api/v1/profile` - Обновить профиль (требует авторизации)

### Товары
- `GET /api/v1/products` - Список товаров (с фильтрами)
- `GET /api/v1/products/:id` - Получить товар по ID
- `POST /api/v1/admin/products` - Создать товар (только админ)
- `PUT /api/v1/admin/products/:id` - Обновить товар (только админ)
- `DELETE /api/v1/admin/products/:id` - Удалить товар (только админ)

### Категории
- `GET /api/v1/categories` - Список категорий
- `POST /api/v1/admin/categories` - Создать категорию (только админ)
- `PUT /api/v1/admin/categories/:id` - Обновить категорию (только админ)
- `DELETE /api/v1/admin/categories/:id` - Удалить категорию (только админ)

### Пользователи (Админ)
- `GET /api/v1/admin/users` - Список пользователей
- `PUT /api/v1/admin/users/:id` - Обновить пользователя
- `DELETE /api/v1/admin/users/:id` - Удалить пользователя

### Контактная форма
- `POST /api/v1/contact` - Отправить контактную форму
- `POST /api/v1/quick-contact` - Быстрая заявка (имя + телефон)
- `GET /api/v1/admin/contacts` - Список обращений (только админ)
- `GET /api/v1/admin/contacts/:id` - Получить обращение (только админ)
- `PUT /api/v1/admin/contacts/:id/read` - Пометить как прочитанное (только админ)
- `DELETE /api/v1/admin/contacts/:id` - Удалить обращение (только админ)

## Примеры запросов

### Регистрация
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Иван Иванов",
    "email": "ivan@example.com",
    "password": "password123",
    "phone": "+998901234567"
  }'
```

### Вход
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "ivan@example.com",
    "password": "password123"
  }'
```

### Получение товаров
```bash
# Все товары
curl http://localhost:8080/api/v1/products

# С фильтрами
curl "http://localhost:8080/api/v1/products?category=1&search=iphone&page=1&limit=10"
```

### Создание товара (требует авторизации админа)
```bash
curl -X POST http://localhost:8080/api/v1/admin/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "iPhone 15",
    "description": "Новый iPhone 15",
    "price": 1000,
    "category_id": 1,
    "brand": "Apple",
    "stock": 10,
    "is_active": true
  }'
```

### Отправка контактной формы
```bash
# Полная контактная форма
curl -X POST http://localhost:8080/api/v1/contact \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Иван Петров",
    "email": "ivan@example.com",
    "phone": "+998901234567",
    "subject": "Консультация",
    "message": "Хочу узнать о наличии iPhone 15"
  }'

# Быстрая заявка
curl -X POST http://localhost:8080/api/v1/quick-contact \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Иван Петров",
    "phone": "+998901234567"
  }'
```

## Подключение к фронтенду React

### Пример конфигурации API клиента

```javascript
// api/client.js
import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api/v1';

// Создание экземпляра axios
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Интерцептор для добавления токена
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

export default api;
```

### Примеры использования в React

```javascript
// services/authService.js
import api from '../api/client';

export const authService = {
  async login(email, password) {
    const response = await api.post('/login', { email, password });
    const { token, user } = response.data;
    
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(user));
    
    return { token, user };
  },

  async register(userData) {
    const response = await api.post('/register', userData);
    const { token, user } = response.data;
    
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(user));
    
    return { token, user };
  },

  logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  }
};

// services/productService.js
import api from '../api/client';

export const productService = {
  async getProducts(params = {}) {
    const response = await api.get('/products', { params });
    return response.data;
  },

  async getProduct(id) {
    const response = await api.get(`/products/${id}`);
    return response.data;
  },

  async getCategories() {
    const response = await api.get('/categories');
    return response.data;
  }
};
```

## Тестовые данные

При первом запуске автоматически создаются:

**Администратор:**
- Email: admin@texnousta.com
- Пароль: password

**Категории:**
- Смартфоны
- Ноутбуки  
- Телевизоры
- Бытовая техника
- Аксессуары

**Товары:**
- iPhone 15 Pro
- Samsung Galaxy S24
- MacBook Pro 16"
- LG OLED TV 55"

## Структура проекта

```
TexnoUsta_Backend/
├── cmd/
│   └── main.go              # Точка входа
├── internal/
│   ├── database/
│   │   └── database.go      # Настройка БД и миграции
│   ├── handlers/
│   │   ├── auth.go          # Обработчики аутентификации
│   │   ├── products.go      # Обработчики товаров
│   │   └── admin.go         # Админские обработчики
│   ├── middleware/
│   │   └── auth.go          # Middleware авторизации
│   └── models/
│       └── models.go        # Модели данных
├── uploads/                 # Загруженные файлы
├── .env                     # Переменные окружения
├── go.mod                   # Go модули
└── README.md
```

## Лицензия

MIT License