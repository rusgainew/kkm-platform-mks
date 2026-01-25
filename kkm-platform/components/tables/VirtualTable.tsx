'use client';

import { useRef, useState } from 'react';
import { useVirtualization } from '@/lib/performance';

interface VirtualTableProps<T> {
  data: T[];
  itemHeight: number;
  containerHeight: number;
  columns: Array<{
    key: keyof T;
    label: string;
    width: string;
  }>;
  renderCell?: (item: T, key: keyof T) => React.ReactNode;
}

/**
 * Виртуализированная таблица для большого количества строк
 * Отрендеривает только видимые строки
 */
export function VirtualTable<T extends { id: string | number }>({
  data,
  itemHeight,
  containerHeight,
  columns,
  renderCell,
}: VirtualTableProps<T>) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [scrollOffset, setScrollOffset] = useState(0);

  const virtualization = useVirtualization(data, itemHeight, containerHeight);

  const handleScroll = (e: React.UIEvent<HTMLDivElement>) => {
    const target = e.target as HTMLDivElement;
    setScrollOffset(target.scrollTop);
    virtualization.onScroll(target.scrollTop);
  };

  return (
    <div className="border rounded-lg overflow-hidden">
      {/* Таблица */}
      <div
        ref={containerRef}
        className="overflow-y-auto"
        style={{ height: containerHeight }}
        onScroll={handleScroll}
      >
        {/* Header */}
        <div className="sticky top-0 bg-gray-50 border-b z-10">
          <div className="flex">
            {columns.map((col) => (
              <div
                key={String(col.key)}
                className="px-4 py-2 font-semibold text-sm"
                style={{ width: col.width }}
              >
                {col.label}
              </div>
            ))}
          </div>
        </div>

        {/* Virtual content */}
        <div
          style={{
            height: virtualization.totalHeight,
            position: 'relative',
          }}
        >
          <div
            style={{
              transform: `translateY(${virtualization.offsetY}px)`,
              position: 'absolute',
              top: 0,
              left: 0,
              right: 0,
            }}
          >
            {virtualization.visibleItems.map((item, index) => (
              <div
                key={item.id}
                className="flex border-b hover:bg-gray-50 transition-colors"
                style={{ height: itemHeight }}
              >
                {columns.map((col) => (
                  <div
                    key={String(col.key)}
                    className="px-4 py-2 flex items-center text-sm"
                    style={{ width: col.width }}
                  >
                    {renderCell ? (
                      renderCell(item, col.key)
                    ) : (
                      <span>{String(item[col.key])}</span>
                    )}
                  </div>
                ))}
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Footer с информацией */}
      <div className="bg-gray-50 border-t px-4 py-2 text-xs text-gray-600">
        Показано {virtualization.visibleItems.length} из {data.length} записей
      </div>
    </div>
  );
}

/**
 * Пример использования
 */
export function VirtualTableExample() {
  // Генерируем 10000 строк для демонстрации
  const largeDataset = Array.from({ length: 10000 }, (_, i) => ({
    id: `row-${i}`,
    name: `Товар ${i + 1}`,
    price: Math.floor(Math.random() * 10000),
    quantity: Math.floor(Math.random() * 100),
    status: ['В наличии', 'На заказ', 'Снято с продажи'][
      Math.floor(Math.random() * 3)
    ],
  }));

  return (
    <VirtualTable
      data={largeDataset}
      itemHeight={40}
      containerHeight={600}
      columns={[
        { key: 'name', label: 'Название', width: '40%' },
        { key: 'price', label: 'Цена', width: '20%' },
        { key: 'quantity', label: 'Количество', width: '20%' },
        { key: 'status', label: 'Статус', width: '20%' },
      ]}
      renderCell={(item, key) => {
        if (key === 'price') {
          return <span className="font-semibold">₽{item.price}</span>;
        }
        if (key === 'status') {
          const statusColors: Record<string, string> = {
            'В наличии': 'bg-green-100 text-green-800',
            'На заказ': 'bg-yellow-100 text-yellow-800',
            'Снято с продажи': 'bg-red-100 text-red-800',
          };
          return (
            <span
              className={`px-2 py-1 rounded text-xs font-medium ${
                statusColors[String(item[key])] || ''
              }`}
            >
              {String(item[key])}
            </span>
          );
        }
        return <span>{String(item[key])}</span>;
      }}
    />
  );
}
