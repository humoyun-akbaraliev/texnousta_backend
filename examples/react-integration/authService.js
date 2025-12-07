// services/authService.js
import api from '../api/client';

export const authService = {
  // Вход в систему
  async login(email, password) {
    try {
      const response = await api.post('/login', { email, password });
      const { token, user } = response.data;
      
      // Сохранение в localStorage
      localStorage.setItem('token', token);
      localStorage.setItem('user', JSON.stringify(user));
      
      return { success: true, token, user };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Ошибка при входе в систему'
      };
    }
  },

  // Регистрация
  async register(userData) {
    try {
      const response = await api.post('/register', userData);
      const { token, user } = response.data;
      
      // Сохранение в localStorage
      localStorage.setItem('token', token);
      localStorage.setItem('user', JSON.stringify(user));
      
      return { success: true, token, user };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Ошибка при регистрации'
      };
    }
  },

  // Получение профиля
  async getProfile() {
    try {
      const response = await api.get('/profile');
      return { success: true, user: response.data.user };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Ошибка при получении профиля'
      };
    }
  },

  // Обновление профиля
  async updateProfile(userData) {
    try {
      const response = await api.put('/profile', userData);
      const { user } = response.data;
      
      // Обновление в localStorage
      localStorage.setItem('user', JSON.stringify(user));
      
      return { success: true, user };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Ошибка при обновлении профиля'
      };
    }
  },

  // Выход из системы
  logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  },

  // Проверка авторизации
  isAuthenticated() {
    return !!localStorage.getItem('token');
  },

  // Получение текущего пользователя
  getCurrentUser() {
    const user = localStorage.getItem('user');
    return user ? JSON.parse(user) : null;
  },

  // Проверка роли админа
  isAdmin() {
    const user = this.getCurrentUser();
    return user && user.role === 'admin';
  }
};