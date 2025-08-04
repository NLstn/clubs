import { FC, useState } from 'react';
import { Input } from '@/components/ui';

interface InviteMemberProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (email: string) => void;
}

const InviteMember: FC<InviteMemberProps> = ({ isOpen, onClose, onSubmit }) => {
  const [email, setEmail] = useState('');

  if (!isOpen) return null;

  return (
    <div className="modal">
      <div className="modal-content">
        <h2>Invite Member</h2>
        <Input
          label="Email"
          id="email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="Enter email"
        />
        <div className="modal-actions">
          <button onClick={() => onSubmit(email)} disabled={!email} className="button-accept">
            Send Invite
          </button>
          <button onClick={onClose} className="button-cancel">Cancel</button>
        </div>
      </div>
    </div>
  );
};

export default InviteMember;