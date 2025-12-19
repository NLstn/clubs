import { FC, useState } from "react";
import api from "../../../../utils/api";
import { Input, Modal, Button, ButtonState } from '@/components/ui';

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
    const [buttonState, setButtonState] = useState<ButtonState>('idle');

    if (!isOpen) return null;

    const handleSubmit = async () => {
        if (!title || !content) {
            setError("Please fill in all fields");
            return;
        }

        setError(null);
        setButtonState('loading');
        
        try {
            await api.post(`/api/v2/News`, { 
                Title: title,
                Content: content,
                ClubID: clubId
            });
            setButtonState('success');
            
            setTimeout(() => {
                setTitle('');
                setContent('');
                setButtonState('idle');
                onSuccess();
                onClose();
            }, 1000);
        } catch (error: unknown) {
            setButtonState('error');
            if (error instanceof Error) {
                setError("Failed to add news: " + error.message);
            } else {
                setError("Failed to add news: Unknown error");
            }
            setTimeout(() => setButtonState('idle'), 3000);
        }
    };

    const handleClose = () => {
        setTitle('');
        setContent('');
        setError(null);
        setButtonState('idle');
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
                    successMessage="News added!"
                    errorMessage="Failed to add news"
                >
                    Add News
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

export default AddNews;