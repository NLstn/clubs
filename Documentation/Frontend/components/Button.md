# Button Component

## Overview

The `Button` component is an enhanced, reusable button component with built-in support for user feedback states including loading, success, and error indicators. It provides visual feedback during asynchronous operations and supports cancellable processes.

## Features

- Multiple visual variants (primary, secondary, accept, cancel, maybe)
- Three size options (small, medium, large)
- Built-in loading state with animated spinner
- Success state with checkmark icon and custom message
- Error state with X icon and custom message
- Cancellable operations with dedicated cancel button
- Counter badge support
- Full width option
- Accessibility compliant with proper ARIA attributes
- TypeScript support

## Import

```tsx
import { Button, ButtonState } from '@/components/ui';
```

## Basic Usage

### Simple Button

```tsx
<Button onClick={handleClick}>
  Click Me
</Button>
```

### Button Variants

```tsx
<Button variant="primary">Primary</Button>
<Button variant="secondary">Secondary</Button>
<Button variant="accept">Accept</Button>
<Button variant="cancel">Cancel</Button>
<Button variant="maybe">Maybe</Button>
```

### Button Sizes

```tsx
<Button size="sm">Small</Button>
<Button size="md">Medium</Button>
<Button size="lg">Large</Button>
```

## Feedback States

### Loading State

Display a loading spinner and disable the button during asynchronous operations:

```tsx
const [buttonState, setButtonState] = useState<ButtonState>('idle');

const handleSave = async () => {
  setButtonState('loading');
  try {
    await saveData();
    setButtonState('success');
    setTimeout(() => setButtonState('idle'), 3000);
  } catch (error) {
    setButtonState('error');
    setTimeout(() => setButtonState('idle'), 3000);
  }
};

<Button 
  state={buttonState}
  onClick={handleSave}
>
  Save
</Button>
```

### Success State

Display a checkmark icon with an optional success message:

```tsx
<Button 
  state="success"
  successMessage="Saved successfully!"
>
  Save
</Button>
```

If no `successMessage` is provided, the button's children text will be displayed with the checkmark icon.

### Error State

Display an X icon with an optional error message:

```tsx
<Button 
  state="error"
  errorMessage="Failed to save!"
>
  Save
</Button>
```

The button remains clickable in the error state, allowing users to retry the operation.

### Cancellable Operations

Add a cancel button during loading state to allow users to abort long-running operations:

```tsx
const [buttonState, setButtonState] = useState<ButtonState>('idle');
const abortControllerRef = useRef<AbortController | null>(null);

const handleLongOperation = async () => {
  setButtonState('loading');
  abortControllerRef.current = new AbortController();
  
  try {
    await longRunningOperation(abortControllerRef.current.signal);
    setButtonState('success');
    setTimeout(() => setButtonState('idle'), 3000);
  } catch (error) {
    if (error.name !== 'AbortError') {
      setButtonState('error');
      setTimeout(() => setButtonState('idle'), 3000);
    }
  }
};

const handleCancel = () => {
  abortControllerRef.current?.abort();
  setButtonState('idle');
};

<Button 
  state={buttonState}
  onCancel={handleCancel}
  onClick={handleLongOperation}
>
  Start Process
</Button>
```

## Additional Features

### Counter Badge

Display a notification counter badge on the button:

```tsx
<Button counter={5}>
  Notifications
</Button>
```

Features:
- Displays "99+" for values over 99
- Automatically sizes for two-digit numbers
- Customizable color
- Hidden during non-idle states

```tsx
<Button counter={25} counterColor="#ff5500">
  Alerts
</Button>
```

### Full Width

Make the button expand to fill its container:

```tsx
<Button fullWidth>
  Full Width Button
</Button>
```

## Complete Example

Here's a complete example showing all features in action:

