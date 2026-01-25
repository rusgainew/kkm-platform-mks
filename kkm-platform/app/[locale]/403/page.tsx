'use client';

import Link from 'next/link';
import { ShieldAlert } from 'lucide-react';

export default function Forbidden() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-800 to-slate-900 flex items-center justify-center px-4">
      <div className="text-center max-w-md">
        <div className="mb-6 flex justify-center">
          <ShieldAlert className="w-24 h-24 text-red-500" />
        </div>
        
        <h1 className="text-5xl font-bold text-white mb-2">403</h1>
        <h2 className="text-2xl font-semibold text-gray-200 mb-4">Доступ запрещен</h2>
        <p className="text-gray-400 mb-8">
          У вас нет прав доступа к этому ресурсу. Пожалуйста, обратитесь к администратору.
        </p>
        
        <div className="flex gap-4 justify-center">
          <Link
            href="/dashboard"
            className="px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white font-semibold rounded-lg transition-colors"
          >
            На главную
          </Link>
          <Link
            href="/"
            className="px-6 py-3 bg-gray-700 hover:bg-gray-600 text-white font-semibold rounded-lg transition-colors"
          >
            На главную страницу
          </Link>
        </div>
      </div>
    </div>
  );
}
