import { useEffect, useState, useCallback } from "react";
import { useParams } from "react-router-dom";
import api from "../../utils/api";
import { useT } from "../../hooks/useTranslation";

interface News {
    id: string;
    title: string;
    content: string;
    created_at: string;
    updated_at: string;
}

const ClubNews = () => {
    const { t } = useT();
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
            return t('dates.notSet');
        }
        
        try {
            const dateTime = new Date(timestamp);
            if (isNaN(dateTime.getTime())) {
                return t('dates.invalidDateTime');
            }
            return dateTime.toLocaleDateString();
        } catch {
            return t('dates.parseError');
        }
    };

    if (loading) return <div>{t('clubs.loading.news')}</div>;
    if (error) return <div className="error">{t('clubs.errors.loadingNews', { error })}</div>;
    if (!news || news.length === 0) return null; // Don't show the section if there's no news

    return (
        <div className="news-section">
            <h3>{t('clubs.latestNews')}</h3>
            <div>
                {news.map(newsItem => (
                    <div key={newsItem.id} className="news-card">
                        <h4 className="news-title">{newsItem.title}</h4>
                        <p className="news-content">{newsItem.content}</p>
                        <small className="news-meta">
                            {t('news.postedOn')} {formatDateTime(newsItem.created_at)}
                        </small>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default ClubNews;