import { useState, useEffect, useCallback } from 'react';
import { useParams } from 'react-router-dom';
import api from '../../../../utils/api';
import { useT } from '../../../../hooks/useTranslation';
import { Input, Table, TableColumn } from '@/components/ui';

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
    const { t } = useT();
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
            setError(t('fines.errors.fetchingTemplates'));
            setLoading(false);
        }
    }, [clubId, t]);

    useEffect(() => {
        // Calling fetchTemplates here is the correct pattern for data fetching
        // eslint-disable-next-line react-hooks/set-state-in-effect
        fetchTemplates();
    }, [fetchTemplates]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!formData.description || formData.amount <= 0) {
            setError(t('fines.validationError'));
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
            setError(t('fines.errors.savingTemplate'));
        }
    };

    const handleEdit = (template: FineTemplate) => {
        setFormData({ description: template.description, amount: template.amount });
        setEditingId(template.id);
        setIsAdding(true);
    };

    const handleDelete = async (templateId: string) => {
        if (!confirm(t('fines.deleteConfirmation'))) return;

        try {
            await api.delete(`/api/v1/clubs/${clubId}/fine-templates/${templateId}`);
            await fetchTemplates();
        } catch (err: Error | unknown) {
            console.error('Error deleting fine template:', err instanceof Error ? err.message : 'Unknown error');
            setError(t('fines.errors.deletingTemplate'));
        }
    };

    const handleCancel = () => {
        setFormData({ description: '', amount: 0 });
        setIsAdding(false);
        setEditingId(null);
        setError(null);
    };

    const columns: TableColumn<FineTemplate>[] = [
        {
            key: 'description',
            header: t('fines.description'),
            render: (template) => template.description
        },
        {
            key: 'amount',
            header: t('fines.amount'),
            render: (template) => `$${template.amount.toFixed(2)}`
        },
        {
            key: 'actions',
            header: t('common.actions'),
            render: (template) => (
                <div className="table-actions">
                    <button 
                        onClick={() => handleEdit(template)}
                        className="action-button edit"
                    >
                        {t('common.edit')}
                    </button>
                    <button 
                        onClick={() => handleDelete(template.id)}
                        className="action-button remove"
                    >
                        {t('common.delete')}
                    </button>
                </div>
            )
        }
    ];

    return (
        <div style={{ marginBottom: '2rem' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
                <h3>{t('fines.fineTemplates')}</h3>
                <button 
                    onClick={() => setIsAdding(true)}
                    disabled={isAdding}
                >
                    {t('fines.addTemplate')}
                </button>
            </div>

            {isAdding && (
                <form onSubmit={handleSubmit} style={{ marginBottom: '1rem', padding: '1rem', border: '1px solid #ddd', borderRadius: '4px' }}>
                    <h4>{editingId ? t('fines.editTemplate') : t('fines.addNewTemplate')}</h4>
                    <Input
                        label={t('fines.description')}
                        id="description"
                        type="text"
                        value={formData.description}
                        onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                        placeholder="Enter fine description"
                        required
                    />
                    <Input
                        label={t('fines.amount')}
                        id="amount"
                        type="number"
                        step="0.01"
                        min="0.01"
                        value={formData.amount}
                        onChange={(e) => setFormData({ ...formData, amount: Number(e.target.value) })}
                        placeholder="Enter amount"
                        required
                    />
                    <div className="form-actions">
                        <button type="submit" className="button-accept">
                            {editingId ? t('fines.updateTemplate') : t('fines.addTemplate')}
                        </button>
                        <button type="button" onClick={handleCancel} className="button-cancel">
                            {t('common.cancel')}
                        </button>
                    </div>
                </form>
            )}

            <Table
                columns={columns}
                data={templates}
                keyExtractor={(template) => template.id}
                loading={loading}
                error={error}
                emptyMessage={t('fines.noTemplates')}
                loadingMessage={t('clubs.loading.fineTemplates')}
                errorMessage={error || undefined}
            />
        </div>
    );
};

export default AdminClubFineTemplateList;