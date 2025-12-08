import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import PageHeader from '../PageHeader';
import '@testing-library/jest-dom';

describe('PageHeader', () => {
  it('renders with children content', () => {
    render(
      <PageHeader>
        <h1>Test Title</h1>
      </PageHeader>
    );
    
    expect(screen.getByText('Test Title')).toBeInTheDocument();
  });

  it('renders with card variant by default', () => {
    const { container } = render(
      <PageHeader>
        <h1>Test Title</h1>
      </PageHeader>
    );
    
    const header = container.querySelector('.page-header');
    expect(header).toHaveClass('page-header--card');
    expect(header).not.toHaveClass('page-header--simple');
  });

  it('renders with simple variant when specified', () => {
    const { container } = render(
      <PageHeader variant="simple">
        <h1>Test Title</h1>
      </PageHeader>
    );
    
    const header = container.querySelector('.page-header');
    expect(header).toHaveClass('page-header--simple');
    expect(header).not.toHaveClass('page-header--card');
  });

  it('renders with action buttons', () => {
    render(
      <PageHeader actions={<button>Edit</button>}>
        <h1>Test Title</h1>
      </PageHeader>
    );
    
    expect(screen.getByText('Edit')).toBeInTheDocument();
  });

  it('renders without actions when not provided', () => {
    const { container } = render(
      <PageHeader>
        <h1>Test Title</h1>
      </PageHeader>
    );
    
    const actions = container.querySelector('.page-header__actions');
    expect(actions).not.toBeInTheDocument();
  });

  it('applies custom className', () => {
    const { container } = render(
      <PageHeader className="custom-header">
        <h1>Test Title</h1>
      </PageHeader>
    );
    
    const header = container.querySelector('.page-header');
    expect(header).toHaveClass('custom-header');
    expect(header).toHaveClass('page-header');
  });

  it('renders multiple action buttons', () => {
    render(
      <PageHeader
        actions={
          <>
            <button>Edit</button>
            <button>Delete</button>
            <button>Share</button>
          </>
        }
      >
        <h1>Test Title</h1>
      </PageHeader>
    );
    
    expect(screen.getByText('Edit')).toBeInTheDocument();
    expect(screen.getByText('Delete')).toBeInTheDocument();
    expect(screen.getByText('Share')).toBeInTheDocument();
  });

  it('renders complex children content', () => {
    render(
      <PageHeader>
        <div className="logo">Logo</div>
        <div className="info">
          <h1>Title</h1>
          <p>Description</p>
        </div>
      </PageHeader>
    );
    
    expect(screen.getByText('Logo')).toBeInTheDocument();
    expect(screen.getByText('Title')).toBeInTheDocument();
    expect(screen.getByText('Description')).toBeInTheDocument();
  });

  it('wraps children in content div', () => {
    const { container } = render(
      <PageHeader>
        <h1>Test Title</h1>
      </PageHeader>
    );
    
    const content = container.querySelector('.page-header__content');
    expect(content).toBeInTheDocument();
    expect(content).toContainHTML('<h1>Test Title</h1>');
  });

  it('wraps actions in actions div when provided', () => {
    const { container } = render(
      <PageHeader actions={<button>Edit</button>}>
        <h1>Test Title</h1>
      </PageHeader>
    );
    
    const actions = container.querySelector('.page-header__actions');
    expect(actions).toBeInTheDocument();
    expect(actions).toContainHTML('<button>Edit</button>');
  });

  it('combines variant and custom className correctly', () => {
    const { container } = render(
      <PageHeader variant="simple" className="custom-class">
        <h1>Test Title</h1>
      </PageHeader>
    );
    
    const header = container.querySelector('.page-header');
    expect(header).toHaveClass('page-header');
    expect(header).toHaveClass('page-header--simple');
    expect(header).toHaveClass('custom-class');
  });

  it('renders card variant with custom className', () => {
    const { container } = render(
      <PageHeader variant="card" className="custom-card">
        <h1>Test Title</h1>
      </PageHeader>
    );
    
    const header = container.querySelector('.page-header');
    expect(header).toHaveClass('page-header');
    expect(header).toHaveClass('page-header--card');
    expect(header).toHaveClass('custom-card');
  });
});
