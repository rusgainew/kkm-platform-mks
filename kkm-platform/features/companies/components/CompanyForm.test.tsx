import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { useRouter } from 'next/navigation';
import CompanyForm from './CompanyForm';
import type { Company } from '@/types/entities';

// Mock next/navigation
vi.mock('next/navigation', () => ({
  useRouter: vi.fn(),
}));

describe('CompanyForm', () => {
  const mockPush = vi.fn();
  const mockBack = vi.fn();
  const mockOnSubmit = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    (useRouter as ReturnType<typeof vi.fn>).mockReturnValue({
      push: mockPush,
      back: mockBack,
    });
  });

  describe('Режим создания', () => {
    it('отображает форму создания компании', () => {
      render(<CompanyForm onSubmit={mockOnSubmit} mode="create" />);
      
      expect(screen.getByText('Создание компании')).toBeInTheDocument();
      expect(screen.getByLabelText(/Название компании/)).toBeInTheDocument();
      expect(screen.getByLabelText(/ИНН/)).toBeInTheDocument();
      expect(screen.getByLabelText(/Юридический адрес/)).toBeInTheDocument();
      expect(screen.getByLabelText(/Телефон/)).toBeInTheDocument();
      expect(screen.getByLabelText(/Email/)).toBeInTheDocument();
    });

    it('показывает ошибки валидации при пустой форме', async () => {
      render(<CompanyForm onSubmit={mockOnSubmit} mode="create" />);
      
      const submitButton = screen.getByText('Создать компанию');
      fireEvent.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText('Название компании обязательно')).toBeInTheDocument();
        expect(screen.getByText('ИНН обязателен')).toBeInTheDocument();
        expect(screen.getByText('Адрес обязателен')).toBeInTheDocument();
        expect(screen.getByText('Телефон обязателен')).toBeInTheDocument();
        expect(screen.getByText('Email обязателен')).toBeInTheDocument();
      });

      expect(mockOnSubmit).not.toHaveBeenCalled();
    });

    it('валидирует ИНН (должен быть 10 или 12 цифр)', async () => {
      render(<CompanyForm onSubmit={mockOnSubmit} mode="create" />);
      
      const tinInput = screen.getByLabelText(/ИНН/) as HTMLInputElement;
      
      // Невалидный ИНН (9 цифр)
      fireEvent.change(tinInput, { target: { value: '123456789' } });
      fireEvent.click(screen.getByText('Создать компанию'));
      
      await waitFor(() => {
        expect(screen.getByText(/ИНН должен содержать 10/)).toBeInTheDocument();
      });
      
      // Валидный ИНН (10 цифр)
      fireEvent.change(tinInput, { target: { value: '1234567890' } });
      await waitFor(() => {
        expect(screen.queryByText(/ИНН должен содержать 10/)).not.toBeInTheDocument();
      });
    });

    it('валидирует email', async () => {
      render(<CompanyForm onSubmit={mockOnSubmit} mode="create" />);
      
      const emailInput = screen.getByLabelText(/Email/) as HTMLInputElement;
      
      // Невалидный email
      fireEvent.change(emailInput, { target: { value: 'invalid-email' } });
      fireEvent.click(screen.getByText('Создать компанию'));
      
      await waitFor(() => {
        expect(screen.getByText('Введите корректный email')).toBeInTheDocument();
      });
      
      // Валидный email
      fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
      await waitFor(() => {
        expect(screen.queryByText('Введите корректный email')).not.toBeInTheDocument();
      });
    });

    it('валидирует телефон', async () => {
      render(<CompanyForm onSubmit={mockOnSubmit} mode="create" />);
      
      const phoneInput = screen.getByLabelText(/Телефон/) as HTMLInputElement;
      
      // Невалидный телефон
      fireEvent.change(phoneInput, { target: { value: '123' } });
      fireEvent.click(screen.getByText('Создать компанию'));
      
      await waitFor(() => {
        expect(screen.getByText(/Введите корректный номер телефона/)).toBeInTheDocument();
      });
      
      // Валидный телефон
      fireEvent.change(phoneInput, { target: { value: '+79991234567' } });
      await waitFor(() => {
        expect(screen.queryByText(/Введите корректный номер телефона/)).not.toBeInTheDocument();
      });
    });

    it('успешно создает компанию с валидными данными', async () => {
      const mockCompany: Company = {
        company_id: '123',
        name: 'Test Company',
        tin: '1234567890',
        address: 'Test Address 123',
        phone: '+79991234567',
        email: 'test@example.com',
        owner_id: 'user1',
        status: 'active',
        created_at: Date.now(),
        updated_at: Date.now(),
      };

      mockOnSubmit.mockResolvedValue(mockCompany);

      render(<CompanyForm onSubmit={mockOnSubmit} mode="create" />);
      
      // Заполнение формы
      fireEvent.change(screen.getByLabelText(/Название компании/), {
        target: { value: 'Test Company' },
      });
      fireEvent.change(screen.getByLabelText(/ИНН/), {
        target: { value: '1234567890' },
      });
      fireEvent.change(screen.getByLabelText(/Юридический адрес/), {
        target: { value: 'Test Address 123' },
      });
      fireEvent.change(screen.getByLabelText(/Телефон/), {
        target: { value: '+79991234567' },
      });
      fireEvent.change(screen.getByLabelText(/Email/), {
        target: { value: 'test@example.com' },
      });

      // Отправка формы
      fireEvent.click(screen.getByText('Создать компанию'));

      await waitFor(() => {
        expect(mockOnSubmit).toHaveBeenCalledWith({
          name: 'Test Company',
          tin: '1234567890',
          address: 'Test Address 123',
          phone: '+79991234567',
          email: 'test@example.com',
        });
      });

      // Проверка сообщения об успехе
      await waitFor(() => {
        expect(screen.getByText('Компания успешно создана')).toBeInTheDocument();
      });

      // Проверка редиректа
      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith('/companies');
      }, { timeout: 2000 });
    });

    it('обрабатывает ошибку при создании', async () => {
      const errorMessage = 'Компания с таким ИНН уже существует';
      mockOnSubmit.mockRejectedValue(new Error(errorMessage));

      render(<CompanyForm onSubmit={mockOnSubmit} mode="create" />);
      
      // Заполнение формы валидными данными
      fireEvent.change(screen.getByLabelText(/Название компании/), {
        target: { value: 'Test Company' },
      });
      fireEvent.change(screen.getByLabelText(/ИНН/), {
        target: { value: '1234567890' },
      });
      fireEvent.change(screen.getByLabelText(/Юридический адрес/), {
        target: { value: 'Test Address 123' },
      });
      fireEvent.change(screen.getByLabelText(/Телефон/), {
        target: { value: '+79991234567' },
      });
      fireEvent.change(screen.getByLabelText(/Email/), {
        target: { value: 'test@example.com' },
      });

      // Отправка формы
      fireEvent.click(screen.getByText('Создать компанию'));

      await waitFor(() => {
        expect(screen.getByText(errorMessage)).toBeInTheDocument();
      });

      expect(mockPush).not.toHaveBeenCalled();
    });
  });

  describe('Режим редактирования', () => {
    const initialData: Company = {
      company_id: '123',
      name: 'Existing Company',
      tin: '1234567890',
      kpp: '123456789',
      ogrn: '1234567890123',
      address: 'Existing Address',
      phone: '+79991234567',
      email: 'existing@example.com',
      website: 'https://example.com',
      owner_id: 'user1',
      status: 'active',
      created_at: Date.now(),
      updated_at: Date.now(),
    };

    it('отображает форму редактирования с начальными данными', () => {
      render(<CompanyForm onSubmit={mockOnSubmit} mode="edit" initialData={initialData} />);
      
      expect(screen.getByText('Редактирование компании')).toBeInTheDocument();
      expect(screen.getByDisplayValue('Existing Company')).toBeInTheDocument();
      expect(screen.getByDisplayValue('1234567890')).toBeInTheDocument();
      expect(screen.getByDisplayValue('123456789')).toBeInTheDocument();
      expect(screen.getByDisplayValue('Existing Address')).toBeInTheDocument();
      expect(screen.getByDisplayValue('+79991234567')).toBeInTheDocument();
      expect(screen.getByDisplayValue('existing@example.com')).toBeInTheDocument();
    });

    it('отображает поле статуса в режиме редактирования', () => {
      render(<CompanyForm onSubmit={mockOnSubmit} mode="edit" initialData={initialData} />);
      
      expect(screen.getByText('Статус компании')).toBeInTheDocument();
      const statusSelect = screen.getByLabelText(/Статус/) as HTMLSelectElement;
      expect(statusSelect).toBeInTheDocument();
      expect(statusSelect.value).toBe('active');
    });

    it('успешно обновляет компанию', async () => {
      const updatedCompany: Company = {
        ...initialData,
        name: 'Updated Company',
      };

      mockOnSubmit.mockResolvedValue(updatedCompany);

      render(<CompanyForm onSubmit={mockOnSubmit} mode="edit" initialData={initialData} />);
      
      // Изменение названия
      const nameInput = screen.getByLabelText(/Название компании/) as HTMLInputElement;
      fireEvent.change(nameInput, { target: { value: 'Updated Company' } });

      // Отправка формы
      fireEvent.click(screen.getByText('Сохранить изменения'));

      await waitFor(() => {
        expect(mockOnSubmit).toHaveBeenCalled();
      });

      // Проверка сообщения об успехе
      await waitFor(() => {
        expect(screen.getByText('Компания успешно обновлена')).toBeInTheDocument();
      });
    });
  });

  describe('Навигация', () => {
    it('возвращается назад при клике на кнопку "Вернуться"', () => {
      render(<CompanyForm onSubmit={mockOnSubmit} />);
      
      const backButton = screen.getByText('Вернуться');
      fireEvent.click(backButton);

      expect(mockBack).toHaveBeenCalled();
    });

    it('переходит к списку компаний при клике на "Отмена"', () => {
      render(<CompanyForm onSubmit={mockOnSubmit} />);
      
      const cancelButton = screen.getByText('Отмена');
      fireEvent.click(cancelButton);

      expect(mockPush).toHaveBeenCalledWith('/companies');
    });
  });

  describe('Опциональные поля', () => {
    it('не включает опциональные поля если они пустые', async () => {
      mockOnSubmit.mockResolvedValue({} as Company);

      render(<CompanyForm onSubmit={mockOnSubmit} mode="create" />);
      
      // Заполнение только обязательных полей
      fireEvent.change(screen.getByLabelText(/Название компании/), {
        target: { value: 'Test Company' },
      });
      fireEvent.change(screen.getByLabelText(/ИНН/), {
        target: { value: '1234567890' },
      });
      fireEvent.change(screen.getByLabelText(/Юридический адрес/), {
        target: { value: 'Test Address 123' },
      });
      fireEvent.change(screen.getByLabelText(/Телефон/), {
        target: { value: '+79991234567' },
      });
      fireEvent.change(screen.getByLabelText(/Email/), {
        target: { value: 'test@example.com' },
      });

      // Отправка формы
      fireEvent.click(screen.getByText('Создать компанию'));

      await waitFor(() => {
        const callArg = mockOnSubmit.mock.calls[0][0];
        expect(callArg).toEqual({
          name: 'Test Company',
          tin: '1234567890',
          address: 'Test Address 123',
          phone: '+79991234567',
          email: 'test@example.com',
        });
        // КПП, ОГРН, Website не должны присутствовать
        expect(callArg).not.toHaveProperty('kpp');
        expect(callArg).not.toHaveProperty('ogrn');
        expect(callArg).not.toHaveProperty('website');
      });
    });

    it('включает опциональные поля если они заполнены', async () => {
      mockOnSubmit.mockResolvedValue({} as Company);

      render(<CompanyForm onSubmit={mockOnSubmit} mode="create" />);
      
      // Заполнение всех полей
      fireEvent.change(screen.getByLabelText(/Название компании/), {
        target: { value: 'Test Company' },
      });
      fireEvent.change(screen.getByLabelText(/ИНН/), {
        target: { value: '1234567890' },
      });
      fireEvent.change(screen.getByLabelText(/КПП/), {
        target: { value: '123456789' },
      });
      fireEvent.change(screen.getByLabelText(/ОГРН/), {
        target: { value: '1234567890123' },
      });
      fireEvent.change(screen.getByLabelText(/Юридический адрес/), {
        target: { value: 'Test Address 123' },
      });
      fireEvent.change(screen.getByLabelText(/Телефон/), {
        target: { value: '+79991234567' },
      });
      fireEvent.change(screen.getByLabelText(/Email/), {
        target: { value: 'test@example.com' },
      });
      fireEvent.change(screen.getByLabelText(/Веб-сайт/), {
        target: { value: 'https://example.com' },
      });

      // Отправка формы
      fireEvent.click(screen.getByText('Создать компанию'));

      await waitFor(() => {
        expect(mockOnSubmit).toHaveBeenCalledWith({
          name: 'Test Company',
          tin: '1234567890',
          kpp: '123456789',
          ogrn: '1234567890123',
          address: 'Test Address 123',
          phone: '+79991234567',
          email: 'test@example.com',
          website: 'https://example.com',
        });
      });
    });
  });
});
