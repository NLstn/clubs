# Toast Component

A reusable toast notification system for displaying temporary messages to users across the application. The toast system provides a consistent way to show success, error, info, and warning messages.

## Overview

The toast system consists of:
- **ToastProvider**: React Context provider that manages toast state
- **ToastContainer**: Container component that renders all active toasts
- **ToastItem**: Individual toast notification component
- **useToast**: React hook for accessing toast functionality

## Features

- ✅ Multiple toast types (success, error, info, warning)
- ✅ Auto-dismiss with configurable duration
- ✅ Manual dismissal via close button
- ✅ Multiple toasts displayed simultaneously
- ✅ Smooth slide-in/out animations
- ✅ Responsive design (desktop and mobile)
- ✅ Accessible (ARIA roles and labels)
- ✅ Internationalization support
- ✅ TypeScript support

## Installation

The toast system is already integrated into the application. The `ToastProvider` is wrapped around the app in `App.tsx`, and the `ToastContainer` is rendered globally.

## Usage

### Basic Usage

Use the `useToast` hook to access toast functions in any component:

```tsx
import { useToast } from '@/hooks/useToast';

function MyComponent() {
  const toast = useToast();

  const handleSuccess = () => {
    toast.success('Operation completed successfully!');
  };

  const handleError = () => {
    toast.error('An error occurred. Please try again.');
  };

  return (
    <div>
      <button onClick={handleSuccess}>Show Success</button>
      <button onClick={handleError}>Show Error</button>
    </div>
  );
}
```

### Toast Methods

The `useToast` hook provides several convenience methods:

```tsx
const toast = useToast();

// Success toast (green)
toast.success('Saved successfully!');

// Error toast (red)
toast.error('Failed to save. Please try again.');

// Info toast (blue)
toast.info('Loading data...');

// Warning toast (orange)
toast.warning('You have unsaved changes.');

// Custom toast with all options
toast.addToast('Custom message', 'success', 3000);
```

### Custom Duration

By default, toasts auto-dismiss after 5000ms (5 seconds). You can customize this:

```tsx
// Show for 3 seconds
toast.success('Quick message', 3000);

// Show for 10 seconds
toast.error('Important error message', 10000);

// Never auto-dismiss (pass 0 or negative number)
toast.info('Permanent message', 0);
```

### With i18n Translations

Use translation keys for internationalized messages:

```tsx
import { useT } from '@/hooks/useTranslation';

function MyComponent() {
  const { t } = useT();
  const toast = useToast();

  const handleSave = async () => {
    try {
      await saveData();
      toast.success(t('toast.success.saved'));
    } catch (error) {
      toast.error(t('toast.error.save'));
    }
  };

  return <button onClick={handleSave}>Save</button>;
}
```

### Common Translation Keys

The following toast translations are available in both English and German:

```tsx
// Success messages
t('toast.success.saved')    // "Saved successfully!"
t('toast.success.created')  // "Created successfully!"
t('toast.success.updated')  // "Updated successfully!"
t('toast.success.deleted')  // "Deleted successfully!"
t('toast.success.copied')   // "Copied to clipboard!"

// Error messages
t('toast.error.generic')    // "An error occurred. Please try again."
t('toast.error.network')    // "Network error. Please check your connection."
t('toast.error.save')       // "Failed to save. Please try again."
t('toast.error.load')       // "Failed to load data. Please try again."
t('toast.error.delete')     // "Failed to delete. Please try again."
t('toast.error.permission') // "You don't have permission..."

// Info messages
t('toast.info.loading')     // "Loading..."
t('toast.info.processing')  // "Processing..."

// Warning messages
t('toast.warning.unsavedChanges') // "You have unsaved changes."
t('toast.warning.confirmAction')  // "Please confirm this action."
```

## Examples

### Form Submission

```tsx
import { useToast } from '@/hooks/useToast';
import { useT } from '@/hooks/useTranslation';

function CreateEventForm() {
  const toast = useToast();
  const { t } = useT();
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (data: EventData) => {
    setIsSubmitting(true);
    
    try {
      await api.post('/api/v2/Events', data);
      toast.success(t('toast.success.created'));
      // Navigate or reset form
    } catch (error) {
      toast.error(t('toast.error.save'));
    } finally {
      setIsSubmitting(false);
    }
  };

  return <form onSubmit={handleSubmit}>...</form>;
}
```

### Delete Confirmation

```tsx
import { useToast } from '@/hooks/useToast';
import { useT } from '@/hooks/useTranslation';

function DeleteButton({ itemId, onDeleted }) {
  const toast = useToast();
  const { t } = useT();

  const handleDelete = async () => {
    if (!confirm('Are you sure you want to delete this item?')) {
      return;
    }

    try {
      await api.delete(`/api/v2/Items/${itemId}`);
      toast.success(t('toast.success.deleted'));
      onDeleted();
    } catch (error) {
      toast.error(t('toast.error.delete'));
    }
  };

  return <button onClick={handleDelete}>Delete</button>;
}
```

