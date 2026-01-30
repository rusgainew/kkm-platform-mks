import React from "react";

export interface SwitchProps {
  /** Состояние переключателя */
  checked?: boolean;
  /** Callback при изменении состояния */
  onChange?: (checked: boolean) => void;
  /** Отключенное состояние */
  disabled?: boolean;
  /** Размер переключателя */
  size?: "sm" | "md" | "lg";
  /** Label текст */
  label?: React.ReactNode;
  /** Позиция label */
  labelPosition?: "left" | "right";
  /** Дополнительные CSS классы */
  className?: string;
  /** ID для связи с label */
  id?: string;
  /** Имя для формы */
  name?: string;
}

const sizeClasses = {
  sm: {
    track: "w-8 h-4",
    thumb: "w-3 h-3",
    translate: "translate-x-4",
  },
  md: {
    track: "w-11 h-6",
    thumb: "w-5 h-5",
    translate: "translate-x-5",
  },
  lg: {
    track: "w-14 h-7",
    thumb: "w-6 h-6",
    translate: "translate-x-7",
  },
};

export const Switch: React.FC<SwitchProps> = ({
  checked = false,
  onChange,
  disabled = false,
  size = "md",
  label,
  labelPosition = "right",
  className = "",
  id,
  name,
}) => {
  const sizeConfig = sizeClasses[size];

  const handleClick = () => {
    if (!disabled && onChange) {
      onChange(!checked);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if ((e.key === " " || e.key === "Enter") && !disabled) {
      e.preventDefault();
      if (onChange) {
        onChange(!checked);
      }
    }
  };

  const switchElement = (
    <button
      type="button"
      role="switch"
      aria-checked={checked}
      aria-disabled={disabled}
      disabled={disabled}
      onClick={handleClick}
      onKeyDown={handleKeyDown}
      id={id}
      name={name}
      className={`
        relative inline-flex items-center rounded-full transition-colors duration-200 ease-in-out
        ${sizeConfig.track}
        ${
          checked
            ? disabled
              ? "bg-blue-300"
              : "bg-blue-600 hover:bg-blue-700"
            : disabled
              ? "bg-gray-300"
              : "bg-gray-400 hover:bg-gray-500"
        }
        ${disabled ? "cursor-not-allowed opacity-50" : "cursor-pointer"}
        focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2
      `}
    >
      <span
        className={`
          inline-block transform rounded-full bg-white shadow-lg transition-transform duration-200 ease-in-out
          ${sizeConfig.thumb}
          ${checked ? sizeConfig.translate : "translate-x-0.5"}
        `}
      />
    </button>
  );

  if (!label) {
    return <div className={className}>{switchElement}</div>;
  }

  return (
    <label
      className={`
        inline-flex items-center gap-2 cursor-pointer
        ${disabled ? "opacity-50 cursor-not-allowed" : ""}
        ${className}
      `}
    >
      {labelPosition === "left" && (
        <span className="text-sm text-gray-700 dark:text-gray-300">
          {label}
        </span>
      )}
      {switchElement}
      {labelPosition === "right" && (
        <span className="text-sm text-gray-700 dark:text-gray-300">
          {label}
        </span>
      )}
    </label>
  );
};
