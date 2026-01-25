'use client';

import { Package } from 'lucide-react';
import Image from 'next/image';
import { Product } from '@/types';

interface ProductCardProps {
  product: Product;
  onAddToCart: (product: Product) => void;
}

export default function ProductCard({ product, onAddToCart }: ProductCardProps) {
  return (
    <button
      onClick={() => onAddToCart(product)}
      className="bg-gray-800 border border-gray-700 rounded-lg p-4 hover:shadow-lg hover:shadow-blue-500/30 hover:border-blue-500/50 transition-all flex flex-col items-center text-center"
    >
      {product.image ? (
        <Image
          src={product.image}
          alt={product.name}
          width={80}
          height={80}
          className="object-cover rounded-md mb-3"
        />
      ) : (
        <div className="w-20 h-20 bg-gray-700 rounded-md flex items-center justify-center mb-3">
          <Package className="w-10 h-10 text-gray-500" />
        </div>
      )}
      <h3 className="font-semibold text-white mb-1 line-clamp-2">{product.name}</h3>
      <p className="text-sm text-gray-400 mb-2">{product.category}</p>
      <p className="text-lg font-bold text-emerald-400">{product.price.toFixed(2)} ₽</p>
    </button>
  );
}
