'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { Users, Plus, Trash2, AlertCircle, CheckCircle, Mail, User } from 'lucide-react';
import { Spinner } from '@/components/ui/Spinner';
import { Button } from '@/components/ui/Button';
import type { Employee, AddMemberRequest } from '@/types/entities';
import { isValidEmail } from '@/types/entities';
import {
  getCompanyMembers,
  addCompanyMember,
  removeCompanyMember,
  updateMemberRole,
} from '@/lib/api/companies';

interface CompanyMembersProps {
  companyId: string;
  companyName: string;
  isOwner?: boolean;
}

export default function CompanyMembers({ companyId, companyName, isOwner = false }: CompanyMembersProps) {
  const [members, setMembers] = useState<Employee[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isAddingMember, setIsAddingMember] = useState(false);
  const [showAddForm, setShowAddForm] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  // Form state
  const [newMemberEmail, setNewMemberEmail] = useState('');
  const [newMemberRole, setNewMemberRole] = useState<string>('employee');
  const [emailError, setEmailError] = useState<string | null>(null);

  // Load members on mount
  const loadMembers = useCallback(async () => {
    try {
      setIsLoading(true);
      setError(null);
      const response = await getCompanyMembers(companyId);
      if (response.success && response.data) {
        setMembers(response.data);
      } else {
        setError(response.error?.message || 'Не удалось загрузить список участников');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка загрузки участников');
      console.error('[CompanyMembers] Load error:', err);
    } finally {
      setIsLoading(false);
    }
  }, [companyId]);

  useEffect(() => {
    loadMembers();
  }, [loadMembers]);

  const handleAddMember = async (e: React.FormEvent) => {
    e.preventDefault();

    // Validation
    if (!newMemberEmail.trim()) {
      setEmailError('Email обязателен');
      return;
    }
    if (!isValidEmail(newMemberEmail)) {
      setEmailError('Введите корректный email');
      return;
    }

    setIsAddingMember(true);
    setError(null);
    setSuccessMessage(null);
    setEmailError(null);

    try {
      const memberData: AddMemberRequest = {
        email: newMemberEmail.trim(),
        role: newMemberRole, // Тип уже string согласно AddMemberRequest
      };

      const response = await addCompanyMember(companyId, memberData);
      
      if (response.success && response.data) {
        setMembers(prev => [...prev, response.data!]);
        setSuccessMessage(`Участник ${newMemberEmail} успешно добавлен`);
        setNewMemberEmail('');
        setNewMemberRole('employee');
        setShowAddForm(false);
        
        // Clear success message after 3 seconds
        setTimeout(() => setSuccessMessage(null), 3000);
      } else {
        setError(response.error?.message || 'Не удалось добавить участника');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка при добавлении участника');
      console.error('[CompanyMembers] Add error:', err);
    } finally {
      setIsAddingMember(false);
    }
  };

  const handleRemoveMember = async (userId: string, userName: string) => {
    if (!confirm(`Вы уверены, что хотите удалить ${userName} из компании?`)) {
      return;
    }

    setError(null);
    setSuccessMessage(null);

    try {
      const response = await removeCompanyMember(companyId, userId);
      
      if (response.success) {
        setMembers(prev => prev.filter(m => m.user_id !== userId));
        setSuccessMessage(`Участник ${userName} удален из компании`);
        
        // Clear success message after 3 seconds
        setTimeout(() => setSuccessMessage(null), 3000);
      } else {
        setError(response.error?.message || 'Не удалось удалить участника');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка при удалении участника');
      console.error('[CompanyMembers] Remove error:', err);
    }
  };

  const handleUpdateRole = async (userId: string, userName: string, newRole: string) => {
    setError(null);
    setSuccessMessage(null);

    try {
      const response = await updateMemberRole(companyId, userId, newRole);
      
      if (response.success && response.data) {
        setMembers(prev => prev.map(m => 
          m.user_id === userId ? response.data! : m
        ));
        setSuccessMessage(`Роль участника ${userName} обновлена`);
        
        // Clear success message after 3 seconds
        setTimeout(() => setSuccessMessage(null), 3000);
      } else {
        setError(response.error?.message || 'Не удалось обновить роль');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка при обновлении роли');
      console.error('[CompanyMembers] Update role error:', err);
    }
  };

  const getRoleName = (role: string): string => {
    const roleMap: Record<string, string> = {
      owner: 'Владелец',
      admin: 'Администратор',
      manager: 'Менеджер',
      employee: 'Сотрудник',
      cashier: 'Кассир',
    };
    return roleMap[role] || role;
  };

  const getRoleColor = (role: string): string => {
    const colorMap: Record<string, string> = {
      owner: 'bg-purple-500/20 text-purple-400 border-purple-500/30',
      admin: 'bg-red-500/20 text-red-400 border-red-500/30',
      manager: 'bg-blue-500/20 text-blue-400 border-blue-500/30',
      employee: 'bg-gray-500/20 text-gray-400 border-gray-500/30',
      cashier: 'bg-green-500/20 text-green-400 border-green-500/30',
    };
    return colorMap[role] || 'bg-gray-500/20 text-gray-400 border-gray-500/30';
  };

  return (
    <div className="bg-gray-800/50 border border-gray-700 rounded-lg p-6">
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-3">
          <Users size={24} className="text-blue-400" />
          <div>
            <h3 className="text-lg font-semibold text-white">Участники компании</h3>
            <p className="text-sm text-gray-400">{companyName}</p>
          </div>
        </div>
        {isOwner && !showAddForm && (
          <button
            onClick={() => setShowAddForm(true)}
            className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
          >
            <Plus size={18} />
            Добавить участника
          </button>
        )}
      </div>

      {/* Success Message */}
      {successMessage && (
        <div className="mb-4 bg-emerald-900/20 border border-emerald-800 rounded-lg p-3 flex items-start gap-2">
          <CheckCircle size={18} className="text-emerald-400 shrink-0 mt-0.5" />
          <p className="text-emerald-300 text-sm">{successMessage}</p>
        </div>
      )}

      {/* Error Message */}
      {error && (
        <div className="mb-4 bg-red-900/20 border border-red-800 rounded-lg p-3 flex items-start gap-2">
          <AlertCircle size={18} className="text-red-400 shrink-0 mt-0.5" />
          <p className="text-red-300 text-sm">{error}</p>
        </div>
      )}

      {/* Add Member Form */}
      {showAddForm && isOwner && (
        <form onSubmit={handleAddMember} className="mb-6 bg-gray-900/50 border border-gray-600 rounded-lg p-4 space-y-4">
          <h4 className="text-sm font-semibold text-white">Добавить нового участника</h4>
          
          <div>
            <label htmlFor="memberEmail" className="block text-sm font-medium text-gray-300 mb-2">
              Email <span className="text-red-400">*</span>
            </label>
            <input
              id="memberEmail"
              type="email"
              value={newMemberEmail}
              onChange={(e) => {
                setNewMemberEmail(e.target.value);
                setEmailError(null);
              }}
              disabled={isAddingMember}
              placeholder="user@example.com"
              className={`w-full px-3 py-2 rounded-lg bg-gray-800 border outline-none transition-colors ${
                emailError
                  ? 'border-red-600 focus:border-red-500'
                  : 'border-gray-600 focus:border-blue-500'
              } text-white placeholder-gray-500 disabled:opacity-50`}
            />
            {emailError && (
              <p className="text-red-400 text-xs mt-1 flex items-center gap-1">
                <AlertCircle size={12} />
                {emailError}
              </p>
            )}
          </div>

          <div>
            <label htmlFor="memberRole" className="block text-sm font-medium text-gray-300 mb-2">
              Роль
            </label>
            <select
              id="memberRole"
              value={newMemberRole}
              onChange={(e) => setNewMemberRole(e.target.value)}
              disabled={isAddingMember}
              className="w-full px-3 py-2 rounded-lg bg-gray-800 border border-gray-600 focus:border-blue-500 outline-none transition-colors text-white disabled:opacity-50"
            >
              <option value="employee">Сотрудник</option>
              <option value="cashier">Кассир</option>
              <option value="manager">Менеджер</option>
              <option value="admin">Администратор</option>
            </select>
          </div>

          <div className="flex gap-2">
            <Button
              type="submit"
              disabled={isAddingMember}
              loading={isAddingMember}
              variant="primary"
              className="flex-1"
            >
              Добавить
            </Button>
            <button
              type="button"
              onClick={() => {
                setShowAddForm(false);
                setNewMemberEmail('');
                setNewMemberRole('employee');
                setEmailError(null);
              }}
              disabled={isAddingMember}
              className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg transition-colors disabled:opacity-50"
            >
              Отмена
            </button>
          </div>
        </form>
      )}

      {/* Members List */}
      {isLoading ? (
        <div className="flex items-center justify-center py-8">
          <Spinner size="lg" label="Загрузка участников..." centered />
        </div>
      ) : members.length === 0 ? (
        <div className="text-center py-8 text-gray-400">
          <Users size={48} className="mx-auto mb-3 opacity-50" />
          <p>Участников пока нет</p>
        </div>
      ) : (
        <div className="space-y-3">
          {members.map((member) => (
            <div
              key={member.user_id}
              className="bg-gray-900/50 border border-gray-600 rounded-lg p-4 flex items-center justify-between hover:border-gray-500 transition-colors"
            >
              <div className="flex items-center gap-4 flex-1">
                <div className="w-10 h-10 rounded-full bg-blue-500/20 flex items-center justify-center">
                  <User size={20} className="text-blue-400" />
                </div>
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-1">
                    <h4 className="text-white font-medium">
                      {member.first_name} {member.last_name}
                    </h4>
                    <span className={`px-2 py-0.5 rounded-full text-xs border ${getRoleColor(member.role)}`}>
                      {getRoleName(member.role)}
                    </span>
                  </div>
                  <div className="flex items-center gap-1 text-sm text-gray-400">
                    <Mail size={14} />
                    {member.email}
                  </div>
                </div>
              </div>

              {/* Actions (only for owner) */}
              {isOwner && member.role !== 'owner' && (
                <div className="flex items-center gap-2">
                  <select
                    value={member.role}
                    onChange={(e) => handleUpdateRole(member.user_id, `${member.first_name} ${member.last_name}`, e.target.value)}
                    className="px-3 py-1.5 bg-gray-800 border border-gray-600 rounded-lg text-sm text-white hover:border-gray-500 transition-colors outline-none"
                  >
                    <option value="employee">Сотрудник</option>
                    <option value="cashier">Кассир</option>
                    <option value="manager">Менеджер</option>
                    <option value="admin">Администратор</option>
                  </select>
                  <button
                    onClick={() => handleRemoveMember(member.user_id, `${member.first_name} ${member.last_name}`)}
                    className="p-2 text-red-400 hover:bg-red-500/20 rounded-lg transition-colors"
                    title="Удалить участника"
                  >
                    <Trash2 size={18} />
                  </button>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
