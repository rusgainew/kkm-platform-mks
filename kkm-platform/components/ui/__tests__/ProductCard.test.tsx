import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import ProductCard from "../ProductCard";
import { Product } from "@/types";

describe("ProductCard", () => {
  const mockProduct: Product = {
    id: "1",
    name: "Test Product",
    category: "Electronics",
    price: 999.99,
    barcode: "1234567890",
    image: "/test-image.jpg",
  };

  const mockProductWithoutImage: Product = {
    id: "2",
    name: "Product Without Image",
    category: "Accessories",
    price: 49.99,
    barcode: "0987654321",
  };

  it("should render product with all details", () => {
    const handleAddToCart = vi.fn();
    render(<ProductCard product={mockProduct} onAddToCart={handleAddToCart} />);

    expect(screen.getByText("Test Product")).toBeInTheDocument();
    expect(screen.getByText("Electronics")).toBeInTheDocument();
    expect(screen.getByText("999.99 ₽")).toBeInTheDocument();
  });

  it("should render product image when image is provided", () => {
    const handleAddToCart = vi.fn();
    const { container } = render(
      <ProductCard product={mockProduct} onAddToCart={handleAddToCart} />,
    );

    const image = container.querySelector("img");
    expect(image).toBeInTheDocument();
    expect(image).toHaveAttribute("alt", "Test Product");
  });

  it("should render placeholder icon when image is not provided", () => {
    const handleAddToCart = vi.fn();
    const { container } = render(
      <ProductCard
        product={mockProductWithoutImage}
        onAddToCart={handleAddToCart}
      />,
    );

    // Check for Package icon placeholder
    const placeholder = container.querySelector(".bg-gray-700");
    expect(placeholder).toBeInTheDocument();

    const icon = container.querySelector("svg");
    expect(icon).toBeInTheDocument();
  });

  it("should call onAddToCart when card is clicked", () => {
    const handleAddToCart = vi.fn();
    render(<ProductCard product={mockProduct} onAddToCart={handleAddToCart} />);

    const card = screen.getByRole("button");
    fireEvent.click(card);

    expect(handleAddToCart).toHaveBeenCalledTimes(1);
    expect(handleAddToCart).toHaveBeenCalledWith(mockProduct);
  });

  it("should format price with 2 decimal places", () => {
    const handleAddToCart = vi.fn();
    const productWithWholePrice: Product = {
      ...mockProduct,
      price: 100,
    };

    render(
      <ProductCard
        product={productWithWholePrice}
        onAddToCart={handleAddToCart}
      />,
    );

    expect(screen.getByText("100.00 ₽")).toBeInTheDocument();
  });

  it("should truncate long product names with line-clamp-2", () => {
    const handleAddToCart = vi.fn();
    const longNameProduct: Product = {
      ...mockProduct,
      name: "This is a very long product name that should be truncated with ellipsis when displayed",
    };

    const { container } = render(
      <ProductCard product={longNameProduct} onAddToCart={handleAddToCart} />,
    );

    const title = screen.getByText(longNameProduct.name);
    expect(title).toHaveClass("line-clamp-2");
  });

  it("should have hover effects applied via CSS classes", () => {
    const handleAddToCart = vi.fn();
    const { container } = render(
      <ProductCard product={mockProduct} onAddToCart={handleAddToCart} />,
    );

    const card = screen.getByRole("button");
    expect(card).toHaveClass("hover:shadow-lg");
    expect(card).toHaveClass("hover:shadow-blue-500/30");
    expect(card).toHaveClass("hover:border-blue-500/50");
  });

  it("should render as a button element", () => {
    const handleAddToCart = vi.fn();
    render(<ProductCard product={mockProduct} onAddToCart={handleAddToCart} />);

    const card = screen.getByRole("button");
    expect(card.tagName).toBe("BUTTON");
  });

  it("should handle products with zero price", () => {
    const handleAddToCart = vi.fn();
    const freeProduct: Product = {
      ...mockProduct,
      price: 0,
    };

    render(<ProductCard product={freeProduct} onAddToCart={handleAddToCart} />);

    expect(screen.getByText("0.00 ₽")).toBeInTheDocument();
  });

  it("should have correct styling classes", () => {
    const handleAddToCart = vi.fn();
    render(<ProductCard product={mockProduct} onAddToCart={handleAddToCart} />);

    const card = screen.getByRole("button");
    expect(card).toHaveClass("bg-gray-800");
    expect(card).toHaveClass("border-gray-700");
    expect(card).toHaveClass("rounded-lg");
  });

  it("should display product name in white color", () => {
    const handleAddToCart = vi.fn();
    render(<ProductCard product={mockProduct} onAddToCart={handleAddToCart} />);

    const name = screen.getByText("Test Product");
    expect(name).toHaveClass("text-white");
  });

  it("should display category in gray color", () => {
    const handleAddToCart = vi.fn();
    render(<ProductCard product={mockProduct} onAddToCart={handleAddToCart} />);

    const category = screen.getByText("Electronics");
    expect(category).toHaveClass("text-gray-400");
  });

  it("should display price in emerald color", () => {
    const handleAddToCart = vi.fn();
    render(<ProductCard product={mockProduct} onAddToCart={handleAddToCart} />);

    const price = screen.getByText("999.99 ₽");
    expect(price).toHaveClass("text-emerald-400");
  });
});
