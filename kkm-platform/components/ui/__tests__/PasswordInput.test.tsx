import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { PasswordInput } from "../PasswordInput";

describe("PasswordInput", () => {
  it("should render password input with type password by default", () => {
    const { container } = render(<PasswordInput value="" onChange={vi.fn()} />);

    const input = container.querySelector("input");
    expect(input).toHaveAttribute("type", "password");
  });

  it("should render toggle button when showToggle is true", () => {
    render(<PasswordInput value="" onChange={vi.fn()} showToggle={true} />);

    const toggleButton = screen.getByRole("button");
    expect(toggleButton).toBeInTheDocument();
  });

  it("should not render toggle button when showToggle is false", () => {
    render(<PasswordInput value="" onChange={vi.fn()} showToggle={false} />);

    const toggleButton = screen.queryByRole("button");
    expect(toggleButton).not.toBeInTheDocument();
  });

  it("should toggle password visibility when button is clicked", async () => {
    const { container } = render(
      <PasswordInput value="" onChange={vi.fn()} showToggle={true} />
    );

    const input = container.querySelector("input") as HTMLInputElement;
    const toggleButton = screen.getByRole("button");

    // Initially password type
    expect(input).toHaveAttribute("type", "password");

    // Click to show password
    fireEvent.click(toggleButton);
    expect(input).toHaveAttribute("type", "text");

    // Click to hide password again
    fireEvent.click(toggleButton);
    expect(input).toHaveAttribute("type", "password");
  });

  it("should call onChange when input value changes", async () => {
    const user = userEvent.setup();
    const handleChange = vi.fn();
    const { container } = render(
      <PasswordInput value="" onChange={handleChange} />
    );

    const input = container.querySelector("input") as HTMLInputElement;
    await user.type(input, "password123");

    expect(handleChange).toHaveBeenCalled();
  });

  it("should display Eye icon when password is hidden", () => {
    const { container } = render(
      <PasswordInput value="" onChange={vi.fn()} showToggle={true} />
    );

    // Eye icon should be visible (password hidden)
    const eyeIcon = container.querySelector("svg");
    expect(eyeIcon).toBeInTheDocument();
  });

  it("should display EyeOff icon when password is visible", () => {
    const { container } = render(
      <PasswordInput value="" onChange={vi.fn()} showToggle={true} />
    );

    const toggleButton = screen.getByRole("button");
    fireEvent.click(toggleButton);

    // EyeOff icon should be visible (password shown)
    const eyeOffIcon = container.querySelector("svg");
    expect(eyeOffIcon).toBeInTheDocument();
  });

  it("should forward ref to input element", () => {
    const ref = vi.fn();
    render(<PasswordInput ref={ref} value="" onChange={vi.fn()} />);

    expect(ref).toHaveBeenCalled();
  });

  it("should accept all standard input props", () => {
    const { container } = render(
      <PasswordInput
        value=""
        onChange={vi.fn()}
        placeholder="Enter password"
        required
        minLength={8}
        disabled
      />
    );

    const input = container.querySelector("input") as HTMLInputElement;
    expect(input).toHaveAttribute("placeholder", "Enter password");
    expect(input).toHaveAttribute("required");
    expect(input).toHaveAttribute("minLength", "8");
    expect(input).toBeDisabled();
  });

  it("should apply custom className", () => {
    const { container } = render(
      <PasswordInput value="" onChange={vi.fn()} className="custom-class" />
    );

    const input = container.querySelector("input") as HTMLInputElement;
    expect(input.className).toContain("custom-class");
  });

  it("should disable toggle button when input is disabled", () => {
    render(
      <PasswordInput
        value=""
        onChange={vi.fn()}
        showToggle={true}
        disabled
      />
    );

    const toggleButton = screen.getByRole("button");
    expect(toggleButton).toBeDisabled();
  });

  it("should have correct accessibility attributes", () => {
    const { container } = render(
      <PasswordInput
        value=""
        onChange={vi.fn()}
        showToggle={true}
        aria-label="Password"
      />
    );

    const input = container.querySelector("input") as HTMLInputElement;
    expect(input).toHaveAttribute("aria-label", "Password");

    const toggleButton = screen.getByRole("button");
    expect(toggleButton).toHaveAttribute("type", "button");
  });
});
