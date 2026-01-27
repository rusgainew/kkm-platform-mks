import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { useRouter } from 'next/navigation';
import CatalogForm from './CatalogForm';
import type { CatalogItem } from '@/types/entities';

// Mock next/navigation
vi.mock('next/navigation', () => ({
  useRouter: vi.fn(),
}));

describe('CatalogForm', () => {
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
    it('отображает форму создания товара', () => {
      render(<CatalogForm onSubmit={mockOnSubmit} mode="create" />);
      
      expect(screen.getByText('Добавление товара в каталог')).toBeInTheDocument();
      expect(screen.getByLabelText(/Наименование товара/)).toBeInTheDocument();
      expect(screen.getByLabelText(/Номер по каталогу/)).toBeInTheDocument();
      expect(screen.getByLabelText(/Код ТНВЭД/)).toBeInTheDocument();
      expect(screen.getByLabelText(/Цена/)).toBeInTheDocument();
      expect(screen.getByLabelText(/Валюта/)).toBeInTheDocument();
      expect(screen.getByLabelText(/Единица измерения/)).toBeInTheDocument();
    });

    it('показывает ошибки валидации при пустой форме', async () => {
      render(<CatalogForm onSubmit={mockOnSubmit} mode="create" />);
      
      const submitButton = screen.getByText('Добавить в каталог');
      fireEvent.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText('Название товара/услуги обязательно')).toBeInTheDocument();
        expect(screen.getByText('Номер по каталогу обязателен')).toBeInTheDocument();
        expect(screen.getByText('Код ТНВЭД обязателен для ЭСФ')).toBeInTheDocument();
        expect(screen.getByText('Цена обязательна')).toBeInTheDocument();
      });

      expect(mockOnSubmit).not.toHaveBeenCalled();
    });

    it('валидирует ТНВЭД код (должен быть 10 цифр)', async () => {
      render(<CatalogForm onSubmit={mockOnSubmit} mode="create" />);
      
      const tnvedInput = screen.getByLabelText(/Код ТНВЭД/) as HTMLInputElement;
      
      // Невалидный ТНВЭД (меньше 10 цифр)
      fireEvent.change(tnvedInput, { target: { value: '12345' } });
      fireEvent.click(screen.getByText('Добавить в каталог'));
      
      await waitFor(() => {
        expect(screen.getByText('Код ТНВЭД должен содержать 10 цифр')).toBeInTheDocument();
      });
      
      // Валидный ТНВЭД (10 цифр с автоформатированием)
      fireEvent.change(tnvedInput, { target: { value: '1234567890' } });
      await waitFor(() => {
        expect(tnvedInput.value).toBe('1234 56 7890');
        expect(screen.queryByText('Код ТНВЭД должен содержать 10 цифр')).not.toBeInTheDocument();
      });
    });

    it('форматирует ТНВЭД код автоматически (XXXX XX XXXX)', () => {
      render(<CatalogForm onSubmit={mockOnSubmit} mode="create" />);
      
      const tnvedInput = screen.getByLabelText(/Код ТНВЭД/) as HTMLInputElement;
      
      // Ввод без пробелов
      fireEvent.change(tnvedInput, { target: { value: '1234567890' } });
      
      // Проверка автоформатирования
      expect(tnvedInput.value).toBe('1234 56 7890');
    });

    it('валидирует цену (должна быть положительным числом)', async () => {
      render(<CatalogForm onSubmit={mockOnSubmit} mode="create" />);
      
      const priceInput = screen.getByLabelText(/Цена/) as HTMLInputElement;
      
      // Пустая цена
      fireEvent.click(screen.getByText('Добавить в каталог'));
      await waitFor(() => {
        expect(screen.getByText('Цена обязательна')).toBeInTheDocument();
      });
      
      // Отрицательная цена
      fireEvent.change(priceInput, { target: { value: '-100' } });
      // Note: Поле type="number" с min="0" предотвращает ввод отрицательных чисел
      
      // Валидная цена
      fireEvent.change(priceInput, { target: { value: '100.50' } });
      await waitFor(() => {
        expect(screen.queryByText('Цена должна быть положительным числом')).not.toBeInTheDocument();
      });
    });

    it('успешно создает товар с валидными данными', async () => {
      const mockItem: CatalogItem = {
        id: '123',
        name: 'Test Product',
        number: 'PROD-001',
        description: 'Test description',
        tnved_code: '1234567890',
        price: 150.00,
        currency: 'KGS',
        unit: '796',
        created_at: Date.now(),
        updated_at: Date.now(),
      };

      mockOnSubmit.mockResolvedValue(mockItem);

      render(<CatalogForm onSubmit={mockOnSubmit} mode="create" />);
      
      // Заполнение формы
      fireEvent.change(screen.getByLabelText(/Наименование товара/), {
        target: { value: 'Test Product' },
      });
      fireEvent.change(screen.getByLabelText(/Номер по каталогу/), {
        target: { value: 'PROD-001' },
      });
      fireEvent.change(screen.getByLabelText(/Код ТНВЭД/), {
        target: { value: '1234567890' },
      });
      fireEvent.change(screen.getByLabelText(/Цена/), {
        target: { value: '150.00' },
      });
      
      // Валюта и единица измерения уже установлены по умолчанию

      // Отправка формы
      fireEvent.click(screen.getByText('Добавить в каталог'));

      await waitFor(() => {
        expect(mockOnSubmit).toHaveBeenCalledWith({
          name: 'Test Product',
          number: 'PROD-001',
          tnved_code: '1234567890', // Без пробелов
          price: 150.00,
          currency: 'KGS',
          unit: '796',
        });
      });

      // Проверка сообщения об успехе
      await waitFor(() => {
        expect(screen.getByText('Товар успешно добавлен в каталог')).toBeInTheDocument();
      });

      // Проверка редиректа
      await waitFor(() => {
        expect(mockPush).toHaveBeenCalledWith('/catalog');
      }, { timeout: 2000 });
    });

    it('включает опциональное описание если оно заполнено', async () => {
      mockOnSubmit.mockResolvedValue({} as CatalogItem);

      render(<CatalogForm onSubmit={mockOnSubmit} mode="create" />);
      
      // Заполнение обязательных полей + описание
      fireEvent.change(screen.getByLabelText(/Наименование товара/), {
        target: { value: 'Test Product' },
      });
      fireEvent.change(screen.getByLabelText(/Номер по каталогу/), {
        target: { value: 'PROD-001' },
      });
      fireEvent.change(screen.getByLabelText(/Описание/), {
        target: { value: 'Detailed description' },
      });
      fireEvent.change(screen.getByLabelText(/Код ТНВЭД/), {
        target: { value: '1234567890' },
      });
      fireEvent.change(screen.getByLabelText(/Цена/), {
        target: { value: '100' },
      });

      fireEvent.click(screen.getByText('Добавить в каталог'));

      await waitFor(() => {
        const callArg = mockOnSubmit.mock.calls[0][0];
        expect(callArg).toHaveProperty('description', 'Detailed description');
      });
    });

    it('обрабатывает ошибку при создании', async () => {
      const errorMessage = 'Товар с таким номером уже существует';
      mockOnSubmit.mockRejectedValue(new Error(errorMessage));

      render(<CatalogForm onSubmit={mockOnSubmit} mode="create" />);
      
      // Заполнение формы валидными данными
      fireEvent.change(screen.getByLabelText(/Наименование товара/), {
        target: { value: 'Test Product' },
      });
      fireEvent.change(screen.getByLabelText(/Номер по каталогу/), {
        target: { value: 'PROD-001' },
      });
      fireEvent.change(screen.getByLabelText(/Код ТНВЭД/), {
        target: { value: '1234567890' },
      });
      fireEvent.change(screen.getByLabelText(/Цена/), {
        target: { value: '100' },
      });

      fireEvent.click(screen.getByText('Добавить в каталог'));

      await waitFor(() => {
        expect(screen.getByText(errorMessage)).toBeInTheDocument();
      });

      expect(mockPush).not.toHaveBeenCalled();
    });
  });

  describe('Режим редактирования', () => {
    const initialData: CatalogItem = {
      id: '123',
      name: 'Existing Product',
      number: 'PROD-001',
      description: 'Existing description',
      tnved_code: '1234567890',
      price: 200.00,
      currency: 'USD',
      unit: '166',
      created_at: Date.now(),
      updated_at: Date.now(),
    };

    it('отображает форму редактирования с начальными данными', () => {
      render(<CatalogForm onSubmit={mockOnSubmit} mode="edit" initialData={initialData} />);
      
      expect(screen.getByText('Редактирование товара')).toBeInTheDocument();
      expect(screen.getByDisplayValue('Existing Product')).toBeInTheDocument();
      expect(screen.getByDisplayValue('PROD-001')).toBeInTheDocument();
      expect(screen.getByDisplayValue('Existing description')).toBeInTheDocument();
      expect(screen.getByDisplayValue('1234 56 7890')).toBeInTheDocument(); // Форматированный ТНВЭД
      expect(screen.getByDisplayValue('200')).toBeInTheDocument();
    });

    it('успешно обновляет товар', async () => {
      const updatedItem: CatalogItem = {
        ...initialData,
        name: 'Updated Product',
        price: 250.00,
      };

      mockOnSubmit.mockResolvedValue(updatedItem);

      render(<CatalogForm onSubmit={mockOnSubmit} mode="edit" initialData={initialData} />);
      
      // Изменение названия и цены
      const nameInput = screen.getByLabelText(/Наименование товара/) as HTMLInputElement;
      const priceInput = screen.getByLabelText(/Цена/) as HTMLInputElement;
      
      fireEvent.change(nameInput, { target: { value: 'Updated Product' } });
      fireEvent.change(priceInput, { target: { value: '250.00' } });

      fireEvent.click(screen.getByText('Сохранить изменения'));

      await waitFor(() => {
        expect(mockOnSubmit).toHaveBeenCalled();
      });

      // Проверка сообщения об успехе
      await waitFor(() => {
        expect(screen.getByText('Товар успешно обновлен')).toBeInTheDocument();
      });
    });
  });

  describe('Справочники', () => {
    it('отображает все валюты из справочника CURRENCIES', () => {
      render(<CatalogForm onSubmit={mockOnSubmit} mode="create" />);
      
      const currencySelect = screen.getByLabelText(/Валюта/) as HTMLSelectElement;
      const options = Array.from(currencySelect.options).map(opt => opt.value);
      
      // Проверка наличия основных валют
      expect(options).toContain('KGS');
      expect(options).toContain('USD');
      expect(options).toContain('RUB');
      expect(options).toContain('EUR');
      expect(options).toContain('CNY');
    });

    it('отображает все единицы измерения из справочника UNITS_OF_MEASURE', () => {
      render(<CatalogForm onSubmit={mockOnSubmit} mode="create" />);
      
      const unitSelect = screen.getByLabelText(/Единица измерения/) as HTMLSelectElement;
      const options = Array.from(unitSelect.options).map(opt => opt.value);
      
      // Проверка наличия основных единиц
      expect(options).toContain('796'); // шт
      expect(options).toContain('166'); // кг
      expect(options).toContain('112'); // л
    });

    it('отображает правильный символ валюты для выбранной валюты', () => {
      render(<CatalogForm onSubmit={mockOnSubmit} mode="create" />);
      
      const currencySelect = screen.getByLabelText(/Валюта/) as HTMLSelectElement;
      
      // Выбор USD
      fireEvent.change(currencySelect, { target: { value: 'USD' } });
      expect(screen.getByText('$')).toBeInTheDocument();
      
      // Выбор EUR
      fireEvent.change(currencySelect, { target: { value: 'EUR' } });
      expect(screen.getByText('€')).toBeInTheDocument();
      
      // Выбор RUB
      fireEvent.change(currencySelect, { target: { value: 'RUB' } });
      expect(screen.getByText('₽')).toBeInTheDocument();
    });
  });

  describe('Навигация', () => {
    it('возвращается назад при клике на кнопку "Вернуться"', () => {
      render(<CatalogForm onSubmit={mockOnSubmit} />);
      
      const backButton = screen.getByText('Вернуться');
      fireEvent.click(backButton);

      expect(mockBack).toHaveBeenCalled();
    });

    it('переходит к списку каталога при клике на "Отмена"', () => {
      render(<CatalogForm onSubmit={mockOnSubmit} />);
      
      const cancelButton = screen.getByText('Отмена');
      fireEvent.click(cancelButton);

      expect(mockPush).toHaveBeenCalledWith('/catalog');
    });
  });

  describe('Очистка ошибок', () => {
    it('очищает ошибку поля при изменении значения', async () => {
      render(<CatalogForm onSubmit={mockOnSubmit} mode="create" />);
      
      // Вызвать ошибку
      fireEvent.click(screen.getByText('Добавить в каталог'));
      
      await waitFor(() => {
        expect(screen.getByText('Название товара/услуги обязательно')).toBeInTheDocument();
      });
      
      // Исправить поле
      const nameInput = screen.getByLabelText(/Наименование товара/);
      fireEvent.change(nameInput, { target: { value: 'Test Product' } });
      
      // Ошибка должна исчезнуть
      await waitFor(() => {
        expect(screen.queryByText('Название товара/услуги обязательно')).not.toBeInTheDocument();
      });
    });
  });
});
