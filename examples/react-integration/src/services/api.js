// examples/react-integration/src/services/api.js

const API_BASE_URL = 'http://localhost:8080/api/v1';

class ApiService {
  constructor() {
    this.baseURL = API_BASE_URL;
    this.token = localStorage.getItem('authToken');
  }

  // Получить заголовки запроса
  getHeaders(requireAuth = false) {
    const headers = {
      'Content-Type': 'application/json',
    };

    if (requireAuth && this.token) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }

    return headers;
  }

  // Базовый метод для запросов
  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    const config = {
      headers: this.getHeaders(options.requireAuth),
      ...options,
    };

    try {
      const response = await fetch(url, config);
      
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || `HTTP error! status: ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error(`API request failed: ${endpoint}`, error);
      throw error;
    }
  }

  // Установить токен
  setToken(token) {
    this.token = token;
    localStorage.setItem('authToken', token);
  }

  // Удалить токен
  clearToken() {
    this.token = null;
    localStorage.removeItem('authToken');
  }

  // ========== АВТОРИЗАЦИЯ ==========
  
  // Регистрация
  async register(userData) {
    const response = await this.request('/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    });
    
    if (response.token) {
      this.setToken(response.token);
    }
    
    return response;
  }

  // Вход
  async login(email, password) {
    const response = await this.request('/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });
    
    if (response.token) {
      this.setToken(response.token);
    }
    
    return response;
  }

  // Выход
  logout() {
    this.clearToken();
  }

  // ========== ТОВАРЫ ==========
  
  // Получить список товаров
  async getProducts(params = {}) {
    const queryString = new URLSearchParams(params).toString();
    const endpoint = queryString ? `/products?${queryString}` : '/products';
    return await this.request(endpoint);
  }

  // Получить товар по ID
  async getProduct(id) {
    return await this.request(`/products/${id}`);
  }

  // Получить рекомендуемые товары
  async getFeaturedProducts(limit = 8) {
    return await this.getProducts({ featured: true, limit });
  }

  // Поиск товаров
  async searchProducts(query, params = {}) {
    return await this.getProducts({ search: query, ...params });
  }

  // ========== КАТЕГОРИИ ==========
  
  // Получить все категории
  async getCategories() {
    return await this.request('/categories');
  }

  // ========== КОНТАКТНЫЕ ФОРМЫ ==========
  
  // Отправить быструю заявку
  async sendQuickContact(name, phone) {
    return await this.request('/quick-contact', {
      method: 'POST',
      body: JSON.stringify({ name, phone }),
    });
  }

  // Отправить полную контактную форму
  async sendContact(contactData) {
    return await this.request('/contact', {
      method: 'POST',
      body: JSON.stringify(contactData),
    });
  }

  // ========== ПРОФИЛЬ ПОЛЬЗОВАТЕЛЯ ==========
  
  // Получить профиль
  async getProfile() {
    return await this.request('/profile', { requireAuth: true });
  }

  // Обновить профиль
  async updateProfile(userData) {
    return await this.request('/profile', {
      method: 'PUT',
      requireAuth: true,
      body: JSON.stringify(userData),
    });
  }

  // ========== АДМИНСКИЕ МЕТОДЫ ==========
  
  // Получить контактные обращения (админ)
  async getContactSubmissions(params = {}) {
    const queryString = new URLSearchParams(params).toString();
    const endpoint = queryString ? `/admin/contacts?${queryString}` : '/admin/contacts';
    return await this.request(endpoint, { requireAuth: true });
  }

  // Пометить обращение как прочитанное (админ)
  async markContactAsRead(contactId) {
    return await this.request(`/admin/contacts/${contactId}/read`, {
      method: 'PUT',
      requireAuth: true,
    });
  }

  // Создать товар (админ)
  async createProduct(productData) {
    return await this.request('/admin/products', {
      method: 'POST',
      requireAuth: true,
      body: JSON.stringify(productData),
    });
  }

  // Обновить товар (админ)
  async updateProduct(productId, productData) {
    return await this.request(`/admin/products/${productId}`, {
      method: 'PUT',
      requireAuth: true,
      body: JSON.stringify(productData),
    });
  }

  // Удалить товар (админ)
  async deleteProduct(productId) {
    return await this.request(`/admin/products/${productId}`, {
      method: 'DELETE',
      requireAuth: true,
    });
  }

  // Создать категорию (админ)
  async createCategory(categoryData) {
    return await this.request('/admin/categories', {
      method: 'POST',
      requireAuth: true,
      body: JSON.stringify(categoryData),
    });
  }

  // Обновить категорию (админ)
  async updateCategory(categoryId, categoryData) {
    return await this.request(`/admin/categories/${categoryId}`, {
      method: 'PUT',
      requireAuth: true,
      body: JSON.stringify(categoryData),
    });
  }

  // Удалить категорию (админ)
  async deleteCategory(categoryId) {
    return await this.request(`/admin/categories/${categoryId}`, {
      method: 'DELETE',
      requireAuth: true,
    });
  }

  // Получить пользователей (админ)
  async getUsers(params = {}) {
    const queryString = new URLSearchParams(params).toString();
    const endpoint = queryString ? `/admin/users?${queryString}` : '/admin/users';
    return await this.request(endpoint, { requireAuth: true });
  }

  // Обновить пользователя (админ)
  async updateUser(userId, userData) {
    return await this.request(`/admin/users/${userId}`, {
      method: 'PUT',
      requireAuth: true,
      body: JSON.stringify(userData),
    });
  }

  // Удалить пользователя (админ)
  async deleteUser(userId) {
    return await this.request(`/admin/users/${userId}`, {
      method: 'DELETE',
      requireAuth: true,
    });
  }
}

// Создать экземпляр API сервиса
const apiService = new ApiService();

export default apiService;

// Экспорт отдельных методов для удобства
export const {
  // Авторизация
  register,
  login,
  logout,
  
  // Товары
  getProducts,
  getProduct,
  getFeaturedProducts,
  searchProducts,
  
  // Категории
  getCategories,
  
  // Контакты
  sendQuickContact,
  sendContact,
  
  // Профиль
  getProfile,
  updateProfile,
  
  // Админка
  getContactSubmissions,
  markContactAsRead,
  createProduct,
  updateProduct,
  deleteProduct,
  createCategory,
  updateCategory,
  deleteCategory,
  getUsers,
  updateUser,
  deleteUser,
} = apiService;