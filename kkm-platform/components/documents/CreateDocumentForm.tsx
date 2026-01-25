"use client";

import React, { useState } from "react";
import { createDocument } from "@/lib/api/documents";
import { useCompanies } from "@/hooks/useCompanies";
import { useAuthStore } from "@/store/authStore";
import { X, Loader2 } from "lucide-react";

interface CreateDocumentFormProps {
  onClose: () => void;
  onSuccess: () => void;
}

export default function CreateDocumentForm({
  onClose,
  onSuccess,
}: CreateDocumentFormProps) {
  const [isHydrated, setIsHydrated] = useState(false);
  
  // Get current user from auth store
  const user = useAuthStore((state) => state.user);
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  
  // Log user info for debugging and handle hydration
  React.useEffect(() => {
    setIsHydrated(true);
    console.log("[CreateDocumentForm] User from store:", user);
    console.log("[CreateDocumentForm] Is authenticated:", isAuthenticated);
  }, [user, isAuthenticated]);
  
  // Load companies from API
  const { companies, isLoading: companiesLoading } = useCompanies();

  const [formData, setFormData] = useState({
    title: "",
    content: "",
    organization_id: "",
  });

  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  // Set first company as default when companies load
  React.useEffect(() => {
    if (companies.length > 0 && !formData.organization_id) {
      setFormData((prev) => ({
        ...prev,
        organization_id: companies[0].id,
      }));
    }
  }, [companies]);

  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setSuccess("");
    setIsLoading(true);

    try {
      console.log("[CreateDocumentForm] handleSubmit - user from store:", user);
      console.log("[CreateDocumentForm] handleSubmit - user?.id:", user?.id);
      
      if (!user?.id) {
        console.warn("[CreateDocumentForm] Пользователь не авторизован");
        throw new Error("Требуется аутентификация");
      }

      if (!formData.title.trim()) {
        throw new Error("Введите название документа");
      }

      if (!formData.organization_id) {
        throw new Error("Выберите компанию");
      }

      const payload = {
        title: formData.title,
        content: formData.content,
        organization_id: formData.organization_id,
        created_by: user.id,
      };

      console.log("[CreateDocumentForm] Submitting:", payload);

      const result = await createDocument(payload);
      console.log("[CreateDocumentForm] Success:", result);

      setSuccess("Документ успешно создан!");
      setFormData({
        title: "",
        content: "",
        organization_id: "company-1",
      });

      setTimeout(() => {
        onSuccess();
        onClose();
      }, 1500);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Ошибка создания";
      console.error("[CreateDocumentForm] Error:", errorMessage);
      console.error("[CreateDocumentForm] Error details - user:", user);
      console.error("[CreateDocumentForm] Error details - isAuthenticated:", isAuthenticated);
      console.error("[CreateDocumentForm] Error details - isHydrated:", isHydrated);
      setError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-md mx-4">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b">
          <h2 className="text-xl font-semibold text-gray-900">
            Создать новый документ
          </h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 transition"
          >
            <X size={24} />
          </button>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="p-6 space-y-4">
          {/* Error Message */}
          {error && (
            <div className="bg-red-50 border border-red-200 rounded-lg p-3">
              <p className="text-sm text-red-700">{error}</p>
            </div>
          )}

          {/* Success Message */}
          {success && (
            <div className="bg-green-50 border border-green-200 rounded-lg p-3">
              <p className="text-sm text-green-700">{success}</p>
            </div>
          )}

          {/* Title */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Название документа *
            </label>
            <input
              type="text"
              name="title"
              value={formData.title}
              onChange={handleChange}
              placeholder="Введите название"
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition"
              disabled={isLoading}
              required
            />
          </div>

          {/* Content */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Содержание документа
            </label>
            <textarea
              name="content"
              value={formData.content}
              onChange={handleChange}
              placeholder="Введите содержание документа (опционально)"
              rows={3}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition resize-none"
              disabled={isLoading}
            />
          </div>

          {/* Organization */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Организация *
            </label>
            <select
              name="organization_id"
              value={formData.organization_id}
              onChange={handleChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition"
              disabled={isLoading || companiesLoading}
            >
              <option value="">
                {companiesLoading ? "Загрузка..." : "Выберите организацию"}
              </option>
              {companies.map((company) => (
                <option key={company.id} value={company.id}>
                  {company.name}
                </option>
              ))}
            </select>
          </div>

          {/* Debug Info */}
          <div className="bg-gray-100 p-2 rounded text-xs text-gray-600 space-y-1">
            <p>Пользователь: {user?.name || "не загружен"}</p>
            <p>ID: {user?.id || "нет"}</p>
            <p>Аутентифицирован: {isAuthenticated ? "да" : "нет"}</p>
            <p>Гидрирован: {isHydrated ? "да" : "нет"}</p>
          </div>

          {/* Buttons */}
          <div className="flex gap-3 pt-4">
            <button
              type="button"
              onClick={onClose}
              disabled={isLoading}
              className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition disabled:opacity-50"
            >
              Отмена
            </button>
            <button
              type="submit"
              disabled={isLoading}
              className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition disabled:opacity-50 flex items-center justify-center gap-2"
            >
              {isLoading && <Loader2 size={16} className="animate-spin" />}
              {isLoading ? "Создание..." : "Создать"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
