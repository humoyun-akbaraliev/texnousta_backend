// components/ContactForm.js
import React, { useState } from 'react';
import { contactService } from '../services/contactService';

// Полная контактная форма
export const ContactForm = ({ onSuccess }) => {
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    phone: '',
    subject: '',
    message: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setSuccess('');

    const result = await contactService.submitContactForm(formData);
    
    if (result.success) {
      setSuccess(result.message);
      setFormData({
        name: '',
        email: '',
        phone: '',
        subject: '',
        message: ''
      });
      if (onSuccess) onSuccess();
    } else {
      setError(result.error);
    }
    
    setLoading(false);
  };

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  return (
    <div className="contact-form">
      <h3>Свяжитесь с нами</h3>
      
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="name">Имя *</label>
          <input
            type="text"
            id="name"
            name="name"
            value={formData.name}
            onChange={handleChange}
            required
            placeholder="Ваше имя"
          />
        </div>

        <div className="form-group">
          <label htmlFor="email">Email</label>
          <input
            type="email"
            id="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
            placeholder="your@email.com"
          />
        </div>

        <div className="form-group">
          <label htmlFor="phone">Телефон *</label>
          <input
            type="tel"
            id="phone"
            name="phone"
            value={formData.phone}
            onChange={handleChange}
            required
            placeholder="+998 90 123 45 67"
          />
        </div>

        <div className="form-group">
          <label htmlFor="subject">Тема *</label>
          <select
            id="subject"
            name="subject"
            value={formData.subject}
            onChange={handleChange}
            required
          >
            <option value="">Выберите тему</option>
            <option value="Консультация">Консультация</option>
            <option value="Заказ">Заказ товара</option>
            <option value="Гарантия">Вопрос по гарантии</option>
            <option value="Ремонт">Ремонт</option>
            <option value="Другое">Другое</option>
          </select>
        </div>

        <div className="form-group">
          <label htmlFor="message">Сообщение *</label>
          <textarea
            id="message"
            name="message"
            value={formData.message}
            onChange={handleChange}
            required
            rows={4}
            placeholder="Расскажите подробнее о вашем вопросе"
          />
        </div>

        {error && <div className="error">{error}</div>}
        {success && <div className="success">{success}</div>}

        <button type="submit" disabled={loading} className="btn btn-primary">
          {loading ? 'Отправка...' : 'Отправить'}
        </button>
      </form>
    </div>
  );
};

// Быстрая форма обратного звонка
export const QuickContactForm = ({ onSuccess, buttonText = "Заказать звонок" }) => {
  const [formData, setFormData] = useState({
    name: '',
    phone: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [isOpen, setIsOpen] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setSuccess('');

    const result = await contactService.submitQuickContact(formData.name, formData.phone);
    
    if (result.success) {
      setSuccess(result.message);
      setFormData({ name: '', phone: '' });
      if (onSuccess) onSuccess();
      
      // Автоматически закрыть форму через 3 секунды
      setTimeout(() => {
        setIsOpen(false);
        setSuccess('');
      }, 3000);
    } else {
      setError(result.error);
    }
    
    setLoading(false);
  };

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  return (
    <>
      <button 
        className="btn btn-callback"
        onClick={() => setIsOpen(true)}
      >
        {buttonText}
      </button>

      {isOpen && (
        <div className="modal-overlay" onClick={() => setIsOpen(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h3>Заказать обратный звонок</h3>
              <button 
                className="close-button"
                onClick={() => setIsOpen(false)}
              >
                ×
              </button>
            </div>

            <div className="modal-body">
              <p>Оставьте свои контактные данные и мы перезвоним вам в течение 15 минут</p>
              
              <form onSubmit={handleSubmit}>
                <div className="form-group">
                  <input
                    type="text"
                    name="name"
                    value={formData.name}
                    onChange={handleChange}
                    required
                    placeholder="Ваше имя"
                  />
                </div>

                <div className="form-group">
                  <input
                    type="tel"
                    name="phone"
                    value={formData.phone}
                    onChange={handleChange}
                    required
                    placeholder="+998 90 123 45 67"
                  />
                </div>

                {error && <div className="error">{error}</div>}
                {success && <div className="success">{success}</div>}

                <div className="form-actions">
                  <button type="submit" disabled={loading} className="btn btn-primary">
                    {loading ? 'Отправка...' : 'Заказать звонок'}
                  </button>
                  <button 
                    type="button" 
                    onClick={() => setIsOpen(false)}
                    className="btn btn-secondary"
                  >
                    Отмена
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}
    </>
  );
};

// Компонент для отображения контактных обращений в админке
export const ContactsList = () => {
  const [contacts, setContacts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [pagination, setPagination] = useState({});
  const [showUnreadOnly, setShowUnreadOnly] = useState(false);

  const fetchContacts = async (page = 1) => {
    setLoading(true);
    const result = await adminContactService.getContacts({
      page,
      unread: showUnreadOnly ? 'true' : undefined
    });
    
    if (result.success) {
      setContacts(result.contacts);
      setPagination(result.pagination);
    }
    
    setLoading(false);
  };

  const handleMarkAsRead = async (id) => {
    const result = await adminContactService.markAsRead(id);
    if (result.success) {
      fetchContacts(pagination.page);
    }
  };

  const handleDelete = async (id) => {
    if (window.confirm('Вы уверены, что хотите удалить это обращение?')) {
      const result = await adminContactService.deleteContact(id);
      if (result.success) {
        fetchContacts(pagination.page);
      }
    }
  };

  useEffect(() => {
    fetchContacts();
  }, [showUnreadOnly]);

  if (loading) {
    return <div className="loading">Загрузка обращений...</div>;
  }

  return (
    <div className="contacts-list">
      <div className="contacts-header">
        <h2>Контактные обращения</h2>
        <label className="checkbox-label">
          <input
            type="checkbox"
            checked={showUnreadOnly}
            onChange={(e) => setShowUnreadOnly(e.target.checked)}
          />
          Показать только непрочитанные
        </label>
      </div>

      <div className="contacts-table">
        <table>
          <thead>
            <tr>
              <th>Дата</th>
              <th>Имя</th>
              <th>Телефон</th>
              <th>Email</th>
              <th>Тема</th>
              <th>Статус</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            {contacts.map(contact => (
              <tr key={contact.id} className={!contact.is_read ? 'unread' : ''}>
                <td>{new Date(contact.created_at).toLocaleDateString()}</td>
                <td>{contact.name}</td>
                <td>{contact.phone}</td>
                <td>{contact.email || '-'}</td>
                <td>{contact.subject}</td>
                <td>
                  {contact.is_read ? (
                    <span className="status-read">Прочитано</span>
                  ) : (
                    <span className="status-unread">Новое</span>
                  )}
                </td>
                <td>
                  {!contact.is_read && (
                    <button onClick={() => handleMarkAsRead(contact.id)}>
                      Прочитано
                    </button>
                  )}
                  <button onClick={() => handleDelete(contact.id)}>
                    Удалить
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};