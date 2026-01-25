'use client';

import React, { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import { AlertCircle, Loader2 } from 'lucide-react';
import { getUserById } from '@/lib/api/users';
import type { ApiUser } from '@/lib/api/users';
import UserEditForm from '@/features/users/components/UserEditForm';

export default function EditUserPage() {
  const params = useParams();
  const userId = params.id as string;
  const [user, setUser] = useState<ApiUser | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!userId) return;

    const loadUser = async () => {
      setIsLoading(true);
      setError(null);

      try {
        const userData = await getUserById(userId);
        setUser(userData);
      } catch (err) {
        const message = err instanceof Error ? err.message : 'Ошибка при загрузке пользователя';
        setError(message);
        console.error('Failed to load user:', err);
      } finally {
        setIsLoading(false);
      }
    };

    loadUser();
  }, [userId]);

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-b from-gray-900 to-black flex items-center justify-center">
        <div className="flex flex-col items-center gap-4">
          <Loader2 size={40} className="text-blue-400 animate-spin" />
          <p className="text-gray-400">Загрузка данных пользователя...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gradient-to-b from-gray-900 to-black flex items-center justify-center px-4">
        <div className="bg-red-900/20 border border-red-800 rounded-lg p-6 max-w-md w-full flex items-start gap-3">
          <AlertCircle size={20} className="text-red-400 flex-shrink-0 mt-0.5" />
          <div>
            <p className="text-red-300 font-medium">Ошибка</p>
            <p className="text-red-300/70 text-sm mt-1">{error}</p>
          </div>
        </div>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="min-h-screen bg-gradient-to-b from-gray-900 to-black flex items-center justify-center px-4">
        <div className="bg-gray-800/50 border border-gray-700 rounded-lg p-6 max-w-md w-full">
          <p className="text-gray-300">Пользователь не найден</p>
        </div>
      </div>
    );
  }

  return <UserEditForm user={user} />;
}
