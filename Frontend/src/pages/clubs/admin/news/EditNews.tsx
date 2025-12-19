import { FC, useState, useEffect } from "react";
import api from "../../../../utils/api";
import { Input, Modal, Button, ButtonState } from '@/components/ui';

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
    const [buttonState, setButtonState] = useState<ButtonState>('idle');

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
        setButtonState('loading');
        
        try {
            await api.patch(`/api/v2/News('${news.ID}')`, { 
                Title: title,
                Content: content
            });
            setButtonState('success');
            
            setTimeout(() => {
                setButtonState('idle');
                onSuccess();
                onClose();
            }, 1000);
        } catch (error: unknown) {
            setButtonState('error');
            if (error instanceof Error) {
                setError("Failed to update news: " + error.message);
            } else {
                setError("Failed to update news: Unknown error");
            }
            setTimeout(() => setButtonState('idle'), 3000);
        }
    };

    const handleClose = () => {
        setError(null);
        setButtonState('idle');
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
                        disabled={buttonState === 'loading'}
                    />

                    <Input
                        label="Content"
                        value={content}
                        onChange={(e) => setContent(e.target.value)}
                        placeholder="News Content"
                        disabled={buttonState === 'loading'}
                        multiline
                        rows={6}
                    />
                </div>
            </Modal.Body>

            <Modal.Actions>
                <Button 
                    variant="accept"
                    onClick={handleSubmit}
                    state={buttonState}
                    successMessage="Updated!"
                    errorMessage="Failed to update"
                >
                    Update News
                </Button>
                <Button 
                    variant="cancel"
                    onClick={handleClose}
                    disabled={buttonState === 'loading'}
                >
                    Cancel
                </Button>
            </Modal.Actions>
        </Modal>
    );
};

export default EditNews;