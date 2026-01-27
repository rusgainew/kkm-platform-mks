import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { ForeignCompanyForm } from "../components/ForeignCompanyForm";
import * as foreignCompaniesApi from "@/lib/api/foreign-companies";

// Мокируем API
vi.mock("@/lib/api/foreign-companies");

// Мокируем toast
vi.mock("react-hot-toast", () => ({
  default: {
    success: vi.fn(),
    error: vi.fn(),
  },
}));

// Тестовые данные
const mockCompany: foreignCompaniesApi.ForeignCompany = {
  id: "test-uuid-123",
  name: "Acme Corporation Ltd",
  tax_id: "12-3456789",
  country: "US",
  address: "123 Main Street, New York, NY 10001, USA",
  contact_email: "contact@acme.com",
  contact_phone: "+1 (555) 123-4567",
  currency: "USD",
  is_active: true,
  created_at: 1700000000,
  updated_at: 1700000000,
};

/**
 * Wrapper для React Query
 */
function createWrapper() {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });

  const Wrapper = ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
  Wrapper.displayName = "QueryWrapper";
  
  return Wrapper;
}

describe("ForeignCompanyForm", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("Rendering", () => {
    it("should render create form with all fields", () => {
      render(<ForeignCompanyForm />, { wrapper: createWrapper() });

      expect(
        screen.getByText("Добавить иностранную компанию")
      ).toBeInTheDocument();
      expect(screen.getByLabelText(/Название компании/)).toBeInTheDocument();
      expect(screen.getByLabelText(/Налоговый номер/)).toBeInTheDocument();
      expect(screen.getByLabelText(/Страна/)).toBeInTheDocument();
      expect(screen.getByLabelText(/Адрес/)).toBeInTheDocument();
      expect(screen.getByLabelText(/Email/)).toBeInTheDocument();
      expect(screen.getByLabelText(/Телефон/)).toBeInTheDocument();
      expect(screen.getByLabelText(/Валюта/)).toBeInTheDocument();
      expect(
        screen.getByRole("button", { name: /Создать компанию/ })
      ).toBeInTheDocument();
    });

    it("should render edit form with company data", () => {
      render(<ForeignCompanyForm company={mockCompany} />, {
        wrapper: createWrapper(),
      });

      expect(
        screen.getByText("Редактировать иностранную компанию")
      ).toBeInTheDocument();
      expect(screen.getByDisplayValue("Acme Corporation Ltd")).toBeInTheDocument();
      expect(screen.getByDisplayValue("12-3456789")).toBeInTheDocument();
      expect(screen.getByDisplayValue("US")).toBeInTheDocument();
      expect(
        screen.getByDisplayValue("123 Main Street, New York, NY 10001, USA")
      ).toBeInTheDocument();
      expect(screen.getByDisplayValue("contact@acme.com")).toBeInTheDocument();
      expect(
        screen.getByDisplayValue("+1 (555) 123-4567")
      ).toBeInTheDocument();
      expect(screen.getByDisplayValue("USD")).toBeInTheDocument();
    });

    it("should show delete button in edit mode", () => {
      render(<ForeignCompanyForm company={mockCompany} />, {
        wrapper: createWrapper(),
      });

      expect(
        screen.getByRole("button", { name: /Удалить/ })
      ).toBeInTheDocument();
    });

    it("should not show delete button in create mode", () => {
      render(<ForeignCompanyForm />, { wrapper: createWrapper() });

      expect(
        screen.queryByRole("button", { name: /Удалить/ })
      ).not.toBeInTheDocument();
    });
  });

  describe("Validation", () => {
    it("should show error for empty company name", async () => {
      const user = userEvent.setup();
      render(<ForeignCompanyForm />, { wrapper: createWrapper() });

      const submitButton = screen.getByRole("button", {
        name: /Создать компанию/,
      });
      await user.click(submitButton);

      await waitFor(() => {
        expect(
          screen.getByText("Название компании обязательно")
        ).toBeInTheDocument();
      });
    });

    it("should show error for too short company name", async () => {
      const user = userEvent.setup();
      render(<ForeignCompanyForm />, { wrapper: createWrapper() });

      const nameInput = screen.getByLabelText(/Название компании/);
      await user.type(nameInput, "A");

      const submitButton = screen.getByRole("button", {
        name: /Создать компанию/,
      });
      await user.click(submitButton);

      await waitFor(() => {
        expect(
          screen.getByText("Название должно содержать минимум 2 символа")
        ).toBeInTheDocument();
      });
    });

    it("should show error for empty tax ID", async () => {
      const user = userEvent.setup();
      render(<ForeignCompanyForm />, { wrapper: createWrapper() });

      const submitButton = screen.getByRole("button", {
        name: /Создать компанию/,
      });
      await user.click(submitButton);

      await waitFor(() => {
        expect(
          screen.getByText("Налоговый номер (PIN/TIN) обязателен")
        ).toBeInTheDocument();
      });
    });

    it("should show error for invalid tax ID format", async () => {
      const user = userEvent.setup();
      render(<ForeignCompanyForm />, { wrapper: createWrapper() });

      const taxIdInput = screen.getByLabelText(/Налоговый номер/);
      await user.type(taxIdInput, "123"); // Too short

      const submitButton = screen.getByRole("button", {
        name: /Создать компанию/,
      });
      await user.click(submitButton);

      await waitFor(() => {
        expect(
          screen.getByText(/Неверный формат налогового номера/)
        ).toBeInTheDocument();
      });
    });

    it("should show error for empty country code", async () => {
      const user = userEvent.setup();
      render(<ForeignCompanyForm />, { wrapper: createWrapper() });

      const submitButton = screen.getByRole("button", {
        name: /Создать компанию/,
      });
      await user.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText("Код страны обязателен")).toBeInTheDocument();
      });
    });

    it("should show error for empty address", async () => {
      const user = userEvent.setup();
      render(<ForeignCompanyForm />, { wrapper: createWrapper() });

      const submitButton = screen.getByRole("button", {
        name: /Создать компанию/,
      });
      await user.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText("Адрес обязателен")).toBeInTheDocument();
      });
    });

    it("should show error for too short address", async () => {
      const user = userEvent.setup();
      render(<ForeignCompanyForm />, { wrapper: createWrapper() });

      const addressInput = screen.getByLabelText(/Адрес/);
      await user.type(addressInput, "123");

      const submitButton = screen.getByRole("button", {
        name: /Создать компанию/,
      });
      await user.click(submitButton);

      await waitFor(() => {
        expect(
          screen.getByText("Адрес должен содержать минимум 5 символов")
        ).toBeInTheDocument();
      });
    });

    it("should show error for empty email", async () => {
      const user = userEvent.setup();
      render(<ForeignCompanyForm />, { wrapper: createWrapper() });

      const submitButton = screen.getByRole("button", {
        name: /Создать компанию/,
      });
      await user.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText("Email обязателен")).toBeInTheDocument();
      });
    });

    it("should show error for invalid email format", async () => {
      const user = userEvent.setup();
      render(<ForeignCompanyForm />, { wrapper: createWrapper() });

      const emailInput = screen.getByLabelText(/Email/);
      await user.type(emailInput, "invalid-email");

      const submitButton = screen.getByRole("button", {
        name: /Создать компанию/,
      });
      await user.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText("Неверный формат email")).toBeInTheDocument();
      });
    });

    it("should show error for empty phone", async () => {
      const user = userEvent.setup();
      render(<ForeignCompanyForm />, { wrapper: createWrapper() });

      const submitButton = screen.getByRole("button", {
        name: /Создать компанию/,
      });
      await user.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText("Телефон обязателен")).toBeInTheDocument();
      });
    });

    it("should show error for invalid phone format", async () => {
      const user = userEvent.setup();
      render(<ForeignCompanyForm />, { wrapper: createWrapper() });

      const phoneInput = screen.getByLabelText(/Телефон/);
      await user.type(phoneInput, "123"); // Too short

      const submitButton = screen.getByRole("button", {
        name: /Создать компанию/,
      });
      await user.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText(/Неверный формат телефона/)).toBeInTheDocument();
      });
    });
  });

  describe("Create Company", () => {
    it("should create company with valid data", async () => {
      const user = userEvent.setup();
      const onSuccess = vi.fn();

      vi.mocked(foreignCompaniesApi.createForeignCompany).mockResolvedValue({
        data: mockCompany,
        success: true,
      });

      render(<ForeignCompanyForm onSuccess={onSuccess} />, {
        wrapper: createWrapper(),
      });

      // Заполняем форму
      await user.type(
        screen.getByLabelText(/Название компании/),
        "Test Company"
      );
      await user.type(screen.getByLabelText(/Налоговый номер/), "12-3456789");

      const countrySelect = screen.getByLabelText(/Страна/);
      await user.selectOptions(countrySelect, "US");

      await user.type(
        screen.getByLabelText(/Адрес/),
        "123 Main Street, New York"
      );
      await user.type(
        screen.getByLabelText(/Email/),
        "test@example.com"
      );
      await user.type(
        screen.getByLabelText(/Телефон/),
        "+1 (555) 123-4567"
      );

      const currencySelect = screen.getByLabelText(/Валюта/);
      await user.selectOptions(currencySelect, "USD");

      // Отправляем форму
      const submitButton = screen.getByRole("button", {
        name: /Создать компанию/,
      });
      await user.click(submitButton);

      await waitFor(() => {
        expect(foreignCompaniesApi.createForeignCompany).toHaveBeenCalledWith({
          name: "Test Company",
          tax_id: "12-3456789",
          country: "US",
          address: "123 Main Street, New York",
          contact_email: "test@example.com",
          contact_phone: "+1 (555) 123-4567",
          currency: "USD",
        });
        expect(onSuccess).toHaveBeenCalled();
      });
    });

    it("should normalize data before submission", async () => {
      const user = userEvent.setup();

      vi.mocked(foreignCompaniesApi.createForeignCompany).mockResolvedValue({
        data: mockCompany,
        success: true,
      });

      render(<ForeignCompanyForm />, { wrapper: createWrapper() });

      // Вводим данные с лишними пробелами и в неправильном регистре
      await user.type(
        screen.getByLabelText(/Название компании/),
        "  Test Company  "
      );
      await user.type(
        screen.getByLabelText(/Налоговый номер/),
        "  abc-123  "
      );

      const countrySelect = screen.getByLabelText(/Страна/);
      await user.selectOptions(countrySelect, "US");

      await user.type(
        screen.getByLabelText(/Адрес/),
        "  123 Main Street  "
      );
      await user.type(
        screen.getByLabelText(/Email/),
        "  TEST@EXAMPLE.COM  "
      );
      await user.type(
        screen.getByLabelText(/Телефон/),
        "  +1 (555) 123-4567  "
      );

      const submitButton = screen.getByRole("button", {
        name: /Создать компанию/,
      });
      await user.click(submitButton);

      await waitFor(() => {
        expect(foreignCompaniesApi.createForeignCompany).toHaveBeenCalledWith(
          expect.objectContaining({
            name: "Test Company",
            tax_id: "ABC-123",
            country: "US",
            address: "123 Main Street",
            contact_email: "test@example.com",
            contact_phone: "+1 (555) 123-4567",
            currency: "USD",
          })
        );
      });
    });
  });

  describe("Update Company", () => {
    it("should update company with modified data", async () => {
      const user = userEvent.setup();
      const onSuccess = vi.fn();

      vi.mocked(foreignCompaniesApi.updateForeignCompany).mockResolvedValue({
        data: { ...mockCompany, name: "Updated Company" },
        success: true,
      });

      render(<ForeignCompanyForm company={mockCompany} onSuccess={onSuccess} />, {
        wrapper: createWrapper(),
      });

      // Изменяем название
      const nameInput = screen.getByLabelText(/Название компании/);
      await user.clear(nameInput);
      await user.type(nameInput, "Updated Company");

      const submitButton = screen.getByRole("button", {
        name: /Сохранить изменения/,
      });
      await user.click(submitButton);

      await waitFor(() => {
        expect(foreignCompaniesApi.updateForeignCompany).toHaveBeenCalledWith(
          mockCompany.id,
          expect.objectContaining({
            name: "Updated Company",
          })
        );
        expect(onSuccess).toHaveBeenCalled();
      });
    });
  });

  describe("Delete Company", () => {
    it("should show delete confirmation modal", async () => {
      const user = userEvent.setup();
      render(<ForeignCompanyForm company={mockCompany} />, {
        wrapper: createWrapper(),
      });

      const deleteButton = screen.getByRole("button", { name: /Удалить/ });
      await user.click(deleteButton);

      await waitFor(() => {
        expect(screen.getByText("Подтвердите удаление")).toBeInTheDocument();
        expect(
          screen.getByText(/Вы уверены, что хотите удалить компанию/)
        ).toBeInTheDocument();
      });
    });

    it("should delete company on confirmation", async () => {
      const user = userEvent.setup();
      const onSuccess = vi.fn();

      vi.mocked(foreignCompaniesApi.deleteForeignCompany).mockResolvedValue({
        success: true,
      });

      render(<ForeignCompanyForm company={mockCompany} onSuccess={onSuccess} />, {
        wrapper: createWrapper(),
      });

      // Открываем модальное окно
      const deleteButton = screen.getByRole("button", { name: /Удалить/ });
      await user.click(deleteButton);

      // Подтверждаем удаление
      const confirmButton = await screen.findByRole("button", {
        name: /Удалить/,
      });
      await user.click(confirmButton);

      await waitFor(() => {
        expect(foreignCompaniesApi.deleteForeignCompany).toHaveBeenCalledWith(
          mockCompany.id
        );
        expect(onSuccess).toHaveBeenCalled();
      });
    });

    it("should cancel deletion on cancel button", async () => {
      const user = userEvent.setup();
      render(<ForeignCompanyForm company={mockCompany} />, {
        wrapper: createWrapper(),
      });

      // Открываем модальное окно
      const deleteButton = screen.getByRole("button", { name: /Удалить/ });
      await user.click(deleteButton);

      // Нажимаем отмену
      const cancelButtons = await screen.findAllByRole("button", {
        name: /Отмена/,
      });
      await user.click(cancelButtons[1]); // Вторая кнопка "Отмена" в модальном окне

      await waitFor(() => {
        expect(
          screen.queryByText("Подтвердите удаление")
        ).not.toBeInTheDocument();
      });
    });
  });

  describe("Form Behavior", () => {
    it("should clear field error on input change", async () => {
      const user = userEvent.setup();
      render(<ForeignCompanyForm />, { wrapper: createWrapper() });

      // Вызываем ошибку
      const submitButton = screen.getByRole("button", {
        name: /Создать компанию/,
      });
      await user.click(submitButton);

      await waitFor(() => {
        expect(
          screen.getByText("Название компании обязательно")
        ).toBeInTheDocument();
      });

      // Вводим текст
      const nameInput = screen.getByLabelText(/Название компании/);
      await user.type(nameInput, "Test");

      // Ошибка должна исчезнуть
      await waitFor(() => {
        expect(
          screen.queryByText("Название компании обязательно")
        ).not.toBeInTheDocument();
      });
    });

    it("should call onCancel when cancel button is clicked", async () => {
      const user = userEvent.setup();
      const onCancel = vi.fn();

      render(<ForeignCompanyForm onCancel={onCancel} />, {
        wrapper: createWrapper(),
      });

      const cancelButton = screen.getByRole("button", { name: /Отмена/ });
      await user.click(cancelButton);

      expect(onCancel).toHaveBeenCalled();
    });

    it("should disable form during submission", async () => {
      const user = userEvent.setup();

      vi.mocked(foreignCompaniesApi.createForeignCompany).mockImplementation(
        () =>
          new Promise((resolve) =>
            setTimeout(
              () =>
                resolve({
                  data: mockCompany,
                  success: true,
                }),
              100
            )
          )
      );

      render(<ForeignCompanyForm />, { wrapper: createWrapper() });

      // Заполняем форму
      await user.type(
        screen.getByLabelText(/Название компании/),
        "Test Company"
      );
      await user.type(screen.getByLabelText(/Налоговый номер/), "12-3456789");

      const countrySelect = screen.getByLabelText(/Страна/);
      await user.selectOptions(countrySelect, "US");

      await user.type(
        screen.getByLabelText(/Адрес/),
        "123 Main Street, New York"
      );
      await user.type(
        screen.getByLabelText(/Email/),
        "test@example.com"
      );
      await user.type(
        screen.getByLabelText(/Телефон/),
        "+1 (555) 123-4567"
      );

      const submitButton = screen.getByRole("button", {
        name: /Создать компанию/,
      });
      await user.click(submitButton);

      // Проверяем, что кнопка отключена
      await waitFor(() => {
        expect(submitButton).toBeDisabled();
        expect(screen.getByText("Сохранение...")).toBeInTheDocument();
      });
    });
  });
});
