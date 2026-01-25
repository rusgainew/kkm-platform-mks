// Error Boundary компонент для перехвата ошибок в React приложении
export { default as ErrorBoundary } from './ErrorBoundary';

// Fallback компоненты для различных секций приложения
export {
  CompaniesErrorFallback,
  UsersErrorFallback,
  DashboardErrorFallback,
  POSErrorFallback,
} from './ErrorFallbacks';

// Утилиты для использования Error Boundary
export { SafeClientComponent, withErrorBoundary } from './SafeClientComponent';

// Тестовые компоненты для разработки (только для разработки)
export {
  ErrorTestComponent,
  EventErrorTestComponent,
  MultiErrorTestComponent,
} from './ErrorTestComponents';
