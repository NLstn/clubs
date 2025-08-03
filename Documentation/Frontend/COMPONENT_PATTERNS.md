# Component Patterns Guide

This document provides detailed specifications and usage guidelines for all UI components in the Clubs Management Application.

## Button Components

### Primary Button

**Purpose**: Main actions, form submissions, primary navigation
**Visual**: Green background with white text

```css
.button-primary {
  background-color: var(--color-primary);
  color: white;
  border: none;
  padding: var(--space-sm) var(--space-md);
  border-radius: var(--border-radius-sm);
  font-size: 1rem;
  font-weight: 500;
  min-height: 44px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.button-primary:hover:not(:disabled) {
  background-color: var(--color-primary-hover);
}
```

**Usage Examples**:
- "Create Club"
- "Save Changes"
- "Join Club"
- "Send Magic Link"

### Secondary Button

**Purpose**: Secondary actions, alternative options
**Visual**: Blue background with white text

```css
.button-secondary {
  background-color: var(--color-secondary);
  color: white;
  /* Same styling as primary but different color */
}
```

**Usage Examples**:
- "Edit"
- "View Details"
- "Login with Keycloak"

### Destructive Button

**Purpose**: Dangerous or irreversible actions
**Visual**: Red background with white text

```css
.button-cancel {
  background-color: var(--color-cancel);
  color: var(--color-cancel-text);
}
```

**Usage Examples**:
- "Delete Club"
- "Remove Member"
- "Cancel"

### Button States

- **Default**: Base styling with appropriate color
- **Hover**: Darker background color
- **Focus**: 4px auto outline for keyboard navigation
- **Disabled**: 60% opacity, no pointer events, not-allowed cursor
- **Loading**: Disabled state with loading indicator

### Button Guidelines

- Use descriptive labels instead of generic terms
- Maintain consistent capitalization (Title Case)
- Ensure minimum 44px touch target (48px on mobile)
- Place primary action on the right in button groups
- Use full-width buttons on mobile for better usability

## Form Components

### Input Fields

**Purpose**: Text input, email, password, and other single-line data entry

```css
.form-group input {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid var(--color-border);
  border-radius: var(--border-radius-md);
  font-size: 1rem;
  background-color: var(--color-background-light);
  color: var(--color-text);
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
  box-sizing: border-box;
}

.form-group input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(76, 175, 80, 0.1);
}
```

### Textarea Fields

**Purpose**: Multi-line text input

```css
.form-group textarea {
  min-height: 100px;
  resize: vertical;
  /* Inherits input styling */
}
```

### Select Dropdowns

**Purpose**: Single selection from predefined options

```css
.form-group select {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid var(--color-border);
  border-radius: var(--border-radius-md);
  background-color: var(--color-background-light);
  color: var(--color-text);
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg...");
  background-repeat: no-repeat;
  background-position: right 12px center;
  background-size: 20px;
  padding-right: 40px;
}
```

### Form Labels

**Purpose**: Clear identification of form fields

```css
.form-group label {
  display: block;
  font-weight: 600;
  font-size: 0.95rem;
  color: var(--color-text);
  margin-bottom: var(--space-xs);
  letter-spacing: 0.025em;
}
```

### Form Validation

**Success State**:
```css
.form-group input.success {
  border-color: var(--color-primary);
}
```

**Error State**:
```css
.form-group input.error {
  border-color: var(--color-cancel);
}

.error-message {
  color: var(--color-cancel);
  font-size: 0.875rem;
  margin-top: var(--space-xs);
}
```

## Card Components

### Basic Card

**Purpose**: Content containers, information display
**Visual**: Light background with subtle border and shadow

```css
.card {
  background-color: white;
  border: 1px solid var(--color-border);
  border-radius: var(--border-radius-lg);
  padding: var(--space-md);
  box-shadow: var(--shadow-sm);
  color: black;
  transition: box-shadow 0.2s ease;
}

.card:hover {
  box-shadow: var(--shadow-md);
}
```

### Event Card

**Purpose**: Displaying event information in dark theme
**Visual**: Dark background consistent with app theme

```css
.event-card {
  background-color: var(--color-background-light);
  border: 1px solid var(--color-border);
  border-radius: var(--border-radius-lg);
  padding: var(--space-md);
  color: var(--color-text);
}
```

### News Card

**Purpose**: Displaying news and announcements

```css
.news-card {
  background-color: var(--color-background-light);
  border: 1px solid var(--color-border);
  border-radius: var(--border-radius-sm);
  padding: var(--space-md);
  color: var(--color-text);
}
```

### Club Card

**Purpose**: Displaying club information in lists
**Visual**: Interactive card with hover effects

```css
.club-card {
  background-color: white;
  border: 1px solid var(--color-border);
  border-radius: var(--border-radius-lg);
  padding: var(--space-md);
  cursor: pointer;
  color: black;
  display: flex;
  flex-direction: column;
  height: 100%;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.club-card:hover {
  border-color: var(--color-primary);
  box-shadow: var(--shadow-sm);
}
```

## Navigation Components

### Header

**Purpose**: Main navigation and branding
**Structure**: Logo + Title + Actions + User Menu

```css
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 var(--space-lg);
  height: 70px;
  background-color: var(--color-background-light);
  box-shadow: var(--shadow-sm);
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  z-index: 1000;
}
```

### User Dropdown

**Purpose**: User account actions and navigation

