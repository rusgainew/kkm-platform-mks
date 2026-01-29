import { describe, it, expect } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Tooltip } from "../Tooltip";

describe("Tooltip", () => {
  describe("Рендеринг", () => {
    it("должен рендерить children", () => {
      render(
        <Tooltip content="Tooltip content">
          <button>Trigger</button>
        </Tooltip>,
      );
      expect(screen.getByText("Trigger")).toBeInTheDocument();
    });

    it("не должен показывать tooltip по умолчанию", () => {
      render(
        <Tooltip content="Tooltip content">
          <button>Trigger</button>
        </Tooltip>,
      );
      expect(screen.queryByRole("tooltip")).not.toBeInTheDocument();
    });

    it("должен показывать tooltip при наведении", async () => {
      const user = userEvent.setup();
      render(
        <Tooltip content="Tooltip content" delay={0}>
          <button>Trigger</button>
        </Tooltip>,
      );

      await user.hover(screen.getByText("Trigger"));
      await waitFor(() =>
        expect(screen.getByRole("tooltip")).toBeInTheDocument(),
      );
    });

    it("должен отображать контент", async () => {
      const user = userEvent.setup();
      render(
        <Tooltip content="Test content" delay={0}>
          <button>Trigger</button>
        </Tooltip>,
      );

      await user.hover(screen.getByText("Trigger"));
      await waitFor(() =>
        expect(screen.getByText("Test content")).toBeInTheDocument(),
      );
    });
  });

  describe("Позиционирование", () => {
    it("top по умолчанию", async () => {
      const user = userEvent.setup();
      render(
        <Tooltip content="Tooltip" delay={0}>
          <button>Trigger</button>
        </Tooltip>,
      );

      await user.hover(screen.getByText("Trigger"));
      await waitFor(() => {
        expect(screen.getByRole("tooltip").className).toContain("bottom-full");
      });
    });

    it("right", async () => {
      const user = userEvent.setup();
      render(
        <Tooltip content="Tooltip" position="right" delay={0}>
          <button>Trigger</button>
        </Tooltip>,
      );

      await user.hover(screen.getByText("Trigger"));
      await waitFor(() => {
        expect(screen.getByRole("tooltip").className).toContain("left-full");
      });
    });
  });

  describe("Disabled", () => {
    it("не должен показывать когда disabled", async () => {
      const user = userEvent.setup();
      render(
        <Tooltip content="Tooltip" disabled delay={0}>
          <button>Trigger</button>
        </Tooltip>,
      );

      await user.hover(screen.getByText("Trigger"));
      await new Promise((resolve) => setTimeout(resolve, 100));
      expect(screen.queryByRole("tooltip")).not.toBeInTheDocument();
    });
  });

  describe("Стилизация", () => {
    it("должен применять maxWidth", async () => {
      const user = userEvent.setup();
      render(
        <Tooltip content="Tooltip" maxWidth="300px" delay={0}>
          <button>Trigger</button>
        </Tooltip>,
      );

      await user.hover(screen.getByText("Trigger"));
      await waitFor(() => {
        expect(screen.getByRole("tooltip")).toHaveStyle({ maxWidth: "300px" });
      });
    });

    it("должен применять className", () => {
      const { container } = render(
        <Tooltip content="Tooltip" className="custom-class">
          <button>Trigger</button>
        </Tooltip>,
      );

      expect(container.firstChild).toHaveClass("custom-class");
    });
  });
});
