import React, { useState } from 'react';
import Layout from "../../components/layout/Layout";
import ProfileSidebar from "./ProfileSidebar";
import { useNotificationPreferences } from '../../hooks/useNotifications';
import { useT } from '../../hooks/useTranslation';
import './ProfileNotificationSettings.css';

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
        <div className="profile-layout">
          <ProfileSidebar />
          <div className="profile-content">
            <div style={{ 
              textAlign: 'center', 
              padding: 'var(--space-xl)',
              color: 'var(--color-text-secondary)'
            }}>
              <p>Loading notification settings...</p>
            </div>
          </div>
        </div>
      </Layout>
    );
  }

  if (error && !preferences) {
    return (
      <Layout title="Notification Settings">
        <div className="profile-layout">
          <ProfileSidebar />
          <div className="profile-content">
            <div style={{ 
              textAlign: 'center', 
              padding: 'var(--space-xl)',
              color: 'var(--color-error-text)'
            }}>
              <p>Failed to load notification settings</p>
            </div>
          </div>
        </div>
      </Layout>
    );
  }

  if (!preferences) {
    return (
      <Layout title="Notification Settings">
        <div className="profile-layout">
          <ProfileSidebar />
          <div className="profile-content">
            <div style={{ 
              textAlign: 'center', 
              padding: 'var(--space-xl)',
              color: 'var(--color-error-text)'
            }}>
              <p>Notification settings not found</p>
            </div>
          </div>
        </div>
      </Layout>
    );
  }

  const notificationTypes = [
    {
      key: 'memberAdded',
      title: t('notifications.types.memberAdded.title') || 'Member Added',
      description: t('notifications.types.memberAdded.description') || 'When you are added to a club',
      inAppKey: 'memberAddedInApp' as keyof typeof preferences,
      emailKey: 'memberAddedEmail' as keyof typeof preferences,
    },
    {
      key: 'eventCreated',
      title: t('notifications.types.eventCreated.title') || 'Event Created',
      description: t('notifications.types.eventCreated.description') || 'When new events are created in your clubs',
      inAppKey: 'eventCreatedInApp' as keyof typeof preferences,
      emailKey: 'eventCreatedEmail' as keyof typeof preferences,
    },
    {
      key: 'fineAssigned',
      title: t('notifications.types.fineAssigned.title') || 'Fine Assigned',
      description: t('notifications.types.fineAssigned.description') || 'When you are assigned a fine',
      inAppKey: 'fineAssignedInApp' as keyof typeof preferences,
      emailKey: 'fineAssignedEmail' as keyof typeof preferences,
    },
    {
      key: 'newsCreated',
      title: t('notifications.types.newsCreated.title') || 'News Created',
      description: t('notifications.types.newsCreated.description') || 'When news is published in your clubs',
      inAppKey: 'newsCreatedInApp' as keyof typeof preferences,
      emailKey: 'newsCreatedEmail' as keyof typeof preferences,
    },
  ];

  return (
    <Layout title="Notification Settings">
      <div className="profile-layout">
        <ProfileSidebar />
        <div className="profile-content">
          <div className="profile-header">
            <h2>Notification Settings</h2>
            <p>Choose how you want to be notified about activities in your clubs.</p>
          </div>

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
                      <label className="toggle-switch">
                        <input
                          type="checkbox"
                          checked={preferences[type.inAppKey] as boolean}
                          onChange={(e) => handleToggle(type.inAppKey, e.target.checked)}
                          disabled={saving}
                        />
                        <span className="slider"></span>
                      </label>
                    </div>
                    
                    <div className="settings-table-cell toggle-cell">
                      <label className="toggle-switch">
                        <input
                          type="checkbox"
                          checked={preferences[type.emailKey] as boolean}
                          onChange={(e) => handleToggle(type.emailKey, e.target.checked)}
                          disabled={saving}
                        />
                        <span className="slider"></span>
                      </label>
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
      </div>
    </Layout>
  );
};

export default ProfileNotificationSettings;
