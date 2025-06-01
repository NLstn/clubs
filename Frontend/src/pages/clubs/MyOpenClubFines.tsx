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
    createdByName: string;
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
                        <th>Reason</th>
                        <th>Amount</th>
                        <th>Created At</th>
                        <th>Created By</th>
                    </tr>
                </thead>
                <tbody>
                    {fines && fines.map((fine) => (
                        <tr key={fine.id}>
                            <td>{fine.reason}</td>
                            <td>{fine.amount}</td>
                            <td>{new Date(fine.createdAt).toLocaleString()}</td>
                            <td>{fine.createdByName}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
};

export default MyOpenClubFines;