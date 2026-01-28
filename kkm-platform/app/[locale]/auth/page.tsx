'use client';

import { useEffect, useMemo, useRef, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { LogIn, UserPlus, Mail, KeyRound, Loader2, ShieldCheck, ArrowLeft, Lock, Eye, EyeOff } from 'lucide-react';
import {
  useLoginMutation,
  useRegisterMutation,
  useForgotPasswordMutation,
} from '@/lib/hooks/useAuthApi';
import { useAuthStore } from '@/store/authStore';
import ForgotPasswordModal from '@/components/auth/ForgotPasswordModal';
import ResetPasswordModal from '@/components/auth/ResetPasswordModal';
import { getRedirectPathByRole } from '@/lib/utils/roleRedirect';
import type { UserRole } from '@/types/entities';

export default function AuthPage() {
  const router = useRouter();
  const user = useAuthStore((s) => s.user);
  const tokens = useAuthStore((s) => s.tokens);
  const loginMutation = useLoginMutation();
  const registerMutation = useRegisterMutation();
  const forgotMutation = useForgotPasswordMutation();

  const [mode, setMode] = useState<'login' | 'register'>('login');
  const [loginForm, setLoginForm] = useState({ email: '', password: '' });
  const [regForm, setRegForm] = useState({ email: '', password: '', confirmPassword: '', firstName: '', lastName: '' });
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [showForgotModal, setShowForgotModal] = useState(false);
  const [showResetModal, setShowResetModal] = useState(false);
  const [message, setMessage] = useState<string | null>(null);
  const hasRedirectedRef = useRef(false);

  // Валидация пароля
  const validatePassword = (password: string) => {
    const checks = {
      length: password.length >= 8,
      uppercase: /[A-Z]/.test(password),
      lowercase: /[a-z]/.test(password),
      number: /[0-9]/.test(password),
      special: /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password),
    };
    return checks;
  };

  const passwordChecks = validatePassword(regForm.password);
  const isPasswordValid = Object.values(passwordChecks).every(Boolean);

  useEffect(() => {
    if (!hasRedirectedRef.current && user && tokens?.accessToken) {
      console.log('[AuthPage] Пользователь аутентифицирован, роль:', user.role);
      hasRedirectedRef.current = true;
      
      // Получаем путь редиректа по роли
      const redirectPath = getRedirectPathByRole(user.role);
      console.log('[AuthPage] Перенаправление на:', redirectPath);
      
      router.replace(redirectPath);
    }
  }, [user, tokens, router]);

  const isLoading = useMemo(
    () => loginMutation.isPending || registerMutation.isPending || forgotMutation.isPending,
    [loginMutation.isPending, registerMutation.isPending, forgotMutation.isPending]
  );

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setMessage(null);
    try {
      const response = await loginMutation.mutateAsync({ email: loginForm.email, password: loginForm.password });
      
      // Получаем роль пользователя из ответа
      const userRole = response?.user?.role as UserRole | undefined;
      console.log('[handleLogin] Ответ:', response);
      console.log('[handleLogin] Роль пользователя из ответа:', userRole);
      
      // Получаем путь редиректа по роли
      const redirectPath = getRedirectPathByRole(userRole);
      
      // Небольшая задержка для обновления state
      setTimeout(() => {
        console.log('[handleLogin] Redirecting to:', redirectPath);
        router.replace(redirectPath);
      }, 100);
    } catch (err) {
      // Parse error message
      let errorMessage = 'Ошибка входа';
      const errMessage = (err as Error)?.message || '';
      if (errMessage) {
        const msg = errMessage;
        if (msg.includes('user not found')) {
          errorMessage = 'Пользователь не найден. Проверьте email или зарегистрируйтесь.';
        } else if (msg.includes('Invalid email or password') || msg.includes('invalid password')) {
          errorMessage = 'Неверный email или пароль.';
        } else if (msg.includes('validation failed')) {
          errorMessage = 'Ошибка валидации. Проверьте данные.';
        } else if (msg.includes('NOT_FOUND')) {
          errorMessage = 'Пользователь не найден. Создайте аккаунт для входа.';
        } else {
          errorMessage = msg;
        }
      }
      setMessage(errorMessage);
    }
  };

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    setMessage(null);
    try {
      await registerMutation.mutateAsync({
        email: regForm.email,
        password: regForm.password,
        first_name: regForm.firstName,
        last_name: regForm.lastName,
      });
      // Reset form and switch to login
      setRegForm({ email: '', password: '', confirmPassword: '', firstName: '', lastName: '' });
      setLoginForm({ email: regForm.email, password: '' });
      setMode('login');
      setMessage('✓ Аккаунт создан! Войдите в систему.');
    } catch (err) {
      // Parse error message
      let errorMessage = 'Ошибка регистрации';
      const errMessage = (err as Error)?.message || '';
      if (errMessage) {
        const msg = errMessage;
        if (msg.includes('already exists')) {
          errorMessage = 'Email уже зарегистрирован.';
        } else if (msg.includes('validation failed')) {
          errorMessage = 'Пароль не соответствует требованиям безопасности.';
        } else if (msg.includes('invalid')) {
          errorMessage = 'Некорректные данные. Проверьте заполнение полей.';
        } else {
          errorMessage = msg;
        }
      }
      setMessage(errorMessage);
    }
  };

  // Forgot password is handled in ForgotPasswordModal component

  return (
    <div className="min-h-screen bg-linear-to-br from-gray-950 via-gray-900 to-gray-950 flex items-center justify-center px-4">
      <div className="absolute top-6 left-6">
        <Link
          href="/"
          className="inline-flex items-center gap-2 text-gray-400 hover:text-white transition-colors"
        >
          <ArrowLeft className="w-4 h-4" />
          На главную
        </Link>
      </div>

      <div className="w-full max-w-4xl grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Left panel */}
        <div className="bg-gray-900/70 border border-gray-800 rounded-2xl p-8 shadow-2xl backdrop-blur-sm">
          <div className="flex items-center gap-3 mb-6">
            <div className="w-10 h-10 bg-linear-to-br from-blue-500 to-emerald-500 rounded-lg flex items-center justify-center shadow-blue-500/40">
              <ShieldCheck className="w-5 h-5 text-white" />
            </div>
            <div>
              <p className="text-sm text-gray-400">KKM Platform</p>
              <h1 className="text-2xl font-bold text-white">Вход в систему</h1>
            </div>
          </div>

          <div className="flex gap-2 mb-8">
            <button
              className={`flex-1 py-2 rounded-lg font-semibold transition-colors ${
                mode === 'login'
                  ? 'bg-linear-to-br from-blue-500 to-emerald-500 text-white shadow-lg shadow-blue-500/30'
                  : 'bg-gray-800 text-gray-300 border border-gray-700'
              }`}
              onClick={() => setMode('login')}
            >
              Вход
            </button>
            <button
              className={`flex-1 py-2 rounded-lg font-semibold transition-colors ${
                mode === 'register'
                  ? 'bg-linear-to-br from-blue-500 to-emerald-500 text-white shadow-lg shadow-blue-500/30'
                  : 'bg-gray-800 text-gray-300 border border-gray-700'
              }`}
              onClick={() => setMode('register')}
            >
              Регистрация
            </button>
          </div>

          {mode === 'login' && (
            <p className="text-xs text-gray-400 mb-4 p-2 rounded border border-gray-700 bg-gray-800/50">
              💡 Если у вас нет аккаунта, создайте его нажав на кнопку "Регистрация" выше.
            </p>
          )}

          {mode === 'login' ? (
            <form className="space-y-4" onSubmit={handleLogin}>
              <div>
                <label className="text-sm text-gray-300 mb-1 block">Email</label>
                <div className="relative">
                  <Mail className="w-4 h-4 text-gray-500 absolute left-3 top-1/2 -translate-y-1/2" />
                  <input
                    type="email"
                    required
                    autoComplete="email"
                    value={loginForm.email}
                    onChange={(e) => setLoginForm({ ...loginForm, email: e.target.value })}
                    className="w-full bg-gray-800 border border-gray-700 rounded-lg py-3 pl-10 pr-3 text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="you@example.com"
                  />
                </div>
              </div>
              <div>
                <label className="text-sm text-gray-300 mb-1 block">Пароль</label>
                <div className="relative">
                  <KeyRound className="w-4 h-4 text-gray-500 absolute left-3 top-1/2 -translate-y-1/2" />
                  <input
                    type="password"
                    required
                    autoComplete="current-password"
                    value={loginForm.password}
                    onChange={(e) => setLoginForm({ ...loginForm, password: e.target.value })}
                    className="w-full bg-gray-800 border border-gray-700 rounded-lg py-3 pl-10 pr-3 text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="••••••••"
                  />
                </div>
              </div>

              <button
                type="submit"
                disabled={isLoading}
                className="w-full py-3 rounded-lg font-semibold text-white bg-linear-to-br from-blue-500 to-emerald-500 shadow-lg shadow-blue-500/30 flex items-center justify-center gap-2 disabled:opacity-60"
              >
                {isLoading ? <Loader2 className="w-5 h-5 animate-spin" /> : <LogIn className="w-5 h-5" />}
                Войти
              </button>
            </form>
          ) : (
            <form className="space-y-4" onSubmit={handleRegister}>
              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="text-sm text-gray-300 mb-1 block">Имя</label>
                  <input
                    type="text"
                    required
                    autoComplete="given-name"
                    value={regForm.firstName}
                    onChange={(e) => setRegForm({ ...regForm, firstName: e.target.value })}
                    className="w-full bg-gray-800 border border-gray-700 rounded-lg py-3 px-3 text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                <div>
                  <label className="text-sm text-gray-300 mb-1 block">Фамилия</label>
                  <input
                    type="text"
                    required
                    autoComplete="family-name"
                    value={regForm.lastName}
                    onChange={(e) => setRegForm({ ...regForm, lastName: e.target.value })}
                    className="w-full bg-gray-800 border border-gray-700 rounded-lg py-3 px-3 text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
              </div>
              <div>
                <label className="text-sm text-gray-300 mb-1 block">Email</label>
                <div className="relative">
                  <Mail className="w-4 h-4 text-gray-500 absolute left-3 top-1/2 -translate-y-1/2" />
                  <input
                    type="email"
                    required
                    autoComplete="email"
                    value={regForm.email}
                    onChange={(e) => setRegForm({ ...regForm, email: e.target.value })}
                    className="w-full bg-gray-800 border border-gray-700 rounded-lg py-3 pl-10 pr-3 text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="you@example.com"
                  />
                </div>
              </div>
              <div>
                <label className="text-sm text-gray-300 mb-1 block">Пароль</label>
                <div className="relative">
                  <KeyRound className="w-4 h-4 text-gray-500 absolute left-3 top-1/2 -translate-y-1/2" />
                  <input
                    type={showPassword ? "text" : "password"}
                    required
                    autoComplete="new-password"
                    value={regForm.password}
                    onChange={(e) => setRegForm({ ...regForm, password: e.target.value })}
                    className="w-full bg-gray-800 border border-gray-700 rounded-lg py-3 pl-10 pr-10 text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="••••••••"
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    disabled={isLoading}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-300 disabled:opacity-50 disabled:cursor-not-allowed"
                    aria-label={showPassword ? "Скрыть пароль" : "Показать пароль"}
                  >
                    {showPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                  </button>
                </div>
                {regForm.password && (
                  <div className="mt-3 p-3 bg-gray-800/50 rounded-lg border border-gray-700 space-y-2">
                    <p className="text-xs text-gray-400 font-semibold mb-2">Требования к паролю:</p>
                    <div className="space-y-1">
                      <div className={`flex items-center gap-2 text-xs ${passwordChecks.length ? 'text-emerald-400' : 'text-gray-500'}`}>
                        <span>{passwordChecks.length ? '✓' : '○'}</span>
                        <span>Минимум 8 символов</span>
                      </div>
                      <div className={`flex items-center gap-2 text-xs ${passwordChecks.uppercase ? 'text-emerald-400' : 'text-gray-500'}`}>
                        <span>{passwordChecks.uppercase ? '✓' : '○'}</span>
                        <span>Заглавная буква (A-Z)</span>
                      </div>
                      <div className={`flex items-center gap-2 text-xs ${passwordChecks.lowercase ? 'text-emerald-400' : 'text-gray-500'}`}>
                        <span>{passwordChecks.lowercase ? '✓' : '○'}</span>
                        <span>Строчная буква (a-z)</span>
                      </div>
                      <div className={`flex items-center gap-2 text-xs ${passwordChecks.number ? 'text-emerald-400' : 'text-gray-500'}`}>
                        <span>{passwordChecks.number ? '✓' : '○'}</span>
                        <span>Цифра (0-9)</span>
                      </div>
                      <div className={`flex items-center gap-2 text-xs ${passwordChecks.special ? 'text-emerald-400' : 'text-gray-500'}`}>
                        <span>{passwordChecks.special ? '✓' : '○'}</span>
                        <span>Спецсимвол (!@#$%^&amp;*)</span>
                      </div>
                    </div>
                  </div>
                )}
              </div>

              <div>
                <label className="text-sm text-gray-300 mb-1 block">Подтвердите пароль</label>
                <div className="relative">
                  <Lock className="w-4 h-4 text-gray-500 absolute left-3 top-1/2 -translate-y-1/2" />
                  <input
                    type={showConfirmPassword ? "text" : "password"}
                    required
                    autoComplete="new-password"
                    value={regForm.confirmPassword}
                    onChange={(e) => setRegForm({ ...regForm, confirmPassword: e.target.value })}
                    className={`w-full bg-gray-800 border rounded-lg py-3 pl-10 pr-10 text-white placeholder-gray-500 focus:outline-none focus:ring-2 ${
                      regForm.confirmPassword && regForm.password !== regForm.confirmPassword
                        ? 'border-red-700 focus:ring-red-500'
                        : 'border-gray-700 focus:ring-blue-500'
                    }`}
                    placeholder="••••••••"
                  />
                  <button
                    type="button"
                    onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                    disabled={isLoading}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-300 disabled:opacity-50 disabled:cursor-not-allowed"
                    aria-label={showConfirmPassword ? "Скрыть пароль" : "Показать пароль"}
                  >
                    {showConfirmPassword ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                  </button>
                </div>
                {regForm.password && regForm.confirmPassword && (
                  <div className="mt-2 flex items-center gap-2 text-xs">
                    {regForm.password === regForm.confirmPassword ? (
                      <>
                        <div className="w-4 h-4 rounded-full bg-emerald-900/30 border border-emerald-700 flex items-center justify-center">
                          <div className="w-2 h-2 bg-emerald-400 rounded-full" />
                        </div>
                        <span className="text-emerald-300">Пароли совпадают</span>
                      </>
                    ) : (
                      <>
                        <div className="w-4 h-4 rounded-full bg-red-900/30 border border-red-700 flex items-center justify-center text-red-400">
                          ✕
                        </div>
                        <span className="text-red-300">Пароли не совпадают</span>
                      </>
                    )}
                  </div>
                )}
              </div>

              <button
                type="submit"
                disabled={isLoading || !isPasswordValid || !regForm.email || !regForm.firstName || !regForm.lastName || regForm.password !== regForm.confirmPassword || !regForm.confirmPassword}
                className="w-full py-3 rounded-lg font-semibold text-white bg-linear-to-br from-blue-500 to-emerald-500 shadow-lg shadow-blue-500/30 flex items-center justify-center gap-2 disabled:opacity-60 disabled:cursor-not-allowed"
              >
                {isLoading ? <Loader2 className="w-5 h-5 animate-spin" /> : <UserPlus className="w-5 h-5" />}
                Создать аккаунт
              </button>
            </form>
          )}

          {message && (
            <div className={`mt-4 p-3 rounded-lg border text-sm ${
              message.includes('✓') || message.includes('Аккаунт создан')
                ? 'border-green-700 bg-green-900/30 text-green-200'
                : 'border-red-700 bg-red-900/30 text-red-200'
            }`}>
              {message}
            </div>
          )}
        </div>

        {/* Right panel */}
        <div className="bg-gray-900/40 border border-gray-800 rounded-2xl p-8 shadow-xl backdrop-blur-sm">
          {/* Test credentials info */}
          <div className="mb-6 p-4 rounded-lg border border-blue-700/50 bg-blue-900/20">
            <div className="flex items-start gap-3">
              <div className="text-blue-400 mt-1">ℹ️</div>
              <div>
                <p className="text-sm font-semibold text-blue-300 mb-2">Тестовые учётные данные</p>
                <p className="text-xs text-gray-400 mb-2">Для быстрого тестирования используйте:</p>
                <code className="text-xs bg-gray-800 p-2 rounded block text-gray-300 mb-2">
                  Email: newadmin@example.com<br/>
                  Пароль: SecurePass123!
                </code>
                <p className="text-xs text-gray-500">Или создайте новый аккаунт в разделе "Регистрация"</p>
              </div>
            </div>
          </div>

          <div className="flex items-center gap-3 mb-4">
            <div className="w-10 h-10 bg-gray-800 rounded-lg flex items-center justify-center border border-gray-700">
              <Lock className="w-5 h-5 text-blue-400" />
            </div>
            <div>
              <p className="text-xs uppercase tracking-wide text-gray-500">Безопасность</p>
              <h3 className="text-xl font-semibold text-white">Восстановление доступа</h3>
            </div>
          </div>

          <div className="space-y-3">
            <p className="text-sm text-gray-400 mb-4">
              Если вы забыли пароль, нажмите на одну из кнопок ниже для восстановления доступа.
            </p>
              <button
                type="button"
                onClick={() => setShowForgotModal(true)}
                className="w-full py-3 rounded-lg font-semibold text-white bg-gray-800 border border-gray-700 hover:border-blue-500 transition-colors flex items-center justify-center gap-2 disabled:opacity-60"
              >
                <Mail className="w-5 h-5" />
                Восстановить пароль
              </button>
              <button
                type="button"
                onClick={() => setShowResetModal(true)}
                className="w-full py-3 rounded-lg font-semibold text-white bg-gray-800 border border-gray-700 hover:border-emerald-500 transition-colors flex items-center justify-center gap-2 disabled:opacity-60"
              >
                <KeyRound className="w-5 h-5" />
                Использовать код восстановления
              </button>
          </div>

          <div className="mt-8 text-gray-400 text-sm space-y-2">
            <p>Доступ по ролям: администратор, менеджер, кассир.</p>
            <p>Токены обновляются через /users/refresh, сессия сохраняется в Zustand.</p>
            <p>Все запросы идут к {process.env.NEXT_PUBLIC_API_URL || 'http://localhost/api/v1'}.</p>
          </div>
        </div>
      </div>

      {/* Modals */}
      <ForgotPasswordModal
        isOpen={showForgotModal}
        onClose={() => setShowForgotModal(false)}
      />
      <ResetPasswordModal
        isOpen={showResetModal}
        onClose={() => setShowResetModal(false)}
      />
    </div>
  );
}
