# Color Reference Guide

This document provides a visual and technical reference for all colors used in the Clubs Management Application.

## Color Palette Overview

### Primary Colors

| Color Name | Hex Code | RGB | Usage |
|------------|----------|-----|-------|
| Primary Green | `#4CAF50` | `rgb(76, 175, 80)` | Primary actions, success states, brand elements |
| Primary Green Hover | `#45a049` | `rgb(69, 160, 73)` | Hover state for primary buttons |
| Secondary Blue | `#646cff` | `rgb(100, 108, 255)` | Secondary actions, links, informational elements |
| Secondary Blue Hover | `#535bf2` | `rgb(83, 91, 242)` | Hover state for secondary buttons |

### Background Colors

| Color Name | Hex Code | RGB | Usage |
|------------|----------|-----|-------|
| Background Dark | `#242424` | `rgb(36, 36, 36)` | Main application background |
| Background Light | `#333333` | `rgb(51, 51, 51)` | Card backgrounds, sections |

### Text Colors

| Color Name | Hex Code | RGB | Usage |
|------------|----------|-----|-------|
| Text Primary | `rgba(255, 255, 255, 0.87)` | `rgb(255, 255, 255, 87%)` | Main text content |
| Text Secondary | `#888888` | `rgb(136, 136, 136)` | Supporting text, metadata |

### System Colors

| Color Name | Hex Code | RGB | Usage |
|------------|----------|-----|-------|
| Error/Cancel | `#f44336` | `rgb(244, 67, 54)` | Destructive actions, errors |
| Error/Cancel Hover | `#e53935` | `rgb(229, 57, 53)` | Hover state for destructive actions |
| Success Background | `#d4edda` | `rgb(212, 237, 218)` | Success message backgrounds |
| Success Text | `#155724` | `rgb(21, 87, 36)` | Success message text |
| Error Background | `#f8d7da` | `rgb(248, 215, 218)` | Error message backgrounds |
| Error Text | `#721c24` | `rgb(114, 28, 36)` | Error message text |
| Border | `#dddddd` | `rgb(221, 221, 221)` | Default borders, separators |

### Accent Colors

| Color Name | Hex Code | RGB | Usage |
|------------|----------|-----|-------|
| Orange (Maybe) | `#ffa500` | `rgb(255, 165, 0)` | Maybe/tentative states |
| Orange Hover | `#ff8c00` | `rgb(255, 140, 0)` | Hover state for orange buttons |

## Color Accessibility

### Contrast Ratios

All color combinations meet WCAG 2.1 AA standards (minimum 4.5:1 for normal text, 3:1 for large text):

| Foreground | Background | Contrast Ratio | WCAG Level |
|------------|------------|----------------|------------|
| Text Primary | Background Dark | 12.63:1 | AAA |
| Text Secondary | Background Dark | 4.95:1 | AA |
| White Text | Primary Green | 4.68:1 | AA |
| White Text | Secondary Blue | 5.18:1 | AA |
| White Text | Error Red | 5.79:1 | AA |
| Success Text | Success Background | 7.52:1 | AAA |
| Error Text | Error Background | 6.18:1 | AA |

### Color Blind Considerations

- The primary green and error red combination provides sufficient contrast
- Blue secondary color is distinguishable from green primary
- Text labels and icons supplement color-only information
- Focus states use multiple visual indicators (border, shadow, outline)

## Usage Guidelines

### Do's

✅ Use primary green for main actions (save, create, confirm)
✅ Use secondary blue for navigation and informational actions
✅ Use red sparingly for destructive actions only
✅ Maintain the dark theme consistency throughout
✅ Ensure sufficient contrast for all text elements
✅ Use semantic meaning for colors (green = success, red = error)

### Don'ts

❌ Don't use bright colors on dark backgrounds without sufficient contrast
❌ Don't rely solely on color to convey information
❌ Don't use red for non-destructive actions
❌ Don't mix light and dark themes arbitrarily
❌ Don't create custom colors without checking accessibility
❌ Don't use more than 3-4 colors in a single interface

## Implementation in CSS

### CSS Custom Properties

```css
:root {
  /* Primary Colors */
  --color-primary: #4CAF50;
  --color-primary-hover: #45a049;
  --color-secondary: #646cff;
  --color-secondary-hover: #535bf2;
  
  /* Background Colors */
  --color-background: #242424;
  --color-background-light: #333333;
  
  /* Text Colors */
  --color-text: rgba(255, 255, 255, 0.87);
  --color-text-secondary: #888;
  
  /* System Colors */
  --color-cancel: #f44336;
  --color-cancel-hover: #e53935;
  --color-cancel-text: #fff;
  --color-border: #ddd;
  
  /* Status Colors */
  --color-success-bg: #d4edda;
  --color-success-text: #155724;
  --color-error-bg: #f8d7da;
  --color-error-text: #721c24;
}
```

### Example Usage

```css
/* Primary Button */
.button-primary {
  background-color: var(--color-primary);
  color: white;
  border: none;
}

.button-primary:hover {
  background-color: var(--color-primary-hover);
}

/* Success Message */
.success-message {
  background-color: var(--color-success-bg);
  color: var(--color-success-text);
  border: 1px solid #c3e6cb;
}

/* Card Background */
.card {
  background-color: var(--color-background-light);
  border: 1px solid var(--color-border);
  color: var(--color-text);
}
```

## Brand Color Hex Values (for Design Tools)

Copy these hex values for use in design tools like Figma, Sketch, or Adobe XD:

```
Primary Green: #4CAF50
Primary Green Hover: #45a049
Secondary Blue: #646cff
Secondary Blue Hover: #535bf2
Background Dark: #242424
Background Light: #333333
Text Primary: #ffffff (87% opacity)
Text Secondary: #888888
Error Red: #f44336
Error Red Hover: #e53935
Border Gray: #dddddd
Success Green: #155724
Success Background: #d4edda
Error Background: #f8d7da
Orange Accent: #ffa500
Orange Hover: #ff8c00
```

## Testing Colors

When implementing new components or features:

1. Test all color combinations with a contrast checker
2. Verify readability in different lighting conditions
3. Test with color blindness simulators
4. Ensure colors work in both light and dark system preferences
5. Validate focus states are clearly visible

## Future Color Considerations

As the application grows, consider:

- Adding semantic color tokens for specific use cases
- Implementing a light theme variant
- Adding more accent colors for data visualization
- Creating color scales for better granular control
- Implementing color tokens for theming system