import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import '@testing-library/jest-dom';
import { Tabs } from '../Tabs';

describe('Tabs Component', () => {
  const mockTabs = [
    { id: 'tab1', label: 'First Tab' },
    { id: 'tab2', label: 'Second Tab' },
    { id: 'tab3', label: 'Third Tab' },
  ];

  it('renders all tabs with the provided tabs array', () => {
    render(
      <Tabs tabs={mockTabs} activeTab="tab1" onTabChange={() => {}}>
        <div>Content</div>
      </Tabs>
    );

    expect(screen.getByRole('button', { name: /first tab/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /second tab/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /third tab/i })).toBeInTheDocument();
  });

  it('highlights the active tab with active class', () => {
    render(
      <Tabs tabs={mockTabs} activeTab="tab2" onTabChange={() => {}}>
        <div>Content</div>
      </Tabs>
    );

    const firstTab = screen.getByRole('button', { name: /first tab/i });
    const secondTab = screen.getByRole('button', { name: /second tab/i });
    const thirdTab = screen.getByRole('button', { name: /third tab/i });

    expect(firstTab).not.toHaveClass('active');
    expect(secondTab).toHaveClass('active');
    expect(thirdTab).not.toHaveClass('active');
  });

  it('calls onTabChange callback when a tab is clicked', () => {
    const handleTabChange = vi.fn();

    render(
      <Tabs tabs={mockTabs} activeTab="tab1" onTabChange={handleTabChange}>
        <div>Content</div>
      </Tabs>
    );

    const secondTab = screen.getByRole('button', { name: /second tab/i });
    fireEvent.click(secondTab);

    expect(handleTabChange).toHaveBeenCalledWith('tab2');
    expect(handleTabChange).toHaveBeenCalledTimes(1);
  });

  it('renders children content in the tab content area', () => {
    const testContent = 'This is my tab content';
    render(
      <Tabs tabs={mockTabs} activeTab="tab1" onTabChange={() => {}}>
        <div>{testContent}</div>
      </Tabs>
    );

    expect(screen.getByText(testContent)).toBeInTheDocument();
  });

  it('renders conditionally filtered tabs', () => {
    const conditionalTabs = [
      { id: 'tab1', label: 'Always Visible' },
      { id: 'tab2', label: 'Sometimes Visible' },
    ];
    const shiftsEnabled = false;
    const tabsToRender = [
      conditionalTabs[0],
      ...(shiftsEnabled ? [conditionalTabs[1]] : []),
    ];

    render(
      <Tabs tabs={tabsToRender} activeTab="tab1" onTabChange={() => {}}>
        <div>Content</div>
      </Tabs>
    );

    expect(screen.getByRole('button', { name: /always visible/i })).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: /sometimes visible/i })).not.toBeInTheDocument();
  });

  it('applies correct CSS classes to tab navigation', () => {
    const { container } = render(
      <Tabs tabs={mockTabs} activeTab="tab1" onTabChange={() => {}}>
        <div>Content</div>
      </Tabs>
    );

    const tabsContainer = container.querySelector('.tabs-container');
    const tabsNav = container.querySelector('.tabs-nav');
    const tabContent = container.querySelector('.tab-content');

    expect(tabsContainer).toBeInTheDocument();
    expect(tabsNav).toBeInTheDocument();
    expect(tabContent).toBeInTheDocument();
  });

  it('all tab buttons have correct type attribute', () => {
    render(
      <Tabs tabs={mockTabs} activeTab="tab1" onTabChange={() => {}}>
        <div>Content</div>
      </Tabs>
    );

    const buttons = screen.getAllByRole('button');
    buttons.forEach((button) => {
      expect(button).toHaveAttribute('type', 'button');
    });
  });

  it('handles multiple tab clicks correctly', () => {
    const handleTabChange = vi.fn();

    render(
      <Tabs tabs={mockTabs} activeTab="tab1" onTabChange={handleTabChange}>
        <div>Content</div>
      </Tabs>
    );

    const secondTab = screen.getByRole('button', { name: /second tab/i });
    const thirdTab = screen.getByRole('button', { name: /third tab/i });

    fireEvent.click(secondTab);
    fireEvent.click(thirdTab);

    expect(handleTabChange).toHaveBeenCalledTimes(2);
    expect(handleTabChange).toHaveBeenNthCalledWith(1, 'tab2');
    expect(handleTabChange).toHaveBeenNthCalledWith(2, 'tab3');
  });

  it('renders with empty children without errors', () => {
    render(
      <Tabs tabs={mockTabs} activeTab="tab1" onTabChange={() => {}} />
    );

    expect(screen.getByRole('button', { name: /first tab/i })).toBeInTheDocument();
  });

  it('renders with single tab', () => {
    const singleTab = [{ id: 'only', label: 'Only Tab' }];
    render(
      <Tabs tabs={singleTab} activeTab="only" onTabChange={() => {}}>
        <div>Content</div>
      </Tabs>
    );

    const button = screen.getByRole('button', { name: /only tab/i });
    expect(button).toBeInTheDocument();
    expect(button).toHaveClass('active');
  });
});
