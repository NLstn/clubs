import { useState, useCallback } from 'react';
import Layout from "../../components/layout/Layout";
import ProfileContentLayout from '../../components/layout/ProfileContentLayout';
import { useAuth } from '../../hooks/useAuth';
import { ODataTable, ODataTableColumn, Button, ConfirmDialog, Input, FormGroup, Modal } from '@/components/ui';
import './Profile.css';

interface APIKey {
  ID: string;
  Name: string;
  KeyPrefix: string;
  IsActive: boolean;
  ExpiresAt?: string;
  CreatedAt: string;
  LastUsedAt?: string;
}

interface CreateAPIKeyResponse {
  ID: string;
  Name: string;
  KeyPrefix: string;
  APIKey: string;
  IsActive: boolean;
  ExpiresAt?: string;
  CreatedAt: string;
}

const ProfileAPIKeys = () => {
  const { api } = useAuth();
  const [message, setMessage] = useState('');
  const [keyToDelete, setKeyToDelete] = useState<string | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [createKeyName, setCreateKeyName] = useState('');
  const [expirationDays, setExpirationDays] = useState('');
  const [isCreating, setIsCreating] = useState(false);
  const [createdKey, setCreatedKey] = useState<CreateAPIKeyResponse | null>(null);
  const [showCreatedKeyModal, setShowCreatedKeyModal] = useState(false);
  const [copiedToClipboard, setCopiedToClipboard] = useState(false);

  const refreshAPIKeys = useCallback(() => {
    // Reload the page or trigger a re-fetch
    // Since ODataTable manages its own state, we need a different approach
    window.location.reload();
  }, []);

  const handleDeleteKey = async () => {
    if (!keyToDelete) return;

    setIsDeleting(true);
    try {
      // Use OData v2 DELETE with entity key
      await api.delete(`/api/v2/APIKeys('${keyToDelete}')`);
      setMessage('API Key deleted successfully');
      setTimeout(() => setMessage(''), 3000);
      refreshAPIKeys(); // Refresh the list
    } catch (error) {
      console.error('Error deleting API key:', error);
      setMessage('Failed to delete API key');
      setTimeout(() => setMessage(''), 3000);
    } finally {
      setIsDeleting(false);
      setKeyToDelete(null);
    }
  };

  const handleCreateKey = async () => {
    if (!createKeyName.trim()) {
      setMessage('Key name is required');
      setTimeout(() => setMessage(''), 3000);
      return;
    }

    setIsCreating(true);
    try {
      const requestBody: { name: string; expiresAt?: string } = {
        name: createKeyName,
      };

      if (expirationDays) {
        const expiresAt = new Date();
        expiresAt.setDate(expiresAt.getDate() + parseInt(expirationDays));
        requestBody.expiresAt = expiresAt.toISOString();
      }

      // Use OData action to create the API key
      const response = await api.post('/api/v2/CreateAPIKey', requestBody);
      const data = response.data;

      // Show the generated plaintext key to the user
      setCreatedKey(data);
      setShowCreatedKeyModal(true);
      setCreateKeyName('');
      setExpirationDays('');
      setShowCreateModal(false);
      setMessage('API Key created successfully');
      // Don't refresh immediately - let user see and copy the key
    } catch (error) {
      console.error('Error creating API key:', error);
      const err = error as { response?: { data?: { error?: { message?: string } } } };
      const errorMessage = err.response?.data?.error?.message || 'Failed to create API key';
      setMessage(errorMessage);
    } finally {
      setIsCreating(false);
    }
  };

  const handleCopyToClipboard = async () => {
    if (!createdKey?.APIKey) return;
    
    try {
      await navigator.clipboard.writeText(createdKey.APIKey);
      setCopiedToClipboard(true);
      setTimeout(() => setCopiedToClipboard(false), 2000);
    } catch (err) {
      console.error('Failed to copy to clipboard:', err);
    }
  };

  const columns: ODataTableColumn<APIKey>[] = [
    {
      key: 'Name',
      header: 'Name',
      sortable: true,
      sortField: 'Name',
      render: (apiKey: APIKey) => apiKey.Name,
    },
    {
      key: 'KeyPrefix',
      header: 'Key Prefix',
      render: (apiKey: APIKey) => <code>{apiKey.KeyPrefix}</code>,
    },
    {
      key: 'IsActive',
      header: 'Status',
      render: (apiKey: APIKey) => (
        <span style={{ color: apiKey.IsActive ? 'green' : 'red' }}>
          {apiKey.IsActive ? 'Active' : 'Inactive'}
        </span>
      ),
    },
    {
      key: 'CreatedAt',
      header: 'Created',
      sortable: true,
      sortField: 'CreatedAt',
      render: (apiKey: APIKey) => new Date(apiKey.CreatedAt).toLocaleDateString(),
    },
    {
      key: 'LastUsedAt',
      header: 'Last Used',
      render: (apiKey: APIKey) => apiKey.LastUsedAt ? new Date(apiKey.LastUsedAt).toLocaleDateString() : 'Never',
    },
    {
      key: 'ExpiresAt',
      header: 'Expires',
      render: (apiKey: APIKey) => apiKey.ExpiresAt ? new Date(apiKey.ExpiresAt).toLocaleDateString() : 'Never',
    },
    {
      key: 'actions',
      header: 'Actions',
      render: (apiKey: APIKey) => (
        <Button
          onClick={() => setKeyToDelete(apiKey.ID)}
          style={{
            padding: '0.5rem 1rem',
            backgroundColor: '#dc3545',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: 'pointer',
          }}
        >
          Delete
        </Button>
      ),
    },
  ];

  return (
    <Layout title="API Keys">
      <ProfileContentLayout title="API Keys">
        <div className="profile-content">
          <h2>API Keys</h2>
          <p>Manage your API keys for programmatic access to the application.</p>

          {message && (
            <div
              style={{
                padding: '1rem',
                marginBottom: '1rem',
                backgroundColor: message.includes('success') ? '#d4edda' : '#f8d7da',
                border: `1px solid ${message.includes('success') ? '#c3e6cb' : '#f5c6cb'}`,
                borderRadius: '4px',
                color: message.includes('success') ? '#155724' : '#721c24',
              }}
            >
              {message}
            </div>
          )}

          <Button
            onClick={() => setShowCreateModal(true)}
            style={{
              marginBottom: '1.5rem',
              padding: '0.75rem 1.5rem',
              backgroundColor: '#007bff',
              color: 'white',
              border: 'none',
              borderRadius: '4px',
              cursor: 'pointer',
            }}
          >
            Create New API Key
          </Button>

          <ODataTable
            endpoint="/api/v2/APIKeys"
            columns={columns}
            keyExtractor={(key: APIKey) => key.ID}
            pageSize={10}
          />

          <Modal
            isOpen={showCreateModal}
            onClose={() => {
              setShowCreateModal(false);
              setCreateKeyName('');
              setExpirationDays('');
            }}
            title="Create New API Key"
          >
            <div style={{ padding: '1.5rem', display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
              <FormGroup>
                <label>Key Name</label>
                <Input
                  type="text"
                  value={createKeyName}
                  onChange={(e) => setCreateKeyName(e.target.value)}
                  placeholder="e.g., My Integration, CI/CD Pipeline"
                />
              </FormGroup>

              <FormGroup>
                <label>Expiration (optional)</label>
                <Input
                  type="number"
                  value={expirationDays}
                  onChange={(e) => setExpirationDays(e.target.value)}
                  placeholder="Days until expiration (leave empty for no expiration)"
                  min="1"
                  max="365"
                />
              </FormGroup>

              <div style={{ display: 'flex', gap: '1rem', justifyContent: 'flex-end', marginTop: '0.5rem' }}>
                <Button
                  onClick={() => {
                    setShowCreateModal(false);
                    setCreateKeyName('');
                    setExpirationDays('');
                  }}
                  style={{
                    padding: '0.5rem 1rem',
                    backgroundColor: '#6c757d',
                    color: 'white',
                    border: 'none',
                    borderRadius: '4px',
                    cursor: 'pointer',
                  }}
                >
                  Cancel
                </Button>
                <Button
                  onClick={handleCreateKey}
                  disabled={isCreating}
                  style={{
                    padding: '0.5rem 1rem',
                    backgroundColor: '#28a745',
                    color: 'white',
                    border: 'none',
                    borderRadius: '4px',
                    cursor: isCreating ? 'not-allowed' : 'pointer',
                    opacity: isCreating ? 0.6 : 1,
                  }}
                >
                  {isCreating ? 'Creating...' : 'Create Key'}
                </Button>
              </div>
            </div>
          </Modal>

          <Modal
            isOpen={showCreatedKeyModal}
            onClose={() => setShowCreatedKeyModal(false)}
            title="API Key Created Successfully"
          >
            <div style={{ padding: '1.5rem', display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
              <div
                style={{
                  padding: '1rem',
                  backgroundColor: '#fff3cd',
                  border: '1px solid #ffc107',
                  borderRadius: '4px',
                  color: '#856404',
                }}
              >
                <strong>Important:</strong> Copy your API key now. You won't be able to see it again!
              </div>

              {createdKey && (
                <>
                  <div>
                    <strong>Key Name:</strong> {createdKey.Name}
                  </div>
                  <div>
                    <strong>Key Prefix:</strong> <code>{createdKey.KeyPrefix}</code>
                  </div>
                  <div>
                    <strong>Full API Key:</strong>
                    <div
                      style={{
                        padding: '0.75rem',
                        backgroundColor: '#f8f9fa',
                        border: '1px solid #dee2e6',
                        borderRadius: '4px',
                        fontFamily: 'monospace',
                        wordBreak: 'break-all',
                        marginTop: '0.5rem',
                      }}
                    >
                      {createdKey.APIKey}
                    </div>
                  </div>
                </>
              )}

              <div style={{ display: 'flex', gap: '1rem', marginTop: '0.5rem' }}>
                <Button
                  onClick={handleCopyToClipboard}
                  style={{
                    padding: '0.5rem 1rem',
                    backgroundColor: '#007bff',
                    color: 'white',
                    border: 'none',
                    borderRadius: '4px',
                    cursor: 'pointer',
                  }}
                >
                  {copiedToClipboard ? '✓ Copied!' : 'Copy to Clipboard'}
                </Button>
                <Button
                  onClick={() => {
                    setShowCreatedKeyModal(false);
                    refreshAPIKeys(); // Refresh after closing
                  }}
                  style={{
                    padding: '0.5rem 1rem',
                    backgroundColor: '#6c757d',
                    color: 'white',
                    border: 'none',
                    borderRadius: '4px',
                    cursor: 'pointer',
                  }}
                >
                  Done
                </Button>
              </div>
            </div>
          </Modal>

          <ConfirmDialog
            isOpen={!!keyToDelete}
            title="Delete API Key"
            message="Are you sure you want to delete this API key? This action cannot be undone."
            onConfirm={handleDeleteKey}
            onClose={() => setKeyToDelete(null)}
            isLoading={isDeleting}
          />
        </div>
      </ProfileContentLayout>
    </Layout>
  );
};

export default ProfileAPIKeys;
