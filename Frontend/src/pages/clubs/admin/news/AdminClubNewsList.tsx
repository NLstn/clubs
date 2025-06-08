import { useEffect, useState, useCallback } from "react";
import { useParams } from "react-router-dom";
import EditNews from "./EditNews";
import AddNews from "./AddNews";
import api from "../../../../utils/api";

interface News {
    id: string;
    title: string;
    content: string;
    created_at: string;
    updated_at: string;
}

const AdminClubNewsList = () => {
    const { id } = useParams();
    const [news, setNews] = useState<News[]>([]);
    const [selectedNews, setSelectedNews] = useState<News | null>(null);
    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [isAddModalOpen, setIsAddModalOpen] = useState(false);
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

    const handleEditNews = (newsItem: News) => {
        setSelectedNews(newsItem);
        setIsEditModalOpen(true);
    };

    const handleCloseEditModal = () => {
        setSelectedNews(null);
        setIsEditModalOpen(false);
    };

    const handleDeleteNews = async (newsId: string) => {
        if (!confirm("Are you sure you want to delete this news post?")) {
            return;
        }

        try {
            await api.delete(`/api/v1/clubs/${id}/news/${newsId}`);
            fetchNews(); // Refresh the list
        } catch (error) {
            console.error("Error deleting news:", error);
            setError(error instanceof Error ? error.message : "Failed to delete news");
        }
    };

    const formatDateTime = (timestamp: string) => {
        if (!timestamp || timestamp === 'null') {
            return 'Not set';
        }
        
        try {
            const dateTime = new Date(timestamp);
            if (isNaN(dateTime.getTime())) {
                return 'Invalid date/time';
            }
            return dateTime.toLocaleString();
        } catch {
            return 'Parse error';
        }
    };

    const truncateContent = (content: string, maxLength: number = 100) => {
        if (content.length <= maxLength) {
            return content;
        }
        return content.substring(0, maxLength) + '...';
    };

    return (
        <div>
            <h3>News</h3>
            {loading && <p>Loading news...</p>}
            {error && <p style={{color: 'red'}}>Error: {error}</p>}
            {!loading && !error && (
                <>
                    <table className="basic-table">
                        <thead>
                            <tr>
                                <th>Title</th>
                                <th>Content</th>
                                <th>Created</th>
                                <th>Updated</th>
                                <th>Actions</th>
                            </tr>
                        </thead>
                        <tbody>
                            {news.length === 0 ? (
                                <tr>
                                    <td colSpan={5} style={{textAlign: 'center', fontStyle: 'italic'}}>
                                        No news posts available
                                    </td>
                                </tr>
                            ) : (
                                news.map(newsItem => (
                                    <tr key={newsItem.id}>
                                        <td>{newsItem.title}</td>
                                        <td>{truncateContent(newsItem.content)}</td>
                                        <td>{formatDateTime(newsItem.created_at)}</td>
                                        <td>{formatDateTime(newsItem.updated_at)}</td>
                                        <td>
                                            <button
                                                onClick={() => handleEditNews(newsItem)}
                                                className="button-accept"
                                                style={{marginRight: '5px'}}
                                            >
                                                Edit
                                            </button>
                                            <button
                                                onClick={() => handleDeleteNews(newsItem.id)}
                                                className="button-cancel"
                                            >
                                                Delete
                                            </button>
                                        </td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                    <button onClick={() => setIsAddModalOpen(true)} className="button-accept">
                        Add News
                    </button>
                </>
            )}
            <EditNews
                isOpen={isEditModalOpen}
                onClose={handleCloseEditModal}
                news={selectedNews}
                clubId={id}
                onSuccess={fetchNews}
            />
            <AddNews 
                isOpen={isAddModalOpen}
                onClose={() => setIsAddModalOpen(false)}
                clubId={id || ''}
                onSuccess={fetchNews}
            />
        </div>
    );
};

export default AdminClubNewsList;