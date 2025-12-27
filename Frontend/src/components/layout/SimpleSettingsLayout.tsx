import { ReactNode } from 'react';
import { useNavigate } from 'react-router-dom';
import './SimpleSettingsLayout.css';

interface SimpleSettingsLayoutProps {
    /** The title/heading for the page */
    title: string;
    /** The main content of the page */
    children: ReactNode;
    /** Show back button (default: true) */
    showBack?: boolean;
    /** Custom back path (default: /settings) */
    backPath?: string;
}

/**
 * SimpleSettingsLayout - A clean layout for settings pages without sidebar
 * 
 * This component provides a simple layout for settings pages with:
 * - Optional back button
 * - Page title
 * - Clean content area
 * 
 * @example
 * ```tsx
 * <SimpleSettingsLayout title="Privacy Settings">
 *   <SettingsList>...</SettingsList>
 * </SimpleSettingsLayout>
 * ```
 */
const SimpleSettingsLayout = ({ 
    title, 
    children, 
    showBack = true, 
    backPath = '/settings' 
}: SimpleSettingsLayoutProps) => {
    const navigate = useNavigate();

    return (
        <div className="simple-settings-layout">
            <div className="simple-settings-header">
                {showBack && (
                    <button
                        onClick={() => navigate(backPath)}
                        className="simple-settings-back-button"
                        aria-label="Back to settings"
                    >
                        ←
                    </button>
                )}
                <h1 className="simple-settings-title">{title}</h1>
            </div>
            <div className="simple-settings-content">
                {children}
            </div>
        </div>
    );
};

export default SimpleSettingsLayout;
