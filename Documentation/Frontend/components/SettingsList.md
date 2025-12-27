# SettingsList Components

## Overview

The SettingsList components provide a smartphone-style settings interface that mimics the native settings apps found on iOS and Android devices. These components are designed to offer a familiar, mobile-friendly experience for managing user preferences and configuration options.

## Components

### SettingsList

Container component for organizing settings sections.

**Props:**
- `children: React.ReactNode` - Settings list content (usually SettingsListSection components)
- `className?: string` - Additional CSS class name

**Usage:**
```tsx
import { SettingsList } from '@/components/ui';

<SettingsList>
  <SettingsListSection title="Preferences">
    {/* Settings items */}
  </SettingsListSection>
</SettingsList>
```

### SettingsListSection

Groups related settings together with an optional header and description.

**Props:**
- `title?: string` - Section title/header
- `description?: string` - Section description shown below title
- `children: React.ReactNode` - Section content (usually SettingsListItem components)

**Usage:**
```tsx
import { SettingsListSection } from '@/components/ui';

<SettingsListSection 
  title="DISPLAY" 
  description="Customize how the app looks"
>
  <SettingsListItem title="Theme" value="Dark" onClick={handleClick} />
  <SettingsListItem title="Text Size" value="Medium" onClick={handleClick} />
</SettingsListSection>
```

### SettingsListItem

Individual setting row in smartphone style.

**Props:**
- `title: string` - Main item title (required)
- `subtitle?: string` - Optional subtitle/description
- `value?: string` - Optional value displayed on the right
- `icon?: React.ReactNode` - Optional icon element
- `control?: React.ReactNode` - Optional control element (toggle, checkbox, etc.)
- `onClick?: () => void` - Click handler for navigable items
- `showChevron?: boolean` - Whether to show chevron indicator (auto-detected if onClick is provided)
- `className?: string` - Custom class name

**Usage Examples:**

**Navigable Item (with chevron):**
```tsx
<SettingsListItem 
  title="Language" 
  value="English" 
  icon="🌐"
  onClick={handleLanguageClick}
/>
```

**Item with Control:**
```tsx
<SettingsListItem 
  title="Notifications" 
  subtitle="Receive push notifications"
  icon="🔔"
  control={<ToggleSwitch checked={enabled} onChange={handleToggle} />}
/>
```

**Item with Subtitle:**
```tsx
<SettingsListItem 
  title="Dark Mode"
  subtitle="Dark theme for reduced eye strain"
  value="✓"
  icon="🌙"
  onClick={handleThemeChange}
/>
```

## Complete Example

```tsx
import { SettingsList, SettingsListSection, SettingsListItem, ToggleSwitch } from '@/components/ui';
import { useState } from 'react';

function UserPreferences() {
  const [notifications, setNotifications] = useState(true);
  const [theme, setTheme] = useState('dark');

  return (
    <SettingsList>
      {/* Language Section */}
      <SettingsListSection 
        title="LANGUAGE"
        description="Select your preferred language"
      >
        <SettingsListItem 
          title="Language" 
          value="English" 
          icon="🌐"
          onClick={() => {/* Navigate to language selection */}}
        />
      </SettingsListSection>

      {/* Appearance Section */}
      <SettingsListSection 
        title="APPEARANCE"
        description="Choose your preferred color theme"
      >
        <SettingsListItem
          title="Light Mode"
          subtitle="Bright theme with light backgrounds"
          value={theme === 'light' ? '✓' : ''}
          icon="☀️"
          onClick={() => setTheme('light')}
        />
        <SettingsListItem
          title="Dark Mode"
          subtitle="Dark theme for reduced eye strain"
          value={theme === 'dark' ? '✓' : ''}
          icon="🌙"
          onClick={() => setTheme('dark')}
        />
        <SettingsListItem
          title="System Setting"
          subtitle="Automatically adjusts to your device theme"
          value={theme === 'system' ? '✓' : ''}
          icon="ℹ️"
          onClick={() => setTheme('system')}
        />
      </SettingsListSection>

      {/* Notifications Section */}
      <SettingsListSection title="NOTIFICATIONS">
        <SettingsListItem 
          title="Push Notifications"
          subtitle="Receive notifications for important updates"
          icon="🔔"
          control={
            <ToggleSwitch 
              checked={notifications} 
              onChange={setNotifications} 
            />
          }
        />
      </SettingsListSection>
    </SettingsList>
  );
}
```

## Design Features

### Visual Style
- **Clean iOS/Android-inspired design** with grouped sections
- **Touch-friendly targets** (minimum 44px height on desktop, 48px on mobile)
- **Proper spacing and borders** for clear visual hierarchy
- **Hover states** for interactive items on desktop
- **Active/pressed states** for mobile tap feedback

### Accessibility
- **Keyboard navigation** support for navigable items (Enter/Space to activate)
- **Proper ARIA roles** (button role for clickable items)
- **Focus indicators** for keyboard users
- **Screen reader friendly** with proper semantic structure

### Responsive Design
- **Mobile-optimized** with larger touch targets on small screens
- **Fluid typography** that scales appropriately
- **Adaptive spacing** that adjusts to viewport size

## Styling

The components use CSS custom properties from the design system:

```css
--color-text
--color-text-secondary
--color-background-light
--color-border
--color-primary
--space-xs, --space-sm, --space-md
--border-radius-sm, --border-radius-md, --border-radius-lg
```

## When to Use

### ✅ Use SettingsList when:
- Building settings/preferences pages for mobile views
- Creating configuration interfaces that need to feel native
- Designing forms with multiple sections on mobile devices
- Implementing list-based navigation menus

### ❌ Don't use SettingsList when:
- Building complex forms with many input fields (use Form components instead)
- Creating data tables (use Table or ODataTable)
- Displaying read-only information (use Card or list components)
- On desktop-only interfaces where the grid layout is more appropriate

## Related Components

- **SettingsSection** - Desktop-oriented settings layout with horizontal controls
- **SettingItem** - Desktop-oriented setting row with side-by-side layout
- **ToggleSwitch** - Common control used within SettingsListItem
- **Card** - For grouping content in non-settings contexts

## Differences from SettingsSection/SettingItem

The `SettingsList` components are specifically designed for **mobile/smartphone interfaces**, while `SettingsSection`/`SettingItem` are better suited for **desktop layouts**:

| Feature | SettingsList | SettingsSection |
|---------|--------------|-----------------|
| Layout | Vertical, stacked | Horizontal, side-by-side |
| Visual Style | iOS/Android native | Desktop settings panel |
| Touch Targets | Extra large (48px) | Standard (44px) |
| Navigation | Chevron indicators | No chevron |
| Best For | Mobile/tablet | Desktop |

## Implementation in ProfilePreferences

The ProfilePreferences page uses responsive design to show different layouts based on viewport size:

- **Desktop (>768px)**: Grid-based theme selector with visual previews
- **Mobile (≤768px)**: SettingsList with smartphone-style interface

This provides an optimal experience on both platforms while maintaining functionality.
