import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Switch } from "../Switch";

describe("Switch", () => {
  describe("Основной рендеринг", () => {
    it("должен рендериться с корректной ролью", () => {
      render(<Switch checked={false} onChange={() => {}} />);
      const switchElement = screen.getByRole("switch");
      expect(switchElement).toBeInTheDocument();
    });

    it("должен отображать checked состояние", () => {
      render(<Switch checked={true} onChange={() => {}} />);
      const switchElement = screen.getByRole("switch");
      expect(switchElement).toHaveAttribute("aria-checked", "true");
    });

    it("должен отображать unchecked состояние", () => {
      render(<Switch checked={false} onChange={() => {}} />);
      const switchElement = screen.getByRole("switch");
      expect(switchElement).toHaveAttribute("aria-checked", "false");
    });

    it("должен рендериться с label", () => {
      render(<Switch label="Test Label" checked={false} onChange={() => {}} />);
      expect(screen.getByText("Test Label")).toBeInTheDocument();
    });

    it("должен рендериться без label", () => {
      const { container } = render(
        <Switch checked={false} onChange={() => {}} />,
      );
      const label = container.querySelector("label");
      expect(label).not.toBeInTheDocument();
    });
  });

  describe("Размеры", () => {
    it("должен применять класс для sm размера", () => {
      const { container } = render(
        <Switch size="sm" checked={false} onChange={() => {}} />,
      );
      const switchButton = container.querySelector("button");
      expect(switchButton?.className).toContain("w-8 h-4");
    });

    it("должен применять класс для md размера (по умолчанию)", () => {
      const { container } = render(
        <Switch checked={false} onChange={() => {}} />,
      );
      const switchButton = container.querySelector("button");
      expect(switchButton?.className).toContain("w-11 h-6");
    });

    it("должен применять класс для lg размера", () => {
      const { container } = render(
        <Switch size="lg" checked={false} onChange={() => {}} />,
      );
      const switchButton = container.querySelector("button");
      expect(switchButton?.className).toContain("w-14 h-7");
    });
  });

  describe("Label позиция", () => {
    it("должен отображать label слева", () => {
      const { container } = render(
        <Switch
          label="Left Label"
          labelPosition="left"
          checked={false}
          onChange={() => {}}
        />,
      );
      const label = container.querySelector("label");
      const text = label?.querySelector("span:first-child");
      expect(text?.textContent).toBe("Left Label");
    });

    it("должен отображать label справа (по умолчанию)", () => {
      render(
        <Switch label="Right Label" checked={false} onChange={() => {}} />,
      );
      expect(screen.getByText("Right Label")).toBeInTheDocument();
    });
  });

  describe("Взаимодействие", () => {
    it("должен вызывать onChange при клике", async () => {
      const user = userEvent.setup();
      const onChange = vi.fn();
      render(<Switch checked={false} onChange={onChange} />);

      const switchElement = screen.getByRole("switch");
      await user.click(switchElement);

      expect(onChange).toHaveBeenCalledTimes(1);
      expect(onChange).toHaveBeenCalledWith(true);
    });

    it("должен переключать состояние с false на true", async () => {
      const user = userEvent.setup();
      const onChange = vi.fn();
      render(<Switch checked={false} onChange={onChange} />);

      const switchElement = screen.getByRole("switch");
      await user.click(switchElement);

      expect(onChange).toHaveBeenCalledWith(true);
    });

    it("должен переключать состояние с true на false", async () => {
      const user = userEvent.setup();
      const onChange = vi.fn();
      render(<Switch checked={true} onChange={onChange} />);

      const switchElement = screen.getByRole("switch");
      await user.click(switchElement);

      expect(onChange).toHaveBeenCalledWith(false);
    });

    it("должен активироваться по Space", async () => {
      const user = userEvent.setup();
      const onChange = vi.fn();
      render(<Switch checked={false} onChange={onChange} />);

      const switchElement = screen.getByRole("switch");
      switchElement.focus();
      await user.keyboard(" ");

      expect(onChange).toHaveBeenCalledWith(true);
    });

    it("должен активироваться по Enter", async () => {
      const user = userEvent.setup();
      const onChange = vi.fn();
      render(<Switch checked={false} onChange={onChange} />);

      const switchElement = screen.getByRole("switch");
      switchElement.focus();
      await user.keyboard("{Enter}");

      expect(onChange).toHaveBeenCalledWith(true);
    });
  });

  describe("Disabled состояние", () => {
    it("должен отображать disabled состояние", () => {
      render(<Switch disabled checked={false} onChange={() => {}} />);
      const switchElement = screen.getByRole("switch");
      expect(switchElement).toHaveAttribute("aria-disabled", "true");
      expect(switchElement).toBeDisabled();
    });

    it("не должен вызывать onChange когда disabled", async () => {
      const user = userEvent.setup();
      const onChange = vi.fn();
      render(<Switch disabled checked={false} onChange={onChange} />);

      const switchElement = screen.getByRole("switch");
      await user.click(switchElement);

      expect(onChange).not.toHaveBeenCalled();
    });

    it("должен применять стили opacity для disabled", () => {
      const { container } = render(
        <Switch disabled checked={false} onChange={() => {}} />,
      );
      const switchButton = container.querySelector("button");
      expect(switchButton?.className).toContain("opacity-50");
    });

    it("должен применять cursor-not-allowed для disabled", () => {
      const { container } = render(
        <Switch disabled checked={false} onChange={() => {}} />,
      );
      const switchButton = container.querySelector("button");
      expect(switchButton?.className).toContain("cursor-not-allowed");
    });
  });

  describe("Accessibility", () => {
    it("должен иметь корректный ARIA атрибут для checked", () => {
      render(<Switch checked={true} onChange={() => {}} />);
      const switchElement = screen.getByRole("switch");
      expect(switchElement).toHaveAttribute("aria-checked", "true");
    });

    it("должен иметь корректный ARIA атрибут для disabled", () => {
      render(<Switch disabled checked={false} onChange={() => {}} />);
      const switchElement = screen.getByRole("switch");
      expect(switchElement).toHaveAttribute("aria-disabled", "true");
    });

    it("должен связывать id с кнопкой", () => {
      render(<Switch id="test-switch" checked={false} onChange={() => {}} />);
      const switchElement = screen.getByRole("switch");
      expect(switchElement).toHaveAttribute("id", "test-switch");
    });

    it("должен иметь name атрибут", () => {
      render(<Switch name="test-name" checked={false} onChange={() => {}} />);
      const switchElement = screen.getByRole("switch");
      expect(switchElement).toHaveAttribute("name", "test-name");
    });
  });

  describe("Дополнительные пропсы", () => {
    it("должен применять кастомный className", () => {
      const { container } = render(
        <Switch className="custom-class" checked={false} onChange={() => {}} />,
      );
      expect(container.firstChild).toHaveClass("custom-class");
    });

    it("должен работать без onChange", () => {
      expect(() => {
        render(<Switch checked={false} />);
      }).not.toThrow();
    });
  });
});
