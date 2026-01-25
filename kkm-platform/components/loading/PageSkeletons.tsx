'use client';

import React from 'react';
import { SkeletonCard, SkeletonTable, SkeletonList, SkeletonChart } from './Skeleton';

/**
 * Skeleton для страницы компаний
 */
export const CompaniesPageSkeleton: React.FC = () => (
  <div className="space-y-6">
    {/* Заголовок и поиск */}
    <div className="flex justify-between items-center">
      <div className="w-48 h-8 rounded bg-gray-800 animate-pulse" />
      <div className="w-64 h-10 rounded bg-gray-800 animate-pulse" />
    </div>

    {/* Фильтры */}
    <div className="flex gap-4">
      <div className="w-32 h-10 rounded bg-gray-800 animate-pulse" />
      <div className="w-32 h-10 rounded bg-gray-800 animate-pulse" />
      <div className="w-32 h-10 rounded bg-gray-800 animate-pulse" />
    </div>

    {/* Таблица компаний */}
    <SkeletonTable rows={6} columns={5} />
  </div>
);

/**
 * Skeleton для страницы пользователей
 */
export const UsersPageSkeleton: React.FC = () => (
  <div className="space-y-6">
    {/* Заголовок и кнопка создания */}
    <div className="flex justify-between items-center">
      <div className="w-48 h-8 rounded bg-gray-800 animate-pulse" />
      <div className="w-40 h-10 rounded bg-gray-800 animate-pulse" />
    </div>

    {/* Фильтры */}
    <div className="flex gap-4">
      <div className="w-40 h-10 rounded bg-gray-800 animate-pulse" />
      <div className="w-40 h-10 rounded bg-gray-800 animate-pulse" />
    </div>

    {/* Список пользователей */}
    <SkeletonList count={8} />
  </div>
);

/**
 * Skeleton для дашборда
 */
export const DashboardSkeleton: React.FC = () => (
  <div className="space-y-6">
    {/* Статистика карточки */}
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      {Array.from({ length: 4 }, (_, i) => (
        <div
          key={i}
          className="p-6 border border-gray-800 rounded-lg bg-gray-900"
        >
          <div className="w-24 h-4 rounded bg-gray-800 animate-pulse mb-3" />
          <div className="w-32 h-8 rounded bg-gray-800 animate-pulse mb-2" />
          <div className="w-20 h-4 rounded bg-gray-800 animate-pulse" />
        </div>
      ))}
    </div>

    {/* Графики */}
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <SkeletonChart height={250} />
      <SkeletonChart height={250} />
    </div>

    {/* Таблица терминалов */}
    <div className="border border-gray-800 rounded-lg p-4 bg-gray-900">
      <div className="w-40 h-6 rounded bg-gray-800 animate-pulse mb-4" />
      <SkeletonTable rows={5} columns={4} />
    </div>
  </div>
);

/**
 * Skeleton для компании карточки
 */
export const CompanySkeleton: React.FC = () => (
  <div className="p-4 border border-gray-800 rounded-lg bg-gray-900">
    {/* Название */}
    <div className="w-3/4 h-6 rounded bg-gray-800 animate-pulse mb-4" />

    {/* Информация */}
    <div className="space-y-3 mb-4">
      <div className="flex justify-between">
        <div className="w-24 h-4 rounded bg-gray-800 animate-pulse" />
        <div className="w-32 h-4 rounded bg-gray-800 animate-pulse" />
      </div>
      <div className="flex justify-between">
        <div className="w-24 h-4 rounded bg-gray-800 animate-pulse" />
        <div className="w-32 h-4 rounded bg-gray-800 animate-pulse" />
      </div>
      <div className="flex justify-between">
        <div className="w-24 h-4 rounded bg-gray-800 animate-pulse" />
        <div className="w-32 h-4 rounded bg-gray-800 animate-pulse" />
      </div>
    </div>

    {/* Кнопки действия */}
    <div className="flex gap-2">
      <div className="flex-1 h-10 rounded bg-gray-800 animate-pulse" />
      <div className="flex-1 h-10 rounded bg-gray-800 animate-pulse" />
    </div>
  </div>
);

/**
 * Skeleton для POS страницы
 */
