import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { Badge } from "../Badge";

describe("Badge", () => {
  it("should render badge with text", () => {
    render(<Badge>Badge Text</Badge>);
    expect(screen.getByText("Badge Text")).toBeInTheDocument();
  });

  it("should apply default variant", () => {
    const { container } = render(<Badge>Badge</Badge>);
    const badge = container.querySelector("span");
    expect(badge?.className).toContain("bg-blue-500/20");
  });

  it("should apply success variant", () => {
    const { container } = render(<Badge variant="success">Success</Badge>);
    const badge = container.querySelector("span");
    expect(badge?.className).toContain("bg-emerald-500/20");
  });

  it("should apply danger variant", () => {
    const { container } = render(<Badge variant="danger">Danger</Badge>);
    const badge = container.querySelector("span");
    expect(badge?.className).toContain("bg-red-500/20");
  });

  it("should apply warning variant", () => {
    const { container } = render(<Badge variant="warning">Warning</Badge>);
    const badge = container.querySelector("span");
    expect(badge?.className).toContain("bg-yellow-500/20");
  });

  it("should apply small size", () => {
    const { container } = render(<Badge size="sm">Small</Badge>);
    const badge = container.querySelector("span");
    expect(badge?.className).toContain("text-xs");
  });

  it("should apply large size", () => {
    const { container } = render(<Badge size="lg">Large</Badge>);
    const badge = container.querySelector("span");
    expect(badge?.className).toContain("text-base");
  });

  it("should render dot indicator", () => {
    const { container } = render(<Badge dot>With Dot</Badge>);
    const dot = container.querySelector(".rounded-full");
    expect(dot).toBeInTheDocument();
  });

  it("should render remove button when removable", () => {
    render(
      <Badge removable onRemove={vi.fn()}>
        Removable
      </Badge>
    );
    const removeButton = screen.getByLabelText("Remove");
    expect(removeButton).toBeInTheDocument();
  });

  it("should call onRemove when remove button is clicked", () => {
    const handleRemove = vi.fn();
    render(
      <Badge removable onRemove={handleRemove}>
        Removable
      </Badge>
    );
    
    const removeButton = screen.getByLabelText("Remove");
    fireEvent.click(removeButton);
    
    expect(handleRemove).toHaveBeenCalledTimes(1);
  });

  it("should not render remove button when not removable", () => {
    render(<Badge>Not Removable</Badge>);
    expect(screen.queryByLabelText("Remove")).not.toBeInTheDocument();
  });

  it("should apply custom className", () => {
    const { container } = render(<Badge className="custom-class">Badge</Badge>);
    const badge = container.querySelector("span");
    expect(badge?.className).toContain("custom-class");
  });

  it("should render with dot and correct color", () => {
    const { container } = render(
      <Badge dot variant="success">
        Success
      </Badge>
    );
    const dot = container.querySelector(".bg-emerald-400");
    expect(dot).toBeInTheDocument();
  });
});