```tsx
import { useState, useRef } from 'react';
import { Button, ButtonState } from '@/components/ui';

function SaveForm() {
  const [buttonState, setButtonState] = useState<ButtonState>('idle');
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  const handleSave = async () => {
    setButtonState('loading');
    
    // Simulate async operation
    timeoutRef.current = setTimeout(async () => {
      try {
        // Your save logic here
        await saveFormData();
        setButtonState('success');
        
        // Reset to idle after 3 seconds
        setTimeout(() => setButtonState('idle'), 3000);
      } catch (error) {
        setButtonState('error');
        setTimeout(() => setButtonState('idle'), 3000);
      }
    }, 2000);
  };

  const handleCancel = () => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
      timeoutRef.current = null;
    }
    setButtonState('idle');
  };

  return (
    <Button
      variant="primary"
      size="md"
      state={buttonState}
      successMessage="Saved successfully!"
      errorMessage="Save failed. Please try again."
      onCancel={handleCancel}
      onClick={handleSave}
    >
      Save Changes
    </Button>
  );
}
```

## Props API

### ButtonProps

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `variant` | `'primary' \| 'secondary' \| 'accept' \| 'cancel' \| 'maybe'` | `'primary'` | Visual variant of the button |
| `size` | `'sm' \| 'md' \| 'lg'` | `'md'` | Size of the button |
| `fullWidth` | `boolean` | `false` | Whether button should fill container width |
| `counter` | `number` | `undefined` | Optional counter badge value |
| `counterColor` | `string` | `'var(--color-cancel)'` | Background color of counter badge |
| `state` | `ButtonState` | `'idle'` | Current state of the button |
| `successMessage` | `string` | `undefined` | Message to display in success state |
| `errorMessage` | `string` | `undefined` | Message to display in error state |
| `onCancel` | `() => void` | `undefined` | Callback when cancel button is clicked |
| `children` | `React.ReactNode` | required | Button content |
| `disabled` | `boolean` | `false` | Whether button is disabled |
| ...rest | `HTMLButtonAttributes` | - | All standard button HTML attributes |

### ButtonState Type

```tsx
type ButtonState = 'idle' | 'loading' | 'success' | 'error';
```

## Behavior Notes

1. **Button Disabled States**:
   - `loading`: Button is disabled, shows spinner
   - `success`: Button is disabled, shows checkmark
   - `error`: Button remains enabled for retry
   - `idle`: Normal button behavior

2. **Counter Badge**:
   - Only displayed when `state='idle'` and `counter > 0`
   - Hidden during loading, success, and error states

3. **Cancel Button**:
   - Only displayed when `state='loading'` and `onCancel` prop is provided
   - Positioned outside the main button to avoid nesting issues
   - Has its own click handler that doesn't trigger the main button's onClick

4. **Accessibility**:
   - `aria-busy="true"` when state is loading
   - Proper focus management
   - Semantic HTML structure
   - All interactive elements are keyboard accessible

## Demo

A live demonstration of all Button features is available at `/demo/button` when running the development server.

## Testing

The Button component includes comprehensive tests covering:
- All variants and sizes
- All state transitions
- Counter badge behavior
- Cancel functionality
- Accessibility features
- Props forwarding

Run tests with:
```bash
npm run test -- src/components/ui/__tests__/Button.test.tsx
```

## Styling

The Button component uses CSS variables from the design system:

```css
--color-primary
--color-primary-hover
--color-secondary
--color-secondary-hover
--color-cancel
--color-cancel-hover
--space-sm, --space-md
--border-radius-sm
```

Custom styles can be applied via the `className` prop:

```tsx
<Button className="my-custom-button">
  Custom Styled
</Button>
```

## Migration from Old Button

If you have existing buttons without state management, no changes are required. The component is fully backward compatible:

```tsx
// Old usage - still works
<Button variant="primary" onClick={handleClick}>
  Click Me
</Button>

// New usage - with states
<Button variant="primary" state={buttonState} onClick={handleClick}>
  Click Me
</Button>
```

## Best Practices

1. **Always manage state**: Use `useState` to manage button states for async operations
2. **Provide feedback messages**: Use `successMessage` and `errorMessage` for better UX
3. **Reset state**: Always reset to 'idle' after showing success/error
4. **Handle cancellation**: Implement proper cleanup in `onCancel` handler
5. **Use appropriate variants**: Use 'accept' for positive actions, 'cancel' for destructive actions
6. **Accessible labels**: Ensure button text clearly indicates the action

## Related Components

- [Modal](./Modal.md) - Modal dialogs often use buttons
- [ConfirmDialog](./ConfirmDialog.md) - Confirmation dialogs with button actions
- [FormGroup](./FormGroup.md) - Form layout with buttons
