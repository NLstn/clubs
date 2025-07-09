import { describe, it, expect, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import { TestI18nProvider } from '../i18n-test-utils';
import LanguageSwitcher from '../../components/LanguageSwitcher';

describe('LanguageSwitcher', () => {
  beforeEach(() => {
    // Reset language to English before each test
    localStorage.removeItem('i18nextLng');
  });

  it('renders language switcher with default English', () => {
    render(
      <TestI18nProvider>
        <LanguageSwitcher />
      </TestI18nProvider>
    );

    expect(screen.getByRole('combobox')).toBeInTheDocument();
    expect(screen.getByDisplayValue('English')).toBeInTheDocument();
  });

  it('allows switching between languages', () => {
    render(
      <TestI18nProvider>
        <LanguageSwitcher />
      </TestI18nProvider>
    );

    const select = screen.getByRole('combobox');
    
    // Switch to German
    fireEvent.change(select, { target: { value: 'de' } });
    expect(select).toHaveValue('de');
    
    // Switch back to English
    fireEvent.change(select, { target: { value: 'en' } });
    expect(select).toHaveValue('en');
  });

  it('displays correct language options', () => {
    render(
      <TestI18nProvider>
        <LanguageSwitcher />
      </TestI18nProvider>
    );

    const englishOption = screen.getByRole('option', { name: 'English' });
    const germanOption = screen.getByRole('option', { name: 'Deutsch' });
    
    expect(englishOption).toBeInTheDocument();
    expect(germanOption).toBeInTheDocument();
  });
});