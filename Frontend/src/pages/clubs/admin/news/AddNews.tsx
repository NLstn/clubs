import { FC, useState } from "react";
import api from "../../../../utils/api";
import { Input, Modal } from '@/components/ui';

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
        <Modal isOpen={isOpen} onClose={handleClose} title="Add News">
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
                <button 
                    onClick={handleSubmit} 
                    className="button-accept"
                    disabled={isSubmitting}
                >
                    {isSubmitting ? (
                        <>
                            <Modal.LoadingSpinner />
                            Adding...
                        </>
                    ) : (
                        'Add News'
                    )}
                </button>
                <button 
                    onClick={handleClose} 
                    className="button-cancel"
                    disabled={isSubmitting}
                >
                    Cancel
                </button>
            </Modal.Actions>
        </Modal>
    );
};

export default AddNews;