/**
 * Utility functions for tracking and managing recently visited clubs
 */

export interface RecentClub {
  id: string;
  name: string;
  visitedAt: number; // timestamp
}

const RECENT_CLUBS_KEY = 'recent_clubs';
const MAX_RECENT_CLUBS = 5;

/**
 * Get recently visited clubs from localStorage
 */
export const getRecentClubs = (): RecentClub[] => {
  try {
    const stored = localStorage.getItem(RECENT_CLUBS_KEY);
    if (!stored) return [];
    
    const clubs = JSON.parse(stored) as RecentClub[];
    // Sort by most recent first
    return clubs.sort((a, b) => b.visitedAt - a.visitedAt).slice(0, MAX_RECENT_CLUBS);
  } catch (error) {
    console.error('Error loading recent clubs:', error);
    return [];
  }
};

/**
 * Add or update a club in the recent clubs list
 */
export const addRecentClub = (clubId: string, clubName: string): void => {
  try {
    const existing = getRecentClubs();
    const now = Date.now();
    
    // Remove existing entry for this club if it exists
    const filtered = existing.filter(club => club.id !== clubId);
    
    // Add the club at the beginning
    const updated: RecentClub[] = [
      { id: clubId, name: clubName, visitedAt: now },
      ...filtered
    ].slice(0, MAX_RECENT_CLUBS);
    
    localStorage.setItem(RECENT_CLUBS_KEY, JSON.stringify(updated));
  } catch (error) {
    console.error('Error saving recent club:', error);
  }
};

/**
 * Clear all recent clubs
 */
export const clearRecentClubs = (): void => {
  try {
    localStorage.removeItem(RECENT_CLUBS_KEY);
  } catch (error) {
    console.error('Error clearing recent clubs:', error);
  }
};