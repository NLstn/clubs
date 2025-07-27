import React, { useState } from 'react';
import { useNotificationPreferences } from '../../hooks/useNotifications';
import { useT } from '../../hooks/useTranslation';
import './NotificationSettings.css';

const NotificationSettings: React.FC = () => {
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
    return <div className="notification-settings-loading">{t('notifications.messages.loading')}</div>;
  }

  if (error && !preferences) {
    return <div className="notification-settings-error">{t('notifications.messages.error', { error })}</div>;
  }

  if (!preferences) {
    return <div className="notification-settings-error">{t('notifications.messages.failedToLoad')}</div>;
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
    <div className="notification-settings">
      <div className="notification-settings-header">
        <h2>{t('notifications.title')}</h2>
        <p>{t('notifications.description')}</p>
      </div>

      {saveMessage && (
        <div className={`save-message ${saveMessage.includes('success') ? 'success' : 'error'}`}>
          {saveMessage}
        </div>
      )}

      <div className="notification-settings-content">
        <div className="settings-table">
          <div className="settings-table-header">
            <div className="settings-table-cell notification-type-header">{t('notifications.notificationType')}</div>
            <div className="settings-table-cell toggle-header">{t('notifications.inApp')}</div>
            <div className="settings-table-cell toggle-header">{t('notifications.email')}</div>
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
            <strong>{t('common.note')}:</strong> {t('notifications.note')}
          </p>
        </div>
      </div>

      {saving && <div className="saving-indicator">{t('notifications.messages.saving')}</div>}
    </div>
  );
};

export default NotificationSettings;