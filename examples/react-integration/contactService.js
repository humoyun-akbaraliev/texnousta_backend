// services/contactService.js
import api from '../api/client';

export const contactService = {
  // Отправка полной контактной формы
  async submitContactForm(formData) {
    try {
      const response = await api.post('/contact', formData);
      return { success: true, message: response.data.message };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Ошибка при отправке формы'
      };
    }
  },

  // Отправка быстрой заявки (только имя и телефон)
  async submitQuickContact(name, phone) {
    try {
      const response = await api.post('/quick-contact', { name, phone });
      return { success: true, message: response.data.message };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Ошибка при отправке заявки'
      };
    }
  }
};

// Админские функции для управления контактными обращениями
export const adminContactService = {
  // Получение списка обращений
  async getContacts(params = {}) {
    try {
      const response = await api.get('/admin/contacts', { params });
      return { 
        success: true, 
        contacts: response.data.contacts,
        pagination: response.data.pagination 
      };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Ошибка при получении обращений'
      };
    }
  },

  // Получение конкретного обращения
  async getContact(id) {
    try {
      const response = await api.get(`/admin/contacts/${id}`);
      return { success: true, contact: response.data.contact };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Ошибка при получении обращения'
      };
    }
  },

  // Пометить как прочитанное
  async markAsRead(id) {
    try {
      const response = await api.put(`/admin/contacts/${id}/read`);
      return { success: true, message: response.data.message };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Ошибка при обновлении статуса'
      };
    }
  },

  // Удаление обращения
  async deleteContact(id) {
    try {
      const response = await api.delete(`/admin/contacts/${id}`);
      return { success: true, message: response.data.message };
    } catch (error) {
      return {
        success: false,
        error: error.response?.data?.error || 'Ошибка при удалении обращения'
      };
    }
  }
};