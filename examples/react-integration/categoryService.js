// services/categoryService.js
import api from '../api/client';

export const categoryService = {
  // Получение списка категорий
  async getCategories() {
    try {
      const response = await api.get('/categories');
      return { success: true, categories: response.data.categories };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Ошибка при получении категорий'
      };
    }
  }
};

// Админские функции для управления категориями
export const adminCategoryService = {
  // Создание категории
  async createCategory(categoryData) {
    try {
      const response = await api.post('/admin/categories', categoryData);
      return { success: true, category: response.data.category };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Ошибка при создании категории'
      };
    }
  },

  // Обновление категории
  async updateCategory(id, categoryData) {
    try {
      const response = await api.put(`/admin/categories/${id}`, categoryData);
      return { success: true, category: response.data.category };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Ошибка при обновлении категории'
      };
    }
  },

  // Удаление категории
  async deleteCategory(id) {
    try {
      await api.delete(`/admin/categories/${id}`);
      return { success: true };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Ошибка при удалении категории'
      };
    }
  }
};