'use client';

import React, { useState } from 'react';
import { FileUp, Upload, File, Trash2, Download } from 'lucide-react';

interface Document {
  id: string;
  name: string;
  size: number;
  uploaded_at: string;
  file_type: string;
}

export default function DocumentsUpload() {
  const [documents] = useState<Document[]>([]);
  const [isDragging, setIsDragging] = useState(false);
  const [isUploading, setIsUploading] = useState(false);

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = () => {
    setIsDragging(false);
  };

  const handleDrop = async (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    // Handle file drop
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i];
  };

  return (
    <div className="bg-gray-900 rounded-lg border border-gray-800 overflow-hidden">
      <div className="p-6 border-b border-gray-800">
        <h2 className="text-2xl font-bold text-white flex items-center gap-2">
          <FileUp className="w-6 h-6 text-indigo-400" />
          Управление документами
        </h2>
      </div>

      <div className="p-6 border-b border-gray-800">
        <div
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
          className={`border-2 border-dashed rounded-lg p-12 text-center transition-colors ${
            isDragging
              ? 'border-indigo-500 bg-indigo-900/10'
              : 'border-gray-700 bg-gray-800/30 hover:border-gray-600'
          }`}
        >
          <Upload className="mx-auto mb-4 text-indigo-400" size={48} />
          <h3 className="text-lg font-semibold text-white mb-2">
            Загрузите документы
          </h3>
          <p className="text-gray-400 mb-4">
            Перетащите файлы сюда или нажмите для выбора
          </p>
          <input
            type="file"
            multiple
            disabled={isUploading}
            className="hidden"
            id="file-upload"
          />
          <label
            htmlFor="file-upload"
            className="inline-block px-6 py-2 bg-indigo-600 text-white rounded-lg cursor-pointer hover:bg-indigo-700 transition-colors disabled:opacity-50"
          >
            {isUploading ? 'Загрузка...' : 'Выбрать файлы'}
          </label>
        </div>
      </div>

      {documents.length > 0 && (
        <div>
          <div className="px-6 py-4 bg-gray-800/50 border-b border-gray-800">
            <h3 className="text-sm font-semibold text-gray-300">
              Загруженные документы ({documents.length})
            </h3>
          </div>
          
          <div className="divide-y divide-gray-800">
            {documents.map((doc) => (
              <div key={doc.id} className="p-6 flex items-center justify-between hover:bg-gray-800/50 transition-colors">
                <div className="flex items-center gap-4">
                  <File className="text-gray-400" size={24} />
                  <div>
                    <p className="text-white font-medium">{doc.name}</p>
                    <p className="text-gray-400 text-sm">
                      {formatFileSize(doc.size)} • {new Date(doc.uploaded_at).toLocaleDateString('ru-RU')}
                    </p>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <button className="p-2 hover:bg-gray-700 text-gray-400 rounded transition-colors">
                    <Download size={18} />
                  </button>
                  <button className="p-2 hover:bg-red-600/20 text-red-400 rounded transition-colors">
                    <Trash2 size={18} />
                  </button>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {documents.length === 0 && (
        <div className="p-8 text-center text-gray-400">
          Документы еще не загружены
        </div>
      )}
    </div>
  );
}
