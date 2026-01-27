/**
 * Unit tests for BankAccountForm component
 * Testing form validation, data input, and submission logic
 *
 * @vitest
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import BankAccountForm from './BankAccountForm';
import type { BankAccount } from '@/types/entities';

// Mock API functions
vi.mock('@/lib/api/bank-accounts', () => ({
  createBankAccount: vi.fn(),
  updateBankAccount: vi.fn(),
}));

const createTestQueryClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });

const renderWithQueryClient = (component: React.ReactElement) => {
  const queryClient = createTestQueryClient();
  return render(
    <QueryClientProvider client={queryClient}>{component}</QueryClientProvider>
  );
};

describe('BankAccountForm', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Rendering', () => {
    it('should render create mode form', () => {
      renderWithQueryClient(<BankAccountForm ownerId="owner-123" />);

      expect(
        screen.getByText('Добавить банковский счет')
      ).toBeInTheDocument();
      expect(screen.getByLabelText(/Номер счета/)).toBeInTheDocument();
      expect(screen.getByLabelText(/Название банка/)).toBeInTheDocument();
      expect(screen.getByLabelText(/БИК банка/)).toBeInTheDocument();
      expect(screen.getByLabelText(/Валюта счета/)).toBeInTheDocument();
      expect(screen.getByText('Создать счет')).toBeInTheDocument();
    });

    it('should render edit mode form with initial data', () => {
      const initialData: BankAccount = {
        id: 'acc-123',
        account_number: '12345678901234567890',
        bank_name: 'ПАО Сбербанк',
        bank_code: '044525225',
        currency: 'RUB',
        owner_id: 'owner-123',
        is_active: true,
        created_at: Date.now(),
        updated_at: Date.now(),
      };

      renderWithQueryClient(<BankAccountForm initialData={initialData} />);

      expect(
        screen.getByText('Редактировать банковский счет')
      ).toBeInTheDocument();
      expect(screen.getByDisplayValue('12345678901234567890')).toBeInTheDocument();
      expect(screen.getByDisplayValue('ПАО Сбербанк')).toBeInTheDocument();
      expect(screen.getByDisplayValue('044525225')).toBeInTheDocument();
      expect(screen.getByText('Обновить счет')).toBeInTheDocument();
    });

    it('should render all currency options', () => {
      renderWithQueryClient(<BankAccountForm ownerId="owner-123" />);

      const currencySelect = screen.getByLabelText(/Валюта счета/) as HTMLSelectElement;
      const options = Array.from(currencySelect.options).map((opt) => opt.value);

      expect(options).toContain('KGS');
      expect(options).toContain('USD');
      expect(options).toContain('RUB');
      expect(options).toContain('EUR');
      expect(options).toContain('CNY');
    });
  });

  describe('Validation', () => {
    it('should show error for empty account number', async () => {
      renderWithQueryClient(<BankAccountForm ownerId="owner-123" />);

      const submitButton = screen.getByText('Создать счет');
      fireEvent.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText('Номер счета обязателен')).toBeInTheDocument();
      });
    });

    it('should show error for invalid account number length', async () => {
      renderWithQueryClient(<BankAccountForm ownerId="owner-123" />);

      const accountInput = screen.getByLabelText(/Номер счета/);
      fireEvent.change(accountInput, { target: { value: '123456789' } });
      fireEvent.blur(accountInput);

      await waitFor(() => {
        expect(
          screen.getByText('Номер счета должен содержать ровно 20 цифр')
        ).toBeInTheDocument();
      });
    });

    it('should accept valid 20-digit account number', async () => {
      renderWithQueryClient(<BankAccountForm ownerId="owner-123" />);

      const accountInput = screen.getByLabelText(/Номер счета/);
      fireEvent.change(accountInput, {
        target: { value: '12345678901234567890' },
      });
      fireEvent.blur(accountInput);

      await waitFor(() => {
        expect(
          screen.queryByText('Номер счета должен содержать ровно 20 цифр')
        ).not.toBeInTheDocument();
      });
    });

    it('should remove non-numeric characters from account number', () => {
      renderWithQueryClient(<BankAccountForm ownerId="owner-123" />);

      const accountInput = screen.getByLabelText(/Номер счета/) as HTMLInputElement;
      fireEvent.change(accountInput, {
        target: { value: '1234-5678-9012-3456-7890' },
      });

      expect(accountInput.value).toBe('12345678901234567890');
    });

    it('should show error for empty bank name', async () => {
      renderWithQueryClient(<BankAccountForm ownerId="owner-123" />);

      const submitButton = screen.getByText('Создать счет');
      fireEvent.click(submitButton);

      await waitFor(() => {
        expect(
          screen.getByText('Название банка обязательно')
        ).toBeInTheDocument();
      });
    });

    it('should show error for short bank name', async () => {
      renderWithQueryClient(<BankAccountForm ownerId="owner-123" />);

      const bankNameInput = screen.getByLabelText(/Название банка/);
      fireEvent.change(bankNameInput, { target: { value: 'AB' } });
      fireEvent.blur(bankNameInput);

      await waitFor(() => {
        expect(
          screen.getByText('Название банка должно содержать минимум 3 символа')
        ).toBeInTheDocument();
      });
    });

    it('should show error for empty bank code', async () => {
      renderWithQueryClient(<BankAccountForm ownerId="owner-123" />);

      const submitButton = screen.getByText('Создать счет');
      fireEvent.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText('БИК банка обязателен')).toBeInTheDocument();
      });
    });

    it('should show error for invalid bank code length', async () => {
      renderWithQueryClient(<BankAccountForm ownerId="owner-123" />);

      const bankCodeInput = screen.getByLabelText(/БИК банка/);
      fireEvent.change(bankCodeInput, { target: { value: '12345' } });
      fireEvent.blur(bankCodeInput);

      await waitFor(() => {
        expect(
          screen.getByText('БИК должен содержать от 6 до 9 цифр')
        ).toBeInTheDocument();
      });
    });

    it('should accept valid 9-digit bank code', async () => {
      renderWithQueryClient(<BankAccountForm ownerId="owner-123" />);

      const bankCodeInput = screen.getByLabelText(/БИК банка/);
      fireEvent.change(bankCodeInput, { target: { value: '044525225' } });
      fireEvent.blur(bankCodeInput);

      await waitFor(() => {
        expect(
          screen.queryByText('БИК должен содержать от 6 до 9 цифр')
        ).not.toBeInTheDocument();
      });
    });

    it('should remove non-numeric characters from bank code', () => {
      renderWithQueryClient(<BankAccountForm ownerId="owner-123" />);

      const bankCodeInput = screen.getByLabelText(/БИК банка/) as HTMLInputElement;
      fireEvent.change(bankCodeInput, { target: { value: '044-525-225' } });

      expect(bankCodeInput.value).toBe('044525225');
    });

    it('should limit bank code to 9 digits', () => {
      renderWithQueryClient(<BankAccountForm ownerId="owner-123" />);

      const bankCodeInput = screen.getByLabelText(/БИК банка/) as HTMLInputElement;
      fireEvent.change(bankCodeInput, { target: { value: '12345678901234' } });

      expect(bankCodeInput.value).toBe('123456789');
    });
  });

  describe('Form interaction', () => {
    it('should update form fields on input', () => {
      renderWithQueryClient(<BankAccountForm ownerId="owner-123" />);

      const accountInput = screen.getByLabelText(/Номер счета/) as HTMLInputElement;
      const bankNameInput = screen.getByLabelText(/Название банка/) as HTMLInputElement;
      const bankCodeInput = screen.getByLabelText(/БИК банка/) as HTMLInputElement;
      const currencySelect = screen.getByLabelText(/Валюта счета/) as HTMLSelectElement;

      fireEvent.change(accountInput, {
        target: { value: '12345678901234567890' },
      });
      fireEvent.change(bankNameInput, { target: { value: 'Тест Банк' } });
      fireEvent.change(bankCodeInput, { target: { value: '044525225' } });
      fireEvent.change(currencySelect, { target: { value: 'USD' } });

      expect(accountInput.value).toBe('12345678901234567890');
      expect(bankNameInput.value).toBe('Тест Банк');
      expect(bankCodeInput.value).toBe('044525225');
      expect(currencySelect.value).toBe('USD');
    });

    it('should toggle active status checkbox', () => {
      renderWithQueryClient(<BankAccountForm ownerId="owner-123" />);

      const checkbox = screen.getByRole('checkbox') as HTMLInputElement;
      expect(checkbox.checked).toBe(true); // default is active

      fireEvent.click(checkbox);
      expect(checkbox.checked).toBe(false);

      fireEvent.click(checkbox);
      expect(checkbox.checked).toBe(true);
    });

    it('should call onCancel when cancel button is clicked', () => {
      const onCancel = vi.fn();
      renderWithQueryClient(
        <BankAccountForm ownerId="owner-123" onCancel={onCancel} />
      );

      const cancelButton = screen.getByText('Отмена');
      fireEvent.click(cancelButton);

      expect(onCancel).toHaveBeenCalledTimes(1);
    });

    it('should disable submit button when form has errors', async () => {
      renderWithQueryClient(<BankAccountForm ownerId="owner-123" />);

      const submitButton = screen.getByText('Создать счет') as HTMLButtonElement;
      fireEvent.click(submitButton);

      await waitFor(() => {
        expect(submitButton.disabled).toBe(true);
      });
    });

    it('should show character counter for account number', () => {
      renderWithQueryClient(<BankAccountForm ownerId="owner-123" />);

      const accountInput = screen.getByLabelText(/Номер счета/);
      fireEvent.change(accountInput, { target: { value: '12345' } });

      expect(screen.getByText(/Текущая длина: 5/)).toBeInTheDocument();
    });
  });

  describe('Information display', () => {
    it('should display informational message about requirements', () => {
      renderWithQueryClient(<BankAccountForm ownerId="owner-123" />);

      expect(
        screen.getByText(/Номер счета должен содержать ровно 20 цифр/)
      ).toBeInTheDocument();
      expect(
        screen.getByText(/БИК банка обычно содержит 9 цифр/)
      ).toBeInTheDocument();
      expect(
        screen.getByText(/Все реквизиты должны быть проверены перед сохранением/)
      ).toBeInTheDocument();
    });

    it('should display hint about bank code format', () => {
      renderWithQueryClient(<BankAccountForm ownerId="owner-123" />);

      expect(
        screen.getByText(/Обычно 9 цифр для РФ, 6-9 цифр для других стран/)
      ).toBeInTheDocument();
    });
  });
});