export const POSSkeleton: React.FC = () => (
  <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 h-screen">
    {/* Товары */}
    <div className="lg:col-span-2">
      <div className="w-32 h-6 rounded bg-gray-800 animate-pulse mb-4" />
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {Array.from({ length: 9 }, (_, i) => (
          <div
            key={i}
            className="p-4 border border-gray-800 rounded-lg bg-gray-900"
          >
            <div className="w-full h-32 rounded bg-gray-800 animate-pulse mb-3" />
            <div className="w-3/4 h-4 rounded bg-gray-800 animate-pulse mb-2" />
            <div className="w-1/2 h-4 rounded bg-gray-800 animate-pulse" />
          </div>
        ))}
      </div>
    </div>

    {/* Корзина */}
    <div className="border border-gray-800 rounded-lg p-4 bg-gray-900">
      <div className="w-20 h-6 rounded bg-gray-800 animate-pulse mb-4" />
      <div className="space-y-3 mb-4">
        {Array.from({ length: 5 }, (_, i) => (
          <div key={i} className="flex justify-between">
            <div className="w-20 h-4 rounded bg-gray-800 animate-pulse" />
            <div className="w-16 h-4 rounded bg-gray-800 animate-pulse" />
          </div>
        ))}
      </div>
      <div className="border-t border-gray-700 pt-4 mb-4">
        <div className="flex justify-between mb-2">
          <div className="w-12 h-4 rounded bg-gray-800 animate-pulse" />
          <div className="w-20 h-4 rounded bg-gray-800 animate-pulse" />
        </div>
      </div>
      <div className="h-10 rounded bg-gray-800 animate-pulse" />
    </div>
  </div>
);

/**
 * Skeleton для модального окна
 */
export const ModalSkeleton: React.FC = () => (
  <div className="space-y-4">
    {/* Заголовок */}
    <div className="w-1/2 h-6 rounded bg-gray-800 animate-pulse" />

    {/* Поля формы */}
    {Array.from({ length: 4 }, (_, i) => (
      <div key={i}>
        <div className="w-24 h-4 rounded bg-gray-800 animate-pulse mb-2" />
        <div className="w-full h-10 rounded bg-gray-800 animate-pulse" />
      </div>
    ))}

    {/* Кнопки */}
    <div className="flex gap-2 justify-end pt-4">
      <div className="w-24 h-10 rounded bg-gray-800 animate-pulse" />
      <div className="w-24 h-10 rounded bg-gray-800 animate-pulse" />
    </div>
  </div>
);

/**
 * Skeleton для карточки статистики
 */
export const StatCardSkeleton: React.FC = () => (
  <div className="p-6 border border-gray-800 rounded-lg bg-gray-900">
    <div className="space-y-3">
      <div className="w-20 h-4 rounded bg-gray-800 animate-pulse" />
      <div className="w-32 h-8 rounded bg-gray-800 animate-pulse" />
      <div className="w-24 h-4 rounded bg-gray-800 animate-pulse" />
    </div>
  </div>
);

/**
 * Skeleton для чарта
 */
export const ChartSkeleton: React.FC = () => (
  <div className="p-6 border border-gray-800 rounded-lg bg-gray-900">
    <div className="w-40 h-6 rounded bg-gray-800 animate-pulse mb-4" />
    <div className="w-full h-64 rounded bg-gray-800 animate-pulse" />
  </div>
);

/**
 * Skeleton для страницы каталога
 */
export const CatalogPageSkeleton: React.FC = () => (
  <div className="space-y-6">
    {/* Заголовок и кнопки */}
    <div className="flex justify-between items-center">
      <div className="w-48 h-8 rounded bg-gray-800 animate-pulse" />
      <div className="flex gap-2">
        <div className="w-32 h-10 rounded bg-gray-800 animate-pulse" />
        <div className="w-32 h-10 rounded bg-gray-800 animate-pulse" />
        <div className="w-32 h-10 rounded bg-gray-800 animate-pulse" />
      </div>
    </div>

    {/* Фильтры */}
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div className="h-10 rounded bg-gray-800 animate-pulse" />
      <div className="h-10 rounded bg-gray-800 animate-pulse" />
      <div className="h-10 rounded bg-gray-800 animate-pulse" />
    </div>

    {/* Таблица товаров */}
    <SkeletonTable rows={8} columns={7} />
  </div>
);
