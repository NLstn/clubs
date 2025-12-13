# ConfirmDialog Component

A reusable confirmation dialog built on top of the Modal component. Used to confirm destructive or important actions before they are executed.

## Features

- Built on top of the Modal component for consistent styling
- Three variants: `danger`, `warning`, and `info`
- Loading state support
- Customizable confirm/cancel button text
- Accessible with proper ARIA labels
- Mobile responsive

## Basic Usage

```tsx
import { ConfirmDialog } from '@/components/ui';
import { useState } from 'react';

function MyComponent() {
  const [isOpen, setIsOpen] = useState(false);
  
  const handleDelete = () => {
    // Perform delete action
    console.log('Item deleted');
    setIsOpen(false);
  };
  
  return (
    <>
      <button onClick={() => setIsOpen(true)}>Delete Item</button>
      
      <ConfirmDialog
        isOpen={isOpen}
        onClose={() => setIsOpen(false)}
        onConfirm={handleDelete}
        title="Delete Item"
        message="Are you sure you want to delete this item? This action cannot be undone."
        variant="danger"
      />
    </>
  );
}
```

## With Loading State

```tsx
function MyComponent() {
  const [isOpen, setIsOpen] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);
  
  const handleDelete = async () => {
    setIsDeleting(true);
    try {
      await api.delete('/api/item');
      setIsOpen(false);
    } catch (error) {
      console.error('Failed to delete:', error);
    } finally {
      setIsDeleting(false);
    }
  };
  
  return (
    <ConfirmDialog
      isOpen={isOpen}
      onClose={() => setIsOpen(false)}
      onConfirm={handleDelete}
      title="Delete Item"
      message="Are you sure you want to delete this item?"
      variant="danger"
      isLoading={isDeleting}
      confirmText="Delete"
      cancelText="Cancel"
    />
  );
}
```

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `isOpen` | `boolean` | Required | Controls whether the dialog is visible |
| `onClose` | `() => void` | Required | Called when the dialog should be closed (Cancel button or backdrop click) |
| `onConfirm` | `() => void` | Required | Called when the user confirms the action |
| `title` | `string` | Required | The title displayed in the dialog header |
| `message` | `string` | Required | The message/question displayed in the dialog body |
| `confirmText` | `string` | `"Confirm"` | Text for the confirm button |
| `cancelText` | `string` | `"Cancel"` | Text for the cancel button |
| `variant` | `"danger" \| "warning" \| "info"` | `"danger"` | Visual style variant of the dialog |
| `isLoading` | `boolean` | `false` | Shows loading spinner and disables buttons during async operations |

## Variants

### Danger (default)
Use for destructive actions like deleting data.
```tsx
<ConfirmDialog
  variant="danger"
  title="Delete Account"
  message="This will permanently delete your account and all associated data."
  // ...
/>
```

### Warning
Use for actions that require caution but aren't destructive.
```tsx
<ConfirmDialog
  variant="warning"
  title="Unsaved Changes"
  message="You have unsaved changes. Are you sure you want to leave?"
  // ...
/>
```

### Info
Use for informational confirmations.
```tsx
<ConfirmDialog
  variant="info"
  title="Confirm Action"
  message="Do you want to proceed with this action?"
  // ...
/>
```

## Examples

### Simple Delete Confirmation
```tsx
const [showConfirm, setShowConfirm] = useState(false);

<ConfirmDialog
  isOpen={showConfirm}
  onClose={() => setShowConfirm(false)}
  onConfirm={() => {
    deleteItem();
    setShowConfirm(false);
  }}
  title="Delete Session"
  message="Are you sure you want to delete this session?"
  variant="danger"
/>
```

### With Async Operation
```tsx
const [showConfirm, setShowConfirm] = useState(false);
const [isDeleting, setIsDeleting] = useState(false);

<ConfirmDialog
  isOpen={showConfirm}
  onClose={() => setShowConfirm(false)}
  onConfirm={async () => {
    setIsDeleting(true);
    try {
      await deleteSession();
      setShowConfirm(false);
    } finally {
      setIsDeleting(false);
    }
  }}
  title="Delete Session"
  message="Are you sure you want to delete this session?"
  variant="danger"
  isLoading={isDeleting}
/>
```

### Custom Button Text
```tsx
<ConfirmDialog
  isOpen={showConfirm}
  onClose={() => setShowConfirm(false)}
  onConfirm={handleLogout}
  title="Log Out"
  message="Are you sure you want to log out from all devices?"
  confirmText="Log Out"
  cancelText="Stay Logged In"
  variant="warning"
/>
```

## Design

The ConfirmDialog follows the application's design system:
- Uses the Modal component as its foundation
- Applies variant-specific header colors
- Consistent spacing and typography
- Mobile responsive
- Dark theme support

## Accessibility

- Proper ARIA labels on buttons
- Keyboard support (Escape to close)
- Focus management
- Screen reader friendly

## Best Practices

1. **Use descriptive titles**: Make it clear what action is being confirmed
2. **Write clear messages**: Explain the consequences of the action
3. **Choose appropriate variants**: Use `danger` for destructive actions
4. **Handle loading states**: Show loading feedback for async operations
5. **Don't overuse**: Only use for important or destructive actions
6. **Provide context**: Include relevant details in the message

## Related Components

- [Modal](./Modal.md) - The underlying modal component
- [Button](./Button.md) - Button component used for actions
