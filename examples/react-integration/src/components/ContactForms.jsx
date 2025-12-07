// examples/react-integration/src/components/ContactForms.jsx
import React, { useState } from 'react';
import apiService from '../services/api';

// Быстрая форма контакта
export const QuickContactForm = ({ onSuccess, className = '' }) => {
  const [formData, setFormData] = useState({
    name: '',
    phone: '',
  });
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState({ type: '', text: '' });

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!formData.name.trim() || !formData.phone.trim()) {
      setMessage({ type: 'error', text: 'Заполните все поля' });
      return;
    }

    try {
      setLoading(true);
      setMessage({ type: '', text: '' });

      const response = await apiService.sendQuickContact(formData.name, formData.phone);
      
      setMessage({ type: 'success', text: response.message });
      setFormData({ name: '', phone: '' });
      
      if (onSuccess) {
        onSuccess(response);
      }
    } catch (error) {
      setMessage({ type: 'error', text: error.message });
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  return (
    <div className={`bg-white p-6 rounded-lg shadow-md ${className}`}>
      <h3 className="text-lg font-semibold mb-4">Быстрая заявка</h3>
      <p className="text-gray-600 text-sm mb-4">
        Оставьте номер телефона и мы перезвоним в течение 15 минут
      </p>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Ваше имя *
          </label>
          <input
            type="text"
            name="name"
            value={formData.name}
            onChange={handleChange}
            placeholder="Введите ваше имя"
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            disabled={loading}
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Номер телефона *
          </label>
          <input
            type="tel"
            name="phone"
            value={formData.phone}
            onChange={handleChange}
            placeholder="+998 90 123 45 67"
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            disabled={loading}
          />
        </div>

        {message.text && (
          <div className={`p-3 rounded-md text-sm ${
            message.type === 'success' 
              ? 'bg-green-100 text-green-700' 
              : 'bg-red-100 text-red-700'
          }`}>
            {message.text}
          </div>
        )}

        <button
          type="submit"
          disabled={loading}
          className="w-full bg-blue-500 text-white py-2 px-4 rounded-md hover:bg-blue-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {loading ? (
            <span className="flex items-center justify-center">
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
              Отправляем...
            </span>
          ) : (
            'Заказать звонок'
          )}
        </button>
      </form>
    </div>
  );
};

// Полная форма контакта
export const ContactForm = ({ onSuccess, className = '' }) => {
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    phone: '',
    subject: '',
    message: '',
  });
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState({ type: '', text: '' });

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    const required = ['name', 'email', 'phone', 'subject', 'message'];
    const missingFields = required.filter(field => !formData[field].trim());
    
    if (missingFields.length > 0) {
      setMessage({ type: 'error', text: 'Заполните все поля' });
      return;
    }

    // Простая валидация email
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(formData.email)) {
      setMessage({ type: 'error', text: 'Введите корректный email' });
      return;
    }

    try {
      setLoading(true);
      setMessage({ type: '', text: '' });

      const response = await apiService.sendContact(formData);
      
      setMessage({ type: 'success', text: response.message });
      setFormData({
        name: '',
        email: '',
        phone: '',
        subject: '',
        message: '',
      });
      
      if (onSuccess) {
        onSuccess(response);
      }
    } catch (error) {
      setMessage({ type: 'error', text: error.message });
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  return (
    <div className={`bg-white p-6 rounded-lg shadow-md ${className}`}>
      <h3 className="text-lg font-semibold mb-4">Связаться с нами</h3>
      <p className="text-gray-600 text-sm mb-6">
        Заполните форму и мы свяжемся с вами в ближайшее время
      </p>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Имя *
            </label>
            <input
              type="text"
              name="name"
              value={formData.name}
              onChange={handleChange}
              placeholder="Ваше имя"
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={loading}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Email *
            </label>
            <input
              type="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              placeholder="your@email.com"
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={loading}
            />
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Телефон *
            </label>
            <input
              type="tel"
              name="phone"
              value={formData.phone}
              onChange={handleChange}
              placeholder="+998 90 123 45 67"
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={loading}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Тема *
            </label>
            <select
              name="subject"
              value={formData.subject}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={loading}
            >
              <option value="">Выберите тему</option>
              <option value="Консультация">Консультация</option>
              <option value="Заказ товара">Заказ товара</option>
              <option value="Техподдержка">Техподдержка</option>
              <option value="Партнерство">Партнерство</option>
              <option value="Другое">Другое</option>
            </select>
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Сообщение *
          </label>
          <textarea
            name="message"
            value={formData.message}
            onChange={handleChange}
            placeholder="Опишите ваш вопрос или пожелание..."
            rows={4}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            disabled={loading}
          />
        </div>

        {message.text && (
          <div className={`p-3 rounded-md text-sm ${
            message.type === 'success' 
              ? 'bg-green-100 text-green-700' 
              : 'bg-red-100 text-red-700'
          }`}>
            {message.text}
          </div>
        )}

        <button
          type="submit"
          disabled={loading}
          className="w-full bg-blue-500 text-white py-3 px-6 rounded-md hover:bg-blue-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {loading ? (
            <span className="flex items-center justify-center">
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
              Отправляем...
            </span>
          ) : (
            'Отправить сообщение'
          )}
        </button>
      </form>
    </div>
  );
};

