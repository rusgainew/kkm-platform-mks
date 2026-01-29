import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "../Tabs";

describe("Tabs", () => {
  describe("Основной рендеринг", () => {
    it("должен рендерить компонент Tabs", () => {
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
        </Tabs>,
      );

      expect(screen.getByRole("tablist")).toBeInTheDocument();
    });

    it("должен показывать активный контент по умолчанию", () => {
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
            <TabsTrigger value="tab2">Tab 2</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
          <TabsContent value="tab2">Content 2</TabsContent>
        </Tabs>,
      );

      expect(screen.getByText("Content 1")).toBeInTheDocument();
      expect(screen.queryByText("Content 2")).not.toBeInTheDocument();
    });

    it("должен рендерить несколько триггеров", () => {
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
            <TabsTrigger value="tab2">Tab 2</TabsTrigger>
            <TabsTrigger value="tab3">Tab 3</TabsTrigger>
          </TabsList>
        </Tabs>,
      );

      expect(screen.getByText("Tab 1")).toBeInTheDocument();
      expect(screen.getByText("Tab 2")).toBeInTheDocument();
      expect(screen.getByText("Tab 3")).toBeInTheDocument();
    });
  });

  describe("Варианты отображения", () => {
    it("должен применять line вариант по умолчанию", () => {
      const { container } = render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
          </TabsList>
        </Tabs>,
      );

      const tablist = container.querySelector('[role="tablist"]');
      expect(tablist?.className).toContain("border-b");
    });

    it("должен применять pills вариант", () => {
      const { container } = render(
        <Tabs defaultValue="tab1" variant="pills">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
          </TabsList>
        </Tabs>,
      );

      const tablist = container.querySelector('[role="tablist"]');
      expect(tablist?.className).toContain("bg-gray-100");
      expect(tablist?.className).toContain("rounded-lg");
    });

    it("должен применять стили line к триггерам", () => {
      render(
        <Tabs defaultValue="tab1" variant="line">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
          </TabsList>
        </Tabs>,
      );

      const trigger = screen.getByRole("tab");
      expect(trigger.className).toContain("border-b-2");
    });

    it("должен применять стили pills к триггерам", () => {
      render(
        <Tabs defaultValue="tab1" variant="pills">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
          </TabsList>
        </Tabs>,
      );

      const trigger = screen.getByRole("tab");
      expect(trigger.className).toContain("rounded-md");
    });
  });

  describe("Переключение вкладок", () => {
    it("должен переключать вкладки при клике", async () => {
      const user = userEvent.setup();
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
            <TabsTrigger value="tab2">Tab 2</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
          <TabsContent value="tab2">Content 2</TabsContent>
        </Tabs>,
      );

      expect(screen.getByText("Content 1")).toBeInTheDocument();
      expect(screen.queryByText("Content 2")).not.toBeInTheDocument();

      await user.click(screen.getByText("Tab 2"));

      expect(screen.queryByText("Content 1")).not.toBeInTheDocument();
      expect(screen.getByText("Content 2")).toBeInTheDocument();
    });

    it("должен вызывать onChange при переключении", async () => {
      const user = userEvent.setup();
      const onChange = vi.fn();
      render(
        <Tabs defaultValue="tab1" onChange={onChange}>
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
            <TabsTrigger value="tab2">Tab 2</TabsTrigger>
          </TabsList>
        </Tabs>,
      );

      await user.click(screen.getByText("Tab 2"));

      expect(onChange).toHaveBeenCalledWith("tab2");
    });

    it("должен работать в контролируемом режиме", async () => {
      const user = userEvent.setup();
      const onChange = vi.fn();
      const { rerender } = render(
        <Tabs defaultValue="tab1" value="tab1" onChange={onChange}>
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
            <TabsTrigger value="tab2">Tab 2</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
          <TabsContent value="tab2">Content 2</TabsContent>
        </Tabs>,
      );

      expect(screen.getByText("Content 1")).toBeInTheDocument();

      await user.click(screen.getByText("Tab 2"));
      expect(onChange).toHaveBeenCalledWith("tab2");

      // В контролируемом режиме нужно обновить value
      rerender(
        <Tabs defaultValue="tab1" value="tab2" onChange={onChange}>
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
            <TabsTrigger value="tab2">Tab 2</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
          <TabsContent value="tab2">Content 2</TabsContent>
        </Tabs>,
      );

      expect(screen.getByText("Content 2")).toBeInTheDocument();
    });
  });

  describe("Клавиатурная навигация", () => {
    it("должен переключаться по Enter", async () => {
      const user = userEvent.setup();
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
            <TabsTrigger value="tab2">Tab 2</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
          <TabsContent value="tab2">Content 2</TabsContent>
        </Tabs>,
      );

      const tab2 = screen.getByText("Tab 2");
      tab2.focus();
      await user.keyboard("{Enter}");

      expect(screen.getByText("Content 2")).toBeInTheDocument();
    });

    it("должен переключаться по Space", async () => {
      const user = userEvent.setup();
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
            <TabsTrigger value="tab2">Tab 2</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
          <TabsContent value="tab2">Content 2</TabsContent>
        </Tabs>,
      );

      const tab2 = screen.getByText("Tab 2");
      tab2.focus();
      await user.keyboard(" ");

      expect(screen.getByText("Content 2")).toBeInTheDocument();
    });
  });

  describe("Disabled состояние", () => {
    it("должен отображать disabled триггер", () => {
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1" disabled>
              Tab 1
            </TabsTrigger>
          </TabsList>
        </Tabs>,
      );

      const trigger = screen.getByRole("tab");
      expect(trigger).toHaveAttribute("aria-disabled", "true");
      expect(trigger).toBeDisabled();
    });

    it("не должен переключаться на disabled вкладку", async () => {
      const user = userEvent.setup();
      const onChange = vi.fn();
      render(
        <Tabs defaultValue="tab1" onChange={onChange}>
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
            <TabsTrigger value="tab2" disabled>
              Tab 2
            </TabsTrigger>
          </TabsList>
        </Tabs>,
      );

      await user.click(screen.getByText("Tab 2"));

      expect(onChange).not.toHaveBeenCalled();
    });

    it("должен применять opacity для disabled", () => {
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1" disabled>
              Tab 1
            </TabsTrigger>
          </TabsList>
        </Tabs>,
      );

      const trigger = screen.getByRole("tab");
      expect(trigger.className).toContain("opacity-50");
    });
  });

  describe("Активное состояние", () => {
    it("должен отмечать активный триггер", () => {
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
            <TabsTrigger value="tab2">Tab 2</TabsTrigger>
          </TabsList>
        </Tabs>,
      );

      const tabs = screen.getAllByRole("tab");
      expect(tabs[0]).toHaveAttribute("aria-selected", "true");
      expect(tabs[1]).toHaveAttribute("aria-selected", "false");
    });

    it("должен применять стили к активному триггеру (line)", () => {
      render(
        <Tabs defaultValue="tab1" variant="line">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
          </TabsList>
        </Tabs>,
      );

      const trigger = screen.getByRole("tab");
      expect(trigger.className).toContain("text-blue-600");
      expect(trigger.className).toContain("border-blue-600");
    });

    it("должен применять стили к активному триггеру (pills)", () => {
      render(
        <Tabs defaultValue="tab1" variant="pills">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
          </TabsList>
        </Tabs>,
      );

      const trigger = screen.getByRole("tab");
      expect(trigger.className).toContain("bg-white");
    });
  });

  describe("TabsList пропсы", () => {
    it("должен применять fullWidth", () => {
      const { container } = render(
        <Tabs defaultValue="tab1">
          <TabsList fullWidth>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
          </TabsList>
        </Tabs>,
      );

      const tablist = container.querySelector('[role="tablist"]');
      expect(tablist?.className).toContain("w-full");
    });

    it("должен применять кастомный className к TabsList", () => {
      const { container } = render(
        <Tabs defaultValue="tab1">
          <TabsList className="custom-tablist">
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
          </TabsList>
        </Tabs>,
      );

      const tablist = container.querySelector('[role="tablist"]');
      expect(tablist).toHaveClass("custom-tablist");
    });
  });

  describe("Accessibility", () => {
    it("должен иметь корректную роль для списка", () => {
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
          </TabsList>
        </Tabs>,
      );

      expect(screen.getByRole("tablist")).toBeInTheDocument();
    });

    it("должен иметь корректную роль для триггеров", () => {
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
            <TabsTrigger value="tab2">Tab 2</TabsTrigger>
          </TabsList>
        </Tabs>,
      );

      const tabs = screen.getAllByRole("tab");
      expect(tabs).toHaveLength(2);
    });

    it("должен иметь корректную роль для контента", () => {
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1">Content 1</TabsContent>
        </Tabs>,
      );

      expect(screen.getByRole("tabpanel")).toBeInTheDocument();
    });
  });

  describe("Дополнительные пропсы", () => {
    it("должен применять кастомный className к Tabs", () => {
      const { container } = render(
        <Tabs defaultValue="tab1" className="custom-tabs">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
          </TabsList>
        </Tabs>,
      );

      const tabsWrapper = container.querySelector(".custom-tabs");
      expect(tabsWrapper).toBeInTheDocument();
    });

    it("должен применять кастомный className к TabsContent", () => {
      render(
        <Tabs defaultValue="tab1">
          <TabsList>
            <TabsTrigger value="tab1">Tab 1</TabsTrigger>
          </TabsList>
          <TabsContent value="tab1" className="custom-content">
            Content 1
          </TabsContent>
        </Tabs>,
      );

      const content = screen.getByRole("tabpanel");
      expect(content).toHaveClass("custom-content");
    });
  });
});
