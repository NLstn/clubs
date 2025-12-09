import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen, act } from '@testing-library/react'
import '@testing-library/jest-dom'
import { ThemeProvider } from '../ThemeProvider'
import { useTheme } from '../../hooks/useTheme'

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
}
Object.defineProperty(window, 'localStorage', { value: localStorageMock })

// Mock window.matchMedia
const mockMatchMedia = vi.fn()
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: mockMatchMedia,
})

// Test component that uses the ThemeContext
const TestComponent = () => {
  const { theme, effectiveTheme, setTheme } = useTheme()
  
  return (
    <div>
      <div data-testid="theme">{theme}</div>
      <div data-testid="effectiveTheme">{effectiveTheme}</div>
      <button onClick={() => setTheme('light')}>Set Light</button>
      <button onClick={() => setTheme('dark')}>Set Dark</button>
      <button onClick={() => setTheme('system')}>Set System</button>
    </div>
  )
}

describe('ThemeProvider', () => {
  let mediaQueryListeners: ((e: MediaQueryListEvent) => void)[] = []

  beforeEach(() => {
    vi.clearAllMocks()
    localStorageMock.getItem.mockReturnValue(null)
    mediaQueryListeners = []

    // Setup matchMedia mock with dark theme preference
    mockMatchMedia.mockImplementation((query: string) => ({
      matches: query === '(prefers-color-scheme: dark)',
      media: query,
      onchange: null,
      addEventListener: vi.fn((event: string, handler: (e: MediaQueryListEvent) => void) => {
        if (event === 'change') {
          mediaQueryListeners.push(handler)
        }
      }),
      removeEventListener: vi.fn(),
      addListener: vi.fn((handler: (e: MediaQueryListEvent) => void) => {
        mediaQueryListeners.push(handler)
      }),
      removeListener: vi.fn(),
      dispatchEvent: vi.fn(),
    }))
  })

  afterEach(() => {
    vi.clearAllMocks()
    document.documentElement.className = ''
  })

  it('initializes with system theme when no localStorage value exists', () => {
    localStorageMock.getItem.mockReturnValue(null)

    render(
      <ThemeProvider>
        <TestComponent />
      </ThemeProvider>
    )

    expect(screen.getByTestId('theme')).toHaveTextContent('system')
    expect(screen.getByTestId('effectiveTheme')).toHaveTextContent('dark') // Based on matchMedia mock
  })

  it('initializes with saved theme from localStorage', () => {
    localStorageMock.getItem.mockReturnValue('light')

    render(
      <ThemeProvider>
        <TestComponent />
      </ThemeProvider>
    )

    expect(screen.getByTestId('theme')).toHaveTextContent('light')
    expect(screen.getByTestId('effectiveTheme')).toHaveTextContent('light')
  })

  it('updates theme and saves to localStorage when setTheme is called', async () => {
    localStorageMock.getItem.mockReturnValue(null)

    render(
      <ThemeProvider>
        <TestComponent />
      </ThemeProvider>
    )

    const lightButton = screen.getByText('Set Light')
    
    await act(async () => {
      lightButton.click()
    })

    expect(screen.getByTestId('theme')).toHaveTextContent('light')
    expect(screen.getByTestId('effectiveTheme')).toHaveTextContent('light')
    expect(localStorageMock.setItem).toHaveBeenCalledWith('user-theme-preference', 'light')
  })

  it('applies correct CSS class to document element', () => {
    localStorageMock.getItem.mockReturnValue('light')

    render(
      <ThemeProvider>
        <TestComponent />
      </ThemeProvider>
    )

    expect(document.documentElement.classList.contains('light-theme')).toBe(true)
    expect(document.documentElement.style.colorScheme).toBe('light')
  })

  it('detects system theme preference changes', async () => {
    localStorageMock.getItem.mockReturnValue('system')

    render(
      <ThemeProvider>
        <TestComponent />
      </ThemeProvider>
    )

    expect(screen.getByTestId('effectiveTheme')).toHaveTextContent('dark')

    // Simulate system theme change to light
    await act(async () => {
      const event = { matches: false } as MediaQueryListEvent
      mediaQueryListeners.forEach(listener => listener(event))
    })

    expect(screen.getByTestId('effectiveTheme')).toHaveTextContent('light')
  })

  it('resolves effectiveTheme from system preference when theme is system', () => {
    localStorageMock.getItem.mockReturnValue('system')

    // Mock light system preference
    mockMatchMedia.mockImplementation((query: string) => ({
      matches: query === '(prefers-color-scheme: light)',
      media: query,
      onchange: null,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      addListener: vi.fn(),
      removeListener: vi.fn(),
      dispatchEvent: vi.fn(),
    }))

    render(
      <ThemeProvider>
        <TestComponent />
      </ThemeProvider>
    )

    expect(screen.getByTestId('theme')).toHaveTextContent('system')
    expect(screen.getByTestId('effectiveTheme')).toHaveTextContent('light')
  })

  it('switches between themes correctly', async () => {
    localStorageMock.getItem.mockReturnValue(null)

    render(
      <ThemeProvider>
        <TestComponent />
      </ThemeProvider>
    )

    // Start with system
    expect(screen.getByTestId('theme')).toHaveTextContent('system')

    // Switch to light
    await act(async () => {
      screen.getByText('Set Light').click()
    })
    expect(screen.getByTestId('theme')).toHaveTextContent('light')
    expect(screen.getByTestId('effectiveTheme')).toHaveTextContent('light')

    // Switch to dark
    await act(async () => {
      screen.getByText('Set Dark').click()
    })
    expect(screen.getByTestId('theme')).toHaveTextContent('dark')
    expect(screen.getByTestId('effectiveTheme')).toHaveTextContent('dark')

    // Switch back to system
    await act(async () => {
      screen.getByText('Set System').click()
    })
    expect(screen.getByTestId('theme')).toHaveTextContent('system')
    expect(screen.getByTestId('effectiveTheme')).toHaveTextContent('dark') // Based on mock
  })
})
