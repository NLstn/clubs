# Frontend Design System & Component Library

## Overview

This document provides a comprehensive guide to the UI design system for the Clubs Management Application. It includes the complete color system, spacing variables, typography guidelines, and component library to ensure consistency, accessibility, and maintainability across the entire application.

## 🎨 Design Principles

1. **Consistency**: Uniform visual language and interaction patterns
2. **Accessibility**: WCAG 2.1 AA compliant design with proper contrast ratios
3. **Responsiveness**: Mobile-first approach with seamless adaptation
4. **Usability**: Intuitive navigation and clear visual hierarchy
5. **Performance**: Optimized for fast loading and smooth interactions
6. **Dark Theme**: Modern dark-themed design for better user experience

## 🌈 Color System

### Primary Colors

| Color Name | Hex Code | RGB | CSS Variable | Usage |
|------------|----------|-----|--------------|-------|
| Primary Green | `#4CAF50` | `rgb(76, 175, 80)` | `--color-primary` | Primary actions, success states, brand elements |
| Primary Green Hover | `#45a049` | `rgb(69, 160, 73)` | `--color-primary-hover` | Hover state for primary buttons |
| Secondary Blue | `#646cff` | `rgb(100, 108, 255)` | `--color-secondary` | Secondary actions, links, informational elements |
| Secondary Blue Hover | `#535bf2` | `rgb(83, 91, 242)` | `--color-secondary-hover` | Hover state for secondary buttons |

### Background Colors

| Color Name | Hex Code | RGB | CSS Variable | Usage |
|------------|----------|-----|--------------|-------|
| Background Dark | `#242424` | `rgb(36, 36, 36)` | `--color-background` | Main application background |
| Background Light | `#333333` | `rgb(51, 51, 51)` | `--color-background-light` | Card backgrounds, sections |

### Text Colors

| Color Name | Hex Code | RGB | CSS Variable | Usage |
|------------|----------|-----|--------------|-------|
| Text Primary | `rgba(255, 255, 255, 0.87)` | `rgb(255, 255, 255, 87%)` | `--color-text` | Main text content |
| Text Secondary | `#888888` | `rgb(136, 136, 136)` | `--color-text-secondary` | Supporting text, metadata |

### System Colors

| Color Name | Hex Code | RGB | CSS Variable | Usage |
|------------|----------|-----|--------------|-------|
| Error/Cancel | `#f44336` | `rgb(244, 67, 54)` | `--color-cancel` | Destructive actions, errors |
| Error/Cancel Hover | `#e53935` | `rgb(229, 57, 53)` | `--color-cancel-hover` | Hover state for destructive actions |
| Success Background | `#d4edda` | `rgb(212, 237, 218)` | `--color-success-bg` | Success message backgrounds |
| Success Text | `#155724` | `rgb(21, 87, 36)` | `--color-success-text` | Success message text |
| Error Background | `#f8d7da` | `rgb(248, 215, 218)` | `--color-error-bg` | Error message backgrounds |
| Error Text | `#721c24` | `rgb(114, 28, 36)` | `--color-error-text` | Error message text |
| Border | `#dddddd` | `rgb(221, 221, 221)` | `--color-border` | Default borders, separators |

### Color Usage Guidelines

