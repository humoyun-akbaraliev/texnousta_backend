// services/productService.js
import api from '../api/client';

export const productService = {
  // Получение списка товаров
  async getProducts(params = {}) {
    try {
      const response = await api.get('/products', { params });
      return { 
        success: true, 
        products: response.data.products,
        pagination: response.data.pagination
      };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Ошибка при получении товаров'
      };
    }
  },

  // Получение товара по ID
  async getProduct(id) {
    try {
      const response = await api.get(`/products/${id}`);
      return { success: true, product: response.data.product };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Товар не найден'
      };
    }
  },

  // Поиск товаров
  async searchProducts(query, filters = {}) {
    const params = {
      search: query,
      ...filters
    };
    return this.getProducts(params);
  },

  // Получение рекомендуемых товаров
  async getFeaturedProducts(limit = 8) {
    const params = {
      featured: 'true',
      limit
    };
    return this.getProducts(params);
  },

  // Получение товаров по категории
  async getProductsByCategory(categoryId, params = {}) {
    const queryParams = {
      category: categoryId,
      ...params
    };
    return this.getProducts(queryParams);
  }
};

// Админские функции для управления товарами
export const adminProductService = {
  // Создание товара
  async createProduct(productData) {
    try {
      const response = await api.post('/admin/products', productData);
      return { success: true, product: response.data.product };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Ошибка при создании товара'
      };
    }
  },

  // Обновление товара
  async updateProduct(id, productData) {
    try {
      const response = await api.put(`/admin/products/${id}`, productData);
      return { success: true, product: response.data.product };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Ошибка при обновлении товара'
      };
    }
  },

  // Удаление товара
  async deleteProduct(id) {
    try {
      await api.delete(`/admin/products/${id}`);
      return { success: true };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Ошибка при удалении товара'
      };
    }
  }
};