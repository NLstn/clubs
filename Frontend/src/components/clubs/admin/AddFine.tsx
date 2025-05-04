import { FC, useState, useEffect } from "react";
import api from "../../../utils/api";
import './AddFine.css';

interface Member {
    id: string;
    name: string;
    role: string;
}

interface AddFineProps {
    isOpen: boolean;
    onClose: () => void;
    clubId: string;
    onSuccess: () => void;
}

const AddFine: FC<AddFineProps> = ({ isOpen, onClose, clubId, onSuccess }) => {
    const [amount, setAmount] = useState<number>(0);
    const [reason, setReason] = useState<string>('');
    const [error, setError] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [members, setMembers] = useState<Member[]>([]);
    const [selectedMember, setSelectedMember] = useState<Member | null>(null);
    const [searchQuery, setSearchQuery] = useState<string>('');

    useEffect(() => {
        const fetchMembers = async () => {
            try {
                const response = await api.get(`/api/v1/clubs/${clubId}/members`);
                setMembers(response.data);
            } catch (error: unknown) {
                if (error instanceof Error) {
                    setError("Failed to fetch members: " + error.message);
                } else {
                    setError("Failed to fetch members: Unknown error");
                }
            }
        };
        if (isOpen) {
            fetchMembers();
        }
    }, [clubId, isOpen]);

    if (!isOpen) return null;

    const handleSubmit = async () => {
        if (!selectedMember) {
            setError("Please select a member");
            return;
        }
        setError(null);
        setIsSubmitting(true);
        try {
            await api.post(`/api/v1/clubs/${clubId}/fines`, { 
                amount, 
                reason,
                userId: selectedMember.id 
            });
            setAmount(0);
            setReason('');
            setSelectedMember(null);
            setSearchQuery('');
            onSuccess();
            onClose();
        } catch (error: unknown) {
            if (error instanceof Error) {
                setError("Failed to add fine: " + error.message);
            } else {
                setError("Failed to add fine: Unknown error");
            }
        } finally {
            setIsSubmitting(false);
        }
    };

    const filteredMembers = members.filter(member =>
        member.name.toLowerCase().includes(searchQuery.toLowerCase())
    );

    return (
        <div className="modal">
            <div className="modal-content">
                <h2>Add Fine</h2>
                {error && <div className="error">{error}</div>}
                <div className="form-group">
                    <label htmlFor="member">Member</label>
                    <div className="member-select">
                        <input
                            id="member"
                            type="text"
                            value={searchQuery}
                            onChange={(e) => {
                                setSearchQuery(e.target.value)
                                setSelectedMember(null);}}
                            placeholder="Search member..."
                            autoComplete="off"
                        />
                        {searchQuery && !selectedMember && (
                            <div className="member-dropdown">
                                {filteredMembers.map(member => (
                                    <div
                                        key={member.id}
                                        className="member-option"
                                        onClick={() => {
                                            setSelectedMember(member);
                                            setSearchQuery(member.name);
                                        }}
                                    >
                                        {member.name}
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                    <label htmlFor="amount">Amount</label>
                    <input
                        id="amount"
                        type="number"
                        value={amount}
                        onChange={(e) => setAmount(Number(e.target.value))}
                        placeholder="Enter amount"
                    />
                    <label htmlFor="reason">Reason</label>
                    <input
                        id="reason"
                        type="text"
                        value={reason}
                        onChange={(e) => setReason(e.target.value)}
                        placeholder="Enter reason"
                    />
                </div>
                <div>
                    <button 
                        onClick={handleSubmit} 
                        disabled={!selectedMember || !amount || !reason || isSubmitting} 
                        className="button-accept"
                    >
                        {isSubmitting ? "Adding..." : "Add Fine"}
                    </button>
                    <button onClick={onClose} className="button-cancel">Cancel</button>
                </div>
            </div>
        </div>
    );
};

export default AddFine;