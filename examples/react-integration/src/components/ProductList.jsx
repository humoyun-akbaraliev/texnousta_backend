// examples/react-integration/src/components/ProductList.jsx
import React, { useState, useEffect } from 'react';
import apiService from '../services/api';

const ProductList = ({ categoryId, featured = false, limit = 12 }) => {
  const [products, setProducts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [pagination, setPagination] = useState(null);
  const [currentPage, setCurrentPage] = useState(1);

  useEffect(() => {
    fetchProducts();
  }, [categoryId, featured, currentPage, limit]);

  const fetchProducts = async () => {
    try {
      setLoading(true);
      setError(null);

      const params = {
        page: currentPage,
        limit,
      };

      if (categoryId) params.category = categoryId;
      if (featured) params.featured = true;

      const response = await apiService.getProducts(params);
      
      setProducts(response.products || []);
      setPagination(response.pagination);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handlePageChange = (page) => {
    setCurrentPage(page);
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center py-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
        <span className="ml-2">Загрузка товаров...</span>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
        <strong>Ошибка:</strong> {error}
        <button 
          onClick={fetchProducts}
          className="ml-4 text-sm underline hover:no-underline"
        >
          Повторить
        </button>
      </div>
    );
  }

  if (products.length === 0) {
    return (
      <div className="text-center py-8 text-gray-500">
        Товары не найдены
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Список товаров */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        {products.map((product) => (
          <ProductCard key={product.id} product={product} />
        ))}
      </div>

      {/* Пагинация */}
      {pagination && pagination.total_pages > 1 && (
        <Pagination 
          pagination={pagination}
          onPageChange={handlePageChange}
        />
      )}
    </div>
  );
};

const ProductCard = ({ product }) => {
  return (
    <div className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow">
      {/* Изображение товара */}
      <div className="h-48 bg-gray-200 flex items-center justify-center">
        {product.image ? (
          <img 
            src={product.image} 
            alt={product.name}
            className="max-h-full max-w-full object-cover"
          />
        ) : (
          <span className="text-gray-400">Нет фото</span>
        )}
      </div>

      {/* Информация о товаре */}
      <div className="p-4">
        {/* Категория */}
        {product.category && (
          <span className="text-xs text-gray-500 uppercase tracking-wide">
            {product.category.name}
          </span>
        )}

        {/* Название */}
        <h3 className="font-semibold text-gray-900 mb-2 line-clamp-2">
          {product.name}
        </h3>

        {/* Бренд и модель */}
        <p className="text-sm text-gray-600 mb-2">
          {product.brand} {product.model}
        </p>

        {/* Цена */}
        <div className="flex items-center space-x-2 mb-3">
          <span className="text-lg font-bold text-blue-600">
            ${product.price}
          </span>
          {product.old_price && product.old_price > product.price && (
            <span className="text-sm text-gray-500 line-through">
              ${product.old_price}
            </span>
          )}
          {product.is_featured && (
            <span className="bg-red-100 text-red-800 text-xs px-2 py-1 rounded">
              ХИТ
            </span>
          )}
        </div>

        {/* Наличие */}
        <div className="mb-3">
          {product.stock > 0 ? (
            <span className="text-green-600 text-sm">
              В наличии ({product.stock} шт.)
            </span>
          ) : (
            <span className="text-red-600 text-sm">Нет в наличии</span>
          )}
        </div>

        {/* Кнопки действий */}
        <div className="flex space-x-2">
          <button className="flex-1 bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 transition-colors">
            Подробнее
          </button>
          <button className="bg-gray-200 text-gray-700 py-2 px-4 rounded hover:bg-gray-300 transition-colors">
            ♡
          </button>
        </div>
      </div>
    </div>
  );
};

const Pagination = ({ pagination, onPageChange }) => {
  const { page: currentPage, total_pages: totalPages } = pagination;

  const getPageNumbers = () => {
    const pages = [];
    const maxVisiblePages = 5;
    
    let start = Math.max(1, currentPage - Math.floor(maxVisiblePages / 2));
    let end = Math.min(totalPages, start + maxVisiblePages - 1);
    
    if (end - start + 1 < maxVisiblePages) {
      start = Math.max(1, end - maxVisiblePages + 1);
    }
    
    for (let i = start; i <= end; i++) {
      pages.push(i);
    }
    
    return pages;
  };

  return (
    <div className="flex justify-center items-center space-x-2">
      {/* Предыдущая страница */}
      <button
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage === 1}
        className="px-3 py-2 text-sm bg-white border border-gray-300 rounded hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        Назад
      </button>

      {/* Номера страниц */}
      {getPageNumbers().map(pageNumber => (
        <button
          key={pageNumber}
          onClick={() => onPageChange(pageNumber)}
          className={`px-3 py-2 text-sm border rounded ${
            pageNumber === currentPage
              ? 'bg-blue-500 text-white border-blue-500'
              : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-50'
          }`}
        >
          {pageNumber}
        </button>
      ))}

      {/* Следующая страница */}
      <button
        onClick={() => onPageChange(currentPage + 1)}
        disabled={currentPage === totalPages}
        className="px-3 py-2 text-sm bg-white border border-gray-300 rounded hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        Вперед
      </button>
    </div>
  );
};

export default ProductList;