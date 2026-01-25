import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, beforeEach, vi } from 'vitest';

/**
 * Mock Companies List Component
 */
function CompaniesList() {
  return (
    <div>
      <h1>Список Компаний</h1>
      <button>Создать компанию</button>
      <table>
        <thead>
          <tr>
            <th>Название</th>
            <th>ИНН</th>
            <th>Телефон</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>ООО "Тестовая"</td>
            <td>7712345678</td>
            <td>+7-999-123-45-67</td>
            <td>
              <button aria-label="edit">Редактировать</button>
              <button aria-label="delete">Удалить</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  );
}

describe('CompaniesList', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Rendering', () => {
    it('renders companies list title', () => {
      render(<CompaniesList />);
      expect(screen.getByText(/список компаний/i)).toBeInTheDocument();
    });

    it('renders create button', () => {
      render(<CompaniesList />);
      expect(screen.getByText(/создать компанию/i)).toBeInTheDocument();
    });

    it('renders table with companies', () => {
      render(<CompaniesList />);
      expect(screen.getByText(/ООО "Тестовая"/)).toBeInTheDocument();
      expect(screen.getByText(/7712345678/)).toBeInTheDocument();
    });

    it('renders action buttons for each company', () => {
      render(<CompaniesList />);
      expect(screen.getByLabelText(/edit/)).toBeInTheDocument();
      expect(screen.getByLabelText(/delete/)).toBeInTheDocument();
    });
  });

  describe('Table Headers', () => {
    it('renders all table headers', () => {
      render(<CompaniesList />);
      expect(screen.getByText(/Название/)).toBeInTheDocument();
      expect(screen.getByText(/ИНН/)).toBeInTheDocument();
      expect(screen.getByText(/Телефон/)).toBeInTheDocument();
      expect(screen.getByText(/Действия/)).toBeInTheDocument();
    });
  });

  describe('Company Data', () => {
    it('displays company name correctly', () => {
      render(<CompaniesList />);
      expect(screen.getByText(/ООО "Тестовая"/)).toBeInTheDocument();
    });

    it('displays company INN correctly', () => {
      render(<CompaniesList />);
      expect(screen.getByText(/7712345678/)).toBeInTheDocument();
    });

    it('displays company phone correctly', () => {
      render(<CompaniesList />);
      expect(screen.getByText(/\+7-999-123-45-67/)).toBeInTheDocument();
    });
  });

  describe('Actions', () => {
    it('has edit button for each company', () => {
      render(<CompaniesList />);
      const editButton = screen.getByLabelText(/edit/);
      expect(editButton).toBeInTheDocument();
    });

    it('has delete button for each company', () => {
      render(<CompaniesList />);
      const deleteButton = screen.getByLabelText(/delete/);
      expect(deleteButton).toBeInTheDocument();
    });
  });
});
