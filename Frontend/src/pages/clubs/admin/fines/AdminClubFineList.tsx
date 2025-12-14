import { useState, useCallback, useMemo } from "react";
import { useParams } from "react-router-dom";
import api from "../../../../utils/api";
import AddFine from "./AddFine";
import AdminClubFineTemplateList from "./AdminClubFineTemplateList";
import { ODataTable, ODataTableColumn, Modal, Button } from '@/components/ui';
import './AdminClubFineList.css';

interface Fine {
    ID: string;
    Amount: number;
    Reason: string;
    CreatedAt: string;
    UpdatedAt: string;
    Paid: boolean;
    User?: {
        FirstName: string;
        LastName: string;
    };
}

const AdminClubFineList = () => {

    const { id } = useParams();

    const [showAllFines, setShowAllFines] = useState(false);
    const [showFineTemplates, setShowFineTemplates] = useState(false);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [refreshKey, setRefreshKey] = useState(0);

    const refreshFines = useCallback(() => {
        setRefreshKey(prev => prev + 1);
    }, []);

    const handleDeleteFine = async (fineId: string) => {
        if (!confirm("Are you sure you want to delete this fine?")) {
            return;
        }

        try {
            await api.delete(`/api/v2/Fines('${fineId}')`);
            refreshFines(); // Refresh the list
        } catch (err) {
            alert("Failed to delete fine: " + err);
        }
    };

    const filter = useMemo(() => {
        const clubFilter = `ClubID eq '${id}'`;
        if (showAllFines) {
            return clubFilter;
        }
        return `${clubFilter} and Paid eq false`;
    }, [id, showAllFines]);

    const columns: ODataTableColumn<Fine>[] = [
        {
            key: 'userName',
            header: 'User',
            render: (fine) => fine.User ? `${fine.User.FirstName} ${fine.User.LastName}` : 'Unknown',
            sortable: true,
            sortField: 'User/FirstName',
        },
        {
            key: 'Amount',
            header: 'Amount',
            render: (fine) => fine.Amount,
            sortable: true,
        },
        {
            key: 'Reason',
            header: 'Reason',
            render: (fine) => fine.Reason,
            sortable: true,
        },
        {
            key: 'CreatedAt',
            header: 'Created At',
            render: (fine) => new Date(fine.CreatedAt).toLocaleString(),
            sortable: true,
        },
        {
            key: 'UpdatedAt',
            header: 'Updated At',
            render: (fine) => new Date(fine.UpdatedAt).toLocaleString(),
            sortable: true,
        },
        {
            key: 'Paid',
            header: 'Paid',
            render: (fine) => fine.Paid ? "Yes" : "No",
            sortable: true,
        },
        {
            key: 'actions',
            header: 'Actions',
            render: (fine) => (
                <Button 
                    variant="cancel"
                    size="sm"
                    onClick={() => handleDeleteFine(fine.ID)}
                >
                    Delete
                </Button>
            )
        }
    ];

    return (
        <div>
            <div className="fines-header">
                <h3>Fines</h3>
                <div className="fines-controls">
                    <label className="checkbox-label">
                        <input
                            type="checkbox"
                            checked={showAllFines}
                            onChange={(e) => setShowAllFines(e.target.checked)}
                        />
                        Show all fines
                    </label>
                    <Button variant="secondary" onClick={() => setShowFineTemplates(true)}>Manage Templates</Button>
                </div>
            </div>
            <ODataTable
                key={`${refreshKey}-${showAllFines}`}
                endpoint="/api/v2/Fines"
                filter={filter}
                expand="User"
                columns={columns}
                keyExtractor={(fine) => fine.ID}
                pageSize={10}
                initialSortField="CreatedAt"
                initialSortDirection="desc"
                emptyMessage="No fines available"
                loadingMessage="Loading fines..."
            />
            <div style={{ marginTop: '20px' }}>
                <Button variant="accept" onClick={() => setIsModalOpen(true)}>
                    Add Fine
                </Button>
            </div>
            <AddFine 
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                clubId={id || ''}
                onSuccess={refreshFines}
            />
            
            <Modal 
                isOpen={showFineTemplates} 
                onClose={() => setShowFineTemplates(false)} 
                title="Manage Fine Templates"
                maxWidth="800px"
            >
                <Modal.Body>
                    <AdminClubFineTemplateList />
                </Modal.Body>
                <Modal.Actions>
                    <Button variant="cancel" onClick={() => setShowFineTemplates(false)}>Close</Button>
                </Modal.Actions>
            </Modal>
        </div>
    );
}

export default AdminClubFineList;