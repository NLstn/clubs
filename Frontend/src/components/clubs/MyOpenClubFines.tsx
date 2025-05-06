import { useCallback, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import api from "../../utils/api";

interface Fine {
    id: string;
    clubId: string;
    amount: number;
    reason: string;
    createdAt: string;
    updatedAt: string;
    paid: boolean;
}

const MyOpenClubFines = () => {

    const { id } = useParams();
    const [fines, setFines] = useState<Fine[]>([]);
    const [error, setError] = useState<string | null>(null);

    const fetchFines = useCallback(async () => {
        try {
            const response = await api.get(`/api/v1/me/fines`);
            if (!response.data) {
                return;
            }
            const filteredFines = response.data.filter((fine: Fine) => fine.clubId === id);
            setFines(filteredFines);
        } catch (err) {
            setError("Failed to fetch fines: " + err);
        }
    }, [id]);

    const payFine = async (fineId: string) => {
        try {
            await api.patch(`/api/v1/clubs/${id}/fines/${fineId}`, { paid: true });
            // Refresh the fines list after successful payment
            await fetchFines();
        } catch (err) {
            setError("Failed to pay fine: " + err);
        }
    };

    useEffect(() => {
        fetchFines();
    }, [fetchFines]);

    return (
        <div>
            <h3>My Open Fines</h3>
            {error && <div className="error">{error}</div>}
            <table>
                <thead>
                    <tr>
                        <th>Amount</th>
                        <th>Reason</th>
                        <th>Created At</th>
                        <th>Updated At</th>
                        <th>Paid</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {fines && fines.map((fine) => (
                        <tr key={fine.id}>
                            <td>{fine.amount}</td>
                            <td>{fine.reason}</td>
                            <td>{new Date(fine.createdAt).toLocaleString()}</td>
                            <td>{new Date(fine.updatedAt).toLocaleString()}</td>
                            <td>{fine.paid ? "Yes" : "No"}</td>
                            <td>
                                {!fine.paid && (
                                    <button 
                                        onClick={() => payFine(fine.id)}
                                        className="button"
                                    >
                                        Pay Fine
                                    </button>
                                )}
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
};

export default MyOpenClubFines;