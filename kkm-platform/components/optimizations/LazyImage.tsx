'use client';

import Image from 'next/image';
import { useRef } from 'react';
import { useIntersectionObserver } from '@/lib/performance';

interface LazyImageProps {
  src: string;
  alt: string;
  width: number;
  height: number;
  className?: string;
  priority?: boolean;
  quality?: number;
}

/**
 * Оптимизированное изображение с ленивой загрузкой
 */
export function LazyImage({
  src,
  alt,
  width,
  height,
  className,
  priority = false,
  quality = 75,
}: LazyImageProps) {
  const ref = useRef<HTMLDivElement>(null);
  const isVisible = useIntersectionObserver(ref as React.RefObject<HTMLDivElement>, { rootMargin: '50px' });

  return (
    <div
      ref={ref}
      className={`relative overflow-hidden bg-gray-100 ${className}`}
      style={{ aspectRatio: `${width}/${height}` }}
    >
      {(isVisible || priority) && (
        <Image
          src={src}
          alt={alt}
          width={width}
          height={height}
          quality={quality}
          priority={priority}
          sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
          className="object-cover w-full h-full"
        />
      )}
    </div>
  );
}

/**
 * Галерея изображений с ленивой загрузкой
 */
export function LazyImageGallery({
  images,
  className = '',
}: {
  images: Array<{ src: string; alt: string; width: number; height: number }>;
  className?: string;
}) {
  return (
    <div className={`grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 ${className}`}>
      {images.map((image, index) => (
        <LazyImage
          key={`${image.src}-${index}`}
          {...image}
          quality={70}
        />
      ))}
    </div>
  );
}
