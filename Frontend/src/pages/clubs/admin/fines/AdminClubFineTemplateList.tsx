import { useState, useCallback, useMemo } from 'react';
import { useParams } from 'react-router-dom';
import api from '../../../../utils/api';
import { useT } from '../../../../hooks/useTranslation';
import { Input, ODataTable, ODataTableColumn, Button } from '@/components/ui';
import { ODataFilter } from '../../../../utils/odata';

interface FineTemplate {
    ID: string;
    ClubID: string;
    Description: string;
    Amount: number;
    CreatedAt: string;
    CreatedBy: string;
    UpdatedAt: string;
    UpdatedBy: string;
}

const AdminClubFineTemplateList = () => {
    const { t } = useT();
    const { id: clubId } = useParams();
    const [refreshKey, setRefreshKey] = useState(0);
    const [error, setError] = useState<string | null>(null);
    const [isAdding, setIsAdding] = useState(false);
    const [editingId, setEditingId] = useState<string | null>(null);
    const [formData, setFormData] = useState({ description: '', amount: 0 });

    const refreshTemplates = useCallback(() => {
        setRefreshKey(prev => prev + 1);
    }, []);

    // Use ODataFilter helpers to safely escape values and prevent filter injection
    const filter = useMemo(() => ODataFilter.eq('ClubID', clubId || ''), [clubId]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!formData.description || formData.amount <= 0) {
            setError(t('fines.validationError'));
            return;
        }

        try {
            if (editingId) {
                // Update existing template
                await api.patch(`/api/v2/FineTemplates('${editingId}')`, {
                    Description: formData.description,
                    Amount: formData.amount
                });
            } else {
                // Create new template
                await api.post(`/api/v2/FineTemplates`, {
                    Description: formData.description,
                    Amount: formData.amount,
                    ClubID: clubId
                });
            }
            setFormData({ description: '', amount: 0 });
            setIsAdding(false);
            setEditingId(null);
            setError(null);
            refreshTemplates();
        } catch (err: Error | unknown) {
            console.error('Error saving fine template:', err instanceof Error ? err.message : 'Unknown error');
            setError(t('fines.errors.savingTemplate'));
        }
    };

    const handleEdit = (template: FineTemplate) => {
        setFormData({ description: template.Description, amount: template.Amount });
        setEditingId(template.ID);
        setIsAdding(true);
    };

    const handleDelete = async (templateId: string) => {
        if (!confirm(t('fines.deleteConfirmation'))) return;

        try {
            await api.delete(`/api/v2/FineTemplates('${templateId}')`);
            refreshTemplates();
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

    const columns: ODataTableColumn<FineTemplate>[] = [
        {
            key: 'Description',
            header: t('fines.description'),
            render: (template) => template.Description,
            sortable: true,
        },
        {
            key: 'Amount',
            header: t('fines.amount'),
            render: (template) => `$${template.Amount.toFixed(2)}`,
            sortable: true,
        },
        {
            key: 'actions',
            header: t('common.actions'),
            render: (template) => (
                <div className="table-actions">
                    <Button
                        size="sm"
                        variant="secondary"
                        onClick={() => handleEdit(template)}
                    >
                        {t('common.edit')}
                    </Button>
                    <Button
                        size="sm"
                        variant="cancel"
                        onClick={() => handleDelete(template.ID)}
                    >
                        {t('common.delete')}
                    </Button>
                </div>
            )
        }
    ];

    return (
        <div style={{ marginBottom: '2rem' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
                <h3>{t('fines.fineTemplates')}</h3>
                <Button
                    size="sm"
                    variant="accept"
                    onClick={() => setIsAdding(true)}
                    disabled={isAdding}
                >
                    {t('fines.addTemplate')}
                </Button>
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
                        <Button type="submit" variant="accept">
                            {editingId ? t('fines.updateTemplate') : t('fines.addTemplate')}
                        </Button>
                        <Button type="button" variant="cancel" onClick={handleCancel}>
                            {t('common.cancel')}
                        </Button>
                    </div>
                </form>
            )}

            {error && <div className="error" style={{ marginBottom: '1rem', color: 'var(--color-cancel)' }}>{error}</div>}

            <ODataTable
                key={refreshKey}
                endpoint="/api/v2/FineTemplates"
                filter={filter}
                columns={columns}
                keyExtractor={(template) => template.ID}
                pageSize={10}
                emptyMessage={t('fines.noTemplates')}
                loadingMessage={t('clubs.loading.fineTemplates')}
                initialSortField="Description"
                initialSortDirection="asc"
            />
        </div>
    );
};

export default AdminClubFineTemplateList;