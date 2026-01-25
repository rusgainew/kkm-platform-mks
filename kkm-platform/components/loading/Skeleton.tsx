'use client';

import React from 'react';

interface SkeletonProps {
  width?: string | number;
  height?: string | number;
  borderRadius?: string;
  className?: string;
  count?: number;
  baseColor?: string;
  highlightColor?: string;
  animated?: boolean;
}

/**
 * Базовый компонент Skeleton Loader
 * Отображает плейсхолдер во время загрузки данных
 * с эффектом shimmer (мерцания)
 */
export const Skeleton: React.FC<SkeletonProps> = ({
  width = '100%',
  height = '20px',
  borderRadius = '4px',
  className = '',
  count = 1,
  baseColor = 'rgb(30 30 30)',
  highlightColor = 'rgb(80 80 80)',
  animated = true,
}) => {
  const widthValue = typeof width === 'number' ? `${width}px` : width;
  const heightValue = typeof height === 'number' ? `${height}px` : height;

  const skeletonStyle: React.CSSProperties = {
    width: widthValue,
    height: heightValue,
    borderRadius,
    backgroundColor: baseColor,
    backgroundImage: animated
      ? `linear-gradient(
          90deg,
          ${baseColor} 0%,
          ${highlightColor} 50%,
          ${baseColor} 100%
        )`
      : undefined,
    backgroundSize: animated ? '200% 100%' : undefined,
    animation: animated ? 'skeleton-loading 2s infinite' : undefined,
  };

  const skeletons = Array.from({ length: count }, (_, i) => (
    <div
      key={i}
      style={{
        ...skeletonStyle,
        marginBottom: i < count - 1 ? '12px' : '0',
      }}
      className={className}
    />
  ));

  return (
    <>
      <style>{`
        @keyframes skeleton-loading {
          0% {
            background-position: 200% 0;
          }
          100% {
            background-position: -200% 0;
          }
        }
      `}</style>
      {count === 1 ? skeletons[0] : <div>{skeletons}</div>}
    </>
  );
};

/**
 * Компонент для скелета текстовой строки
 */
export const SkeletonText: React.FC<SkeletonProps> = (props) => (
  <Skeleton height={props.height || '16px'} {...props} />
);

/**
 * Компонент для скелета аватара
 */
export const SkeletonAvatar: React.FC<{
  size?: number;
  className?: string;
}> = ({ size = 40, className = '' }) => (
  <Skeleton
    width={size}
    height={size}
    borderRadius="50%"
    className={className}
  />
);

/**
 * Компонент для скелета кнопки
 */
export const SkeletonButton: React.FC<{
  width?: string | number;
  className?: string;
}> = ({ width = '120px', className = '' }) => (
  <Skeleton width={width} height={40} borderRadius="6px" className={className} />
);

/**
 * Компонент для скелета карточки
 */
export const SkeletonCard: React.FC<{
  className?: string;
  count?: number;
}> = ({ className = '', count = 1 }) => (
  <div className={className}>
    {Array.from({ length: count }, (_, i) => (
      <div key={i} className="p-4 border border-gray-800 rounded-lg bg-gray-900 mb-4">
        <Skeleton width="80%" height={24} className="mb-3" />
        <Skeleton width="100%" height={16} className="mb-2" />
        <Skeleton width="95%" height={16} className="mb-4" />
        <div className="flex gap-2">
          <Skeleton width="30%" height={40} borderRadius="6px" />
          <Skeleton width="30%" height={40} borderRadius="6px" />
        </div>
      </div>
    ))}
  </div>
);

/**
 * Компонент для скелета таблицы
 */
export const SkeletonTable: React.FC<{
  rows?: number;
  columns?: number;
  className?: string;
}> = ({ rows = 5, columns = 4, className = '' }) => (
  <div className={`w-full ${className}`}>
    {Array.from({ length: rows }, (_, rowIndex) => (
      <div
        key={rowIndex}
        className="flex gap-4 p-4 border-b border-gray-800 last:border-b-0"
      >
        {Array.from({ length: columns }, (_, colIndex) => (
          <Skeleton
            key={colIndex}
            width={`${100 / columns}%`}
            height={20}
            borderRadius="4px"
          />
        ))}
      </div>
    ))}
  </div>
);

/**
 * Компонент для скелета чарта/графика
 */
export const SkeletonChart: React.FC<{
  width?: string | number;
  height?: string | number;
  className?: string;
}> = ({ width = '100%', height = 300, className = '' }) => (
  <div className={`p-4 border border-gray-800 rounded-lg bg-gray-900 ${className}`}>
    <Skeleton width="40%" height={24} className="mb-4" />
    <Skeleton width={width} height={height} borderRadius="8px" />
  </div>
);

/**
 * Компонент для скелета списка
 */
export const SkeletonList: React.FC<{
  count?: number;
  className?: string;
}> = ({ count = 5, className = '' }) => (
  <div className={className}>
    {Array.from({ length: count }, (_, i) => (
      <div
        key={i}
        className="flex items-center gap-4 p-4 border-b border-gray-800 last:border-b-0"
      >
        <SkeletonAvatar size={40} />
        <div className="flex-1">
          <Skeleton width="60%" height={16} className="mb-2" />
          <Skeleton width="40%" height={12} />
        </div>
      </div>
    ))}
  </div>
);
