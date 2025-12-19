import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import Layout from '../../components/layout/Layout';
import { Input, Button, FormGroup, ButtonState } from '@/components/ui';

const CreateClub = () => {
    const navigate = useNavigate();
    const { api } = useAuth();
    const [clubName, setClubName] = useState('');
    const [description, setDescription] = useState('');
    const [createButtonState, setCreateButtonState] = useState<ButtonState>('idle');
    const [message, setMessage] = useState('');

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setCreateButtonState('loading');
        setMessage('');
        
        try {
            // OData v2: POST to Clubs entity set
            const response = await api.post('/api/v2/Clubs', { 
                Name: clubName, 
                Description: description 
            });
            const createdClub = response.data;
            setCreateButtonState('success');
            setMessage('Club created successfully!');
            setTimeout(() => {
                // OData returns ID property with capital I
                navigate(`/clubs/${createdClub.ID || createdClub.id}`);
            }, 1000);
        } catch (error) {
            setCreateButtonState('error');
            setMessage('Error creating club');
            console.error(error);
            setTimeout(() => setCreateButtonState('idle'), 3000);
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
                        disabled={createButtonState === 'loading'}
                    />
                    <FormGroup>
                        <label htmlFor="description">Description:</label>
                        <textarea
                            id="description"
                            value={description}
                            onChange={(e) => setDescription(e.target.value)}
                            autoComplete="off"
                            disabled={createButtonState === 'loading'}
                        />
                    </FormGroup>
                    <div className="form-actions">
                        <Button 
                            type="submit" 
                            variant="accept"
                            state={createButtonState}
                            successMessage="Club created successfully!"
                            errorMessage="Error creating club"
                        >
                            Create Club
                        </Button>
                        <Button 
                            type="button" 
                            variant="cancel" 
                            onClick={() => navigate('/')}
                            disabled={createButtonState === 'loading'}
                        >
                            Cancel
                        </Button>
                    </div>
                </form>
            </div>
        </Layout>
    );
};

export default CreateClub;
