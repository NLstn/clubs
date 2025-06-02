import { useState, useEffect, useCallback } from 'react';
import { useParams } from 'react-router-dom';
import api from '../../../../utils/api';

interface FineTemplate {
    id: string;
    club_id: string;
    description: string;
    amount: number;
    created_at: string;
    created_by: string;
    updated_at: string;
    updated_by: string;
}

const AdminClubFineTemplateList = () => {
    const { id: clubId } = useParams();
    const [templates, setTemplates] = useState<FineTemplate[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [isAdding, setIsAdding] = useState(false);
    const [editingId, setEditingId] = useState<string | null>(null);
    const [formData, setFormData] = useState({ description: '', amount: 0 });

    const fetchTemplates = useCallback(async () => {
        try {
            const response = await api.get(`/api/v1/clubs/${clubId}/fine-templates`);
            setTemplates(response.data);
            setLoading(false);
        } catch (err: Error | unknown) {
            console.error('Error fetching fine templates:', err instanceof Error ? err.message : 'Unknown error');
            setError('Error fetching fine templates');
            setLoading(false);
        }
    }, [clubId]);

    useEffect(() => {
        fetchTemplates();
    }, [fetchTemplates]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!formData.description || formData.amount <= 0) {
            setError('Please provide a valid description and amount');
            return;
        }

        try {
            if (editingId) {
                // Update existing template
                await api.put(`/api/v1/clubs/${clubId}/fine-templates/${editingId}`, formData);
            } else {
                // Create new template
                await api.post(`/api/v1/clubs/${clubId}/fine-templates`, formData);
            }
            setFormData({ description: '', amount: 0 });
            setIsAdding(false);
            setEditingId(null);
            setError(null);
            await fetchTemplates();
        } catch (err: Error | unknown) {
            console.error('Error saving fine template:', err instanceof Error ? err.message : 'Unknown error');
            setError('Error saving fine template');
        }
    };

    const handleEdit = (template: FineTemplate) => {
        setFormData({ description: template.description, amount: template.amount });
        setEditingId(template.id);
        setIsAdding(true);
    };

    const handleDelete = async (templateId: string) => {
        if (!confirm('Are you sure you want to delete this fine template?')) return;

        try {
            await api.delete(`/api/v1/clubs/${clubId}/fine-templates/${templateId}`);
            await fetchTemplates();
        } catch (err: Error | unknown) {
            console.error('Error deleting fine template:', err instanceof Error ? err.message : 'Unknown error');
            setError('Error deleting fine template');
        }
    };

    const handleCancel = () => {
        setFormData({ description: '', amount: 0 });
        setIsAdding(false);
        setEditingId(null);
        setError(null);
    };

    if (loading) return <div>Loading fine templates...</div>;

    return (
        <div style={{ marginBottom: '2rem' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
                <h3>Fine Templates</h3>
                <button 
                    onClick={() => setIsAdding(true)}
                    disabled={isAdding}
                >
                    Add Template
                </button>
            </div>

            {error && <div className="error">{error}</div>}

            {isAdding && (
                <form onSubmit={handleSubmit} style={{ marginBottom: '1rem', padding: '1rem', border: '1px solid #ddd', borderRadius: '4px' }}>
                    <h4>{editingId ? 'Edit Template' : 'Add New Template'}</h4>
                    <div className="form-group">
                        <label htmlFor="description">Description</label>
                        <input
                            id="description"
                            type="text"
                            value={formData.description}
                            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                            placeholder="Enter fine description"
                            required
                        />
                    </div>
                    <div className="form-group">
                        <label htmlFor="amount">Amount</label>
                        <input
                            id="amount"
                            type="number"
                            step="0.01"
                            min="0.01"
                            value={formData.amount}
                            onChange={(e) => setFormData({ ...formData, amount: Number(e.target.value) })}
                            placeholder="Enter amount"
                            required
                        />
                    </div>
                    <div className="form-actions">
                        <button type="submit" className="button-accept">
                            {editingId ? 'Update' : 'Add'} Template
                        </button>
                        <button type="button" onClick={handleCancel} className="button-cancel">
                            Cancel
                        </button>
                    </div>
                </form>
            )}

            {templates.length === 0 ? (
                <p>No fine templates found. Add one to get started!</p>
            ) : (
                <table>
                    <thead>
                        <tr>
                            <th>Description</th>
                            <th>Amount</th>
                            <th>Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        {templates.map((template) => (
                            <tr key={template.id}>
                                <td>{template.description}</td>
                                <td>${template.amount.toFixed(2)}</td>
                                <td>
                                    <button 
                                        onClick={() => handleEdit(template)}
                                        style={{ marginRight: '0.5rem' }}
                                    >
                                        Edit
                                    </button>
                                    <button 
                                        onClick={() => handleDelete(template.id)}
                                        className="button-cancel"
                                    >
                                        Delete
                                    </button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            )}
        </div>
    );
};

export default AdminClubFineTemplateList;