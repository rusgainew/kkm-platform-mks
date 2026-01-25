'use client';

import React from 'react';
import DocumentsUpload from '@/features/documents/components/DocumentsUpload';

export default function DocumentsUploadPage() {
  return (
    <div className="min-h-screen bg-gray-950 py-8 px-4 sm:px-6 lg:px-8">
      <div className="max-w-7xl mx-auto">
        <DocumentsUpload />
      </div>
    </div>
  );
}
