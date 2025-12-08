import React, { useId } from 'react';
import './ToggleSwitch.css';

interface ToggleSwitchProps {
  checked: boolean;
  onChange: (checked: boolean) => void;
  disabled?: boolean;
  label?: string;
  id?: string;
}

export const ToggleSwitch: React.FC<ToggleSwitchProps> = ({
  checked,
  onChange,
  disabled = false,
  label,
  id,
}) => {
  const generatedId = useId();
  const toggleId = id || `toggle-${generatedId}`;

  return (
    <label className="toggle-switch" htmlFor={toggleId}>
      <input
        id={toggleId}
        type="checkbox"
        checked={checked}
        onChange={(e) => onChange(e.target.checked)}
        disabled={disabled}
      />
      <span className="slider"></span>
      {label && <span className="toggle-label">{label}</span>}
    </label>
  );
};
