# UI Components Library

This document provides an overview of all reusable UI components available in the application.

## Overview

The UI components are located in `Frontend/src/components/ui/` and are designed to be reusable across the entire application. Each component follows consistent design patterns and supports theming through CSS custom properties.

## Design System

The components follow these design principles:
- **Dark Theme**: All components are designed for dark backgrounds
- **Consistent Spacing**: Uses CSS custom properties for consistent spacing (`--space-sm`, `--space-md`, `--space-lg`)
- **Responsive Design**: Components adapt to different screen sizes
- **Accessibility**: Components are built with accessibility in mind
- **TypeScript**: All components are fully typed with TypeScript

## Available Components

### Data Display

- **[Table](./components/Table.md)** - A flexible, reusable table component with support for loading states, error handling, and custom rendering

### Form Controls

- **[TypeAheadDropdown](./components/TypeAheadDropdown.md)** - A type-ahead dropdown component with search functionality and autocomplete suggestions

## Component Structure

Each component follows this structure:
```
ComponentName.tsx    // Main component file
ComponentName.css    // Component-specific styles
```

## Usage Guidelines

1. **Import Components**: Import components from their specific files in the ui folder:
   ```tsx
   import Table from '@/components/ui/Table';
   import TypeAheadDropdown from '@/components/ui/TypeAheadDropdown';
   ```

2. **TypeScript**: All components are fully typed. Use the exported interfaces for proper type safety.

3. **Styling**: Components come with default styles but can be customized through:
   - CSS custom properties (recommended)
   - Additional CSS classes via `className` prop
   - Inline styles (for specific cases)

4. **Responsive Design**: Use the provided responsive classes when needed:
   - `.hide-mobile` - Hides on screens < 768px
   - `.hide-small` - Hides on screens < 480px

## CSS Custom Properties

The components use these CSS custom properties:
- `--space-sm` - Small spacing
- `--space-md` - Medium spacing
- `--space-lg` - Large spacing
- `--border-radius-sm` - Small border radius
- `--border-radius-md` - Medium border radius
- `--shadow-sm` - Small shadow

## Contributing

When adding new components:

1. **Create the component** in `Frontend/src/components/ui/`
2. **Follow naming conventions**: PascalCase for component files
3. **Include TypeScript interfaces** for all props
4. **Add CSS file** with matching name for styling
5. **Add responsive design** considerations
6. **Create comprehensive documentation** in `Documentation/Frontend/components/`
7. **Update this main index file** with component description
8. **Add unit tests** in `__tests__` folder
9. **Export interfaces** that other components might need

### Component Template

```tsx
// ComponentName.tsx
import React from 'react';
import './ComponentName.css';

export interface ComponentNameProps {
    // Define props here
}

function ComponentName({ }: ComponentNameProps) {
    return (
        <div className="component-name">
            {/* Component content */}
        </div>
    );
}

export default ComponentName;
```

### CSS Template

```css
/* ComponentName.css */
.component-name {
    /* Base styles using CSS custom properties */
    padding: var(--space-md);
    border-radius: var(--border-radius-sm);
    background: #1a1a1a;
    color: rgba(255, 255, 255, 0.9);
}

/* Responsive design */
@media (max-width: 768px) {
    .component-name {
        padding: var(--space-sm);
    }
}
```

## Color Scheme

Components use the following color scheme:
- Background: `#1a1a1a`
- Secondary Background: `#333`
- Borders: `#444`, `#555`
- Text: `rgba(255, 255, 255, 0.85-0.95)`
- Hover States: `#2a2a2a`
- Error: `#f44336`
- Success: `#4caf50`
- Warning: `#ff9800`
- Info: `#2196f3`
