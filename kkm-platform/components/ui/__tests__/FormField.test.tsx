import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Input, Textarea, Select } from "../FormField";

describe("Input", () => {
  it("should render input element", () => {
    render(<Input value="" onChange={vi.fn()} />);

    const input = screen.getByRole("textbox");
    expect(input).toBeInTheDocument();
  });

  it("should call onChange when value changes", async () => {
    const user = userEvent.setup();
    const handleChange = vi.fn();
    render(<Input value="" onChange={handleChange} />);

    const input = screen.getByRole("textbox");
    await user.type(input, "test");

    expect(handleChange).toHaveBeenCalled();
  });

  it("should apply placeholder", () => {
    render(<Input value="" onChange={vi.fn()} placeholder="Enter text" />);

    const input = screen.getByPlaceholderText("Enter text");
    expect(input).toBeInTheDocument();
  });

  it("should be disabled when disabled prop is true", () => {
    render(<Input value="" onChange={vi.fn()} disabled />);

    const input = screen.getByRole("textbox");
    expect(input).toBeDisabled();
  });

  it("should apply custom className", () => {
    render(<Input value="" onChange={vi.fn()} className="custom-class" />);

    const input = screen.getByRole("textbox");
    expect(input.className).toContain("custom-class");
  });

  it("should have required attribute when required prop is true", () => {
    render(<Input value="" onChange={vi.fn()} required />);

    const input = screen.getByRole("textbox");
    expect(input).toHaveAttribute("required");
  });

  it("should support different input types", () => {
    render(<Input type="email" value="" onChange={vi.fn()} />);

    const input = screen.getByRole("textbox");
    expect(input).toHaveAttribute("type", "email");
  });
});

describe("Textarea", () => {
  it("should render textarea element", () => {
    render(<Textarea value="" onChange={vi.fn()} />);

    const textarea = screen.getByRole("textbox");
    expect(textarea.tagName).toBe("TEXTAREA");
  });

  it("should call onChange when value changes", async () => {
    const user = userEvent.setup();
    const handleChange = vi.fn();
    render(<Textarea value="" onChange={handleChange} />);

    const textarea = screen.getByRole("textbox");
    await user.type(textarea, "test");

    expect(handleChange).toHaveBeenCalled();
  });

  it("should apply rows attribute", () => {
    render(<Textarea value="" onChange={vi.fn()} rows={5} />);

    const textarea = screen.getByRole("textbox");
    expect(textarea).toHaveAttribute("rows", "5");
  });

  it("should apply placeholder", () => {
    render(
      <Textarea
        value=""
        onChange={vi.fn()}
        placeholder="Enter description"
      />
    );

    const textarea = screen.getByPlaceholderText("Enter description");
    expect(textarea).toBeInTheDocument();
  });

  it("should be disabled when disabled prop is true", () => {
    render(<Textarea value="" onChange={vi.fn()} disabled />);

    const textarea = screen.getByRole("textbox");
    expect(textarea).toBeDisabled();
  });

  it("should apply custom className", () => {
    render(<Textarea value="" onChange={vi.fn()} className="custom-class" />);

    const textarea = screen.getByRole("textbox");
    expect(textarea.className).toContain("custom-class");
  });

  it("should have required attribute when required prop is true", () => {
    render(<Textarea value="" onChange={vi.fn()} required />);

    const textarea = screen.getByRole("textbox");
    expect(textarea).toHaveAttribute("required");
  });
});

describe("Select", () => {
  it("should render select element", () => {
    render(
      <Select value="" onChange={vi.fn()}>
        <option value="1">Option 1</option>
        <option value="2">Option 2</option>
      </Select>
    );

    const select = screen.getByRole("combobox");
    expect(select).toBeInTheDocument();
  });

  it("should render all options", () => {
    render(
      <Select value="" onChange={vi.fn()}>
        <option value="1">Option 1</option>
        <option value="2">Option 2</option>
        <option value="3">Option 3</option>
      </Select>
    );

    expect(screen.getByText("Option 1")).toBeInTheDocument();
    expect(screen.getByText("Option 2")).toBeInTheDocument();
    expect(screen.getByText("Option 3")).toBeInTheDocument();
  });

  it("should call onChange when selection changes", async () => {
    const user = userEvent.setup();
    const handleChange = vi.fn();
    render(
      <Select value="1" onChange={handleChange}>
        <option value="1">Option 1</option>
        <option value="2">Option 2</option>
      </Select>
    );

    const select = screen.getByRole("combobox");
    await user.selectOptions(select, "2");

    expect(handleChange).toHaveBeenCalled();
  });

  it("should render placeholder option when provided", () => {
    render(
      <Select value="" onChange={vi.fn()}>
        <option value="" disabled>
          Select an option
        </option>
        <option value="1">Option 1</option>
      </Select>
    );

    expect(screen.getByText("Select an option")).toBeInTheDocument();
  });

  it("should be disabled when disabled prop is true", () => {
    render(
      <Select value="" onChange={vi.fn()} disabled>
        <option value="1">Option 1</option>
      </Select>
    );

    const select = screen.getByRole("combobox");
    expect(select).toBeDisabled();
  });

  it("should apply custom className", () => {
    render(
      <Select value="" onChange={vi.fn()} className="custom-class">
        <option value="1">Option 1</option>
      </Select>
    );

    const select = screen.getByRole("combobox");
    expect(select.className).toContain("custom-class");
  });

  it("should have required attribute when required prop is true", () => {
    render(
      <Select value="" onChange={vi.fn()} required>
        <option value="1">Option 1</option>
      </Select>
    );

    const select = screen.getByRole("combobox");
    expect(select).toHaveAttribute("required");
  });

  it("should display selected value", () => {
    render(
      <Select value="2" onChange={vi.fn()}>
        <option value="1">Option 1</option>
        <option value="2">Option 2</option>
      </Select>
    );

    const select = screen.getByRole("combobox") as HTMLSelectElement;
    expect(select.value).toBe("2");
  });

  it("should disable placeholder option", () => {
    const { container } = render(
      <Select value="" onChange={vi.fn()}>
        <option value="" disabled>
          Select
        </option>
        <option value="1">Option 1</option>
      </Select>
    );

    const placeholderOption = container.querySelector("option[value='']");
    expect(placeholderOption).toHaveAttribute("disabled");
  });
});
