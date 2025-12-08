import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import Layout from '../../components/layout/Layout';
import { Input, Button, FormGroup } from '@/components/ui';

const CreateClub = () => {
    const navigate = useNavigate();
    const { api } = useAuth();
    const [clubName, setClubName] = useState('');
    const [description, setDescription] = useState('');
    const [message, setMessage] = useState('');

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            const response = await api.post('/api/v1/clubs', { name: clubName, description });
            const createdClub = response.data;
            setMessage('Club created successfully!');
            setTimeout(() => {
                navigate(`/clubs/${createdClub.id}`);
            }, 1000);
        } catch (error) {
            setMessage('Error creating club');
            console.error(error);
        }
    };

    return (
        <Layout title="Create New Club">
            <div>
                <h2>Create New Club</h2>

                {message && <p className={`message ${message.includes('Error') ? 'error' : 'success'}`}>
                    {message}
                </p>}

                <form onSubmit={handleSubmit}>
                    <Input
                        label="Club Name:"
                        id="clubName"
                        type="text"
                        value={clubName}
                        onChange={(e) => setClubName(e.target.value)}
                        autoComplete="off"
                        required
                    />
                    <FormGroup>
                        <label htmlFor="description">Description:</label>
                        <textarea
                            id="description"
                            value={description}
                            onChange={(e) => setDescription(e.target.value)}
                            autoComplete="off"
                        />
                    </FormGroup>
                    <div className="form-actions">
                        <Button type="submit" variant="accept">Create Club</Button>
                        <Button type="button" variant="cancel" onClick={() => navigate('/')}>Cancel</Button>
                    </div>
                </form>
            </div>
        </Layout>
    );
};

export default CreateClub;
