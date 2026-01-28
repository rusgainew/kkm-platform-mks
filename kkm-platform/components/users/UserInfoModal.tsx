'use client';

import { Mail, User, Shield, Clock, LockKeyhole, CheckCircle2 } from 'lucide-react';
import type { ApiUser } from '@/lib/api/users';
import { Modal, ModalButton } from '@/components/ui';

interface UserInfoModalProps {
  user: ApiUser | null;
  isOpen: boolean;
  onClose: () => void;
}

const getRolePermissions = (role: string): string[] => {
  const permissions: Record<string, string[]> = {
    manager: [
      'Просмотр пользователей',
      'Управление инвентарем',
      'Просмотр отчетов',
      'Работа с одним магазином',
    ],
    user: [
      'Просмотр пользователей',
      'Базовые операции',
    ],
    admin: [
      'Полный доступ',
      'Управление всеми пользователями',
      'Управление магазинами',
      'Просмотр отчетов',
      'Управление инвентарем',
      'Управление сотрудниками',
    ],
    cashier: [
      'Просмотр пользователей',
      'Работа с кассой',
    ],
  };
  return permissions[role] || [];
};

const getRoleLabel = (role: string) => {
  const labels: Record<string, string> = {
    user: 'Пользователь',
    manager: 'Менеджер',
    admin: 'Администратор',
    cashier: 'Кассир',
    viewer: 'Просмотр',
  };
  return labels[role] || role;
};

const getRoleBadgeColor = (role: string) => {
  const colors: Record<string, string> = {
    user: 'bg-gray-800/50 text-gray-300 border border-gray-700',
    manager: 'bg-blue-900/20 text-blue-300 border border-blue-800',
    admin: 'bg-red-900/20 text-red-300 border border-red-800',
    cashier: 'bg-green-900/20 text-green-300 border border-green-800',
    viewer: 'bg-gray-800/50 text-gray-300 border border-gray-700',
  };
  return colors[role] || colors.viewer;
};

const getStatusBadge = (isActive: boolean) => {
  if (isActive) {
    return {
      label: 'Активен',
      color: 'bg-emerald-900/20 text-emerald-300 border border-emerald-800',
    };
  }
  return {
    label: 'Неактивен',
    color: 'bg-gray-800/50 text-gray-300 border border-gray-700',
  };
};

export default function UserInfoModal({ user, isOpen, onClose }: UserInfoModalProps) {
  if (!isOpen || !user) return null;

  const statusBadge = getStatusBadge(user.is_active ?? true);
  const createdTs = typeof user.created_at === 'string' ? parseInt(user.created_at) : user.created_at;
  const updatedTs = typeof user.updated_at === 'string' ? parseInt(user.updated_at) : user.updated_at;
  const createdDate = new Date(createdTs * 1000).toLocaleDateString('ru-RU', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
  const updatedDate = new Date(updatedTs * 1000).toLocaleDateString('ru-RU', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title="Информация о пользователе"
      size="xl"
      footer={
        <ModalButton onClick={onClose} variant="primary" className="w-full">
          Закрыть
        </ModalButton>
      }
    >
      <div className="space-y-8 max-h-[60vh] overflow-y-auto">
          {/* Avatar and Basic Info */}
          <div className="flex items-start gap-6">
            <div className="w-24 h-24 rounded-full bg-linear-to-br from-blue-500 to-emerald-500 flex items-center justify-center shrink-0">
              <span className="text-white font-bold text-4xl">
                {user.first_name.charAt(0)}{user.last_name.charAt(0)}
              </span>
            </div>
            <div className="flex-1">
              <h3 className="text-2xl font-bold text-white mb-2">
                {user.first_name} {user.last_name}
              </h3>
              <p className="text-gray-400 mb-4">{user.email}</p>
              <div className="flex gap-3 flex-wrap">
                <span className={`inline-flex px-3 py-1 rounded-full text-sm font-semibold ${getRoleBadgeColor(user.role)}`}>
                  {getRoleLabel(user.role)}
                </span>
                <span className={`inline-flex px-3 py-1 rounded-full text-sm font-semibold ${statusBadge.color}`}>
                  {statusBadge.label}
                </span>
              </div>
            </div>
          </div>

          {/* Divider */}
          <div className="border-t border-gray-800" />

          {/* Contact Information */}
          <div>
            <h4 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
              <Mail className="w-5 h-5 text-blue-400" />
              Контактные данные
            </h4>
            <div className="space-y-3 ml-7">
              <div>
                <p className="text-sm text-gray-400">Email</p>
                <p className="text-white font-medium">{user.email}</p>
              </div>
              <div>
                <p className="text-sm text-gray-400">ID пользователя</p>
                <p className="text-white font-mono text-sm break-all">{user.id}</p>
              </div>
            </div>
          </div>

          {/* Personal Information */}
          <div>
            <h4 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
              <User className="w-5 h-5 text-emerald-400" />
              Персональные данные
            </h4>
            <div className="space-y-3 ml-7">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-sm text-gray-400">Имя</p>
                  <p className="text-white font-medium">{user.first_name || '—'}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-400">Фамилия</p>
                  <p className="text-white font-medium">{user.last_name || '—'}</p>
                </div>
              </div>
            </div>
          </div>

          {/* Role & Status */}
          <div>
            <h4 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
              <Shield className="w-5 h-5 text-purple-400" />
              Роль и статус
            </h4>
            <div className="space-y-3 ml-7">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-sm text-gray-400">Роль</p>
                  <p className="text-white font-medium">{getRoleLabel(user.role)}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-400">Статус</p>
                  <p className={`font-medium ${user.is_active ? 'text-emerald-400' : 'text-gray-400'}`}>
                    {user.is_active ? 'Активен' : 'Неактивен'}
                  </p>
                </div>
              </div>
            </div>
          </div>

          {/* Dates */}
          <div>
            <h4 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
              <Clock className="w-5 h-5 text-orange-400" />
              Даты
            </h4>
            <div className="space-y-3 ml-7">
              <div>
                <p className="text-sm text-gray-400">Дата создания</p>
                <p className="text-white font-medium">{createdDate}</p>
              </div>
              <div>
                <p className="text-sm text-gray-400">Последнее обновление</p>
                <p className="text-white font-medium">{updatedDate}</p>
              </div>
            </div>
          </div>

          {/* Divider */}
          <div className="border-t border-gray-800" />

          {/* Permissions */}
          <div>
            <h4 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
              <LockKeyhole className="w-5 h-5 text-cyan-400" />
              Разрешения и доступ
            </h4>
            <div className="ml-7">
              <div className="mb-4">
                <p className="text-sm text-gray-400 mb-2">Доступные операции:</p>
                <div className="space-y-2">
                  {getRolePermissions(user.role).map((permission, idx) => (
                    <div key={idx} className="flex items-center gap-2">
                      <CheckCircle2 className="w-4 h-4 text-emerald-400 shrink-0" />
                      <span className="text-sm text-gray-300">{permission}</span>
                    </div>
                  ))}
                </div>
              </div>
              <div className="pt-4 border-t border-gray-800">
                <p className="text-xs text-gray-500">
                  Разрешения определяются в соответствии с назначенной ролью пользователя
                </p>
              </div>
            </div>
          </div>
        </div>
    </Modal>
  );
}