// Компонент для отображения контактных обращений (для админки)
export const ContactSubmissionsList = ({ className = '' }) => {
  const [contacts, setContacts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState('all'); // all, unread, read

  useEffect(() => {
    fetchContacts();
  }, [filter]);

  const fetchContacts = async () => {
    try {
      setLoading(true);
      
      const params = {};
      if (filter === 'unread') params.unread = true;
      
      const response = await apiService.getContactSubmissions(params);
      setContacts(response.contacts || []);
    } catch (error) {
      console.error('Ошибка загрузки обращений:', error);
    } finally {
      setLoading(false);
    }
  };

  const markAsRead = async (contactId) => {
    try {
      await apiService.markContactAsRead(contactId);
      setContacts(contacts.map(contact => 
        contact.id === contactId 
          ? { ...contact, is_read: true }
          : contact
      ));
    } catch (error) {
      console.error('Ошибка при пометке как прочитанное:', error);
    }
  };

  if (loading) {
    return <div className="text-center py-8">Загрузка...</div>;
  }

  return (
    <div className={className}>
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-xl font-semibold">Контактные обращения</h2>
        
        <select
          value={filter}
          onChange={(e) => setFilter(e.target.value)}
          className="px-3 py-2 border border-gray-300 rounded-md"
        >
          <option value="all">Все обращения</option>
          <option value="unread">Непрочитанные</option>
        </select>
      </div>

      <div className="space-y-4">
        {contacts.length === 0 ? (
          <p className="text-gray-500 text-center py-8">Обращений нет</p>
        ) : (
          contacts.map((contact) => (
            <div
              key={contact.id}
              className={`p-4 border rounded-lg ${
                contact.is_read ? 'bg-gray-50' : 'bg-white border-blue-200'
              }`}
            >
              <div className="flex justify-between items-start mb-2">
                <div>
                  <h3 className="font-semibold">{contact.name}</h3>
                  <p className="text-sm text-gray-600">
                    {contact.email} • {contact.phone}
                  </p>
                </div>
                <div className="text-right">
                  <p className="text-xs text-gray-500">
                    {new Date(contact.created_at).toLocaleString()}
                  </p>
                  {!contact.is_read && (
                    <button
                      onClick={() => markAsRead(contact.id)}
                      className="text-xs text-blue-600 hover:underline mt-1"
                    >
                      Пометить прочитанным
                    </button>
                  )}
                </div>
              </div>
              
              {contact.subject && (
                <p className="text-sm font-medium text-gray-700 mb-2">
                  Тема: {contact.subject}
                </p>
              )}
              
              {contact.message && (
                <p className="text-sm text-gray-700">
                  {contact.message}
                </p>
              )}
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default { QuickContactForm, ContactForm, ContactSubmissionsList };