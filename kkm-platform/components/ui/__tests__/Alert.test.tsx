import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { Alert } from "../Alert";
import { AlertCircle } from "lucide-react";

describe("Alert", () => {
  it("should render alert with text", () => {
    render(<Alert>Alert message</Alert>);
    expect(screen.getByText("Alert message")).toBeInTheDocument();
  });

  it("should have role alert", () => {
    render(<Alert>Message</Alert>);
    expect(screen.getByRole("alert")).toBeInTheDocument();
  });

  it("should apply info variant by default", () => {
    const { container } = render(<Alert>Info</Alert>);
    const alert = container.querySelector('[role="alert"]');
    expect(alert?.className).toContain("bg-blue-500/10");
  });

  it("should apply success variant", () => {
    const { container } = render(<Alert variant="success">Success</Alert>);
    const alert = container.querySelector('[role="alert"]');
    expect(alert?.className).toContain("bg-emerald-500/10");
  });

  it("should apply warning variant", () => {
    const { container } = render(<Alert variant="warning">Warning</Alert>);
    const alert = container.querySelector('[role="alert"]');
    expect(alert?.className).toContain("bg-yellow-500/10");
  });

  it("should apply danger variant", () => {
    const { container } = render(<Alert variant="danger">Danger</Alert>);
    const alert = container.querySelector('[role="alert"]');
    expect(alert?.className).toContain("bg-red-500/10");
  });

  it("should render title when provided", () => {
    render(<Alert title="Alert Title">Message</Alert>);
    expect(screen.getByText("Alert Title")).toBeInTheDocument();
  });

  it("should render default icon for each variant", () => {
    const { container } = render(<Alert variant="info">Info</Alert>);
    const icon = container.querySelector("svg");
    expect(icon).toBeInTheDocument();
  });

  it("should render custom icon when provided", () => {
    render(
      <Alert icon={<AlertCircle data-testid="custom-icon" />}>
        Message
      </Alert>
    );
    expect(screen.getByTestId("custom-icon")).toBeInTheDocument();
  });

  it("should hide icon when hideIcon is true", () => {
    const { container } = render(<Alert hideIcon>Message</Alert>);
    const icon = container.querySelector("svg");
    expect(icon).not.toBeInTheDocument();
  });

  it("should render dismiss button when dismissible", () => {
    render(
      <Alert dismissible onDismiss={vi.fn()}>
        Message
      </Alert>
    );
    expect(screen.getByLabelText("Dismiss alert")).toBeInTheDocument();
  });

  it("should call onDismiss when dismiss button is clicked", () => {
    const handleDismiss = vi.fn();
    render(
      <Alert dismissible onDismiss={handleDismiss}>
        Message
      </Alert>
    );
    
    const dismissButton = screen.getByLabelText("Dismiss alert");
    fireEvent.click(dismissButton);
    
    expect(handleDismiss).toHaveBeenCalledTimes(1);
  });

  it("should not render dismiss button when not dismissible", () => {
    render(<Alert>Message</Alert>);
    expect(screen.queryByLabelText("Dismiss alert")).not.toBeInTheDocument();
  });

  it("should apply custom className", () => {
    const { container } = render(<Alert className="custom-class">Message</Alert>);
    const alert = container.querySelector('[role="alert"]');
    expect(alert?.className).toContain("custom-class");
  });

  it("should render with both title and message", () => {
    render(
      <Alert title="Important" variant="warning">
        Please review the details carefully.
      </Alert>
    );
    
    expect(screen.getByText("Important")).toBeInTheDocument();
    expect(screen.getByText("Please review the details carefully.")).toBeInTheDocument();
  });

  it("should support all HTML div attributes", () => {
    render(
      <Alert data-testid="test-alert" id="alert-1">
        Message
      </Alert>
    );
    
    const alert = screen.getByTestId("test-alert");
    expect(alert).toHaveAttribute("id", "alert-1");
  });
});