- **Primary Green (#4CAF50)**: Use for primary actions, success states, and brand elements
- **Secondary Blue (#646cff)**: Use for secondary actions, links, and informational elements
- **Red (#f44336)**: Use sparingly for destructive actions and error states
- **Dark Backgrounds**: Maintain the dark theme throughout the application
- **Text Contrast**: Ensure minimum 4.5:1 contrast ratio for accessibility

## 📏 Spacing System

### Spacing Variables

The application uses a consistent 8px-based spacing system with responsive adjustments for mobile devices:

| Variable | Desktop Value | Mobile Value | Usage |
|----------|---------------|--------------|-------|
| `--space-xs` | `0.5rem` (8px) | `0.4rem` (6.4px) | Small gaps, icon spacing |
| `--space-sm` | `1rem` (16px) | `0.8rem` (12.8px) | Component padding, small margins |
| `--space-md` | `1.5rem` (24px) | `1.2rem` (19.2px) | Section spacing, form groups |
| `--space-lg` | `2rem` (32px) | `1.6rem` (25.6px) | Page sections, large components |
| `--space-xl` | `3rem` (48px) | `2.4rem` (38.4px) | Page margins, major sections |

### Border Radius Variables

| Variable | Value | Usage |
|----------|-------|-------|
| `--border-radius-sm` | `4px` | Small components, tags |
| `--border-radius-md` | `6px` | Buttons, input fields |
| `--border-radius-lg` | `8px` | Cards, modals |
| `--border-radius-circle` | `50%` | Circular elements, avatars |

### Shadow Variables

| Variable | Value | Usage |
|----------|-------|-------|
| `--shadow-sm` | `0 2px 6px rgba(0, 0, 0, 0.1)` | Small elevation |
| `--shadow-md` | `0 4px 12px rgba(0, 0, 0, 0.15)` | Medium elevation, dropdowns |

## 📝 Typography

### Font Family
```css
font-family: Inter, system-ui, Avenir, Helvetica, Arial, sans-serif;
```

### Font Hierarchy

| Element | Desktop | Tablet | Mobile | Weight | Usage |
|---------|---------|--------|--------|--------|-------|
| **H1** | 3.2rem | 2.2rem | 1.8rem | 600-700 | Page titles, main headings |
| **H2** | 1.8rem | 1.5rem | 1.3rem | 600 | Section titles, dashboard sections |
| **H3** | 1.2rem | 1.2rem | 1.2rem | 600 | Subsection titles, card headers |
| **H4** | 1.1rem | 1.1rem | 1.1rem | 500-600 | Component titles, form sections |
| **Body** | 1rem | 1rem | 1rem | 400 | Regular text content |
| **Small** | 0.9rem | 0.9rem | 0.9rem | 400 | Supporting text, metadata |

**Line Height**: 1.5 for optimal readability across all text elements.

## 🧱 Component Library

### Available Components

#### Core UI Components
- **[Input](./components/Input.md)** - Reusable input component with variants, sizes, error handling, and accessibility features
- **[Table](./components/Table.md)** - Flexible, reusable table component with loading states and error handling
- **[ODataTable](./components/ODataTable.md)** - Table with server-side pagination and sorting via OData v4 queries
- **[TypeAheadDropdown](./components/TypeAheadDropdown.md)** - Type-ahead dropdown with search functionality and autocomplete
- **[Modal](./components/Modal.md)** - Enhanced modal component with compound pattern, error handling, and loading states

#### Layout Components
- **Header** - Main navigation header with user menu and notifications
- **GlobalSearch** - Application-wide search functionality
- **NotificationDropdown** - User notification center
- **RecentClubsDropdown** - Quick access to recently visited clubs

#### Utility Components
- **ProtectedRoute** - Authentication wrapper for protected pages
- **LanguageSwitcher** - Interface language selection
- **CookieConsent** - GDPR compliance cookie consent banner

### Button Components

#### Primary Button
```css
.button-primary {
  background-color: var(--color-primary);
  color: white;
  border: none;
  padding: var(--space-sm) var(--space-md);
  border-radius: var(--border-radius-md);
  font-size: 1rem;
  font-weight: 500;
  min-height: 44px;
  cursor: pointer;
  transition: background-color 0.2s;
}
```
**Usage**: Main actions, form submissions, primary navigation

#### Secondary Button
```css
.button-secondary {
  background-color: var(--color-secondary);
  color: white;
  /* Same styling as primary but different color */
}
```
**Usage**: Secondary actions, alternative options

#### Destructive Button
```css
.button-cancel {
  background-color: var(--color-cancel);
  color: white;
  /* Same styling as primary but red color */
}
```
**Usage**: Dangerous or irreversible actions

### Form Components

#### Input Component
The reusable Input component provides consistent styling and behavior across the application:

```tsx
import { Input } from '@/components/ui';

// Basic usage
<Input label="Username" placeholder="Enter username" />

// With error state
<Input 
  label="Email" 
  error="Please enter a valid email" 
  type="email" 
/>

// Different variants and sizes
<Input variant="outline" size="lg" label="Large Outline Input" />
```

**Features:**
- Three variants: `default`, `outline`, `filled`
- Three sizes: `sm`, `md`, `lg`
- Built-in error handling and helper text
- Full accessibility support with proper labeling
- TypeScript support with all HTML input attributes

#### Input Fields
```css
.form-group input {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid var(--color-border);
  border-radius: var(--border-radius-md);
  font-size: 1rem;
  background-color: var(--color-background-light);
  color: var(--color-text);
}
```

#### Form Groups
```css
.form-group {
  margin-bottom: var(--space-md);
}

.form-group label {
  display: block;
  margin-bottom: var(--space-xs);
  font-weight: 500;
  color: var(--color-text);
}
```

## 📱 Responsive Design

### Breakpoints
- **Mobile**: ≤480px (single column, large touch targets)
- **Tablet**: 481px - 768px (moderate layouts)
- **Desktop**: >768px (full layouts with hover states)

### Mobile-First Philosophy
- Start with mobile design and enhance for larger screens
- Touch-friendly interactions (minimum 44px targets)
- Readable text without zooming (16px minimum)
- Stacked layouts for better mobile usability

### Mobile Utility Classes
```css
.mobile-hide { display: none !important; }       /* Hide on mobile */
.mobile-full-width { width: 100% !important; }   /* Full width on mobile */
.mobile-center { text-align: center !important; } /* Center text on mobile */
.mobile-stack { flex-direction: column !important; } /* Stack elements on mobile */
```

## ♿ Accessibility Guidelines

### WCAG 2.1 AA Compliance
- **Contrast**: Minimum 4.5:1 for normal text, 3:1 for large text
- **Keyboard Navigation**: Full keyboard accessibility with proper focus indicators
- **Focus Management**: Clear focus indicators and logical tab order
- **Screen Readers**: Semantic HTML and proper ARIA usage
- **Color Independence**: Information not conveyed by color alone
- **Touch Targets**: Minimum 44px for touch targets (48px on mobile)

### Implementation Guidelines
- Use semantic HTML elements
- Provide alternative text for images
- Ensure proper heading hierarchy
- Test with keyboard navigation
- Verify screen reader compatibility

## 🛠️ Implementation Guidelines

### CSS Architecture

#### Root `index.css`
The root `src/index.css` file should contain **only**:
- **CSS Variables**: All design tokens (colors, spacing, typography, etc.)
- **Global Resets**: Basic HTML element resets (body, html, *, etc.)
- **Base Typography**: Font family, line height, and base text styling
- **Universal Box Model**: Universal box-sizing rules

**What NOT to include in `index.css`:**
- Component-specific styles
- Layout styles
- Button styles
- Form styles
- Page-specific styles

#### Component Styles
- **Reusable Components**: Create components in `Frontend/src/components/ui/` for commonly used UI patterns (buttons, inputs, cards, modals, etc.)
- **Component CSS Files**: Each component should have its own CSS file (e.g., `Button.tsx` → `Button.css`)
- **Page/Component-Specific Styles**: Custom styles that are only used in specific pages or components should be in their own CSS files, not in `index.css`

#### CSS Best Practices
- **Use CSS Variables**: Always use CSS variables for colors, spacing, and other design tokens
- **Mobile-First**: Write responsive CSS starting with mobile
- **BEM Methodology**: Follow consistent naming conventions for component classes
- **Performance**: Optimize selectors and minimize specificity

### Component Development
1. **Create components** in `Frontend/src/components/ui/`
2. **Follow naming conventions**: PascalCase for component files
3. **Include TypeScript interfaces** for all props
4. **Add CSS file** with matching name for styling
5. **Implement responsive design** considerations
6. **Add comprehensive documentation** in `Documentation/Frontend/components/`
7. **Include unit tests** in `__tests__` folder

### Component Usage
```tsx
// Import components from UI folder
import { Input, Table, TypeAheadDropdown } from '@/components/ui';

// Use TypeScript interfaces for type safety
interface MyComponentProps {
  data: TableData[];
  onSelect: (item: DropdownItem) => void;
}

// Example usage of Input component
<Input
  label="Search"
  placeholder="Search for clubs..."
  value={searchTerm}
  onChange={(e) => setSearchTerm(e.target.value)}
  helperText="Start typing to search"
/>
```

### Internationalization (i18n)

The application uses **i18next** for internationalization, supporting multiple languages (English and German).

#### Translation Hook
Use the `useT` hook to access the translation function:

```tsx
import { useT } from '@/hooks/useTranslation';

function MyComponent() {
  const { t } = useT();
  
  return <h1>{t('common.welcome')}</h1>;
}
```

#### **IMPORTANT: Do NOT Use Fallback Texts**

**Never** use fallback text with the `||` operator when calling the translation function `t()`:

```tsx
// ❌ WRONG - Do not use fallback text
<h1>{t('common.welcome') || 'Welcome'}</h1>
<span>{t('common.name') || 'Name'}</span>

// ✅ CORRECT - Use translation key only
<h1>{t('common.welcome')}</h1>
<span>{t('common.name')}</span>
```

**Rationale:**
- The `t()` function automatically returns the translation key if no translation is found
- Fallback texts hide missing translations, making it impossible to identify untranslated strings
- This makes it difficult to maintain complete translations across all supported languages
- By removing fallbacks, missing translations become immediately visible during development and testing

#### Translation Best Practices

1. **Use descriptive translation keys**: `clubs.members.addButton` instead of `button1`
2. **Organize keys hierarchically**: Group related translations under common prefixes
3. **Test in all languages**: Verify translations in both English and German
4. **Check for missing translations**: Look for displayed translation keys (e.g., "common.welcome") in the UI
5. **Add translations before committing**: Ensure all new translation keys have entries in both language files

#### Translation Files Location
- English: `Frontend/src/i18n/locales/en.json`
- German: `Frontend/src/i18n/locales/de.json`

### File Structure
```
Frontend/src/components/
├── ui/                     # Reusable UI components
│   ├── Input.tsx
│   ├── Input.css
│   ├── Table.tsx
│   ├── Table.css
│   ├── TypeAheadDropdown.tsx
│   ├── TypeAheadDropdown.css
│   └── index.ts           # Component exports
├── layout/                 # Layout components
│   ├── Header.tsx
│   ├── GlobalSearch.tsx
│   └── NotificationDropdown.tsx
├── auth/                   # Authentication components
│   └── ProtectedRoute.tsx
└── __tests__/             # Component tests
```

## 🔄 Maintenance & Contributing

### Design System Updates
- Regular review of component usage and effectiveness
- User feedback integration for continuous improvement
- Performance monitoring and optimization
- Accessibility audits and enhancements

### Contributing Guidelines
1. Follow established patterns when creating new components
2. Document any new patterns or exceptions
3. Test thoroughly across devices and accessibility tools
4. Update documentation when making changes
5. Ensure all new components follow the design system

### Component Template
```tsx
// ComponentName.tsx
import React from 'react';
import './ComponentName.css';

export interface ComponentNameProps {
  // Define props with proper TypeScript types
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

## 📊 Implementation Status

### ✅ Completed Features
- Comprehensive color system with accessibility compliance
- Responsive typography and spacing systems
- Core component library (Input, Table, TypeAheadDropdown)
- Mobile-first responsive design
- Dark theme implementation
- Accessibility features (focus management, contrast, keyboard navigation)
- Layout components (Header, Search, Notifications)
- Reusable form components with consistent styling

### 🔄 Future Enhancements
- Additional reusable UI components
- Enhanced data visualization components
- Formal component library documentation (Storybook)
- Design token management system
- Performance optimizations

## 📞 Support

For questions about the design system or component usage:
1. Review this documentation and component-specific guides
2. Check existing component patterns for similar use cases
3. Consider accessibility implications of any proposed changes
4. Document any new patterns or variations

---

**Note**: This design system is a living document that evolves with the application. Regular reviews ensure it continues to serve user needs effectively while maintaining consistency and accessibility standards.
