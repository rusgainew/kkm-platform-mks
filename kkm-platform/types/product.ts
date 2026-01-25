export interface Product {
  id: number;
  name: string;
  price: number;
  category: string;
  image?: string;
}

export interface ProductCategory {
  id: string;
  name: string;
}
