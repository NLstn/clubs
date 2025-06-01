import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import TypeAheadDropdown from '../TypeAheadDropdown'

interface TestOption {
  id: string
  label: string
}

const mockOptions: TestOption[] = [
  { id: '1', label: 'Option 1' },
  { id: '2', label: 'Option 2' },
  { id: '3', label: 'Another Option' },
]

describe('TypeAheadDropdown', () => {
  it('renders with placeholder text', () => {
    const mockOnChange = vi.fn()
    const mockOnSearch = vi.fn()

    render(
      <TypeAheadDropdown
        options={mockOptions}
        value={null}
        onChange={mockOnChange}
        onSearch={mockOnSearch}
        placeholder="Search for options..."
      />
    )

    expect(screen.getByPlaceholderText('Search for options...')).toBeInTheDocument()
  })

  it('renders with label when provided', () => {
    const mockOnChange = vi.fn()
    const mockOnSearch = vi.fn()

    render(
      <TypeAheadDropdown
        options={mockOptions}
        value={null}
        onChange={mockOnChange}
        onSearch={mockOnSearch}
        label="Choose Option"
        id="test-dropdown"
      />
    )

    expect(screen.getByLabelText('Choose Option')).toBeInTheDocument()
  })

  it('calls onSearch when input value changes', () => {
    const mockOnChange = vi.fn()
    const mockOnSearch = vi.fn()

    render(
      <TypeAheadDropdown
        options={mockOptions}
        value={null}
        onChange={mockOnChange}
        onSearch={mockOnSearch}
      />
    )

    const input = screen.getByRole('textbox')
    fireEvent.change(input, { target: { value: 'Option' } })

    expect(mockOnSearch).toHaveBeenCalledWith('Option')
  })

  it('shows dropdown options when typing and no value is selected', () => {
    const mockOnChange = vi.fn()
    const mockOnSearch = vi.fn()

    render(
      <TypeAheadDropdown
        options={mockOptions}
        value={null}
        onChange={mockOnChange}
        onSearch={mockOnSearch}
      />
    )

    const input = screen.getByRole('textbox')
    fireEvent.change(input, { target: { value: 'Option' } })

    expect(screen.getByText('Option 1')).toBeInTheDocument()
    expect(screen.getByText('Option 2')).toBeInTheDocument()
    expect(screen.getByText('Another Option')).toBeInTheDocument()
  })

  it('selects option when clicked', () => {
    const mockOnChange = vi.fn()
    const mockOnSearch = vi.fn()

    render(
      <TypeAheadDropdown
        options={mockOptions}
        value={null}
        onChange={mockOnChange}
        onSearch={mockOnSearch}
      />
    )

    const input = screen.getByRole('textbox')
    fireEvent.change(input, { target: { value: 'Option' } })
    
    fireEvent.click(screen.getByText('Option 1'))

    expect(mockOnChange).toHaveBeenCalledWith(mockOptions[0])
  })

  it('shows selected value in input', () => {
    const mockOnChange = vi.fn()
    const mockOnSearch = vi.fn()

    render(
      <TypeAheadDropdown
        options={mockOptions}
        value={mockOptions[0]}
        onChange={mockOnChange}
        onSearch={mockOnSearch}
      />
    )

    const input = screen.getByDisplayValue('Option 1')
    expect(input).toBeInTheDocument()
  })

  it('calls onChange with null when input is cleared', () => {
    const mockOnChange = vi.fn()
    const mockOnSearch = vi.fn()

    render(
      <TypeAheadDropdown
        options={mockOptions}
        value={mockOptions[0]}
        onChange={mockOnChange}
        onSearch={mockOnSearch}
      />
    )

    const input = screen.getByRole('textbox')
    fireEvent.change(input, { target: { value: '' } })

    expect(mockOnChange).toHaveBeenCalledWith(null)
  })

  it('opens dropdown on focus', () => {
    const mockOnChange = vi.fn()
    const mockOnSearch = vi.fn()

    render(
      <TypeAheadDropdown
        options={mockOptions}
        value={null}
        onChange={mockOnChange}
        onSearch={mockOnSearch}
      />
    )

    const input = screen.getByRole('textbox')
    fireEvent.change(input, { target: { value: 'test' } })
    fireEvent.focus(input)

    expect(screen.getByText('Option 1')).toBeInTheDocument()
  })
})