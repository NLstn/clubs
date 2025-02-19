import { useState, useEffect } from 'react';
import axios from 'axios';
import './AdminDashboard.css';

interface Club {
    id: number;
    name: string;
    description: string;
}

const AdminDashboard = () => {
    const [clubs, setClubs] = useState<Club[]>([]);
    const [showCreateForm, setShowCreateForm] = useState(false);
    const [clubName, setClubName] = useState('');
    const [description, setDescription] = useState('');
    const [message, setMessage] = useState('');

    useEffect(() => {
        fetchClubs();
    }, []);

    const fetchClubs = async () => {
        try {
            const response = await axios.get('http://localhost:8080/api/v1/clubs');
            setClubs(response.data);
        } catch (error) {
            setMessage('Error fetching clubs');
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await axios.post('http://localhost:8080/api/v1/clubs', { name: clubName, description });
            setMessage('Club created successfully!');
            setClubName('');
            setDescription('');
            setShowCreateForm(false);
            fetchClubs(); // Refresh the list
        } catch (error) {
            setMessage('Error creating club');
        }
    };

    return (
        <div className="admin-dashboard">
            <div className="dashboard-header">
                <h2>Admin Dashboard</h2>
                <button className="create-button" onClick={() => setShowCreateForm(true)}>
                    Create New Club
                </button>
            </div>
            
            {message && <p className={`message ${message.includes('Error') ? 'error' : 'success'}`}>
                {message}
            </p>}

            <div className="clubs-list">
                {clubs.map(club => (
                    <div key={club.id} className="club-item">
                        <h4>{club.name}</h4>
                        <p>{club.description}</p>
                    </div>
                ))}
            </div>

            {showCreateForm && (
                <>
                    <div className="modal-overlay" onClick={() => setShowCreateForm(false)} />
                    <div className="modal">
                        <h3>Create New Club</h3>
                        <form onSubmit={handleSubmit}>
                            <div className="form-group">
                                <label>Club Name:</label>
                                <input
                                    type="text"
                                    value={clubName}
                                    onChange={(e) => setClubName(e.target.value)}
                                    required
                                />
                            </div>
                            <div className="form-group">
                                <label>Description:</label>
                                <textarea
                                    value={description}
                                    onChange={(e) => setDescription(e.target.value)}
                                />
                            </div>
                            <div className="form-actions">
                                <button type="button" onClick={() => setShowCreateForm(false)}>Cancel</button>
                                <button type="submit">Create Club</button>
                            </div>
                        </form>
                    </div>
                </>
            )}
        </div>
    );
};

export default AdminDashboard;