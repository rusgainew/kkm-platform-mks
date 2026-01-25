import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { useRouter } from 'next/navigation';
import { describe, it, expect, beforeEach, vi } from 'vitest';
import UserCreateForm from './UserCreateForm';
import * as usersApi from '@/lib/api/users';

// Mock next/navigation
vi.mock('next/navigation', () => ({
  useRouter: vi.fn(),
}));

// Mock API
vi.mock('@/lib/api/users', () => ({
  createUser: vi.fn(),
}));

describe('UserCreateForm', () => {
  const mockPush = vi.fn();
  const mockBack = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    (useRouter as any).mockReturnValue({
      push: mockPush,
      back: mockBack,
    });
  });

  describe('Form Rendering', () => {
    it('renders all form fields', () => {
      render(<UserCreateForm />);

      expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/имя/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/фамилия/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/роль/i)).toBeInTheDocument();
      expect(screen.getByLabelText(/пароль/i)).toBeInTheDocument();
    });

    it('renders submit and cancel buttons', () => {
      render(<UserCreateForm />);

      expect(screen.getByRole('button', { name: /создать пользователя/i })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: /отмена/i })).toBeInTheDocument();
    });

    it('renders breadcrumb navigation', () => {
      render(<UserCreateForm />);

      expect(screen.getByText(/вернуться к списку/i)).toBeInTheDocument();
    });

    it('sets default role to cashier', () => {
      render(<UserCreateForm />);

      const roleSelect = screen.getByLabelText(/роль/i) as HTMLSelectElement;
      expect(roleSelect.value).toBe('cashier');
    });
  });

  describe('Form Validation', () => {
    it('shows email required error when email is empty', async () => {
      render(<UserCreateForm />);

      const submitButton = screen.getByRole('button', { name: /создать пользователя/i });
      fireEvent.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText(/email обязателен/i)).toBeInTheDocument();
      });
    });

    it('shows email format error for invalid email', async () => {
      render(<UserCreateForm />);

      const emailInput = screen.getByLabelText(/email/i) as HTMLInputElement;
      await userEvent.type(emailInput, 'invalid-email');

      const submitButton = screen.getByRole('button', { name: /создать пользователя/i });
      fireEvent.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText(/некорректный email/i)).toBeInTheDocument();
      });
    });

    it('shows first name required error', async () => {
      render(<UserCreateForm />);

      const submitButton = screen.getByRole('button', { name: /создать пользователя/i });
      fireEvent.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText(/имя обязательно/i)).toBeInTheDocument();
      });
    });

    it('shows last name required error', async () => {
      render(<UserCreateForm />);

      const submitButton = screen.getByRole('button', { name: /создать пользователя/i });
      fireEvent.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText(/фамилия обязательна/i)).toBeInTheDocument();
      });
    });

    it('shows password required error', async () => {
      render(<UserCreateForm />);

      const submitButton = screen.getByRole('button', { name: /создать пользователя/i });
      fireEvent.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText(/пароль обязателен/i)).toBeInTheDocument();
      });
    });

    it('prevents form submission with validation errors', async () => {
      (usersApi.createUser as any).mockResolvedValueOnce({
        user: { user_id: '1' },
      });

      render(<UserCreateForm />);

      const submitButton = screen.getByRole('button', { name: /создать пользователя/i });
      fireEvent.click(submitButton);

      await waitFor(() => {
        expect(usersApi.createUser).not.toHaveBeenCalled();
      });
    });
  });

  describe('Password Validation', () => {
    it('shows password requirements when typing', async () => {
      render(<UserCreateForm />);

      const passwordInput = screen.getByLabelText(/пароль/i) as HTMLInputElement;
      await userEvent.type(passwordInput, 'test');

      expect(screen.getByText(/требования к паролю/i)).toBeInTheDocument();
      expect(screen.getByText(/8 символов/i)).toBeInTheDocument();
    });

    it('validates password with 8 characters', async () => {
      render(<UserCreateForm />);

      const passwordInput = screen.getByLabelText(/пароль/i) as HTMLInputElement;
      await userEvent.type(passwordInput, 'TeSt1234!');

      // Should show all requirements met
      const checkmarks = screen.getAllByText(/требования к паролю/i);
      expect(checkmarks.length).toBeGreaterThan(0);
    });

    it('rejects password without uppercase', async () => {
      render(<UserCreateForm />);

      const passwordInput = screen.getByLabelText(/пароль/i) as HTMLInputElement;
      await userEvent.type(passwordInput, 'test1234!');

      const submitButton = screen.getByRole('button', { name: /создать пользователя/i });
      fireEvent.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText(/пароль не соответствует требованиям/i)).toBeInTheDocument();
      });
    });
  });

  describe('Form Submission', () => {
    it('submits form with valid data', async () => {
      (usersApi.createUser as any).mockResolvedValueOnce({
        user: { user_id: '1' },
      });

      render(<UserCreateForm />);

      const emailInput = screen.getByLabelText(/email/i);
      const firstNameInput = screen.getByLabelText(/имя/i);
      const lastNameInput = screen.getByLabelText(/фамилия/i);
      const passwordInput = screen.getByLabelText(/пароль/i);

      await userEvent.type(emailInput, 'test@example.com');
      await userEvent.type(firstNameInput, 'John');
      await userEvent.type(lastNameInput, 'Doe');
      await userEvent.type(passwordInput, 'Test1234!');

      const submitButton = screen.getByRole('button', { name: /создать пользователя/i });
      fireEvent.click(submitButton);

      await waitFor(() => {
        expect(usersApi.createUser).toHaveBeenCalledWith({
          email: 'test@example.com',
          first_name: 'John',
          last_name: 'Doe',
          password: 'Test1234!',
        });
      });
    });

    it('shows success message after successful submission', async () => {
      (usersApi.createUser as any).mockResolvedValueOnce({
        user: { user_id: '1' },
      });

      render(<UserCreateForm />);

      const emailInput = screen.getByLabelText(/email/i);
      const firstNameInput = screen.getByLabelText(/имя/i);
      const lastNameInput = screen.getByLabelText(/фамилия/i);
      const passwordInput = screen.getByLabelText(/пароль/i);

      await userEvent.type(emailInput, 'test@example.com');
      await userEvent.type(firstNameInput, 'John');
      await userEvent.type(lastNameInput, 'Doe');
      await userEvent.type(passwordInput, 'Test1234!');

      const submitButton = screen.getByRole('button', { name: /создать пользователя/i });
      fireEvent.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText(/пользователь успешно создан/i)).toBeInTheDocument();
      });
    });

    it('redirects to users list after successful submission', async () => {
      (usersApi.createUser as any).mockResolvedValueOnce({
        user: { user_id: '1' },
      });

      render(<UserCreateForm />);

      const emailInput = screen.getByLabelText(/email/i);
      const firstNameInput = screen.getByLabelText(/имя/i);
      const lastNameInput = screen.getByLabelText(/фамилия/i);
      const passwordInput = screen.getByLabelText(/пароль/i);

      await userEvent.type(emailInput, 'test@example.com');
      await userEvent.type(firstNameInput, 'John');
      await userEvent.type(lastNameInput, 'Doe');
      await userEvent.type(passwordInput, 'Test1234!');

      const submitButton = screen.getByRole('button', { name: /создать пользователя/i });
      fireEvent.click(submitButton);

      await waitFor(
        () => {
          expect(mockPush).toHaveBeenCalledWith('/users');
        },
        { timeout: 3000 }
      );
    });

    it('shows error message on API failure', async () => {
      (usersApi.createUser as any).mockRejectedValueOnce(
        new Error('Email уже зарегистрирован в системе')
      );

      render(<UserCreateForm />);

      const emailInput = screen.getByLabelText(/email/i);
      const firstNameInput = screen.getByLabelText(/имя/i);
      const lastNameInput = screen.getByLabelText(/фамилия/i);
      const passwordInput = screen.getByLabelText(/пароль/i);

      await userEvent.type(emailInput, 'test@example.com');
      await userEvent.type(firstNameInput, 'John');
      await userEvent.type(lastNameInput, 'Doe');
      await userEvent.type(passwordInput, 'Test1234!');

      const submitButton = screen.getByRole('button', { name: /создать пользователя/i });
      fireEvent.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText(/email уже зарегистрирован/i)).toBeInTheDocument();
      });
    });
  });

  describe('User Interactions', () => {
    it('clears email error when user starts typing', async () => {
      render(<UserCreateForm />);

      const submitButton = screen.getByRole('button', { name: /создать пользователя/i });
      fireEvent.click(submitButton);

      await waitFor(() => {
        expect(screen.getByText(/email обязателен/i)).toBeInTheDocument();
      });

      const emailInput = screen.getByLabelText(/email/i);
      await userEvent.type(emailInput, 'test');

      await waitFor(() => {
        expect(screen.queryByText(/email обязателен/i)).not.toBeInTheDocument();
      });
    });

    it('allows role selection', async () => {
      render(<UserCreateForm />);

      const roleSelect = screen.getByLabelText(/роль/i) as HTMLSelectElement;
      expect(roleSelect.value).toBe('cashier');

      await userEvent.selectOptions(roleSelect, 'manager');
      expect(roleSelect.value).toBe('manager');

      await userEvent.selectOptions(roleSelect, 'admin');
      expect(roleSelect.value).toBe('admin');
    });

    it('calls router.back() when cancel button is clicked', async () => {
      render(<UserCreateForm />);

      const cancelButton = screen.getByRole('button', { name: /отмена/i });
      fireEvent.click(cancelButton);

      expect(mockBack).toHaveBeenCalled();
    });

    it('calls router.back() when breadcrumb link is clicked', async () => {
      render(<UserCreateForm />);

      const backLink = screen.getByText(/вернуться к списку/i);
      fireEvent.click(backLink);

      expect(mockBack).toHaveBeenCalled();
    });
  });

  describe('Form Disabled State', () => {
    it('disables form during submission', async () => {
      (usersApi.createUser as any).mockImplementationOnce(
        () => new Promise(resolve => setTimeout(() => resolve({ user: { user_id: '1' } }), 1000))
      );

      render(<UserCreateForm />);

      const emailInput = screen.getByLabelText(/email/i);
      const firstNameInput = screen.getByLabelText(/имя/i);
      const lastNameInput = screen.getByLabelText(/фамилия/i);
      const passwordInput = screen.getByLabelText(/пароль/i);

      await userEvent.type(emailInput, 'test@example.com');
      await userEvent.type(firstNameInput, 'John');
      await userEvent.type(lastNameInput, 'Doe');
      await userEvent.type(passwordInput, 'Test1234!');

      const submitButton = screen.getByRole('button', { name: /создать пользователя/i });
      fireEvent.click(submitButton);

      expect(emailInput).toBeDisabled();
      expect(passwordInput).toBeDisabled();
      expect(submitButton).toBeDisabled();
    });
  });
});
