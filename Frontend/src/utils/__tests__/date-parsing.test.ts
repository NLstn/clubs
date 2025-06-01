import { describe, it, expect } from 'vitest'

describe('Date Parsing', () => {
  it('should correctly parse RFC3339 formatted dates from API responses', () => {
    // Use actual RFC3339 timestamps
    const testTimestamp1 = '2025-06-01T18:57:44Z'
    const testTimestamp2 = '2025-06-01T18:57:45Z'
    
    // Mock API response with the new camelCase field names
    const mockFineResponse = {
      id: 'test-id',
      clubId: 'club-id',
      clubName: 'Test Club',
      reason: 'Test fine',
      amount: 25.0,
      createdAt: testTimestamp1,
      updatedAt: testTimestamp2,
      paid: false,
      createdByName: 'Test User'
    }

    // Test that dates can be parsed correctly (this is what happens in the components)
    const createdDate = new Date(mockFineResponse.createdAt)
    const updatedDate = new Date(mockFineResponse.updatedAt)

    // Verify the dates are valid (not "Invalid Date")
    expect(createdDate.toString()).not.toBe('Invalid Date')
    expect(updatedDate.toString()).not.toBe('Invalid Date')
    
    // Verify the dates can be converted to locale strings (as used in the components)
    const createdLocaleString = createdDate.toLocaleString()
    const updatedLocaleString = updatedDate.toLocaleString()
    
    expect(createdLocaleString).not.toBe('Invalid Date')
    expect(updatedLocaleString).not.toBe('Invalid Date')
    
    // Verify the dates are parsed correctly (without being too strict about formatting)
    expect(createdDate.getUTCFullYear()).toBe(2025)
    expect(createdDate.getUTCMonth()).toBe(5) // June (0-indexed)
    expect(createdDate.getUTCDate()).toBe(1)
    expect(createdDate.getUTCHours()).toBe(18)
    expect(createdDate.getUTCMinutes()).toBe(57)
    expect(createdDate.getUTCSeconds()).toBe(44)
  })

  it('should handle the old snake_case format as invalid', () => {
    // Mock the old API response format that was causing issues
    const mockOldResponse = {
      id: 'test-id',
      created_at: '2025-06-01T18:57:44Z',  // snake_case
      updated_at: '2025-06-01T18:57:45Z'   // snake_case
    }

    // With the TypeScript interface expecting camelCase, these will be undefined
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const fine = mockOldResponse as any
    
    // Simulate what would happen in the component with the old format
    const createdDate = new Date(fine.createdAt)  // undefined
    const updatedDate = new Date(fine.updatedAt)  // undefined

    // These should be Invalid Date since the fields are undefined
    expect(createdDate.toString()).toBe('Invalid Date')
    expect(updatedDate.toString()).toBe('Invalid Date')
  })
})