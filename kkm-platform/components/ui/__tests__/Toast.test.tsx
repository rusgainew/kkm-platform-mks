import { describe, it, expect, vi } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { ToastProvider, useToast } from "../Toast";

const TestComponent = () => {
  const { showToast, dismissAll } = useToast();

  return (
    <div>
      <button onClick={() => showToast({ message: "Test toast", duration: 0 })}>
        Show Toast
      </button>
      <button
        onClick={() =>
          showToast({ message: "Success", variant: "success", duration: 0 })
        }
      >
        Show Success
      </button>
      <button
        onClick={() =>
          showToast({ message: "Error", variant: "error", duration: 0 })
        }
      >
        Show Error
      </button>
      <button
        onClick={() =>
          showToast({ title: "Title", message: "With title", duration: 0 })
        }
      >
        Show With Title
      </button>
      <button onClick={dismissAll}>Dismiss All</button>
    </div>
  );
};

describe("Toast", () => {
  describe("ToastProvider", () => {
    it("должен рендерить дочерние элементы", () => {
      render(
        <ToastProvider>
          <div>Test Content</div>
        </ToastProvider>,
      );

      expect(screen.getByText("Test Content")).toBeInTheDocument();
    });

    it("должен выбрасывать ошибку вне провайдера", () => {
      const spy = vi.spyOn(console, "error").mockImplementation(() => {});
      expect(() => render(<TestComponent />)).toThrow(
        "useToast must be used within ToastProvider",
      );
      spy.mockRestore();
    });
  });

  describe("Отображение", () => {
    it("должен показывать toast", async () => {
      const user = userEvent.setup();
      render(
        <ToastProvider>
          <TestComponent />
        </ToastProvider>,
      );

      await user.click(screen.getByText("Show Toast"));
      await waitFor(() =>
        expect(screen.getByText("Test toast")).toBeInTheDocument(),
      );
    });

    it("должен показывать с заголовком", async () => {
      const user = userEvent.setup();
      render(
        <ToastProvider>
          <TestComponent />
        </ToastProvider>,
      );

      await user.click(screen.getByText("Show With Title"));
      await waitFor(() => {
        expect(screen.getByText("Title")).toBeInTheDocument();
        expect(screen.getByText("With title")).toBeInTheDocument();
      });
    });

    it("должен показывать несколько", async () => {
      const user = userEvent.setup();
      render(
        <ToastProvider>
          <TestComponent />
        </ToastProvider>,
      );

      await user.click(screen.getByText("Show Toast"));
      await user.click(screen.getByText("Show Success"));

      await waitFor(() => expect(screen.getAllByRole("alert")).toHaveLength(2));
    });
  });

  describe("Варианты", () => {
    it("success", async () => {
      const user = userEvent.setup();
      render(
        <ToastProvider>
          <TestComponent />
        </ToastProvider>,
      );

      await user.click(screen.getByText("Show Success"));
      await waitFor(() =>
        expect(screen.getByRole("alert").className).toContain("bg-green-50"),
      );
    });

    it("error", async () => {
      const user = userEvent.setup();
      render(
        <ToastProvider>
          <TestComponent />
        </ToastProvider>,
      );

      await user.click(screen.getByText("Show Error"));
      await waitFor(() =>
        expect(screen.getByRole("alert").className).toContain("bg-red-50"),
      );
    });
  });

  describe("Закрытие", () => {
    it("при клике", async () => {
      const user = userEvent.setup();
      render(
        <ToastProvider>
          <TestComponent />
        </ToastProvider>,
      );

      await user.click(screen.getByText("Show Toast"));
      await waitFor(() =>
        expect(screen.getByRole("alert")).toBeInTheDocument(),
      );

      await user.click(screen.getByLabelText("Закрыть"));
      await waitFor(() =>
        expect(screen.queryByRole("alert")).not.toBeInTheDocument(),
      );
    });

    it("всех тостов", async () => {
      const user = userEvent.setup();
      render(
        <ToastProvider>
          <TestComponent />
        </ToastProvider>,
      );

      await user.click(screen.getByText("Show Toast"));
      await user.click(screen.getByText("Show Success"));
      await waitFor(() => expect(screen.getAllByRole("alert")).toHaveLength(2));

      await user.click(screen.getByText("Dismiss All"));
      await waitFor(() =>
        expect(screen.queryByRole("alert")).not.toBeInTheDocument(),
      );
    });
  });

  describe("Позиционирование", () => {
    it("top-right", () => {
      const { container } = render(
        <ToastProvider>
          <div>Content</div>
        </ToastProvider>,
      );

      const toastContainer = container.querySelector(".fixed");
      expect(toastContainer?.className).toContain("top-4");
      expect(toastContainer?.className).toContain("right-4");
    });
  });
});
