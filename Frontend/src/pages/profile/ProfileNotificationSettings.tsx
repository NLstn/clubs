import React, { useState } from 'react';
import Layout from "../../components/layout/Layout";
import ProfileContentLayout from '../../components/layout/ProfileContentLayout';
import { useNotificationPreferences } from '../../hooks/useNotifications';
import { useT } from '../../hooks/useTranslation';
import { ToggleSwitch } from '@/components/ui';
import './ProfileNotificationSettings.css';
import './Profile.css';

const ProfileNotificationSettings: React.FC = () => {
  const { t } = useT();
  const {
    preferences,
    loading,
    error,
    updatePreferences,
  } = useNotificationPreferences();

  const [saving, setSaving] = useState(false);
  const [saveMessage, setSaveMessage] = useState<string | null>(null);

  const handleToggle = async (field: string, value: boolean) => {
    if (!preferences) return;
    
    try {
      setSaving(true);
      setSaveMessage(null);
      
      await updatePreferences({ [field]: value });
      
      setSaveMessage(t('notifications.messages.saveSuccess'));
      setTimeout(() => setSaveMessage(null), 3000);
    } catch {
      setSaveMessage(t('notifications.messages.saveError'));
      setTimeout(() => setSaveMessage(null), 5000);
    } finally {
      setSaving(false);
    }
  };

  if (loading && !preferences) {
    return (
      <Layout title="Notification Settings">
        <ProfileContentLayout title="Notification Settings">
          <div style={{ 
            textAlign: 'center', 
            padding: 'var(--space-xl)',
            color: 'var(--color-text-secondary)'
          }}>
            <p>Loading notification settings...</p>
          </div>
        </ProfileContentLayout>
      </Layout>
    );
  }

  if (error && !preferences) {
    return (
      <Layout title="Notification Settings">
        <ProfileContentLayout title="Notification Settings">
          <div style={{ 
            textAlign: 'center', 
            padding: 'var(--space-xl)',
            color: 'var(--color-error-text)'
          }}>
            <p>Failed to load notification settings</p>
          </div>
        </ProfileContentLayout>
      </Layout>
    );
  }

  if (!preferences) {
    return (
      <Layout title="Notification Settings">
        <ProfileContentLayout title="Notification Settings">
          <div style={{ 
            textAlign: 'center', 
            padding: 'var(--space-xl)',
            color: 'var(--color-error-text)'
          }}>
            <p>Notification settings not found</p>
          </div>
        </ProfileContentLayout>
      </Layout>
    );
  }

  const notificationTypes = [
    {
      key: 'memberAdded',
      title: t('notifications.types.memberAdded.title'),
      description: t('notifications.types.memberAdded.description'),
      inAppKey: 'memberAddedInApp' as keyof typeof preferences,
      emailKey: 'memberAddedEmail' as keyof typeof preferences,
    },
    {
      key: 'eventCreated',
      title: t('notifications.types.eventCreated.title'),
      description: t('notifications.types.eventCreated.description'),
      inAppKey: 'eventCreatedInApp' as keyof typeof preferences,
      emailKey: 'eventCreatedEmail' as keyof typeof preferences,
    },
    {
      key: 'fineAssigned',
      title: t('notifications.types.fineAssigned.title'),
      description: t('notifications.types.fineAssigned.description'),
      inAppKey: 'fineAssignedInApp' as keyof typeof preferences,
      emailKey: 'fineAssignedEmail' as keyof typeof preferences,
    },
    {
      key: 'newsCreated',
      title: t('notifications.types.newsCreated.title'),
      description: t('notifications.types.newsCreated.description'),
      inAppKey: 'newsCreatedInApp' as keyof typeof preferences,
      emailKey: 'newsCreatedEmail' as keyof typeof preferences,
    },
  ];

  return (
    <Layout title="Notification Settings">
      <ProfileContentLayout title="Notification Settings">
        {saveMessage && (
          <div className={`save-message ${saveMessage.includes('success') ? 'success' : 'error'}`}>
            {saveMessage}
          </div>
        )}

        <div className="profile-container" style={{ maxWidth: '800px' }}>
            <div className="notification-settings-content">
              <div className="settings-table">
                <div className="settings-table-header">
                  <div className="settings-table-cell notification-type-header">Notification Type</div>
                  <div className="settings-table-cell toggle-header">In-App</div>
                  <div className="settings-table-cell toggle-header">Email</div>
                </div>

                {notificationTypes.map((type) => (
                  <div key={type.key} className="settings-table-row">
                    <div className="settings-table-cell notification-type-cell">
                      <div className="notification-type-info">
                        <h4>{type.title}</h4>
                        <p>{type.description}</p>
                      </div>
                    </div>
                    
                    <div className="settings-table-cell toggle-cell">
                      <ToggleSwitch
                        checked={preferences[type.inAppKey] as boolean}
                        onChange={(checked) => handleToggle(type.inAppKey, checked)}
                        disabled={saving}
                      />
                    </div>
                    
                    <div className="settings-table-cell toggle-cell">
                      <ToggleSwitch
                        checked={preferences[type.emailKey] as boolean}
                        onChange={(checked) => handleToggle(type.emailKey, checked)}
                        disabled={saving}
                      />
                    </div>
                  </div>
                ))}
              </div>

              <div className="notification-settings-footer">
                <p className="settings-note">
                  <strong>Note:</strong> In-app notifications appear in the notification bell icon. Email notifications are sent to your registered email address.
                </p>
              </div>
            </div>
          {saving && (
            <div className="saving-indicator" style={{
              textAlign: 'center',
              padding: 'var(--space-md)',
              color: 'var(--color-text-secondary)'
            }}>
              Saving...
            </div>
          )}
        </div>
      </ProfileContentLayout>
    </Layout>
  );
};

export default ProfileNotificationSettings;
