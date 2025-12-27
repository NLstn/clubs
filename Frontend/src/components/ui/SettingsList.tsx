import React from 'react';
import './SettingsList.css';

interface SettingsListProps {
    /** Settings list content (usually SettingsListSection components) */
    children: React.ReactNode;
    /** Additional CSS class name */
    className?: string;
}

/**
 * SettingsList - Smartphone-style settings list container
 * 
 * This component mimics the settings interface found on iOS and Android devices,
 * providing a clean, mobile-friendly layout for settings pages.
 * 
 * @example
 * ```tsx
 * <SettingsList>
 *   <SettingsListSection title="Preferences">
 *     <SettingsListItem title="Language" value="English" onClick={handleLanguageClick} />
 *     <SettingsListItem title="Theme" value="Dark" onClick={handleThemeClick} />
 *   </SettingsListSection>
 * </SettingsList>
 * ```
 */
export const SettingsList: React.FC<SettingsListProps> = ({ children, className = '' }) => {
    return (
        <div className={`settings-list ${className}`}>
            {children}
        </div>
    );
};

interface SettingsListSectionProps {
    /** Section title/header */
    title?: string;
    /** Section description shown below title */
    description?: string;
    /** Section content (usually SettingsListItem components) */
    children: React.ReactNode;
}

/**
 * SettingsListSection - Grouped section within a settings list
 * 
 * Groups related settings together with an optional header and description.
 * 
 * @example
 * ```tsx
 * <SettingsListSection title="Display" description="Customize how the app looks">
 *   <SettingsListItem title="Theme" value="Dark" onClick={handleClick} />
 *   <SettingsListItem title="Text Size" value="Medium" onClick={handleClick} />
 * </SettingsListSection>
 * ```
 */
export const SettingsListSection: React.FC<SettingsListSectionProps> = ({ 
    title, 
    description,
    children 
}) => {
    return (
        <div className="settings-list-section">
            {(title || description) && (
                <div className="settings-list-section-header">
                    {title && <h3 className="settings-list-section-title">{title}</h3>}
                    {description && <p className="settings-list-section-description">{description}</p>}
                </div>
            )}
            <div className="settings-list-section-content">
                {children}
            </div>
        </div>
    );
};

interface SettingsListItemProps {
    /** Main item title */
    title: string;
    /** Optional subtitle/description */
    subtitle?: string;
    /** Optional value displayed on the right */
    value?: string;
    /** Optional icon element */
    icon?: React.ReactNode;
    /** Optional control element (toggle, checkbox, etc.) */
    control?: React.ReactNode;
    /** Click handler for navigable items */
    onClick?: () => void;
    /** Whether to show chevron indicator (auto-detected if onClick is provided) */
    showChevron?: boolean;
    /** Custom class name */
    className?: string;
}

/**
 * SettingsListItem - Individual setting row in smartphone style
 * 
 * Represents a single setting with optional icon, title, subtitle, value, and control.
 * Can be tappable (with chevron) or contain inline controls.
 * 
 * @example
 * ```tsx
 * // Navigable item
 * <SettingsListItem 
 *   title="Language" 
 *   value="English" 
 *   onClick={handleClick}
 * />
 * 
 * // Item with control
 * <SettingsListItem 
 *   title="Notifications" 
 *   subtitle="Receive push notifications"
 *   control={<ToggleSwitch checked={true} onChange={handleToggle} />}
 * />
 * ```
 */
export const SettingsListItem: React.FC<SettingsListItemProps> = ({
    title,
    subtitle,
    value,
    icon,
    control,
    onClick,
    showChevron,
    className = ''
}) => {
    const isNavigable = onClick !== undefined;
    const shouldShowChevron = showChevron ?? (isNavigable && !control);

    const handleClick = () => {
        if (onClick) {
            onClick();
        }
    };

    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (onClick && (e.key === 'Enter' || e.key === ' ')) {
            e.preventDefault();
            onClick();
        }
    };

    return (
        <div
            className={`settings-list-item ${isNavigable ? 'settings-list-item-navigable' : ''} ${className}`}
            onClick={handleClick}
            onKeyDown={handleKeyDown}
            role={isNavigable ? 'button' : undefined}
            tabIndex={isNavigable ? 0 : undefined}
        >
            {icon && <div className="settings-list-item-icon">{icon}</div>}
            
            <div className="settings-list-item-content">
                <div className="settings-list-item-text">
                    <div className="settings-list-item-title">{title}</div>
                    {subtitle && <div className="settings-list-item-subtitle">{subtitle}</div>}
                </div>
                
                {value && !control && (
                    <div className="settings-list-item-value">{value}</div>
                )}
                
                {control && (
                    <div className="settings-list-item-control">{control}</div>
                )}
                
                {shouldShowChevron && (
                    <div className="settings-list-item-chevron">
                        <svg width="8" height="13" viewBox="0 0 8 13" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <path d="M1 1L6.5 6.5L1 12" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                        </svg>
                    </div>
                )}
            </div>
        </div>
    );
};

export type { SettingsListProps, SettingsListSectionProps, SettingsListItemProps };
