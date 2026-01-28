'use client';

import { useAuthStore } from '@/store/authStore';
import { ROLE_CONFIGS } from '@/types/entities';
import { Shield, ChevronDown } from 'lucide-react';
import { useState } from 'react';

export default function RoleSelector() {
  const user = useAuthStore((state) => state.user);
  const switchRole = useAuthStore((state) => state.switchRole);
  const [isOpen, setIsOpen] = useState(false);

  if (!user) return null;

  const currentRoleConfig = ROLE_CONFIGS[user.role];

  return (
    <div className="relative">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center gap-2 px-4 py-2 bg-gray-800 hover:bg-gray-700 rounded-lg border border-gray-700 transition-colors"
      >
        <Shield className="w-5 h-5 text-blue-400" />
        <div className="text-left">
          <p className="text-sm font-semibold text-white">{currentRoleConfig.label}</p>
          <p className="text-xs text-gray-400">{user.name}</p>
        </div>
        <ChevronDown className={`w-4 h-4 text-gray-400 transition-transform ${isOpen ? 'rotate-180' : ''}`} />
      </button>

      {isOpen && (
        <div className="absolute top-full mt-2 right-0 w-64 bg-gray-800 border border-gray-700 rounded-lg shadow-2xl z-50">
          <div className="p-2">
            <p className="text-xs font-semibold text-gray-400 uppercase px-3 py-2">Сменить роль (demo)</p>
            {Object.values(ROLE_CONFIGS).map((roleConfig) => (
              <button
                key={roleConfig.role}
                onClick={() => {
                  switchRole(roleConfig.role as any);
                  setIsOpen(false);
                }}
                className={`w-full text-left px-3 py-2 rounded-lg transition-colors ${
                  user.role === roleConfig.role
                    ? 'bg-blue-500/20 text-blue-400'
                    : 'text-gray-300 hover:bg-gray-700'
                }`}
              >
                <p className="font-semibold">{roleConfig.label}</p>
                <p className="text-xs text-gray-400 mt-1">
                  {roleConfig.permissions.length} разрешений
                </p>
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
