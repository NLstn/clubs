<div align="center">
  <img src="../../assets/logo.png" alt="Clubs Logo" width="120"/>
  
  # Input Component
  
  **Reusable form input with consistent styling**
</div>

---

# Input Component

A reusable input component that can be used across the application. This component provides consistent styling and behavior for text inputs.

## Features

- **Accessibility**: Proper label association and keyboard navigation
- **Variants**: Multiple visual styles (default, outline, filled)
- **Sizes**: Three size options (sm, md, lg)
- **Error handling**: Built-in error state and error message display
- **Helper text**: Optional helper text for additional guidance
- **Flexible**: Supports all standard HTML input attributes

## Usage

### Basic Usage

```tsx
import { Input } from '@/components/ui';

<Input
  label="Username"
  placeholder="Enter your username"
  onChange={(e) => setUsername(e.target.value)}
/>
```

### With Error State

```tsx
<Input
  label="Email"
  type="email"
  error="Please enter a valid email address"
  value={email}
  onChange={(e) => setEmail(e.target.value)}
/>
```

### Different Variants

```tsx
{/* Default variant */}
<Input variant="default" label="Default" />

{/* Outline variant */}
<Input variant="outline" label="Outline" />

{/* Filled variant */}
<Input variant="filled" label="Filled" />
```

### Different Sizes

```tsx
{/* Small */}
<Input size="sm" label="Small Input" />

{/* Medium (default) */}
<Input size="md" label="Medium Input" />

{/* Large */}
<Input size="lg" label="Large Input" />
```

### With Helper Text

```tsx
<Input
  label="Password"
  type="password"
  helperText="Password must be at least 8 characters long"
/>
```

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `label` | `string` | - | Label text for the input |
| `error` | `string` | - | Error message to display |
| `helperText` | `string` | - | Helper text to display below the input |
| `variant` | `'default' \| 'outline' \| 'filled'` | `'default'` | Visual variant of the input |
| `size` | `'sm' \| 'md' \| 'lg'` | `'md'` | Size of the input |
| `...props` | `InputHTMLAttributes<HTMLInputElement>` | - | All standard HTML input attributes |

## Integration with TypeAheadDropdown

The `TypeAheadDropdown` component has been updated to use this Input component internally, ensuring consistent styling across all input fields in the application.

## Styling

The component uses CSS classes that can be customized:

- `.input-container` - Container wrapper
- `.input-label` - Label styling
- `.input-base` - Base input styles
- `.input-{variant}` - Variant-specific styles
- `.input-{size}` - Size-specific styles
- `.input-error` - Error state styles
- `.input-error-text` - Error message styles
- `.input-helper-text` - Helper text styles

## Accessibility

- Proper label association using `htmlFor` and `id`
- Focus management with visible focus indicators
- Support for screen readers
- Keyboard navigation support
