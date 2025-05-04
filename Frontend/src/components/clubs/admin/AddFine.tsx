import { FC, useState, useEffect } from "react";
import api from "../../../utils/api";
import TypeAheadDropdown from "../../TypeAheadDropdown";

interface Member {
    id: string;
    name: string;
    role: string;
    userId: string;
}

interface MemberOption {
    id: string;
    label: string;
    member: Member;
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
    const [selectedOption, setSelectedOption] = useState<MemberOption | null>(null);
    const [memberOptions, setMemberOptions] = useState<MemberOption[]>([]);

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
        if (!selectedOption) {
            setError("Please select a member");
            return;
        }
        setError(null);
        setIsSubmitting(true);
        try {
            await api.post(`/api/v1/clubs/${clubId}/fines`, { 
                amount, 
                reason,
                userId: selectedOption.member.userId
            });
            setAmount(0);
            setReason('');
            setSelectedOption(null);
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

    const handleSearch = (query: string) => {
        const filtered = members.map(member => ({
            id: member.id,
            label: member.name,
            member: member
        })).filter(option =>
            option.label.toLowerCase().includes(query.toLowerCase())
        );
        setMemberOptions(filtered);
    };

    return (
        <div className="modal">
            <div className="modal-content">
                <h2>Add Fine</h2>
                {error && <div className="error">{error}</div>}
                <div className="form-group">
                    <TypeAheadDropdown<MemberOption>
                        options={memberOptions}
                        value={selectedOption}
                        onChange={setSelectedOption}
                        onSearch={handleSearch}
                        placeholder="Search member..."
                        id="member"
                        label="Member"
                    />
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
                        disabled={!selectedOption || !amount || !reason || isSubmitting} 
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