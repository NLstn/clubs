import { FC, useState } from 'react';
import { Input, Modal, Button, ButtonState } from '@/components/ui';

interface InviteMemberProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (email: string) => Promise<void>;
}

const InviteMember: FC<InviteMemberProps> = ({ isOpen, onClose, onSubmit }) => {
  const [email, setEmail] = useState('');
  const [buttonState, setButtonState] = useState<ButtonState>('idle');

  const handleSubmit = async () => {
    setButtonState('loading');
    try {
      await onSubmit(email);
      setButtonState('success');
      setTimeout(() => {
        setButtonState('idle');
        setEmail('');
        onClose();
      }, 1000);
    } catch (error) {
      setButtonState('error');
      setTimeout(() => setButtonState('idle'), 3000);
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Invite Member">
      <Modal.Body>
        <div className="modal-form-section">
          <Input
            label="Email"
            id="email"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="Enter email"
            disabled={buttonState === 'loading'}
          />
        </div>
      </Modal.Body>
      <Modal.Actions>
        <Button 
          variant="accept" 
          onClick={handleSubmit} 
          disabled={!email}
          state={buttonState}
          successMessage="Invite sent!"
          errorMessage="Failed to send invite"
        >
          Send Invite
        </Button>
        <Button 
          variant="cancel" 
          onClick={onClose}
          disabled={buttonState === 'loading'}
        >
          Cancel
        </Button>
      </Modal.Actions>
    </Modal>
  );
};

export default InviteMember;