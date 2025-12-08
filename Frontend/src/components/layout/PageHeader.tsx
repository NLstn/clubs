import React from 'react';
import './PageHeader.css';

export interface PageHeaderProps {
    /** Content to display in the header (typically includes title, description, icon/avatar) */
    children: React.ReactNode;
    /** Optional action buttons or elements to display on the right side */
    actions?: React.ReactNode;
    /** Visual variant of the header */
    variant?: 'card' | 'simple';
    /** Additional CSS class name */
    className?: string;
}

/**
 * PageHeader - A reusable header component for detail pages
 * 
 * Creates a consistent header pattern across Club, Team, Profile, Event, and Admin pages.
 * Provides a flex container with content on the left and actions on the right.
 * 
 * @example
 * // Card variant (default) - used for Club/Team/Profile pages
 * <PageHeader actions={<Button>Edit</Button>}>
 *   <Avatar src={logo} />
 *   <div>
 *     <h1>Title</h1>
 *     <p>Description</p>
 *   </div>
 * </PageHeader>
 * 
 * @example
 * // Simple variant - used for Event pages with breadcrumbs
 * <PageHeader variant="simple" actions={<Button>Delete</Button>}>
 *   <Breadcrumb>Home > Events > Event Name</Breadcrumb>
 * </PageHeader>
 */
export const PageHeader: React.FC<PageHeaderProps> = ({
    children,
    actions,
    variant = 'card',
    className = '',
}) => {
    const headerClass = variant === 'simple' 
        ? 'page-header page-header--simple' 
        : 'page-header page-header--card';

    return (
        <div className={`${headerClass} ${className}`.trim()}>
            <div className="page-header__content">
                {children}
            </div>
            {actions && (
                <div className="page-header__actions">
                    {actions}
                </div>
            )}
        </div>
    );
};

export default PageHeader;
