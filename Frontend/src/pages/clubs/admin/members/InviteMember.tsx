import { FC, useState } from 'react';
import { Input, Modal, Button } from '@/components/ui';

interface InviteMemberProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (email: string) => void;
}

const InviteMember: FC<InviteMemberProps> = ({ isOpen, onClose, onSubmit }) => {
  const [email, setEmail] = useState('');

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
          />
        </div>
      </Modal.Body>
      <Modal.Actions>
        <Button variant="accept" onClick={() => onSubmit(email)} disabled={!email}>
          Send Invite
        </Button>
        <Button variant="cancel" onClick={onClose}>Cancel</Button>
      </Modal.Actions>
    </Modal>
  );
};

export default InviteMember;