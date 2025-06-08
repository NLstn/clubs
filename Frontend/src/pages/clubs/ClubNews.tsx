import { useEffect, useState, useCallback } from "react";
import { useParams } from "react-router-dom";
import api from "../../utils/api";

interface News {
    id: string;
    title: string;
    content: string;
    created_at: string;
    updated_at: string;
}

const ClubNews = () => {
    const { id } = useParams();
    const [news, setNews] = useState<News[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const fetchNews = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await api.get(`/api/v1/clubs/${id}/news`);
            setNews(response.data || []);
        } catch (error) {
            console.error("Error fetching news:", error);
            setError(error instanceof Error ? error.message : "Failed to fetch news");
            setNews([]);
        } finally {
            setLoading(false);
        }
    }, [id]);

    useEffect(() => {
        fetchNews();
    }, [fetchNews]);

    const formatDateTime = (timestamp: string) => {
        if (!timestamp || timestamp === 'null') {
            return 'Not set';
        }
        
        try {
            const dateTime = new Date(timestamp);
            if (isNaN(dateTime.getTime())) {
                return 'Invalid date/time';
            }
            return dateTime.toLocaleDateString();
        } catch {
            return 'Parse error';
        }
    };

    if (loading) return <div>Loading news...</div>;
    if (error) return <div style={{color: 'red'}}>Error loading news: {error}</div>;
    if (!news || news.length === 0) return null; // Don't show the section if there's no news

    return (
        <div style={{ marginTop: '30px' }}>
            <h3>Latest News</h3>
            <div>
                {news.map(newsItem => (
                    <div key={newsItem.id} style={{ 
                        border: '1px solid #ddd', 
                        padding: '15px', 
                        marginBottom: '15px', 
                        borderRadius: '5px',
                        backgroundColor: 'var(--background-color)',
                    }}>
                        <h4 style={{ margin: '0 0 10px 0' }}>{newsItem.title}</h4>
                        <p style={{ margin: '0 0 10px 0', lineHeight: '1.5' }}>{newsItem.content}</p>
                        <small style={{ color: '#666' }}>
                            Posted on {formatDateTime(newsItem.created_at)}
                        </small>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default ClubNews;