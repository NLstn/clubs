import { useState, useCallback } from "react";
import { useParams } from "react-router-dom";
import EditNews from "./EditNews";
import AddNews from "./AddNews";
import { ODataTable, ODataTableColumn, Button, ButtonState } from '@/components/ui';
import api from "../../../../utils/api";

interface News {
    ID: string;
    Title: string;
    Content: string;
    CreatedAt: string;
    UpdatedAt: string;
}

const AdminClubNewsList = () => {
    const { id } = useParams();
    const [selectedNews, setSelectedNews] = useState<News | null>(null);
    const [isEditModalOpen, setIsEditModalOpen] = useState(false);
    const [isAddModalOpen, setIsAddModalOpen] = useState(false);
    const [refreshKey, setRefreshKey] = useState(0);
    const [deleteStates, setDeleteStates] = useState<Record<string, ButtonState>>({});

    const refreshNews = useCallback(() => {
        setRefreshKey(prev => prev + 1);
    }, []);

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

        setDeleteStates(prev => ({ ...prev, [newsId]: 'loading' }));

        try {
            await api.delete(`/api/v2/News('${newsId}')`);
            setDeleteStates(prev => ({ ...prev, [newsId]: 'success' }));
            
            setTimeout(() => {
                refreshNews(); // Refresh the list
                setDeleteStates(prev => {
                    const newState = { ...prev };
                    delete newState[newsId];
                    return newState;
                });
            }, 1000);
        } catch (error) {
            console.error("Error deleting news:", error);
            setDeleteStates(prev => ({ ...prev, [newsId]: 'error' }));
            
            setTimeout(() => {
                setDeleteStates(prev => {
                    const newState = { ...prev };
                    delete newState[newsId];
                    return newState;
                });
            }, 3000);
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

    const columns: ODataTableColumn<News>[] = [
        {
            key: 'Title',
            header: 'Title',
            render: (newsItem) => newsItem.Title,
            sortable: true,
            sortField: 'Title'
        },
        {
            key: 'Content',
            header: 'Content',
            render: (newsItem) => truncateContent(newsItem.Content)
        },
        {
            key: 'CreatedAt',
            header: 'Created',
            render: (newsItem) => formatDateTime(newsItem.CreatedAt),
            sortable: true,
            sortField: 'CreatedAt'
        },
        {
            key: 'UpdatedAt',
            header: 'Updated',
            render: (newsItem) => formatDateTime(newsItem.UpdatedAt),
            sortable: true,
            sortField: 'UpdatedAt'
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
                        onClick={() => handleDeleteNews(newsItem.ID)}
                        state={deleteStates[newsItem.ID] || 'idle'}
                        successMessage="Deleted!"
                        errorMessage="Failed"
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
            <ODataTable
                key={refreshKey}
                endpoint="/api/v2/News"
                filter={`ClubID eq '${id}'`}
                columns={columns}
                keyExtractor={(newsItem) => newsItem.ID}
                pageSize={10}
                initialSortField="CreatedAt"
                initialSortDirection="desc"
                emptyMessage="No news posts available"
                loadingMessage="Loading news..."
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
                onSuccess={refreshNews}
            />
            <AddNews 
                isOpen={isAddModalOpen}
                onClose={() => setIsAddModalOpen(false)}
                clubId={id || ''}
                onSuccess={refreshNews}
            />
        </div>
    );
};

export default AdminClubNewsList;