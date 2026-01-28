import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { Button } from "../Button";
import { Plus, X } from "lucide-react";

describe("Button", () => {
  it("should render button with text", () => {
    render(<Button>Click me</Button>);
    expect(screen.getByText("Click me")).toBeInTheDocument();
  });

  it("should call onClick when clicked", () => {
    const handleClick = vi.fn();
    render(<Button onClick={handleClick}>Click</Button>);
    
    fireEvent.click(screen.getByText("Click"));
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it("should apply primary variant by default", () => {
    render(<Button>Button</Button>);
    const button = screen.getByRole("button");
    expect(button.className).toContain("bg-blue-600");
  });

  it("should apply secondary variant", () => {
    render(<Button variant="secondary">Button</Button>);
    const button = screen.getByRole("button");
    expect(button.className).toContain("bg-gray-600");
  });

  it("should apply danger variant", () => {
    render(<Button variant="danger">Button</Button>);
    const button = screen.getByRole("button");
    expect(button.className).toContain("bg-red-600");
  });

  it("should apply success variant", () => {
    render(<Button variant="success">Button</Button>);
    const button = screen.getByRole("button");
    expect(button.className).toContain("bg-emerald-600");
  });

  it("should apply small size", () => {
    render(<Button size="sm">Button</Button>);
    const button = screen.getByRole("button");
    expect(button.className).toContain("px-3");
  });

  it("should apply large size", () => {
    render(<Button size="lg">Button</Button>);
    const button = screen.getByRole("button");
    expect(button.className).toContain("px-6");
  });

  it("should be disabled when disabled prop is true", () => {
    render(<Button disabled>Button</Button>);
    expect(screen.getByRole("button")).toBeDisabled();
  });

  it("should show loading spinner when loading", () => {
    const { container } = render(<Button loading>Button</Button>);
    const spinner = container.querySelector(".animate-spin");
    expect(spinner).toBeInTheDocument();
  });

  it("should be disabled when loading", () => {
    render(<Button loading>Button</Button>);
    expect(screen.getByRole("button")).toBeDisabled();
  });

  it("should render left icon", () => {
    render(
      <Button leftIcon={<Plus data-testid="left-icon" />}>
        Button
      </Button>
    );
    expect(screen.getByTestId("left-icon")).toBeInTheDocument();
  });

  it("should render right icon", () => {
    render(
      <Button rightIcon={<X data-testid="right-icon" />}>
        Button
      </Button>
    );
    expect(screen.getByTestId("right-icon")).toBeInTheDocument();
  });

  it("should hide icons when loading", () => {
    render(
      <Button loading leftIcon={<Plus data-testid="left-icon" />}>
        Button
      </Button>
    );
    expect(screen.queryByTestId("left-icon")).not.toBeInTheDocument();
  });

  it("should apply fullWidth class", () => {
    render(<Button fullWidth>Button</Button>);
    const button = screen.getByRole("button");
    expect(button.className).toContain("w-full");
  });

  it("should forward ref", () => {
    const ref = vi.fn();
    render(<Button ref={ref}>Button</Button>);
    expect(ref).toHaveBeenCalled();
  });

  it("should apply custom className", () => {
    render(<Button className="custom-class">Button</Button>);
    const button = screen.getByRole("button");
    expect(button.className).toContain("custom-class");
  });

  it("should support all button HTML attributes", () => {
    render(
      <Button type="submit" name="test" value="value">
        Button
      </Button>
    );
    const button = screen.getByRole("button");
    expect(button).toHaveAttribute("type", "submit");
    expect(button).toHaveAttribute("name", "test");
  });
});
