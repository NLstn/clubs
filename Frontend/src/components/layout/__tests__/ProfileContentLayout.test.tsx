import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import ProfileContentLayout from '../ProfileContentLayout';
import '@testing-library/jest-dom';

// Mock ProfileSidebar
vi.mock('../../../pages/profile/ProfileSidebar', () => ({
  default: () => <div data-testid="profile-sidebar">Profile Sidebar</div>
}));

describe('ProfileContentLayout', () => {
  it('renders with title and children', () => {
    render(
      <ProfileContentLayout title="Test Title">
        <div>Test Content</div>
      </ProfileContentLayout>
    );
    
    expect(screen.getByText('Test Title')).toBeInTheDocument();
    expect(screen.getByText('Test Content')).toBeInTheDocument();
    expect(screen.getByTestId('profile-sidebar')).toBeInTheDocument();
  });

  it('renders simple title when no rich header props provided', () => {
    const { container } = render(
      <ProfileContentLayout title="Simple Title">
        <div>Content</div>
      </ProfileContentLayout>
    );
    
    const simpleTitle = container.querySelector('.profile-content-title-simple');
    expect(simpleTitle).toBeInTheDocument();
    expect(simpleTitle).toHaveTextContent('Simple Title');
    
    const richHeader = container.querySelector('.profile-content-header');
    expect(richHeader).not.toBeInTheDocument();
  });

  it('renders rich header when actions provided', () => {
    const { container } = render(
      <ProfileContentLayout 
        title="Rich Title"
        actions={<button>Edit</button>}
      >
        <div>Content</div>
      </ProfileContentLayout>
    );
    
    const richHeader = container.querySelector('.profile-content-header');
    expect(richHeader).toBeInTheDocument();
    
    const richTitle = container.querySelector('.profile-content-title');
    expect(richTitle).toBeInTheDocument();
    expect(richTitle).toHaveTextContent('Rich Title');
    
    const simpleTitle = container.querySelector('.profile-content-title-simple');
    expect(simpleTitle).not.toBeInTheDocument();
    
    expect(screen.getByText('Edit')).toBeInTheDocument();
  });

  it('renders rich header when headerContent provided', () => {
    const { container } = render(
      <ProfileContentLayout 
        title="Rich Title"
        headerContent={<div className="avatar">Avatar</div>}
      >
        <div>Content</div>
      </ProfileContentLayout>
    );
    
    const richHeader = container.querySelector('.profile-content-header');
    expect(richHeader).toBeInTheDocument();
    
    expect(screen.getByText('Avatar')).toBeInTheDocument();
  });

  it('renders rich header when subtitle provided', () => {
    const { container } = render(
      <ProfileContentLayout 
        title="Rich Title"
        subtitle="test@example.com"
      >
        <div>Content</div>
      </ProfileContentLayout>
    );
    
    const richHeader = container.querySelector('.profile-content-header');
    expect(richHeader).toBeInTheDocument();
    
    const subtitle = container.querySelector('.profile-content-subtitle');
    expect(subtitle).toBeInTheDocument();
    expect(subtitle).toHaveTextContent('test@example.com');
  });

  it('renders rich header with all props', () => {
    const { container } = render(
      <ProfileContentLayout 
        title="Full Rich Title"
        subtitle="user@example.com"
        headerContent={<div className="avatar">UC</div>}
        actions={
          <>
            <button>Edit</button>
            <button>Save</button>
          </>
        }
      >
        <div>Content</div>
      </ProfileContentLayout>
    );
    
    const richHeader = container.querySelector('.profile-content-header');
    expect(richHeader).toBeInTheDocument();
    
    expect(screen.getByText('Full Rich Title')).toBeInTheDocument();
    expect(screen.getByText('user@example.com')).toBeInTheDocument();
    expect(screen.getByText('UC')).toBeInTheDocument();
    expect(screen.getByText('Edit')).toBeInTheDocument();
    expect(screen.getByText('Save')).toBeInTheDocument();
  });

  it('does not render subtitle when not provided in rich header', () => {
    const { container } = render(
      <ProfileContentLayout 
        title="Title"
        actions={<button>Edit</button>}
      >
        <div>Content</div>
      </ProfileContentLayout>
    );
    
    const subtitle = container.querySelector('.profile-content-subtitle');
    expect(subtitle).not.toBeInTheDocument();
  });

  it('applies correct layout structure', () => {
    const { container } = render(
      <ProfileContentLayout title="Test">
        <div>Content</div>
      </ProfileContentLayout>
    );
    
    const layout = container.querySelector('.profile-content-layout');
    expect(layout).toBeInTheDocument();
    
    const main = container.querySelector('.profile-content-main');
    expect(main).toBeInTheDocument();
  });

  it('renders children inside main content area', () => {
    const { container } = render(
      <ProfileContentLayout title="Test">
        <div className="test-child">Child Content</div>
      </ProfileContentLayout>
    );
    
    const main = container.querySelector('.profile-content-main');
    const child = container.querySelector('.test-child');
    
    expect(main).toContainElement(child as HTMLElement);
  });
});
