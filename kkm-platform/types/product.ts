export interface Product {
  id: number | string;
  name: string;
  price: number;
  category: string;
  barcode?: string;
  image?: string;
}

export interface ProductCategory {
  id: string;
  name: string;
}
