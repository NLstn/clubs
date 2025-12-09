import { FC, useState, useEffect, useCallback } from 'react';
import api from '../../../../utils/api';
import { Modal, Button, FormGroup } from '@/components/ui';

interface Shift {
    id: string;
    startTime: string;
    endTime: string;
}

interface Member {
    id: string;
    userId: string;
    name: string;
    role: string;
}

interface ShiftMember {
    id: string;
    name: string;
}

interface EditShiftProps {
    isOpen: boolean;
    onClose: () => void;
    shift: Shift | null;
    clubId: string | undefined;
}

const EditShift: FC<EditShiftProps> = ({ isOpen, onClose, shift, clubId }) => {
    const [shiftMembers, setShiftMembers] = useState<ShiftMember[]>([]);
    const [availableMembers, setAvailableMembers] = useState<Member[]>([]);
    const [selectedMemberId, setSelectedMemberId] = useState('');
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState(false);

    const fetchShiftMembers = useCallback(async () => {
        if (!shift || !clubId) return;

        try {
            setLoading(true);
            // OData v2: Query ShiftMembers for this shift with expanded Member and User
            const response = await api.get(`/api/v2/ShiftMembers?$filter=ShiftID eq '${shift.id}'&$expand=Member($expand=User)`);
            setShiftMembers(response.data);
        } catch {
            setError('Failed to fetch shift members');
        } finally {
            setLoading(false);
        }
    }, [shift, clubId]);

    const fetchAvailableMembers = useCallback(async () => {
        if (!clubId) return;

        try {
            // OData v2: Query Members for this club with expanded User
            const response = await api.get(`/api/v2/Members?$filter=ClubID eq '${clubId}'&$expand=User`);
            setAvailableMembers(response.data);
        } catch {
            setError('Failed to fetch club members');
        }
    }, [clubId]);

    useEffect(() => {
        if (isOpen && shift && clubId) {
            fetchShiftMembers();
            fetchAvailableMembers();
        }
    }, [isOpen, shift, clubId, fetchShiftMembers, fetchAvailableMembers]);

    const addMemberToShift = async () => {
        if (!selectedMemberId || !shift || !clubId) return;

        try {
            // OData v2: Use AddMember action on Shift entity
            await api.post(`/api/v2/Shifts('${shift.id}')/AddMember`, {
                memberId: selectedMemberId
            });

            // Refresh shift members
            await fetchShiftMembers();
            setSelectedMemberId('');
            setError(null);
        } catch {
            setError('Failed to add member to shift');
        }
    };

    const removeMemberFromShift = async (memberId: string) => {
        if (!shift || !clubId) return;

        try {
            // OData v2: Use RemoveMember action on Shift entity
            await api.post(`/api/v2/Shifts('${shift.id}')/RemoveMember`, {
                memberId
            });

            // Refresh shift members
            await fetchShiftMembers();
            setError(null);
        } catch {
            setError('Failed to remove member from shift');
        }
    };

    const getAvailableMembersForSelection = () => {
        // Filter out members who are already assigned to this shift
        return availableMembers.filter(member =>
            !(shiftMembers && shiftMembers.some(shiftMember => shiftMember.id === member.userId))
        );
    };

    if (!isOpen || !shift) return null;

    return (
        <Modal isOpen={isOpen && !!shift} onClose={onClose} title="Edit Shift" maxWidth="700px">
            <Modal.Error error={error} />
            
            <Modal.Body>
                <div className="modal-form-section">
                    <div>
                        <p><strong>Start Time:</strong> {new Date(shift.startTime).toLocaleString()}</p>
                        <p><strong>End Time:</strong> {new Date(shift.endTime).toLocaleString()}</p>
                    </div>

                    <div>
                        <h3>Assigned Members</h3>
                        {loading ? (
                            <p>Loading...</p>
                        ) : (
                            <div>
                                <table className="basic-table">
                                    <thead>
                                        <tr>
                                            <th>Name</th>
                                            <th>Actions</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {shiftMembers && shiftMembers.map(member => (
                                            <tr key={member.id}>
                                                <td>{member.name}</td>
                                                <td>
                                                    <Button
                                                        size="sm"
                                                        variant="cancel"
                                                        onClick={() => removeMemberFromShift(member.id)}
                                                        aria-label="Remove member"
                                                    >
                                                        ×
                                                    </Button>
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        )}
                    </div>

                    <div>
                        <h3>Add Member to Shift</h3>
                        <FormGroup>
                            <label htmlFor="member-select">Select Member</label>
                            <select
                                id="member-select"
                                value={selectedMemberId}
                                onChange={(e) => {
                                    setSelectedMemberId(e.target.value)
                                }}
                            >
                                <option value="">Choose a member...</option>
                                {getAvailableMembersForSelection().map(member => (
                                    <option key={member.userId} value={member.userId}>
                                        {member.name}
                                    </option>
                                ))}
                            </select>
                        </FormGroup>
                        <Button
                            onClick={addMemberToShift}
                            disabled={!selectedMemberId}
                            variant="accept"
                            size="sm"
                        >
                            Add Member
                        </Button>
                    </div>
                </div>
            </Modal.Body>

            <Modal.Actions>
                <Button onClick={onClose} variant="cancel">
                    Close
                </Button>
            </Modal.Actions>
        </Modal>
    );
};

export default EditShift;