import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import Button from '../Button'

describe('Button Component', () => {
  it('renders with children text', () => {
    render(<Button>Click me</Button>)
    
    expect(screen.getByTestId('custom-button')).toBeInTheDocument()
    expect(screen.getByText('Click me')).toBeInTheDocument()
  })

  it('calls onClick when clicked', () => {
    const handleClick = vi.fn()
    render(<Button onClick={handleClick}>Click me</Button>)
    
    fireEvent.click(screen.getByTestId('custom-button'))
    
    expect(handleClick).toHaveBeenCalledTimes(1)
  })

  it('is disabled when disabled prop is true', () => {
    render(<Button disabled>Click me</Button>)
    
    const button = screen.getByTestId('custom-button')
    expect(button).toBeDisabled()
  })

  it('applies correct CSS class for variant', () => {
    render(<Button variant="secondary">Click me</Button>)
    
    const button = screen.getByTestId('custom-button')
    expect(button).toHaveClass('btn-secondary')
  })

  it('applies primary variant by default', () => {
    render(<Button>Click me</Button>)
    
    const button = screen.getByTestId('custom-button')
    expect(button).toHaveClass('btn-primary')
  })
})