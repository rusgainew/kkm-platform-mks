import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { Spinner, SpinnerOverlay } from "../Spinner";

describe("Spinner", () => {
  it("should render spinner", () => {
    const { container } = render(<Spinner />);
    const spinner = container.querySelector(".animate-spin");
    expect(spinner).toBeInTheDocument();
  });

  it("should apply small size", () => {
    const { container } = render(<Spinner size="sm" />);
    const spinner = container.querySelector(".w-4");
    expect(spinner).toBeInTheDocument();
  });

  it("should apply medium size by default", () => {
    const { container } = render(<Spinner />);
    const spinner = container.querySelector(".w-6");
    expect(spinner).toBeInTheDocument();
  });

  it("should apply large size", () => {
    const { container } = render(<Spinner size="lg" />);
    const spinner = container.querySelector(".w-8");
    expect(spinner).toBeInTheDocument();
  });

  it("should apply extra large size", () => {
    const { container } = render(<Spinner size="xl" />);
    const spinner = container.querySelector(".w-12");
    expect(spinner).toBeInTheDocument();
  });

  it("should apply primary variant by default", () => {
    const { container } = render(<Spinner />);
    const spinner = container.querySelector(".text-blue-500");
    expect(spinner).toBeInTheDocument();
  });

  it("should apply secondary variant", () => {
    const { container } = render(<Spinner variant="secondary" />);
    const spinner = container.querySelector(".text-gray-400");
    expect(spinner).toBeInTheDocument();
  });

  it("should apply white variant", () => {
    const { container } = render(<Spinner variant="white" />);
    const spinner = container.querySelector(".text-white");
    expect(spinner).toBeInTheDocument();
  });

  it("should render with label", () => {
    render(<Spinner label="Loading..." />);
    expect(screen.getByText("Loading...")).toBeInTheDocument();
  });

  it("should center spinner when centered prop is true", () => {
    const { container } = render(<Spinner centered />);
    const wrapper = container.querySelector(".flex.items-center.justify-center");
    expect(wrapper).toBeInTheDocument();
    expect(wrapper?.className).toContain("min-h-[200px]");
  });

  it("should not center spinner by default", () => {
    const { container } = render(<Spinner />);
    const centered = container.querySelector(".min-h-\\[200px\\]");
    expect(centered).not.toBeInTheDocument();
  });

  it("should apply custom className", () => {
    const { container } = render(<Spinner className="custom-class" />);
    const wrapper = container.querySelector(".custom-class");
    expect(wrapper).toBeInTheDocument();
  });
});

describe("SpinnerOverlay", () => {
  it("should render overlay with spinner", () => {
    const { container } = render(<SpinnerOverlay />);
    
    const overlay = container.querySelector(".fixed.inset-0");
    expect(overlay).toBeInTheDocument();
    
    const spinner = container.querySelector(".animate-spin");
    expect(spinner).toBeInTheDocument();
  });

  it("should render with label", () => {
    render(<SpinnerOverlay label="Processing..." />);
    expect(screen.getByText("Processing...")).toBeInTheDocument();
  });

  it("should have backdrop blur", () => {
    const { container } = render(<SpinnerOverlay />);
    const overlay = container.querySelector(".backdrop-blur-sm");
    expect(overlay).toBeInTheDocument();
  });

  it("should use large white spinner", () => {
    const { container } = render(<SpinnerOverlay />);
    const spinner = container.querySelector(".w-8.text-white");
    expect(spinner).toBeInTheDocument();
  });

  it("should have z-50 for proper stacking", () => {
    const { container } = render(<SpinnerOverlay />);
    const overlay = container.querySelector(".z-50");
    expect(overlay).toBeInTheDocument();
  });
});
