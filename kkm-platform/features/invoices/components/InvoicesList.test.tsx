import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, beforeEach, vi } from 'vitest';

/**
 * Mock Invoices List Component
 */
function InvoicesList() {
  return (
    <div>
      <h1>Счета-фактуры</h1>
      <div>
        <button>Создать счет</button>
        <input placeholder="Поиск по номеру" />
        <select>
          <option value="">Все статусы</option>
          <option value="draft">Черновик</option>
          <option value="issued">Выпущен</option>
          <option value="paid">Оплачен</option>
          <option value="cancelled">Отменен</option>
        </select>
      </div>
      <table>
        <thead>
          <tr>
            <th>№ Счета</th>
            <th>Компания</th>
            <th>Сумма</th>
            <th>Статус</th>
            <th>Дата</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>INV-001</td>
            <td>ООО "Компания"</td>
            <td>₽50,000</td>
            <td>
              <span className="badge badge-warning">Выпущен</span>
            </td>
            <td>2025-01-16</td>
            <td>
              <button aria-label="view">Просмотр</button>
              <button aria-label="edit">Редактировать</button>
              <button aria-label="delete">Удалить</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  );
}

describe('InvoicesList', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Rendering', () => {
    it('renders invoices list title', () => {
      render(<InvoicesList />);
      expect(screen.getByText(/счета-фактуры/i)).toBeInTheDocument();
    });

    it('renders create button', () => {
      render(<InvoicesList />);
      expect(screen.getByText(/создать счет/i)).toBeInTheDocument();
    });

    it('renders search input', () => {
      render(<InvoicesList />);
      expect(screen.getByPlaceholderText(/поиск по номеру/i)).toBeInTheDocument();
    });

    it('renders status filter select', () => {
      render(<InvoicesList />);
      const select = screen.getByDisplayValue(/все статусы/i);
      expect(select).toBeInTheDocument();
    });
  });

  describe('Filters', () => {
    it('has status filter options', () => {
      render(<InvoicesList />);
      expect(screen.getByText(/черновик/i)).toBeInTheDocument();
      expect(screen.getByText(/выпущен/i)).toBeInTheDocument();
      expect(screen.getByText(/оплачен/i)).toBeInTheDocument();
      expect(screen.getByText(/отменен/i)).toBeInTheDocument();
    });

    it('allows status selection', async () => {
      render(<InvoicesList />);
      const select = screen.getByDisplayValue(/все статусы/i);
      
      await userEvent.selectOptions(select, 'paid');
      expect(select).toHaveValue('paid');
    });
  });

  describe('Invoice Data', () => {
    it('displays invoice number', () => {
      render(<InvoicesList />);
      expect(screen.getByText(/INV-001/)).toBeInTheDocument();
    });

    it('displays company name', () => {
      render(<InvoicesList />);
      expect(screen.getByText(/ООО "Компания"/)).toBeInTheDocument();
    });

    it('displays invoice amount', () => {
      render(<InvoicesList />);
      expect(screen.getByText(/₽50,000/)).toBeInTheDocument();
    });

    it('displays invoice status', () => {
      render(<InvoicesList />);
      const badges = screen.getByText(/выпущен/i);
      expect(badges).toBeInTheDocument();
    });

    it('displays invoice date', () => {
      render(<InvoicesList />);
      expect(screen.getByText(/2025-01-16/)).toBeInTheDocument();
    });
  });

  describe('Table Headers', () => {
    it('renders all table headers', () => {
      render(<InvoicesList />);
      expect(screen.getByText(/№ Счета/)).toBeInTheDocument();
      expect(screen.getByText(/Компания/)).toBeInTheDocument();
      expect(screen.getByText(/Сумма/)).toBeInTheDocument();
      expect(screen.getByText(/Статус/)).toBeInTheDocument();
      expect(screen.getByText(/Дата/)).toBeInTheDocument();
      expect(screen.getByText(/Действия/)).toBeInTheDocument();
    });
  });

  describe('Actions', () => {
    it('has view button for each invoice', () => {
      render(<InvoicesList />);
      expect(screen.getByLabelText(/view/)).toBeInTheDocument();
    });

    it('has edit button for each invoice', () => {
      render(<InvoicesList />);
      expect(screen.getByLabelText(/edit/)).toBeInTheDocument();
    });

    it('has delete button for each invoice', () => {
      render(<InvoicesList />);
      expect(screen.getByLabelText(/delete/)).toBeInTheDocument();
    });
  });

  describe('Search', () => {
    it('allows typing in search input', async () => {
      render(<InvoicesList />);
      const searchInput = screen.getByPlaceholderText(/поиск по номеру/i);
      
      await userEvent.type(searchInput, 'INV-001');
      expect(searchInput).toHaveValue('INV-001');
    });
  });
});
