// components/ProductList.js
import React from 'react';
import { useProducts } from '../hooks/useAuth';

const ProductList = ({ categoryId, searchQuery, featured }) => {
  const { products, loading, error, pagination, fetchProducts } = useProducts({
    category: categoryId,
    search: searchQuery,
    featured: featured ? 'true' : undefined
  });

  if (loading) {
    return <div className="loading">Загрузка товаров...</div>;
  }

  if (error) {
    return <div className="error">Ошибка: {error}</div>;
  }

  return (
    <div className="product-list">
      <div className="products-grid">
        {products.map(product => (
          <ProductCard key={product.id} product={product} />
        ))}
      </div>
      
      {pagination && (
        <Pagination 
          currentPage={pagination.page}
          totalPages={pagination.total_pages}
          onPageChange={(page) => fetchProducts({ page })}
        />
      )}
    </div>
  );
};

// components/ProductCard.js
const ProductCard = ({ product }) => {
  const { addToCart } = useCart();

  return (
    <div className="product-card">
      <div className="product-image">
        <img src={product.image || '/placeholder.jpg'} alt={product.name} />
        {product.old_price && (
          <div className="discount-badge">
            -{Math.round(((product.old_price - product.price) / product.old_price) * 100)}%
          </div>
        )}
      </div>
      
      <div className="product-info">
        <h3 className="product-name">{product.name}</h3>
        <p className="product-brand">{product.brand}</p>
        
        <div className="product-prices">
          <span className="current-price">${product.price}</span>
          {product.old_price && (
            <span className="old-price">${product.old_price}</span>
          )}
        </div>
        
        <div className="product-actions">
          <button 
            className="btn btn-primary"
            onClick={() => addToCart(product)}
            disabled={product.stock === 0}
          >
            {product.stock === 0 ? 'Нет в наличии' : 'В корзину'}
          </button>
        </div>
      </div>
    </div>
  );
};

// components/LoginForm.js
import { useState } from 'react';
import { useAuth } from '../hooks/useAuth';

const LoginForm = () => {
  const [formData, setFormData] = useState({
    email: '',
    password: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  
  const { login } = useAuth();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    const result = await login(formData.email, formData.password);
    
    if (!result.success) {
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
    <form onSubmit={handleSubmit} className="login-form">
      <div className="form-group">
        <label htmlFor="email">Email:</label>
        <input
          type="email"
          id="email"
          name="email"
          value={formData.email}
          onChange={handleChange}
          required
        />
      </div>
      
      <div className="form-group">
        <label htmlFor="password">Пароль:</label>
        <input
          type="password"
          id="password"
          name="password"
          value={formData.password}
          onChange={handleChange}
          required
        />
      </div>
      
      {error && <div className="error">{error}</div>}
      
      <button type="submit" disabled={loading}>
        {loading ? 'Вход...' : 'Войти'}
      </button>
    </form>
  );
};

export { ProductList, ProductCard, LoginForm };