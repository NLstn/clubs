import { ReactNode } from 'react';
import ProfileSidebar from '../../pages/profile/ProfileSidebar';
import './ProfileContentLayout.css';

interface ProfileContentLayoutProps {
    /** The title/heading for the page */
    title: string;
    /** The main content of the page */
    children: ReactNode;
    /** Optional actions to display in the header (e.g., buttons) */
    actions?: ReactNode;
    /** Optional custom header content (e.g., avatar, user info) */
    headerContent?: ReactNode;
    /** Optional subtitle text (e.g., email address) */
    subtitle?: string;
}

/**
 * ProfileContentLayout - A consistent layout wrapper for all profile pages.
 * 
 * This component provides a standardized structure for profile pages with:
 * - Consistent sidebar navigation
 * - Consistent page title styling
 * - Proper content area with consistent spacing
 * - Optional header content and actions for rich headers
 * 
 * @example
 * Simple title only:
 * ```tsx
 * <ProfileContentLayout title="Active Sessions">
 *   <Table data={sessions} />
 * </ProfileContentLayout>
 * ```
 * 
 * Rich header with content and actions:
 * ```tsx
 * <ProfileContentLayout 
 *   title="Github Copilot"
 *   headerContent={<div className="profile-avatar">GC</div>}
 *   actions={<Button>Edit Profile</Button>}
 * >
 *   <FormGroup>...</FormGroup>
 * </ProfileContentLayout>
 * ```
 */
const ProfileContentLayout = ({ title, children, actions, headerContent, subtitle }: ProfileContentLayoutProps) => {
    const hasRichHeader = headerContent || actions || subtitle;
    
    return (
        <div className="profile-content-layout">
            <ProfileSidebar />
            <div className="profile-content-main">
                {hasRichHeader && (
                    <div className="profile-content-header">
                        {headerContent}
                        <div className="profile-content-header-text">
                            <h1 className="profile-content-title">{title}</h1>
                            {subtitle && <p className="profile-content-subtitle">{subtitle}</p>}
                        </div>
                        {actions && (
                            <div className="profile-content-actions">
                                {actions}
                            </div>
                        )}
                    </div>
                )}
                {children}
            </div>
        </div>
    );
};

export default ProfileContentLayout;
