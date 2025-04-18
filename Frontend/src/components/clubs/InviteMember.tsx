import React from 'react';

interface InviteMemberProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (email: string) => void;
}

const InviteMember: React.FC<InviteMemberProps> = ({ isOpen, onClose, onSubmit }) => {
  const [email, setEmail] = React.useState('');

  if (!isOpen) return null;

  return (
    <div className="modal-overlay">
      <div className="modal">
        <h2>Invite Member</h2>
        <input
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