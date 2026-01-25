import BulkUserImport from '@/features/users/components/BulkUserImport';
import Link from 'next/link';
import { ArrowLeft } from 'lucide-react';

export const metadata = {
  title: 'Импорт пользователей | KKM',
  description: 'Загрузите пользователей из CSV файла',
};

export default function ImportUsersPage() {
  return (
    <div className="min-h-screen bg-gray-950 p-4 md:p-8">
      <div className="w-full max-w-2xl mx-auto">
        {/* Breadcrumb */}
        <Link
          href="/users"
          className="flex items-center gap-2 text-blue-400 hover:text-blue-300 mb-6 transition-colors"
        >
          <ArrowLeft className="w-4 h-4" />
          Вернуться к списку пользователей
        </Link>

        {/* Card */}
        <div className="bg-gray-900 border border-gray-800 rounded-lg shadow-lg">
          {/* Header */}
          <div className="p-6 border-b border-gray-800">
            <h1 className="text-2xl font-bold text-white">Импорт пользователей</h1>
            <p className="text-gray-400 text-sm mt-1">
              Загрузите CSV файл с данными новых пользователей для создания их в массовом порядке
            </p>
          </div>

          {/* Content */}
          <div className="p-6">
            <BulkUserImport />
          </div>
        </div>
      </div>
    </div>
  );
}
