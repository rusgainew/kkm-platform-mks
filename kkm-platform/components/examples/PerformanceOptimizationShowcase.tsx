'use client';

import { useEffect, useRef, useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import {
  useDebounce,
  useThrottle,
  useIntersectionObserver,
} from '@/lib/performance';
import { apiCache } from '@/lib/cache';
import { LazyImage } from '@/components/optimizations/LazyImage';
import { VirtualTable } from '@/components/tables/VirtualTable';

/**
 * Полный пример всех оптимизаций вместе
 */
export function PerformanceOptimizationShowcase() {
  // 1. Состояние для демонстрации
  const [search, setSearch] = useState('');
  const [scrollPos, setScrollPos] = useState(0);
  const containerRef = useRef<HTMLDivElement>(null);

  // 2. Debounce для поиска
  const debouncedSearch = useDebounce(search, 300);

  // 3. Throttle для скролла
  const throttledScroll = useThrottle(scrollPos, 100);

  // 4. Lazy loading с Intersection Observer
  const chartsRef = useRef<HTMLDivElement>(null);
  const chartsVisible = useIntersectionObserver(chartsRef as React.RefObject<HTMLDivElement>, {
    rootMargin: '100px',
  });

  // 5. React Query с автокешем
  const { data: searchResults = [] } = useQuery({
    queryKey: ['users-search', debouncedSearch],
    queryFn: async () => {
      // Проверить кеш
      const cached = apiCache.get(`users-search-${debouncedSearch}`);
      if (cached) return cached;

      // API запрос
      const response = await fetch(
        `/api/users?search=${debouncedSearch}`
      );
      const data = await response.json();

      // Сохранить в кеш на 5 минут
      apiCache.set(`users-search-${debouncedSearch}`, data);
      return data;
    },
    enabled: debouncedSearch.length > 2,
    staleTime: 5 * 60 * 1000, // 5 минут
  });

  // 6. Обработка скролла
  const handleContainerScroll = (e: React.UIEvent<HTMLDivElement>) => {
    const target = e.target as HTMLDivElement;
    setScrollPos(target.scrollTop);
  };

  return (
    <div className="max-w-6xl mx-auto space-y-8 p-6">
      <h1 className="text-3xl font-bold">Performance Optimization Showcase</h1>

      {/* 1. Debounce пример */}
      <section className="space-y-4">
        <h2 className="text-2xl font-bold">1. Debounced Search</h2>
        <p className="text-gray-600">
          Поиск с задержкой 300ms - API запрос отправляется только после того
          как пользователь прекратит вводить. Результат кешируется на 5 минут.
        </p>
        <input
          type="text"
          placeholder="Поиск пользователей... (debounced)"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="w-full px-4 py-2 border rounded-lg"
        />
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {searchResults.map((result: any) => (
            <div
              key={result.id}
              className="p-4 border rounded-lg hover:shadow-lg transition"
            >
              <div className="font-semibold">{result.name}</div>
              <div className="text-sm text-gray-600">{result.email}</div>
            </div>
          ))}
        </div>
      </section>

      {/* 2. Throttle пример */}
      <section className="space-y-4">
        <h2 className="text-2xl font-bold">2. Throttled Scroll Position</h2>
        <p className="text-gray-600">
          Позиция скролла обновляется максимум 1 раз в 100ms. Полезно для
          отслеживания скролла без перегрузки обработчиков.
        </p>
        <div className="text-sm text-gray-600">
          Текущая позиция: {throttledScroll.toFixed(0)}px
        </div>
      </section>

      {/* 3. Lazy Image Loading */}
      <section className="space-y-4">
        <h2 className="text-2xl font-bold">3. Lazy Image Loading</h2>
        <p className="text-gray-600">
          Изображения загружаются только когда они видны на экране или на расстоянии
          50px. Качество оптимизировано (75%).
        </p>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {[
            'https://via.placeholder.com/400x300?text=Image+1',
            'https://via.placeholder.com/400x300?text=Image+2',
            'https://via.placeholder.com/400x300?text=Image+3',
          ].map((src, i) => (
            <LazyImage
              key={i}
              src={src}
              alt={`Demo image ${i + 1}`}
              width={400}
              height={300}
              className="rounded-lg"
            />
          ))}
        </div>
      </section>

      {/* 4. Intersection Observer для ленивой загрузки компонентов */}
      <section className="space-y-4">
        <h2 className="text-2xl font-bold">4. Lazy Component Loading</h2>
        <p className="text-gray-600">
          Этот раздел загружается только когда он видна на экране (intersectionObserver).
        </p>
        <div ref={chartsRef} className="border rounded-lg p-6">
          {chartsVisible ? (
            <div className="space-y-4">
              <div className="h-64 bg-gradient-to-r from-blue-400 to-blue-600 rounded-lg flex items-center justify-center text-white font-bold">
                Диаграмма доходов (загружена при скролле)
              </div>
              <div className="h-48 bg-gradient-to-r from-green-400 to-green-600 rounded-lg flex items-center justify-center text-white font-bold">
                Диаграмма статусов (загружена при скролле)
              </div>
            </div>
          ) : (
            <div className="h-64 bg-gray-200 rounded-lg flex items-center justify-center text-gray-500">
              Прокрути до этого раздела чтобы загрузить компоненты...
            </div>
          )}
        </div>
      </section>

      {/* 5. Virtual Table */}
      <section className="space-y-4">
        <h2 className="text-2xl font-bold">5. Virtual Table (10000 строк)</h2>
        <p className="text-gray-600">
          Виртуализированная таблица отрендеривает только видимые строки.
          Эффективно работает с 10000+ строк без фризов.
        </p>
        <VirtualTable
          data={Array.from({ length: 10000 }, (_, i) => ({
            id: `${i}`,
            product: `Товар ${i + 1}`,
            price: Math.floor(Math.random() * 50000),
            stock: Math.floor(Math.random() * 1000),
            status: ['В наличии', 'Заканчивается', 'Нет'][
              Math.floor(Math.random() * 3)
            ],
          }))}
          itemHeight={40}
          containerHeight={400}
          columns={[
            { key: 'product', label: 'Товар', width: '40%' },
            { key: 'price', label: 'Цена', width: '20%' },
            { key: 'stock', label: 'Склад', width: '20%' },
            { key: 'status', label: 'Статус', width: '20%' },
          ]}
          renderCell={(item: any, key) => {
            if (key === 'price') {
              return <span className="font-semibold">₽{item.price}</span>;
            }
            if (key === 'status') {
              const colors: Record<string, string> = {
                'В наличии': 'text-green-600',
                Заканчивается: 'text-yellow-600',
                Нет: 'text-red-600',
              };
              return (
                <span className={`font-medium ${colors[item.status] || ''}`}>
                  {item.status}
                </span>
              );
            }
            return <>{item[key]}</>;
          }}
        />
      </section>

      {/* 6. Performance Stats */}
      <section className="space-y-4">
        <h2 className="text-2xl font-bold">6. Cache Stats</h2>
        <p className="text-gray-600">
          Информация о состоянии кеша API запросов
        </p>
        <div className="bg-gray-100 p-4 rounded-lg font-mono text-sm">
          <div>Cache размер: {apiCache.size()} записей</div>
          <div>Кешировано: {searchResults.length > 0 ? 'Да' : 'Нет'}</div>
          <div>Последний поиск: {debouncedSearch || '(пусто)'}</div>
        </div>
      </section>

      {/* Tips */}
      <section className="bg-blue-50 border border-blue-200 rounded-lg p-6">
        <h3 className="font-bold text-lg mb-2">💡 Tips for Performance:</h3>
        <ul className="space-y-2 text-sm">
          <li>✅ Используй debounce для поиска и фильтров</li>
          <li>✅ Используй throttle для частых событий (скролл, resize)</li>
          <li>✅ Кеш API результаты с React Query или локальным кешем</li>
          <li>✅ Ленивую загрузку для изображений и тяжелых компонентов</li>
          <li>✅ Виртуализацию для больших списков (10000+ элементов)</li>
          <li>✅ Dynamic imports для кода разделения на chunks</li>
          <li>✅ Мемоизацию для дорогих вычислений (useMemo)</li>
          <li>✅ Code splitting в next.config для оптимального bundling</li>
        </ul>
      </section>
    </div>
  );
}
