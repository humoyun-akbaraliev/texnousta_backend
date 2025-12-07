// examples/react-integration/src/App.jsx
import React, { useState, useEffect } from 'react';
import ProductList from './components/ProductList';
import { QuickContactForm, ContactForm, ContactSubmissionsList } from './components/ContactForms';
import apiService from './services/api';

function App() {
  const [categories, setCategories] = useState([]);
  const [selectedCategory, setSelectedCategory] = useState(null);
  const [user, setUser] = useState(null);
  const [activeTab, setActiveTab] = useState('products');

  useEffect(() => {
    loadCategories();
    checkAuth();
  }, []);

  const loadCategories = async () => {
    try {
      const response = await apiService.getCategories();
      setCategories(response.categories || []);
    } catch (error) {
      console.error('Ошибка загрузки категорий:', error);
    }
  };

  const checkAuth = async () => {
    try {
      const response = await apiService.getProfile();
      setUser(response.user);
    } catch (error) {
      // Пользователь не авторизован
    }
  };

  const handleLogin = async () => {
    try {
      const response = await apiService.login('admin@texnousta.com', 'password');
      setUser(response.user);
      alert('Успешная авторизация!');
    } catch (error) {
      alert('Ошибка авторизации: ' + error.message);
    }
  };

  const handleLogout = () => {
    apiService.logout();
    setUser(null);
  };

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Навигация */}
      <nav className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <h1 className="text-xl font-bold text-gray-900">TexnoUsta</h1>
            </div>

            <div className="flex items-center space-x-4">
              <div className="flex space-x-2">
                <button
                  onClick={() => setActiveTab('products')}
                  className={`px-3 py-2 rounded-md text-sm font-medium ${
                    activeTab === 'products'
                      ? 'bg-blue-100 text-blue-700'
                      : 'text-gray-500 hover:text-gray-700'
                  }`}
                >
                  Товары
                </button>
                <button
                  onClick={() => setActiveTab('contact')}
                  className={`px-3 py-2 rounded-md text-sm font-medium ${
                    activeTab === 'contact'
                      ? 'bg-blue-100 text-blue-700'
                      : 'text-gray-500 hover:text-gray-700'
                  }`}
                >
                  Контакты
                </button>
                {user?.role === 'admin' && (
                  <button
                    onClick={() => setActiveTab('admin')}
                    className={`px-3 py-2 rounded-md text-sm font-medium ${
                      activeTab === 'admin'
                        ? 'bg-red-100 text-red-700'
                        : 'text-gray-500 hover:text-gray-700'
                    }`}
                  >
                    Админка
                  </button>
                )}
              </div>

              {user ? (
                <div className="flex items-center space-x-3">
                  <span className="text-sm text-gray-700">
                    Привет, {user.name}!
                  </span>
                  <button
                    onClick={handleLogout}
                    className="text-sm text-red-600 hover:text-red-800"
                  >
                    Выйти
                  </button>
                </div>
              ) : (
                <button
                  onClick={handleLogin}
                  className="bg-blue-500 text-white px-4 py-2 rounded-md text-sm hover:bg-blue-600"
                >
                  Войти (админ)
                </button>
              )}
            </div>
          </div>
        </div>
      </nav>

      {/* Основное содержимое */}
      <main className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
        {/* Вкладка товаров */}
        {activeTab === 'products' && (
          <div className="space-y-6">
            <div className="bg-white p-6 rounded-lg shadow-sm">
              <h2 className="text-lg font-semibold mb-4">Каталог товаров</h2>
              
              {/* Фильтр по категориям */}
              <div className="flex flex-wrap gap-2 mb-6">
                <button
                  onClick={() => setSelectedCategory(null)}
                  className={`px-4 py-2 rounded-md text-sm ${
                    selectedCategory === null
                      ? 'bg-blue-500 text-white'
                      : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                  }`}
                >
                  Все категории
                </button>
                {categories.map((category) => (
                  <button
                    key={category.id}
                    onClick={() => setSelectedCategory(category.id)}
                    className={`px-4 py-2 rounded-md text-sm ${
                      selectedCategory === category.id
                        ? 'bg-blue-500 text-white'
                        : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                    }`}
                  >
                    {category.name}
                  </button>
                ))}
              </div>
            </div>

            {/* Рекомендуемые товары */}
            <div className="bg-white p-6 rounded-lg shadow-sm">
              <h3 className="text-lg font-semibold mb-4">🔥 Рекомендуемые товары</h3>
              <ProductList featured={true} limit={4} />
            </div>

            {/* Все товары */}
            <div className="bg-white p-6 rounded-lg shadow-sm">
              <h3 className="text-lg font-semibold mb-4">
                {selectedCategory 
                  ? `Товары в категории: ${categories.find(c => c.id === selectedCategory)?.name}`
                  : 'Все товары'
                }
              </h3>
              <ProductList categoryId={selectedCategory} />
            </div>
          </div>
        )}

        {/* Вкладка контактов */}
        {activeTab === 'contact' && (
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <QuickContactForm
              onSuccess={(response) => alert('Заявка принята! ID: ' + response.id)}
            />
            <ContactForm
              onSuccess={(response) => alert('Сообщение отправлено! ID: ' + response.id)}
            />
          </div>
        )}

        {/* Админская панель */}
        {activeTab === 'admin' && user?.role === 'admin' && (
          <div className="space-y-6">
            <div className="bg-white p-6 rounded-lg shadow-sm">
              <h2 className="text-lg font-semibold mb-4">Админская панель</h2>
              <p className="text-gray-600">
                Добро пожаловать в админскую панель! Здесь вы можете управлять товарами,
                категориями и просматривать обращения клиентов.
              </p>
            </div>

            <ContactSubmissionsList />
          </div>
        )}
      </main>
    </div>
  );
}

export default App;