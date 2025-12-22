import React from 'react';
import './SettingsSection.css';

interface SettingsSectionProps {
    /** Title of the settings section */
    title: string;
    /** Optional description for the section */
    description?: string;
    /** Section content (usually setting items) */
    children: React.ReactNode;
}

/**
 * SettingsSection - A reusable container for grouping related settings
 * 
 * This component provides a consistent layout for settings sections across
 * the application (e.g., club settings, user settings).
 * 
 * @example
 * ```tsx
 * <SettingsSection title="Features" description="Enable or disable features">
 *   <SettingItem title="Fines" description="Enable fines management">
 *     <ToggleSwitch checked={enabled} onChange={handleToggle} />
 *   </SettingItem>
 * </SettingsSection>
 * ```
 */
const SettingsSection: React.FC<SettingsSectionProps> = ({
    title,
    description,
    children
}) => {
    return (
        <div className="settings-section">
            <div className="settings-section-header">
                <h3 className="settings-section-title">{title}</h3>
                {description && (
                    <p className="settings-section-description">{description}</p>
                )}
            </div>
            <div className="settings-section-content">
                {children}
            </div>
        </div>
    );
};

interface SettingItemProps {
    /** Title of the setting */
    title: string;
    /** Optional description explaining what the setting does */
    description?: string;
    /** Control element (toggle, checkbox, etc.) */
    children: React.ReactNode;
}

/**
 * SettingItem - A single setting row with label and control
 * 
 * @example
 * ```tsx
 * <SettingItem title="Enable Feature" description="Allow users to use this feature">
 *   <ToggleSwitch checked={enabled} onChange={handleToggle} />
 * </SettingItem>
 * ```
 */
const SettingItem: React.FC<SettingItemProps> = ({
    title,
    description,
    children
}) => {
    return (
        <div className="setting-item">
            <div className="setting-info">
                <h4 className="setting-title">{title}</h4>
                {description && (
                    <p className="setting-description">{description}</p>
                )}
            </div>
            <div className="setting-control">
                {children}
            </div>
        </div>
    );
};

export { SettingsSection, SettingItem };
export type { SettingsSectionProps, SettingItemProps };
