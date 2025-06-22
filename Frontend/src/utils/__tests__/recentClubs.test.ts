import { describe, it, expect, beforeEach, vi } from 'vitest';
import { getRecentClubs, addRecentClub, clearRecentClubs } from '../recentClubs';

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
};

Object.defineProperty(window, 'localStorage', {
  value: localStorageMock,
});

describe('recentClubs utilities', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('getRecentClubs', () => {
    it('returns empty array when no stored clubs', () => {
      localStorageMock.getItem.mockReturnValue(null);
      
      const result = getRecentClubs();
      
      expect(result).toEqual([]);
      expect(localStorageMock.getItem).toHaveBeenCalledWith('recent_clubs');
    });

    it('returns stored clubs sorted by most recent first', () => {
      const storedClubs = [
        { id: '1', name: 'Club A', visitedAt: 1000 },
        { id: '2', name: 'Club B', visitedAt: 2000 },
        { id: '3', name: 'Club C', visitedAt: 1500 },
      ];
      
      localStorageMock.getItem.mockReturnValue(JSON.stringify(storedClubs));
      
      const result = getRecentClubs();
      
      expect(result).toEqual([
        { id: '2', name: 'Club B', visitedAt: 2000 },
        { id: '3', name: 'Club C', visitedAt: 1500 },
        { id: '1', name: 'Club A', visitedAt: 1000 },
      ]);
    });

    it('limits results to 5 clubs maximum', () => {
      const storedClubs = Array.from({ length: 7 }, (_, i) => ({
        id: `${i}`,
        name: `Club ${i}`,
        visitedAt: i * 1000,
      }));
      
      localStorageMock.getItem.mockReturnValue(JSON.stringify(storedClubs));
      
      const result = getRecentClubs();
      
      expect(result).toHaveLength(5);
    });

    it('handles JSON parse errors gracefully', () => {
      localStorageMock.getItem.mockReturnValue('invalid json');
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
      
      const result = getRecentClubs();
      
      expect(result).toEqual([]);
      expect(consoleSpy).toHaveBeenCalled();
      
      consoleSpy.mockRestore();
    });
  });

  describe('addRecentClub', () => {
    it('adds new club to empty list', () => {
      localStorageMock.getItem.mockReturnValue(null);
      const mockDate = 1234567890;
      vi.spyOn(Date, 'now').mockReturnValue(mockDate);
      
      addRecentClub('club1', 'Club One');
      
      expect(localStorageMock.setItem).toHaveBeenCalledWith(
        'recent_clubs',
        JSON.stringify([{ id: 'club1', name: 'Club One', visitedAt: mockDate }])
      );
    });

    it('updates existing club position and timestamp', () => {
      const existingClubs = [
        { id: 'club1', name: 'Club One', visitedAt: 1000 },
        { id: 'club2', name: 'Club Two', visitedAt: 2000 },
      ];
      
      localStorageMock.getItem.mockReturnValue(JSON.stringify(existingClubs));
      const mockDate = 3000;
      vi.spyOn(Date, 'now').mockReturnValue(mockDate);
      
      addRecentClub('club1', 'Club One Updated');
      
      expect(localStorageMock.setItem).toHaveBeenCalledWith(
        'recent_clubs',
        JSON.stringify([
          { id: 'club1', name: 'Club One Updated', visitedAt: mockDate },
          { id: 'club2', name: 'Club Two', visitedAt: 2000 },
        ])
      );
    });

    it('limits stored clubs to maximum of 5', () => {
      const existingClubs = Array.from({ length: 5 }, (_, i) => ({
        id: `club${i}`,
        name: `Club ${i}`,
        visitedAt: i * 1000,
      }));
      
      localStorageMock.getItem.mockReturnValue(JSON.stringify(existingClubs));
      const mockDate = 6000;
      vi.spyOn(Date, 'now').mockReturnValue(mockDate);
      
      addRecentClub('club5', 'Club Five');
      
      const setItemCall = localStorageMock.setItem.mock.calls[0];
      const storedClubs = JSON.parse(setItemCall[1]);
      
      expect(storedClubs).toHaveLength(5);
      expect(storedClubs[0]).toEqual({ id: 'club5', name: 'Club Five', visitedAt: mockDate });
    });

    it('handles localStorage errors gracefully', () => {
      localStorageMock.getItem.mockImplementation(() => {
        throw new Error('localStorage error');
      });
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
      
      expect(() => addRecentClub('club1', 'Club One')).not.toThrow();
      expect(consoleSpy).toHaveBeenCalled();
      
      consoleSpy.mockRestore();
    });
  });

  describe('clearRecentClubs', () => {
    it('removes the recent clubs key from localStorage', () => {
      clearRecentClubs();
      
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('recent_clubs');
    });

    it('handles localStorage errors gracefully', () => {
      localStorageMock.removeItem.mockImplementation(() => {
        throw new Error('localStorage error');
      });
      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
      
      expect(() => clearRecentClubs()).not.toThrow();
      expect(consoleSpy).toHaveBeenCalled();
      
      consoleSpy.mockRestore();
    });
  });
});