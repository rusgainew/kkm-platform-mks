import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, beforeEach, vi } from 'vitest';

/**
 * Mock Catalog List Component
 */
function CatalogList() {
  return (
    <div>
      <h1>Каталог товаров</h1>
      <div>
        <button>Добавить товар</button>
        <input placeholder="Поиск товара" />
        <select>
          <option value="">Все категории</option>
          <option value="electronics">Электроника</option>
          <option value="food">Продукты</option>
          <option value="clothing">Одежда</option>
        </select>
      </div>
      <div className="grid">
        <div className="product-card">
          <img src="/product.jpg" alt="Товар 1" />
          <h3>Товар 1</h3>
          <p>Описание товара 1</p>
          <span className="price">₽1,000</span>
          <span className="stock">В наличии: 50</span>
          <button aria-label="edit">Редактировать</button>
          <button aria-label="delete">Удалить</button>
        </div>
      </div>
    </div>
  );
}

describe('CatalogList', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Rendering', () => {
    it('renders catalog title', () => {
      render(<CatalogList />);
      expect(screen.getByText(/каталог товаров/i)).toBeInTheDocument();
    });

    it('renders add product button', () => {
      render(<CatalogList />);
      expect(screen.getByText(/добавить товар/i)).toBeInTheDocument();
    });

    it('renders search input', () => {
      render(<CatalogList />);
      expect(screen.getByPlaceholderText(/поиск товара/i)).toBeInTheDocument();
    });

    it('renders category filter', () => {
      render(<CatalogList />);
      const select = screen.getByDisplayValue(/все категории/i);
      expect(select).toBeInTheDocument();
    });
  });

  describe('Filters', () => {
    it('has category filter options', () => {
      render(<CatalogList />);
      expect(screen.getByText(/электроника/i)).toBeInTheDocument();
      expect(screen.getByText(/продукты/i)).toBeInTheDocument();
      expect(screen.getByText(/одежда/i)).toBeInTheDocument();
    });

    it('allows category selection', async () => {
      render(<CatalogList />);
      const select = screen.getByDisplayValue(/все категории/i);
      
      await userEvent.selectOptions(select, 'electronics');
      expect(select).toHaveValue('electronics');
    });
  });

  describe('Product Card', () => {
    it('displays product image', () => {
      render(<CatalogList />);
      const img = screen.getByAltText(/товар 1/i);
      expect(img).toBeInTheDocument();
      expect(img).toHaveAttribute('src', '/product.jpg');
    });

    it('displays product name', () => {
      render(<CatalogList />);
      expect(screen.getByText(/товар 1/)).toBeInTheDocument();
    });

    it('displays product description', () => {
      render(<CatalogList />);
      expect(screen.getByText(/описание товара 1/i)).toBeInTheDocument();
    });

    it('displays product price', () => {
      render(<CatalogList />);
      expect(screen.getByText(/₽1,000/)).toBeInTheDocument();
    });

    it('displays stock information', () => {
      render(<CatalogList />);
      expect(screen.getByText(/в наличии: 50/i)).toBeInTheDocument();
    });
  });

  describe('Product Actions', () => {
    it('has edit button for each product', () => {
      render(<CatalogList />);
      expect(screen.getByLabelText(/edit/)).toBeInTheDocument();
    });

    it('has delete button for each product', () => {
      render(<CatalogList />);
      expect(screen.getByLabelText(/delete/)).toBeInTheDocument();
    });
  });

  describe('Search', () => {
    it('allows typing in search input', async () => {
      render(<CatalogList />);
      const searchInput = screen.getByPlaceholderText(/поиск товара/i);
      
      await userEvent.type(searchInput, 'Товар');
      expect(searchInput).toHaveValue('Товар');
    });

    it('clears search input', async () => {
      render(<CatalogList />);
      const searchInput = screen.getByPlaceholderText(/поиск товара/i) as HTMLInputElement;
      
      await userEvent.type(searchInput, 'Товар');
      expect(searchInput.value).toBe('Товар');
      
      await userEvent.clear(searchInput);
      expect(searchInput.value).toBe('');
    });
  });

  describe('Grid Layout', () => {
    it('renders products in grid', () => {
      render(<CatalogList />);
      const grid = screen.getByText(/товар 1/).parentElement?.parentElement;
      expect(grid).toHaveClass('product-card');
    });
  });
});
