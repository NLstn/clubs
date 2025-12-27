import React, { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import { useT } from '../../hooks/useTranslation';
import { getRecentClubs, removeRecentClub, RecentClub } from '../../utils/recentClubs';
import './GlobalSearch.css';

interface SearchResult {
  type: 'club' | 'event';
  id: string;
  name: string;
  description?: string;
  club_id?: string;
  club_name?: string;
  start_time?: string;
  end_time?: string;
}

interface SearchResponse {
  clubs: SearchResult[];
  events: SearchResult[];
}

const GlobalSearch: React.FC = () => {
  const { t } = useT();
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<SearchResponse>({ clubs: [], events: [] });
  const [isOpen, setIsOpen] = useState(false);
  const [isFocused, setIsFocused] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [recentClubs, setRecentClubs] = useState<RecentClub[]>([]);
  const searchRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const navigate = useNavigate();
  const { api } = useAuth();

  // Load recent clubs when focused
  useEffect(() => {
    if (isFocused && !query.trim()) {
      setRecentClubs(getRecentClubs());
    }
  }, [isFocused, query]);

  // Close dropdown when clicking outside
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (searchRef.current && !searchRef.current.contains(event.target as Node)) {
        setIsOpen(false);
        setIsFocused(false);
      }
    }
    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  // Debounced search
  useEffect(() => {
    const timeoutId = setTimeout(() => {
      if (query.trim()) {
        performSearch(query.trim());
      } else {
        setResults({ clubs: [], events: [] });
        // Keep isOpen true if focused to show recent clubs
        if (!isFocused) {
          setIsOpen(false);
        }
      }
    }, 300);

    return () => clearTimeout(timeoutId);
  }, [query, isFocused]); // eslint-disable-line react-hooks/exhaustive-deps

  const performSearch = async (searchQuery: string) => {
    if (!searchQuery) return;

    setIsLoading(true);
    try {
      // OData v2: Use SearchGlobal function
      const response = await api.get(`/api/v2/SearchGlobal(query='${encodeURIComponent(searchQuery)}')`);
      const data = response.data;
      // Map OData response to match expected format
      interface ODataClub { id?: string; ID?: string; name?: string; Name?: string; description?: string; Description?: string; }
      interface ODataEvent { 
        id?: string; 
        ID?: string; 
        name?: string; 
        Name?: string; 
        description?: string; 
        Description?: string; 
        clubId?: string; 
        ClubID?: string; 
        club_id?: string;
        clubName?: string; 
        ClubName?: string; 
        club_name?: string;
        startTime?: string; 
        StartTime?: string; 
        start_time?: string;
        endTime?: string; 
        EndTime?: string;
        end_time?: string;
      }
      const mappedClubs = (data.Clubs || data.clubs || []).map((c: ODataClub) => ({
        type: 'club' as const,
        id: c.ID || c.id || '',
        name: c.Name || c.name || '',
        description: c.Description || c.description
      }));
      const mappedEvents = (data.Events || data.events || []).map((e: ODataEvent) => ({
        type: 'event' as const,
        id: e.ID || e.id,
        name: e.Name || e.name,
        description: e.Description || e.description,
        club_id: e.ClubID || e.clubId || e.club_id,
        club_name: e.ClubName || e.clubName || e.club_name,
        start_time: e.StartTime || e.startTime || e.start_time,
        end_time: e.EndTime || e.endTime || e.end_time
      }));
      setResults({ clubs: mappedClubs, events: mappedEvents });
      setIsOpen(true);
    } catch (error) {
      console.error('Search failed:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleResultClick = (result: SearchResult) => {
    if (result.type === 'club') {
      navigate(`/clubs/${result.id}`);
    } else if (result.type === 'event') {
      navigate(`/clubs/${result.club_id}`); // Navigate to club page for events
    }
    setIsOpen(false);
    setIsFocused(false);
    setQuery('');
  };

  const handleRecentClubClick = async (club: RecentClub) => {
    setIsOpen(false);
    setIsFocused(false);
    setQuery('');
    
    try {
      // First try to check if the club exists (OData v2)
      await api.get(`/api/v2/Clubs('${club.id}')`);
      navigate(`/clubs/${club.id}`);
    } catch {
      // If club doesn't exist, remove it from recent clubs
      console.warn(`Club ${club.id} not found, removing from recent clubs`);
      removeRecentClub(club.id);
      setRecentClubs(getRecentClubs()); // Refresh the list
      
      // Still navigate to the club page, which will show the ClubNotFound component
      navigate(`/clubs/${club.id}`);
    }
  };

  const handleRemoveRecentClub = (e: React.MouseEvent, clubId: string) => {
    e.stopPropagation();
    removeRecentClub(clubId);
    setRecentClubs(getRecentClubs());
  };

  const handleFocus = () => {
    setIsFocused(true);
    setIsOpen(true);
    setRecentClubs(getRecentClubs());
  };

  const handleViewAllClubs = () => {
    navigate('/clubs');
    setIsOpen(false);
    setIsFocused(false);
    setQuery('');
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const totalResults = (results?.clubs?.length || 0) + (results?.events?.length || 0);
  const showRecentClubs = isFocused && !query.trim();
  const showSearchResults = isOpen && query.trim();

  return (
    <div className={`global-search ${isFocused ? 'global-search-focused' : ''}`} ref={searchRef}>
      <div className="search-input-container">
        <input
          ref={inputRef}
          type="text"
          placeholder={t('common.search')}
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onFocus={handleFocus}
          className="search-input"
        />
        <div className="search-icon">
          {isLoading ? (
            <div className="search-loading">⟳</div>
          ) : (
            '🔍'
          )}
        </div>
      </div>

      {/* Show recent clubs when focused with empty query */}
      {showRecentClubs && (
        <div className="search-dropdown">
          <div className="search-section">
            <div className="search-section-header">
              <span className="search-section-title">{t('recentClubs.title')}</span>
            </div>
            {recentClubs.length > 0 ? (
              <>
                {recentClubs.map((club) => (
                  <div
                    key={`recent-${club.id}`}
                    className="search-result-item"
                    onClick={() => handleRecentClubClick(club)}
                  >
                    <div className="search-result-type">Club</div>
                    <div className="search-result-content">
                      <div className="search-result-title">{club.name}</div>
                    </div>
                    <button
                      className="search-result-remove"
                      onClick={(e) => handleRemoveRecentClub(e, club.id)}
                      title={t('recentClubs.removeFromRecent')}
                      aria-label={`${t('recentClubs.removeFromRecent')}: ${club.name}`}
                    >
                      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path d="M18 6L6 18M6 6l12 12" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                      </svg>
                    </button>
                  </div>
                ))}
              </>
            ) : (
              <div className="search-no-results">
                {t('recentClubs.noRecentClubs')}
              </div>
            )}
            <button
              className="search-view-all"
              onClick={handleViewAllClubs}
            >
              {t('recentClubs.viewAllClubs')}
            </button>
          </div>
        </div>
      )}

      {/* Show search results when there's a query */}
      {showSearchResults && (
        <div className="search-dropdown">
          {totalResults === 0 ? (
            <div className="search-no-results">
              No results found for "{query}"
            </div>
          ) : (
            <>
              {results?.clubs && results.clubs.length > 0 && (
                <div className="search-section">
                  <div className="search-section-header">
                    <span className="search-section-title">Clubs ({results?.clubs?.length || 0})</span>
                  </div>
                  {results.clubs?.map((club) => (
                    <div
                      key={`club-${club.id}`}
                      className="search-result-item"
                      onClick={() => handleResultClick(club)}
                    >
                      <div className="search-result-type">Club</div>
                      <div className="search-result-content">
                        <div className="search-result-title">{club.name}</div>
                        {club.description && (
                          <div className="search-result-description">
                            {club.description.length > 100
                              ? `${club.description.substring(0, 100)}...`
                              : club.description
                            }
                          </div>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              )}

              {results?.events && results.events.length > 0 && (
                <div className="search-section">
                  <div className="search-section-header">
                    <span className="search-section-title">Events ({results?.events?.length || 0})</span>
                  </div>
                  {results.events?.map((event) => (
                    <div
                      key={`event-${event.id}`}
                      className="search-result-item"
                      onClick={() => handleResultClick(event)}
                    >
                      <div className="search-result-type">Event</div>
                      <div className="search-result-content">
                        <div className="search-result-title">{event.name}</div>
                        <div className="search-result-meta">
                          <div className="search-result-club">{event.club_name}</div>
                          {event.start_time && (
                            <div className="search-result-date">
                              {formatDate(event.start_time)}
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </>
          )}
        </div>
      )}
    </div>
  );
};

export default GlobalSearch;
