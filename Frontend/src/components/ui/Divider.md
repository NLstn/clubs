# Divider Component

A versatile divider component for visually separating content sections in the UI.

## Features

- **Simple Divider**: Minimal horizontal line for clean separation
- **Text Divider**: Divider with centered text (e.g., "OR")
- **Spacing Control**: Three spacing options (sm, md, lg)
- **Consistent Styling**: Uses design system color tokens

## Usage

### Basic Divider

A simple horizontal line divider:

```tsx
import { Divider } from '@/components/ui';

<Divider />
```

### Divider with Text

Useful for separating alternative actions (e.g., login options):

```tsx
<Divider text="OR" />
```

### Custom Spacing

Control the margin around the divider:

```tsx
<Divider spacing="sm" />  // Small spacing
<Divider spacing="md" />  // Medium spacing (default)
<Divider spacing="lg" />  // Large spacing
```

### With Custom Class

Add additional styling as needed:

```tsx
<Divider className="my-custom-class" />
```

## Props

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `text` | `string` | `undefined` | Optional text to display in the center of the divider |
| `spacing` | `'sm' \| 'md' \| 'lg'` | `'md'` | Spacing variant for top and bottom margins |
| `className` | `string` | `''` | Additional CSS class name |

## Examples

### Login Page Pattern

```tsx
<form onSubmit={handleSubmit}>
  <Input name="email" type="email" />
  <Button type="submit">Send Magic Link</Button>
</form>

<Divider text="OR" />

<Button onClick={handleOAuth}>Login with OAuth</Button>
```

### Section Separator

```tsx
<div className="content-section">
  <h3>Recent Activities</h3>
  <ActivityList />
</div>

<Divider spacing="lg" />

<div className="content-section">
  <h3>Upcoming Events</h3>
  <EventList />
</div>
```

### List Separator

```tsx
<div className="list-container">
  {recentItems.map(item => (
    <div key={item.id}>
      <ListItem item={item} />
      <Divider spacing="sm" />
    </div>
  ))}
</div>
```

## Design Tokens Used

- `--color-border`: Divider line color
- `--color-text-secondary`: Text color for centered text
- `--space-sm`, `--space-md`, `--space-lg`, `--space-xl`: Spacing values

## Accessibility

- Dividers are purely visual and use CSS for rendering
- Text dividers maintain semantic meaning through text content
- Proper contrast ratios maintained with design system colors

## Related Components

- **Card**: Uses borders for section separation
- **Modal**: Uses dividers between header, body, and footer
- **Table**: Uses borders for row separation
