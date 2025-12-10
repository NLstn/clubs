import { FC, useState, useEffect } from "react";
import api from "../../../../utils/api";
import { Input, Modal, Button } from '@/components/ui';

interface News {
    ID: string;
    Title: string;
    Content: string;
    CreatedAt: string;
    UpdatedAt: string;
}

interface EditNewsProps {
    isOpen: boolean;
    onClose: () => void;
    news: News | null;
    clubId: string | undefined;
    onSuccess: () => void;
}

const EditNews: FC<EditNewsProps> = ({ isOpen, onClose, news, onSuccess }) => {
    const [title, setTitle] = useState<string>('');
    const [content, setContent] = useState<string>('');
    const [error, setError] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);

    useEffect(() => {
        if (news) {
            setTitle(news.Title);
            setContent(news.Content);
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
            await api.patch(`/api/v2/News('${news.ID}')`, { 
                Title: title,
                Content: content
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
        <Modal isOpen={isOpen && !!news} onClose={handleClose} title="Edit News">
            <Modal.Error error={error} />
            
            <Modal.Body>
                <div className="modal-form-section">
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
                </div>
            </Modal.Body>

            <Modal.Actions>
                <Button 
                    variant="accept"
                    onClick={handleSubmit}
                    disabled={isSubmitting}
                >
                    {isSubmitting ? (
                        <>
                            <Modal.LoadingSpinner />
                            Updating...
                        </>
                    ) : (
                        'Update News'
                    )}
                </Button>
                <Button 
                    variant="cancel"
                    onClick={handleClose}
                    disabled={isSubmitting}
                >
                    Cancel
                </Button>
            </Modal.Actions>
        </Modal>
    );
};

export default EditNews;