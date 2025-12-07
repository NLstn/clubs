import { useEffect, useState, useCallback } from "react";
import { useParams } from "react-router-dom";
import EditNews from "./EditNews";
import AddNews from "./AddNews";
import { Table, TableColumn, Button } from '@/components/ui';
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

    const columns: TableColumn<News>[] = [
        {
            key: 'title',
            header: 'Title',
            render: (newsItem) => newsItem.title
        },
        {
            key: 'content',
            header: 'Content',
            render: (newsItem) => truncateContent(newsItem.content)
        },
        {
            key: 'created_at',
            header: 'Created',
            render: (newsItem) => formatDateTime(newsItem.created_at)
        },
        {
            key: 'updated_at',
            header: 'Updated',
            render: (newsItem) => formatDateTime(newsItem.updated_at)
        },
        {
            key: 'actions',
            header: 'Actions',
            render: (newsItem) => (
                <div>
                    <Button
                        variant="accept"
                        size="sm"
                        onClick={() => handleEditNews(newsItem)}
                        style={{marginRight: '5px'}}
                    >
                        Edit
                    </Button>
                    <Button
                        variant="cancel"
                        size="sm"
                        onClick={() => handleDeleteNews(newsItem.id)}
                    >
                        Delete
                    </Button>
                </div>
            )
        }
    ];

    return (
        <div>
            <h3>News</h3>
            <Table
                columns={columns}
                data={news}
                keyExtractor={(newsItem) => newsItem.id}
                loading={loading}
                error={error}
                emptyMessage="No news posts available"
                loadingMessage="Loading news..."
                errorMessage={error || "Error loading news"}
            />
            <Button 
                variant="accept"
                onClick={() => setIsAddModalOpen(true)}
                style={{ marginTop: '16px' }}
            >
                Add News
            </Button>
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