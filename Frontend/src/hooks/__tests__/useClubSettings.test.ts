import { describe, it, expect, beforeEach, vi } from 'vitest';
import { renderHook } from '@testing-library/react';
import MockAdapter from 'axios-mock-adapter';
import axios from 'axios';
import { useClubSettings } from '../useClubSettings';

// Create a shared axios instance and mock
const axiosInstance = axios.create({ baseURL: 'http://localhost:8080' });
const mockAxios = new MockAdapter(axiosInstance);

// Mock the api module to return our mocked axios instance
vi.mock('../utils/api', () => ({
  default: axiosInstance
}));

describe('useClubSettings', () => {
  beforeEach(() => {
    mockAxios.reset();
  });

  it('should return default settings initially', () => {
    const { result } = renderHook(() => useClubSettings('123'));

    expect(result.current.settings).toBeNull();
    expect(result.current.loading).toBe(true);
    expect(result.current.error).toBeNull();
  });
});
