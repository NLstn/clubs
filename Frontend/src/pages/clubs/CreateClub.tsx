import { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import { useT } from '../../hooks/useTranslation';
import Layout from '../../components/layout/Layout';
import { Input, Button, FormGroup, ButtonState } from '@/components/ui';

const CreateClub = () => {
    const navigate = useNavigate();
    const { api } = useAuth();
    const { t } = useT();
    const [clubName, setClubName] = useState('');
    const [description, setDescription] = useState('');
    const [createButtonState, setCreateButtonState] = useState<ButtonState>('idle');
    const [message, setMessage] = useState('');
    const timeoutRef = useRef<number | undefined>(undefined);

    useEffect(() => {
        // Cleanup timeout on unmount
        return () => {
            if (timeoutRef.current) {
                clearTimeout(timeoutRef.current);
            }
        };
    }, []);

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
            setMessage(t('createClub.successMessage'));
            timeoutRef.current = window.setTimeout(() => {
                // OData returns ID property with capital I
                navigate(`/clubs/${createdClub.ID || createdClub.id}`);
            }, 1000);
        } catch (error) {
            setCreateButtonState('error');
            setMessage(t('createClub.errorMessage'));
            console.error(error);
            timeoutRef.current = window.setTimeout(() => setCreateButtonState('idle'), 3000);
        }
    };

    return (
        <Layout title={t('createClub.title')}>
            <div>
                <h2>{t('createClub.title')}</h2>

                {message && <p className={`message ${message.includes(t('common.error')) || message.includes('Error') ? 'error' : 'success'}`}>
                    {message}
                </p>}

                <form onSubmit={handleSubmit}>
                    <Input
                        label={t('createClub.clubName')}
                        id="clubName"
                        type="text"
                        value={clubName}
                        onChange={(e) => setClubName(e.target.value)}
                        autoComplete="off"
                        required
                        disabled={createButtonState === 'loading'}
                    />
                    <FormGroup>
                        <label htmlFor="description">{t('createClub.description')}</label>
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
                            successMessage={t('createClub.successMessage')}
                            errorMessage={t('createClub.errorMessage')}
                        >
                            {t('createClub.createClub')}
                        </Button>
                        <Button 
                            type="button" 
                            variant="cancel" 
                            onClick={() => navigate('/')}
                            disabled={createButtonState === 'loading'}
                        >
                            {t('common.cancel')}
                        </Button>
                    </div>
                </form>
            </div>
        </Layout>
    );
};

export default CreateClub;