```css
.userIcon {
  width: 40px;
  height: 40px;
  border-radius: var(--border-radius-circle);
  background-color: var(--color-primary);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-weight: bold;
}

.dropdown {
  position: absolute;
  right: 0;
  top: 50px;
  background-color: white;
  box-shadow: var(--shadow-md);
  border-radius: var(--border-radius-md);
  width: 140px;
  z-index: 100;
  border: 1px solid var(--color-border);
}
```

### Tab Navigation

**Purpose**: Section navigation within pages

```css
.tabs-nav {
  display: flex;
  border-bottom: 2px solid var(--color-border);
  margin-bottom: var(--space-lg);
}

.tab-button {
  background: none;
  border: none;
  padding: var(--space-sm) var(--space-lg);
  cursor: pointer;
  font-size: 1rem;
  color: var(--color-text-secondary);
  border-bottom: 2px solid transparent;
  transition: all 0.2s ease;
}

.tab-button.active {
  color: var(--color-primary);
  border-bottom-color: var(--color-primary);
  background-color: var(--color-background-light);
}
```

## Table Components

### Basic Table

**Purpose**: Displaying tabular data

```css
table {
  width: 100%;
  border-collapse: collapse;
  margin: var(--space-lg) 0;
}

th, td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid var(--color-border);
}

th {
  background-color: var(--color-background-light);
  font-weight: bold;
  color: var(--color-text);
}

td {
  color: var(--color-text);
}

tr:hover {
  background-color: var(--color-background-light);
}
```

### Responsive Table Wrapper

**Purpose**: Horizontal scrolling on mobile

```css
.table-responsive {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  border: 1px solid var(--color-border);
  border-radius: var(--border-radius-sm);
}
```

## Modal Components

### Modal Overlay

**Purpose**: Focused interactions and confirmations

```css
.modal {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background-color: var(--color-background-light);
  padding: var(--space-lg);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-md);
  width: 90%;
  max-width: 500px;
  color: var(--color-text);
}
```

## Message Components

### Success Message

**Purpose**: Positive feedback and confirmations

```css
.success-message {
  background-color: var(--color-success-bg);
  color: var(--color-success-text);
  padding: 12px 16px;
  border-radius: var(--border-radius-md);
  border: 1px solid #c3e6cb;
  margin-bottom: var(--space-md);
  font-weight: 500;
}
```

### Error Message

**Purpose**: Error feedback and warnings

```css
.error-message {
  background-color: var(--color-error-bg);
  color: var(--color-error-text);
  padding: 12px 16px;
  border-radius: var(--border-radius-md);
  border: 1px solid #f5c6cb;
  margin-bottom: var(--space-md);
  font-weight: 500;
}
```

### Empty State

**Purpose**: Informing users when no content is available

```css
.empty-state {
  background-color: var(--color-background-light);
  border: 1px dashed var(--color-border);
  border-radius: var(--border-radius-lg);
  padding: var(--space-lg);
  text-align: center;
  color: var(--color-text-secondary);
}
```

## Badge Components

### Club Badge

**Purpose**: Displaying club association

```css
.club-badge {
  background-color: var(--color-primary);
  color: white;
  padding: 4px 12px;
  border-radius: var(--border-radius-sm);
  font-size: 0.8rem;
  cursor: pointer;
  transition: background-color 0.2s;
  white-space: nowrap;
}

.club-badge:hover {
  background-color: var(--color-primary-hover);
}
```

### Activity Type Badge

**Purpose**: Categorizing activity feed items

```css
.activity-type-badge {
  background-color: var(--color-primary);
  color: white;
  padding: var(--space-xs) var(--space-sm);
  border-radius: var(--border-radius-sm);
  font-size: 0.8rem;
  font-weight: 500;
  text-transform: capitalize;
}
```

## Toggle Components

### Toggle Switch

**Purpose**: Boolean settings and preferences

```css
.toggle-switch {
  position: relative;
  display: inline-block;
  width: 60px;
  height: 34px;
}

.toggle-switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0; left: 0; right: 0; bottom: 0;
  background-color: #ccc;
  transition: .4s;
  border-radius: 34px;
}

.slider:before {
  position: absolute;
  content: "";
  height: 26px;
  width: 26px;
  left: 4px;
  bottom: 4px;
  background-color: white;
  transition: .4s;
  border-radius: 50%;
}

input:checked + .slider {
  background-color: var(--color-primary);
}

input:checked + .slider:before {
  transform: translateX(26px);
}
```

## Layout Components

### Main Layout

**Purpose**: Consistent page structure

```css
.layout {
  width: 100%;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-content {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 90px var(--space-lg) var(--space-lg);
  box-sizing: border-box;
  flex: 1;
}
```

### Grid Layouts

**Purpose**: Responsive content grids

```css
.clubs-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

/* Mobile responsive */
@media (max-width: 480px) {
  .clubs-list {
    grid-template-columns: 1fr;
  }
}
```

## Component Usage Guidelines

### When to Use Each Component

- **Primary Buttons**: Main actions, form submissions
- **Secondary Buttons**: Alternative actions, navigation
- **Cards**: Content containers, clickable items
- **Tables**: Structured data display
- **Modals**: Focused interactions, confirmations
- **Badges**: Labels, categories, status indicators
- **Forms**: Data collection and editing

### Accessibility Considerations

- All interactive elements must be keyboard accessible
- Focus states should be clearly visible
- Color should not be the only way to convey information
- Text should meet contrast ratio requirements
- Screen readers should be able to understand component purpose

### Responsive Behavior

- Components should adapt gracefully to different screen sizes
- Touch targets should be minimum 44px (48px on mobile)
- Text should remain readable at all sizes
- Complex components may need alternative layouts for mobile

---

This component guide should be used alongside the main UI Design Guideline for comprehensive design system implementation.