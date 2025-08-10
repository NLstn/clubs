# Button Component

A reusable button component with customizable variants, sizes, and an optional counter display.

## Features

- **Multiple Variants**: `primary`, `secondary`, `accept`, `cancel`, `maybe`
- **Different Sizes**: `sm`, `md`, `lg`
- **Full Width Option**: Can span the full width of its container
- **Counter Badge**: Display a counter in the top-right corner
- **Accessibility**: Proper focus states and touch-friendly sizing
- **Responsive**: Mobile-optimized touch targets

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `variant` | `'primary' \| 'secondary' \| 'accept' \| 'cancel' \| 'maybe'` | `'primary'` | Visual style variant |
| `size` | `'sm' \| 'md' \| 'lg'` | `'md'` | Button size |
| `fullWidth` | `boolean` | `false` | Whether button should take full width |
| `counter` | `number` | `undefined` | Number to display in counter badge (hidden if 0 or undefined) |
| `counterColor` | `string` | `'var(--color-cancel)'` | Background color for the counter badge |
| `children` | `React.ReactNode` | - | Button content |

All standard HTML button props are also supported.

## Usage Examples

### Basic Button
```tsx
import { Button } from '@/components/ui';

<Button onClick={handleClick}>
  Click Me
</Button>
```

### Button with Counter
```tsx
<Button 
  variant="primary" 
  counter={5}
  onClick={handleViewRequests}
>
  View Join Requests
</Button>
```

### Different Variants
```tsx
<Button variant="accept">Save</Button>
<Button variant="cancel">Delete</Button>
<Button variant="maybe">Maybe Later</Button>
```

### Different Sizes
```tsx
<Button size="sm">Small</Button>
<Button size="md">Medium</Button>
<Button size="lg">Large</Button>
```

### Full Width
```tsx
<Button fullWidth>Full Width Button</Button>
```

### Custom Counter Color
```tsx
<Button 
  counter={3}
  counterColor="var(--color-primary)"
>
  Notifications
</Button>
```

## Counter Behavior

- The counter is only displayed when the `counter` prop is provided and greater than 0
- Numbers greater than 99 are displayed as "99+"
- The counter appears as a small badge in the top-right corner of the button
- The counter has a subtle shadow and border for better visibility

## Styling

The component uses CSS custom properties (CSS variables) for consistent theming:
- `--color-primary` and `--color-primary-hover` for primary buttons
- `--color-cancel` and `--color-cancel-hover` for cancel buttons
- Spacing uses `--space-*` variables
- Border radius uses `--border-radius-*` variables

## Accessibility

- Proper focus states with visible focus rings
- Touch-friendly minimum heights (44px on mobile, 48px on phones)
- Disabled state with reduced opacity and no-cursor
- Semantic button element for screen readers
