import { FC, useState, useEffect } from "react";
import api from "../../../../utils/api";
import { Input } from '@/components/ui';

interface News {
    id: string;
    title: string;
    content: string;
    created_at: string;
    updated_at: string;
}

interface EditNewsProps {
    isOpen: boolean;
    onClose: () => void;
    news: News | null;
    clubId: string | undefined;
    onSuccess: () => void;
}

const EditNews: FC<EditNewsProps> = ({ isOpen, onClose, news, clubId, onSuccess }) => {
    const [title, setTitle] = useState<string>('');
    const [content, setContent] = useState<string>('');
    const [error, setError] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);

    useEffect(() => {
        if (news) {
            setTitle(news.title);
            setContent(news.content);
        }
    }, [news]);

    if (!isOpen || !news) return null;

    const handleSubmit = async () => {
        if (!title || !content) {
            setError("Please fill in all fields");
            return;
        }

        setError(null);
        setIsSubmitting(true);
        
        try {
            await api.put(`/api/v1/clubs/${clubId}/news/${news.id}`, { 
                title,
                content
            });
            onSuccess();
            onClose();
        } catch (error: unknown) {
            if (error instanceof Error) {
                setError("Failed to update news: " + error.message);
            } else {
                setError("Failed to update news: Unknown error");
            }
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleClose = () => {
        setError(null);
        onClose();
    };

    return (
        <div className="modal">
            <div className="modal-content">
                <h2>Edit News</h2>
                {error && <p style={{ color: 'red' }}>{error}</p>}
                
                <Input
                    label="Title"
                    id="newsTitle"
                    type="text"
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                    placeholder="News Title"
                    disabled={isSubmitting}
                />

                <Input
                    label="Content"
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                    placeholder="News Content"
                    disabled={isSubmitting}
                    multiline
                    rows={6}
                />

                <div className="modal-actions">
                    <button 
                        onClick={handleSubmit} 
                        className="button-accept"
                        disabled={isSubmitting}
                    >
                        {isSubmitting ? 'Updating...' : 'Update News'}
                    </button>
                    <button 
                        onClick={handleClose} 
                        className="button-cancel"
                        disabled={isSubmitting}
                    >
                        Cancel
                    </button>
                </div>
            </div>
        </div>
    );
};

export default EditNews;