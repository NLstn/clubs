import { useState } from 'react';
import Layout from "../../components/layout/Layout";
import PageHeader from '../../components/layout/PageHeader';
import ProfileSidebar from "./ProfileSidebar";
import { useTheme } from "../../hooks/useTheme";
import { ThemeMode } from "../../context/ThemeContext";
import { FormGroup } from '@/components/ui';
import './Profile.css';

const ProfilePreferences = () => {
    const { theme, setTheme, effectiveTheme } = useTheme();
    const [message, setMessage] = useState('');

    const handleThemeChange = (newTheme: ThemeMode) => {
        setTheme(newTheme);
        setMessage('Theme preference saved successfully!');
        setTimeout(() => setMessage(''), 3000);
    };

    const themeOptions: { value: ThemeMode; label: string; description: string }[] = [
        {
            value: 'light',
            label: 'Light Mode',
            description: 'Always use light theme'
        },
        {
            value: 'dark',
            label: 'Dark Mode',
            description: 'Always use dark theme'
        },
        {
            value: 'system',
            label: 'Use System Setting',
            description: 'Automatically switch based on your device settings'
        }
    ];

    return (
        <Layout>
            <div className="profile-container">
                <ProfileSidebar />
                <div className="profile-content">
                    <PageHeader>
                        <h1>User Preferences</h1>
                    </PageHeader>
                    
                    {message && (
                        <div className="success-message" style={{
                            padding: 'var(--space-sm)',
                            backgroundColor: 'var(--color-success-bg)',
                            color: 'var(--color-success-text)',
                            borderRadius: 'var(--border-radius-md)',
                            marginBottom: 'var(--space-md)'
                        }}>
                            {message}
                        </div>
                    )}

                    <div className="settings-section">
                        <h3>Appearance</h3>
                        
                        <FormGroup>
                            <label style={{ 
                                display: 'block', 
                                marginBottom: 'var(--space-sm)',
                                fontWeight: 500
                            }}>
                                Theme Preference
                            </label>
                            
                            <div style={{ 
                                display: 'flex', 
                                flexDirection: 'column', 
                                gap: 'var(--space-sm)' 
                            }}>
                                {themeOptions.map((option) => (
                                    <div
                                        key={option.value}
                                        onClick={() => handleThemeChange(option.value)}
                                        style={{
                                            padding: 'var(--space-md)',
                                            border: `2px solid ${theme === option.value ? 'var(--color-primary)' : 'var(--color-border)'}`,
                                            borderRadius: 'var(--border-radius-md)',
                                            cursor: 'pointer',
                                            transition: 'all 0.2s',
                                            backgroundColor: theme === option.value ? 'var(--color-success-bg)' : 'var(--color-background-light)'
                                        }}
                                        onMouseEnter={(e) => {
                                            if (theme !== option.value) {
                                                e.currentTarget.style.borderColor = 'var(--color-border-light)';
                                            }
                                        }}
                                        onMouseLeave={(e) => {
                                            if (theme !== option.value) {
                                                e.currentTarget.style.borderColor = 'var(--color-border)';
                                            }
                                        }}
                                    >
                                        <div style={{ 
                                            display: 'flex', 
                                            alignItems: 'center',
                                            gap: 'var(--space-sm)'
                                        }}>
                                            <input
                                                type="radio"
                                                name="theme"
                                                value={option.value}
                                                checked={theme === option.value}
                                                onChange={() => handleThemeChange(option.value)}
                                                style={{ cursor: 'pointer' }}
                                            />
                                            <div>
                                                <div style={{ 
                                                    fontWeight: 500,
                                                    marginBottom: '4px'
                                                }}>
                                                    {option.label}
                                                </div>
                                                <div style={{ 
                                                    fontSize: '0.9rem',
                                                    color: 'var(--color-text-secondary)'
                                                }}>
                                                    {option.description}
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                ))}
                            </div>
                            
                            <div style={{ 
                                marginTop: 'var(--space-md)',
                                padding: 'var(--space-sm)',
                                backgroundColor: 'var(--color-background-medium)',
                                borderRadius: 'var(--border-radius-md)',
                                fontSize: '0.9rem'
                            }}>
                                <strong>Current theme:</strong> {effectiveTheme === 'light' ? '☀️ Light' : '🌙 Dark'}
                                {theme === 'system' && ' (from system)'}
                            </div>
                        </FormGroup>
                    </div>
                </div>
            </div>
        </Layout>
    );
};

export default ProfilePreferences;
