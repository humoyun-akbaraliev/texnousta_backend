# 🛍️ TexnoUsta API Integration Examples

Полное руководство по интеграции TexnoUsta Go API с React фронтендом.

## Установка зависимостей

```bash
npm install axios
```

## Переменные окружения React (.env)

```env
REACT_APP_API_URL=http://localhost:8080/api/v1
```

## Структура файлов

```
src/
├── api/
│   └── client.js           # Настройка axios
├── services/
│   ├── authService.js      # Сервис авторизации
│   ├── productService.js   # Сервис товаров
│   └── categoryService.js  # Сервис категорий
├── hooks/
│   └── useAuth.js         # Хуки для работы с API
├── components/
│   ├── ProductList.js     # Компонент списка товаров
│   ├── ProductCard.js     # Компонент карточки товара
│   └── LoginForm.js       # Компонент формы входа
└── App.js                 # Основной компонент
```

## Использование в App.js

```javascript
import React from 'react';
import { AuthProvider, useAuth } from './hooks/useAuth';
import { ProductList, LoginForm } from './components';

function App() {
  return (
    <AuthProvider>
      <div className="App">
        <Header />
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/products" element={<ProductsPage />} />
          <Route path="/products/:id" element={<ProductDetailPage />} />
        </Routes>
      </div>
    </AuthProvider>
  );
}

const HomePage = () => {
  return (
    <div>
      <h1>Добро пожаловать в TexnoUsta</h1>
      <ProductList featured={true} />
    </div>
  );
};

const ProductsPage = () => {
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('');

  return (
    <div>
      <div className="filters">
        <input 
          type="text" 
          placeholder="Поиск товаров..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
        <CategoryFilter 
          selectedCategory={selectedCategory}
          onCategoryChange={setSelectedCategory}
        />
      </div>
      <ProductList 
        searchQuery={searchQuery}
        categoryId={selectedCategory}
      />
    </div>
  );
};

export default App;
```

## Защищенные маршруты

```javascript
// components/ProtectedRoute.js
import { useAuth } from '../hooks/useAuth';
import { Navigate } from 'react-router-dom';

const ProtectedRoute = ({ children, adminOnly = false }) => {
  const { isAuthenticated, isAdmin, loading } = useAuth();

  if (loading) {
    return <div>Загрузка...</div>;
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (adminOnly && !isAdmin) {
    return <Navigate to="/" replace />;
  }

  return children;
};

export default ProtectedRoute;
```

## CSS стили (базовые)

```css
/* styles.css */
.product-list {
  padding: 20px;
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.product-card {
  border: 1px solid #ddd;
  border-radius: 8px;
  overflow: hidden;
  transition: transform 0.2s;
}

.product-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.product-image {
  position: relative;
  height: 200px;
  overflow: hidden;
}

.product-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.discount-badge {
  position: absolute;
  top: 10px;
  right: 10px;
  background: #ff4757;
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: bold;
}

.product-info {
  padding: 15px;
}

.product-name {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 5px;
}

.product-brand {
  color: #666;
  font-size: 14px;
  margin-bottom: 10px;
}

.product-prices {
  margin-bottom: 15px;
}

.current-price {
  font-size: 18px;
  font-weight: bold;
  color: #2ed573;
}

.old-price {
  font-size: 14px;
  color: #666;
  text-decoration: line-through;
  margin-left: 10px;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
  transition: background-color 0.2s;
}

.btn-primary {
  background-color: #3742fa;
  color: white;
}

.btn-primary:hover {
  background-color: #2f3542;
}

.btn:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.login-form {
  max-width: 400px;
  margin: 0 auto;
  padding: 20px;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}

.form-group input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.error {
  color: #ff4757;
  background-color: #ffe3e3;
  padding: 10px;
  border-radius: 4px;
  margin: 10px 0;
}

.loading {
  text-align: center;
  padding: 20px;
  color: #666;
}
```

## Запуск

1. Убедитесь, что Go API запущено на порту 8080
2. Запустите React приложение:
   ```bash
   npm start
   ```
3. Откройте http://localhost:3000

API автоматически будет доступно по адресу http://localhost:8080/api/v1

## Тестирование интеграции

Для тестирования можете использовать тестовые данные:
- Админ: admin@texnousta.com / password
- Или создать нового пользователя через регистрацию