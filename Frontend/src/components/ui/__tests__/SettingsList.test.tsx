import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { SettingsList, SettingsListSection, SettingsListItem } from '../SettingsList';

describe('SettingsList', () => {
    it('renders children correctly', () => {
        render(
            <SettingsList>
                <div>Test Content</div>
            </SettingsList>
        );
        expect(screen.getByText('Test Content')).toBeInTheDocument();
    });

    it('applies custom className', () => {
        const { container } = render(
            <SettingsList className="custom-class">
                <div>Content</div>
            </SettingsList>
        );
        expect(container.firstChild).toHaveClass('settings-list');
        expect(container.firstChild).toHaveClass('custom-class');
    });
});

describe('SettingsListSection', () => {
    it('renders with title and description', () => {
        render(
            <SettingsListSection title="Test Section" description="Test Description">
                <div>Content</div>
            </SettingsListSection>
        );
        expect(screen.getByText('Test Section')).toBeInTheDocument();
        expect(screen.getByText('Test Description')).toBeInTheDocument();
        expect(screen.getByText('Content')).toBeInTheDocument();
    });

    it('renders without header when no title or description provided', () => {
        const { container } = render(
            <SettingsListSection>
                <div>Content</div>
            </SettingsListSection>
        );
        expect(container.querySelector('.settings-list-section-header')).not.toBeInTheDocument();
    });

    it('renders with title only', () => {
        render(
            <SettingsListSection title="Test Section">
                <div>Content</div>
            </SettingsListSection>
        );
        expect(screen.getByText('Test Section')).toBeInTheDocument();
        expect(screen.queryByText('Test Description')).not.toBeInTheDocument();
    });

    it('renders with description only', () => {
        render(
            <SettingsListSection description="Test Description">
                <div>Content</div>
            </SettingsListSection>
        );
        expect(screen.getByText('Test Description')).toBeInTheDocument();
    });
});

