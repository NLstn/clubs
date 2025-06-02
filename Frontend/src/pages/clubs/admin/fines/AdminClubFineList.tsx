import { useEffect, useState, useCallback } from "react";
import { useParams } from "react-router-dom";
import api from "../../../../utils/api";
import AddFine from "./AddFine";

interface Fine {
    id: string;
    userName: string;
    amount: number;
    reason: string;
    createdAt: string;
    updatedAt: string;
    paid: boolean;
}

const AdminClubFineList = () => {

    const { id } = useParams();

    const [fines, setFines] = useState<Fine[]>([]);
    const [error, setError] = useState<string | null>(null);
    const [isModalOpen, setIsModalOpen] = useState(false);

    const fetchFines = useCallback(async () => {
        try {
            const response = await api.get(`/api/v1/clubs/${id}/fines`);
            setFines(response.data);
        } catch (err) {
            setError("Failed to fetch fines: " + err);
        }
    }, [id]);

    useEffect(() => {
        fetchFines();
    }, [fetchFines]);

    return (
        <div>
            <h3>Fines</h3>
            {error && <div className="error">{error}</div>}
            <table>
                <thead>
                    <tr>
                        <th>User</th>
                        <th>Amount</th>
                        <th>Reason</th>
                        <th>Created At</th>
                        <th>Updated At</th>
                        <th>Paid</th>
                    </tr>
                </thead>
                <tbody>
                    {fines && fines.map((fine) => (
                        <tr key={fine.id}>
                            <td>{fine.userName}</td>
                            <td>{fine.amount}</td>
                            <td>{fine.reason}</td>
                            <td>{new Date(fine.createdAt).toLocaleString()}</td>
                            <td>{new Date(fine.updatedAt).toLocaleString()}</td>
                            <td>{fine.paid ? "Yes" : "No"}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
            <button onClick={() => setIsModalOpen(true)} className="button-accept">
                Add Fine
            </button>
            <AddFine 
                isOpen={isModalOpen}
                onClose={() => setIsModalOpen(false)}
                clubId={id || ''}
                onSuccess={fetchFines}
            />
        </div>
    );
}

export default AdminClubFineList;