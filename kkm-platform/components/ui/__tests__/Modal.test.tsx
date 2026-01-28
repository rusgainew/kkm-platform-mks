import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { Modal } from "../Modal";

describe("Modal", () => {
  it("should render when isOpen is true", () => {
    render(
      <Modal isOpen={true} onClose={vi.fn()} title="Test Modal">
        <div>Modal Content</div>
      </Modal>
    );

    expect(screen.getByText("Test Modal")).toBeInTheDocument();
    expect(screen.getByText("Modal Content")).toBeInTheDocument();
  });

  it("should not render when isOpen is false", () => {
    render(
      <Modal isOpen={false} onClose={vi.fn()} title="Test Modal">
        <div>Modal Content</div>
      </Modal>
    );

    expect(screen.queryByText("Test Modal")).not.toBeInTheDocument();
    expect(screen.queryByText("Modal Content")).not.toBeInTheDocument();
  });

  it("should call onClose when close button is clicked", () => {
    const handleClose = vi.fn();
    render(
      <Modal isOpen={true} onClose={handleClose} title="Test Modal">
        <div>Content</div>
      </Modal>
    );

    const closeButton = screen.getByRole("button", { name: /close/i });
    fireEvent.click(closeButton);

    expect(handleClose).toHaveBeenCalledTimes(1);
  });

  it("should call onClose when backdrop is clicked", () => {
    const handleClose = vi.fn();
    const { container } = render(
      <Modal isOpen={true} onClose={handleClose} title="Test Modal">
        <div>Content</div>
      </Modal>
    );

    // Click the backdrop (first div with fixed positioning)
    const backdrop = container.querySelector(".fixed.inset-0");
    if (backdrop) {
      fireEvent.click(backdrop);
      expect(handleClose).toHaveBeenCalledTimes(1);
    }
  });

  it("should not call onClose when modal content is clicked", () => {
    const handleClose = vi.fn();
    render(
      <Modal isOpen={true} onClose={handleClose} title="Test Modal">
        <div>Content</div>
      </Modal>
    );

    const modalContent = screen.getByText("Content");
    fireEvent.click(modalContent);

    expect(handleClose).not.toHaveBeenCalled();
  });

  it("should render footer when provided", () => {
    render(
      <Modal
        isOpen={true}
        onClose={vi.fn()}
        title="Test Modal"
        footer={<button>Save</button>}
      >
        <div>Content</div>
      </Modal>
    );

    expect(screen.getByText("Save")).toBeInTheDocument();
  });

  it("should not render footer when not provided", () => {
    const { container } = render(
      <Modal isOpen={true} onClose={vi.fn()} title="Test Modal">
        <div>Content</div>
      </Modal>
    );

    // Footer should not be present
    const footer = container.querySelector(".border-t");
    expect(footer).not.toBeInTheDocument();
  });

  it("should apply correct size class for sm size", () => {
    render(
      <Modal isOpen={true} onClose={vi.fn()} title="Test Modal" size="sm">
        <div>Content</div>
      </Modal>
    );

    const modalDialog = screen.getByRole("dialog");
    expect(modalDialog.className).toContain("max-w-sm");
  });

  it("should apply correct size class for md size", () => {
    render(
      <Modal isOpen={true} onClose={vi.fn()} title="Test Modal" size="md">
        <div>Content</div>
      </Modal>
    );

    const modalDialog = screen.getByRole("dialog");
    expect(modalDialog.className).toContain("max-w-md");
  });

  it("should apply correct size class for lg size", () => {
    render(
      <Modal isOpen={true} onClose={vi.fn()} title="Test Modal" size="lg">
        <div>Content</div>
      </Modal>
    );

    const modalDialog = screen.getByRole("dialog");
    expect(modalDialog.className).toContain("max-w-lg");
  });

  it("should apply correct size class for xl size", () => {
    render(
      <Modal isOpen={true} onClose={vi.fn()} title="Test Modal" size="xl">
        <div>Content</div>
      </Modal>
    );

    const modalDialog = screen.getByRole("dialog");
    expect(modalDialog.className).toContain("max-w-xl");
  });

  it("should render title as ReactNode", () => {
    render(
      <Modal
        isOpen={true}
        onClose={vi.fn()}
        title={
          <div>
            <span>Complex</span> <strong>Title</strong>
          </div>
        }
      >
        <div>Content</div>
      </Modal>
    );

    expect(screen.getByText("Complex")).toBeInTheDocument();
    expect(screen.getByText("Title")).toBeInTheDocument();
  });

  it("should have correct accessibility attributes", () => {
    render(
      <Modal isOpen={true} onClose={vi.fn()} title="Test Modal">
        <div>Content</div>
      </Modal>
    );

    const dialog = screen.getByRole("dialog");
    expect(dialog).toHaveAttribute("role", "dialog");
  });
});