describe('SettingsListItem', () => {
    it('renders with title only', () => {
        render(<SettingsListItem title="Test Item" />);
        expect(screen.getByText('Test Item')).toBeInTheDocument();
    });

    it('renders with title and subtitle', () => {
        render(<SettingsListItem title="Test Item" subtitle="Test Subtitle" />);
        expect(screen.getByText('Test Item')).toBeInTheDocument();
        expect(screen.getByText('Test Subtitle')).toBeInTheDocument();
    });

    it('renders with value', () => {
        render(<SettingsListItem title="Test Item" value="Test Value" />);
        expect(screen.getByText('Test Item')).toBeInTheDocument();
        expect(screen.getByText('Test Value')).toBeInTheDocument();
    });

    it('renders with icon', () => {
        render(<SettingsListItem title="Test Item" icon="🔔" />);
        expect(screen.getByText('Test Item')).toBeInTheDocument();
        expect(screen.getByText('🔔')).toBeInTheDocument();
    });

    it('adds aria-hidden to icon container', () => {
        const { container } = render(<SettingsListItem title="Test Item" icon="🔔" />);
        const iconContainer = container.querySelector('.settings-list-item-icon');
        expect(iconContainer).toHaveAttribute('aria-hidden', 'true');
    });

    it('renders with control element', () => {
        const control = <button>Toggle</button>;
        render(<SettingsListItem title="Test Item" control={control} />);
        expect(screen.getByText('Test Item')).toBeInTheDocument();
        expect(screen.getByText('Toggle')).toBeInTheDocument();
    });

    it('calls onClick when clicked', () => {
        const handleClick = vi.fn();
        render(<SettingsListItem title="Test Item" onClick={handleClick} />);
        
        const item = screen.getByRole('button');
        fireEvent.click(item);
        
        expect(handleClick).toHaveBeenCalledTimes(1);
    });

    it('shows chevron when onClick is provided', () => {
        const { container } = render(<SettingsListItem title="Test Item" onClick={() => {}} />);
        expect(container.querySelector('.settings-list-item-chevron')).toBeInTheDocument();
    });

    it('does not show chevron when onClick is not provided', () => {
        const { container } = render(<SettingsListItem title="Test Item" />);
        expect(container.querySelector('.settings-list-item-chevron')).not.toBeInTheDocument();
    });

    it('does not show chevron when control is provided even with onClick', () => {
        const { container } = render(
            <SettingsListItem 
                title="Test Item" 
                onClick={() => {}} 
                control={<button>Control</button>}
            />
        );
        expect(container.querySelector('.settings-list-item-chevron')).not.toBeInTheDocument();
    });

    it('can force chevron display with showChevron prop', () => {
        const { container } = render(
            <SettingsListItem title="Test Item" showChevron={true} />
        );
        expect(container.querySelector('.settings-list-item-chevron')).toBeInTheDocument();
    });

    it('adds navigable class when onClick is provided', () => {
        const { container } = render(<SettingsListItem title="Test Item" onClick={() => {}} />);
        expect(container.querySelector('.settings-list-item-navigable')).toBeInTheDocument();
    });

    it('adds button role when onClick is provided', () => {
        render(<SettingsListItem title="Test Item" onClick={() => {}} />);
        expect(screen.getByRole('button')).toBeInTheDocument();
    });

    it('does not add button role when onClick is not provided', () => {
        render(<SettingsListItem title="Test Item" />);
        expect(screen.queryByRole('button')).not.toBeInTheDocument();
    });

    it('is keyboard accessible with Enter key', () => {
        const handleClick = vi.fn();
        render(<SettingsListItem title="Test Item" onClick={handleClick} />);
        
        const item = screen.getByRole('button');
        fireEvent.keyDown(item, { key: 'Enter' });
        
        expect(handleClick).toHaveBeenCalledTimes(1);
    });

    it('is keyboard accessible with Space key', () => {
        const handleClick = vi.fn();
        render(<SettingsListItem title="Test Item" onClick={handleClick} />);
        
        const item = screen.getByRole('button');
        fireEvent.keyDown(item, { key: ' ' });
        
        expect(handleClick).toHaveBeenCalledTimes(1);
    });

    it('has tabIndex 0 when navigable', () => {
        render(<SettingsListItem title="Test Item" onClick={() => {}} />);
        const item = screen.getByRole('button');
        expect(item).toHaveAttribute('tabIndex', '0');
    });

    it('applies custom className', () => {
        const { container } = render(<SettingsListItem title="Test Item" className="custom-class" />);
        expect(container.firstChild).toHaveClass('settings-list-item');
        expect(container.firstChild).toHaveClass('custom-class');
    });

    it('renders complex icon element', () => {
        const icon = <svg data-testid="test-svg"><path /></svg>;
        render(<SettingsListItem title="Test Item" icon={icon} />);
        expect(screen.getByTestId('test-svg')).toBeInTheDocument();
    });

    it('renders all elements together', () => {
        const { container } = render(
            <SettingsListItem
                title="Test Item"
                subtitle="Test Subtitle"
                value="Test Value"
                icon="🔔"
                onClick={() => {}}
            />
        );
        expect(screen.getByText('Test Item')).toBeInTheDocument();
        expect(screen.getByText('Test Subtitle')).toBeInTheDocument();
        expect(screen.getByText('Test Value')).toBeInTheDocument();
        expect(screen.getByText('🔔')).toBeInTheDocument();
        expect(container.querySelector('.settings-list-item-chevron')).toBeInTheDocument();
    });
});

describe('SettingsList integration', () => {
    it('renders complete settings interface', () => {
        const handleLanguageClick = vi.fn();
        const handleThemeClick = vi.fn();
        
        render(
            <SettingsList>
                <SettingsListSection title="LANGUAGE" description="Select your language">
                    <SettingsListItem
                        title="Language"
                        value="English"
                        icon="🌐"
                        onClick={handleLanguageClick}
                    />
                </SettingsListSection>
                <SettingsListSection title="APPEARANCE">
                    <SettingsListItem
                        title="Dark Mode"
                        subtitle="Dark theme for reduced eye strain"
                        value="✓"
                        icon="🌙"
                        onClick={handleThemeClick}
                    />
                </SettingsListSection>
            </SettingsList>
        );

        expect(screen.getByText('LANGUAGE')).toBeInTheDocument();
        expect(screen.getByText('Select your language')).toBeInTheDocument();
        expect(screen.getByText('Language')).toBeInTheDocument();
        expect(screen.getByText('English')).toBeInTheDocument();
        
        expect(screen.getByText('APPEARANCE')).toBeInTheDocument();
        expect(screen.getByText('Dark Mode')).toBeInTheDocument();
        expect(screen.getByText('Dark theme for reduced eye strain')).toBeInTheDocument();
        expect(screen.getByText('✓')).toBeInTheDocument();
    });
});