### API Error Handling

```tsx
import { useToast } from '@/hooks/useToast';
import { useT } from '@/hooks/useTranslation';
import axios from 'axios';

function DataFetcher() {
  const toast = useToast();
  const { t } = useT();

  const fetchData = async () => {
    try {
      const response = await api.get('/api/v2/Data');
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        if (error.response?.status === 403) {
          toast.error(t('toast.error.permission'));
        } else if (error.response?.status >= 500) {
          toast.error(t('toast.error.network'));
        } else {
          toast.error(t('toast.error.load'));
        }
      } else {
        toast.error(t('toast.error.generic'));
      }
      throw error;
    }
  };

  return ...;
}
```

### Copy to Clipboard

```tsx
import { useToast } from '@/hooks/useToast';
import { useT } from '@/hooks/useTranslation';

function CopyButton({ text }) {
  const toast = useToast();
  const { t } = useT();

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(text);
      toast.success(t('toast.success.copied'));
    } catch (error) {
      toast.error('Failed to copy to clipboard');
    }
  };

  return <button onClick={handleCopy}>Copy</button>;
}
```

## Styling

The toast system uses CSS variables from the design system for consistent theming:

- **Success**: `--color-success-text`, `--color-success-bg`
- **Error**: `--color-error-text`, `--color-error-bg`
- **Warning**: `--color-maybe`
- **Info**: `--color-secondary`

The toasts are fully responsive and adapt to mobile screens with different animations.

## Accessibility

The toast system follows accessibility best practices:

- Uses `role="alert"` for screen reader announcements
- Provides `aria-live="polite"` for non-intrusive updates
- Includes `aria-label` for close buttons
- Keyboard accessible (close buttons can be focused and activated)

## API Reference

### useToast Hook

```typescript
interface ToastContextType {
  toasts: Toast[];
  addToast: (message: string, type?: ToastType, duration?: number) => void;
  removeToast: (id: string) => void;
  success: (message: string, duration?: number) => void;
  error: (message: string, duration?: number) => void;
  info: (message: string, duration?: number) => void;
  warning: (message: string, duration?: number) => void;
}
```

### Toast Types

```typescript
type ToastType = 'success' | 'error' | 'info' | 'warning';

interface Toast {
  id: string;
  message: string;
  type: ToastType;
  duration?: number;
}
```

## Best Practices

1. **Use appropriate toast types**: Match the toast type to the message content
2. **Keep messages concise**: Toasts should be brief and easy to scan
3. **Use translations**: Always use i18n translation keys for user-facing text
4. **Don't overuse**: Toasts are for temporary notifications, not permanent information
5. **Consider duration**: Important errors may need longer display times
6. **Provide context**: Include enough information for users to understand the message
7. **Avoid technical jargon**: Use user-friendly language in toast messages

## Migration Guide

To migrate from existing alert/error patterns to toasts:

### Before
```tsx
// Using alert()
alert('Invite link copied to clipboard!');

// Using setError state
setError("Failed to add event: " + error.message);
```

### After
```tsx
// Using toast
import { useToast } from '@/hooks/useToast';
import { useT } from '@/hooks/useTranslation';

const toast = useToast();
const { t } = useT();

// Success
toast.success(t('toast.success.copied'));

// Error
toast.error(t('toast.error.save'));
```

## Testing

The toast system includes comprehensive test coverage. To test components that use toasts:

```tsx
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { ToastProvider } from '@/context/ToastProvider';
import MyComponent from './MyComponent';

test('shows success toast', async () => {
  const user = userEvent.setup();
  
  render(
    <ToastProvider>
      <MyComponent />
    </ToastProvider>
  );

  await user.click(screen.getByText('Save'));
  
  expect(screen.getByText('Saved successfully!')).toBeInTheDocument();
});
```

## Troubleshooting

### Toast not appearing
- Ensure your component is wrapped in `ToastProvider`
- Check that you're importing `useToast` correctly
- Verify the toast message is not empty

### Toast not auto-dismissing
- Check that you haven't passed `0` or negative duration
- Ensure timers are working (not mocked in tests without proper cleanup)

### Multiple toasts overlapping
- This is expected behavior; toasts stack vertically
- Consider reducing the number of simultaneous toasts if needed

---

**Related Components:**
- [Modal](./Modal.md) - For blocking user interactions
- [ConfirmDialog](./ConfirmDialog.md) - For confirmation prompts
- [Input](./Input.md) - Form inputs with error states
