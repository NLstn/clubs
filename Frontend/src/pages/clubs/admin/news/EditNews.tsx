import { FC, useState, useEffect, useRef } from "react";
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
    const [formData, setFormData] = useState({ title: '', content: '', newsId: '' });
    const [error, setError] = useState<string | null>(null);
    const [buttonState, setButtonState] = useState<ButtonState>('idle');
    const timeoutRef = useRef<number | undefined>(undefined);

    // Derive form data from news prop
    const title = news && formData.newsId === news.ID ? formData.title : (news?.Title || '');
    const content = news && formData.newsId === news.ID ? formData.content : (news?.Content || '');

    const setTitle = (newTitle: string) => {
        setFormData(prev => ({ ...prev, title: newTitle, newsId: news?.ID || '' }));
    };

    const setContent = (newContent: string) => {
        setFormData(prev => ({ ...prev, content: newContent, newsId: news?.ID || '' }));
    };

    /* eslint-disable react-hooks/set-state-in-effect */
    useEffect(() => {
        // Reset form data when modal closes or news changes
        // This is a legitimate use of setState in effect for controlled form synchronization
        if (!isOpen || !news) {
            setFormData({ title: '', content: '', newsId: '' });
            setError(null);
            setButtonState('idle');
        } else if (news.ID !== formData.newsId) {
            setFormData({ title: news.Title, content: news.Content, newsId: news.ID });
        }
    }, [isOpen, news]); // eslint-disable-line react-hooks/exhaustive-deps
    /* eslint-enable react-hooks/set-state-in-effect */

    useEffect(() => {
        // Cleanup timeout on unmount
        return () => {
            if (timeoutRef.current) {
                clearTimeout(timeoutRef.current);
            }
        };
    }, []);

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
            
            timeoutRef.current = window.setTimeout(() => {
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
            timeoutRef.current = window.setTimeout(() => setButtonState('idle'), 3000);
        }
    };

    const handleClose = () => {
        setFormData({ title: '', content: '', newsId: '' });
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