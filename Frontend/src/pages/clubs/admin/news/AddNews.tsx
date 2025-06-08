import { FC, useState } from "react";
import api from "../../../../utils/api";

interface AddNewsProps {
    isOpen: boolean;
    onClose: () => void;
    clubId: string;
    onSuccess: () => void;
}

const AddNews: FC<AddNewsProps> = ({ isOpen, onClose, clubId, onSuccess }) => {
    const [title, setTitle] = useState<string>('');
    const [content, setContent] = useState<string>('');
    const [error, setError] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);

    if (!isOpen) return null;

    const handleSubmit = async () => {
        if (!title || !content) {
            setError("Please fill in all fields");
            return;
        }

        setError(null);
        setIsSubmitting(true);
        
        try {
            await api.post(`/api/v1/clubs/${clubId}/news`, { 
                title,
                content
            });
            setTitle('');
            setContent('');
            onSuccess();
            onClose();
        } catch (error: unknown) {
            if (error instanceof Error) {
                setError("Failed to add news: " + error.message);
            } else {
                setError("Failed to add news: Unknown error");
            }
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleClose = () => {
        setTitle('');
        setContent('');
        setError(null);
        onClose();
    };

    return (
        <div className="modal">
            <div className="modal-content">
                <h2>Add News</h2>
                {error && <p style={{ color: 'red' }}>{error}</p>}
                
                <div className="form-group">
                    <label htmlFor="newsTitle">Title</label>
                    <input
                        id="newsTitle"
                        type="text"
                        value={title}
                        onChange={(e) => setTitle(e.target.value)}
                        placeholder="News Title"
                        disabled={isSubmitting}
                    />
                </div>

                <div className="form-group">
                    <label htmlFor="newsContent">Content</label>
                    <textarea
                        id="newsContent"
                        value={content}
                        onChange={(e) => setContent(e.target.value)}
                        placeholder="News Content"
                        disabled={isSubmitting}
                        rows={6}
                        style={{ resize: 'vertical' }}
                    />
                </div>

                <div className="modal-actions">
                    <button 
                        onClick={handleSubmit} 
                        className="button-accept"
                        disabled={isSubmitting}
                    >
                        {isSubmitting ? 'Adding...' : 'Add News'}
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

export default AddNews;