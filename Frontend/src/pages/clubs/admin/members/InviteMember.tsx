import { FC, useState } from 'react';

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
        <div className="form-group">
          <label htmlFor="email">Email</label>
          <input
            id="email"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="Enter email"
          />
        </div>
        <div>
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