'use client';

import React from 'react';
import BankAccountsList from '@/features/bank-accounts/components/BankAccountsList';

export default function BankAccountsPage() {
  return (
    <div className="min-h-screen bg-gray-950 py-8 px-4 sm:px-6 lg:px-8">
      <div className="max-w-7xl mx-auto">
        <BankAccountsList />
      </div>
    </div>
  );
}
