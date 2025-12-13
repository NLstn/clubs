import React, { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
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
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<SearchResponse>({ clubs: [], events: [] });
  const [isOpen, setIsOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const searchRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const navigate = useNavigate();
  const { api } = useAuth();

  // Close dropdown when clicking outside
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (searchRef.current && !searchRef.current.contains(event.target as Node)) {
        setIsOpen(false);
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
        setIsOpen(false);
      }
    }, 300);

    return () => clearTimeout(timeoutId);
  }, [query]); // eslint-disable-line react-hooks/exhaustive-deps

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

  return (
    <div className="global-search" ref={searchRef}>
      <div className="search-input-container">
        <input
          ref={inputRef}
          type="text"
          placeholder="Search clubs and events..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
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

      {isOpen && query.trim() && (
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
