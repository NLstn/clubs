import { useState } from 'react';
import Layout from "../../components/layout/Layout";
import ProfileContentLayout from '../../components/layout/ProfileContentLayout';
import LanguageSwitcher from "../../components/LanguageSwitcher";
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
        <Layout title="Preferences">
            <ProfileContentLayout title="Preferences">
                {message && (
                    <div className="success-message">
                        {message}
                    </div>
                )}

                <div className="profile-content-sections">
                    <div className="content-section">
                        <h3>Language</h3>
                                <FormGroup>
                                    <label>Preferred Language</label>
                                    <div style={{ marginTop: 'var(--space-xs)' }}>
                                        <LanguageSwitcher />
                                    </div>
                                </FormGroup>
                            </div>

                            <div className="content-section">
                                <h3>Appearance</h3>
                                
                                <FormGroup>
                                    <label>Select Theme</label>
                                    <p style={{ 
                                        fontSize: '0.9rem', 
                                        color: 'var(--color-text-secondary)',
                                        marginTop: 'var(--space-xs)',
                                        marginBottom: 'var(--space-md)'
                                    }}>
                                        Pick a theme to customize the appearance.
                                    </p>
                                    
                                    <div style={{ 
                                        display: 'grid', 
                                        gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
                                        gap: 'var(--space-md)',
                                        marginTop: 'var(--space-sm)'
                                    }}>
                                        {themeOptions.map((option) => (
                                            <div
                                                key={option.value}
                                                onClick={() => handleThemeChange(option.value)}
                                                style={{
                                                    cursor: 'pointer',
                                                    border: `3px solid ${theme === option.value ? 'var(--color-primary)' : 'var(--color-border)'}`,
                                                    borderRadius: 'var(--border-radius-md)',
                                                    overflow: 'hidden',
                                                    transition: 'all 0.2s',
                                                    backgroundColor: 'var(--color-background)',
                                                    position: 'relative'
                                                }}
                                            >
                                                {/* Preview area */}
                                                <div style={{
                                                    height: '140px',
                                                    background: option.value === 'light' 
                                                        ? 'linear-gradient(to bottom, #f5f5f5 0%, #ffffff 100%)'
                                                        : option.value === 'dark'
                                                        ? 'linear-gradient(to bottom, #1a1a1a 0%, #0d0d0d 100%)'
                                                        : 'linear-gradient(135deg, #f5f5f5 0%, #f5f5f5 49%, #1a1a1a 51%, #0d0d0d 100%)',
                                                    display: 'flex',
                                                    alignItems: 'center',
                                                    justifyContent: 'center',
                                                    padding: 'var(--space-md)',
                                                    position: 'relative'
                                                }}>
                                                    {/* Mock UI elements */}
                                                    <div style={{
                                                        width: '100%',
                                                        display: 'flex',
                                                        flexDirection: 'column',
                                                        gap: '8px'
                                                    }}>
                                                        {/* Header bar */}
                                                        <div style={{
                                                            height: '24px',
                                                            background: option.value === 'light' ? '#e0e0e0' : '#2a2a2a',
                                                            borderRadius: '4px',
                                                            display: 'flex',
                                                            alignItems: 'center',
                                                            padding: '0 8px',
                                                            gap: '4px'
                                                        }}>
                                                            <div style={{
                                                                width: '12px',
                                                                height: '12px',
                                                                borderRadius: '50%',
                                                                background: '#4CAF50'
                                                            }}></div>
                                                            <div style={{
                                                                flex: 1,
                                                                height: '8px',
                                                                background: option.value === 'light' ? '#d0d0d0' : '#3a3a3a',
                                                                borderRadius: '2px'
                                                            }}></div>
                                                        </div>
                                                        {/* Content cards */}
                                                        <div style={{
                                                            display: 'grid',
                                                            gridTemplateColumns: '1fr 1fr',
                                                            gap: '6px'
                                                        }}>
                                                            {[1, 2, 3, 4].map(i => (
                                                                <div key={i} style={{
                                                                    height: '40px',
                                                                    background: option.value === 'light' ? '#ffffff' : '#1e1e1e',
                                                                    border: `1px solid ${option.value === 'light' ? '#d0d0d0' : '#333'}`,
                                                                    borderRadius: '4px',
                                                                    padding: '6px',
                                                                    display: 'flex',
                                                                    flexDirection: 'column',
                                                                    gap: '3px'
                                                                }}>
                                                                    <div style={{
                                                                        height: '4px',
                                                                        background: i === 1 ? '#4CAF50' : i === 2 ? '#2196F3' : i === 3 ? '#FF9800' : '#E91E63',
                                                                        borderRadius: '2px',
                                                                        width: '80%'
                                                                    }}></div>
                                                                    <div style={{
                                                                        height: '3px',
                                                                        background: option.value === 'light' ? '#e0e0e0' : '#333',
                                                                        borderRadius: '2px',
                                                                        width: '60%'
                                                                    }}></div>
                                                                </div>
                                                            ))}
                                                        </div>
                                                    </div>
                                                    
                                                    {/* System setting indicator */}
                                                    {option.value === 'system' && (
                                                        <div style={{
                                                            position: 'absolute',
                                                            top: '8px',
                                                            right: '8px',
                                                            background: 'var(--color-background)',
                                                            borderRadius: '50%',
                                                            width: '28px',
                                                            height: '28px',
                                                            display: 'flex',
                                                            alignItems: 'center',
                                                            justifyContent: 'center',
                                                            fontSize: '16px',
                                                            boxShadow: '0 2px 4px rgba(0,0,0,0.2)'
                                                        }}>
                                                            ℹ️
                                                        </div>
                                                    )}
                                                </div>
                                                
                                                {/* Label area */}
                                                <div style={{
                                                    padding: 'var(--space-sm) var(--space-md)',
                                                    background: theme === option.value ? 'var(--color-background-light)' : 'var(--color-background)',
                                                    borderTop: `1px solid ${theme === option.value ? 'var(--color-primary)' : 'var(--color-border)'}`
                                                }}>
                                                    <div style={{
                                                        display: 'flex',
                                                        alignItems: 'center',
                                                        justifyContent: 'space-between',
                                                        marginBottom: '4px'
                                                    }}>
                                                        <div style={{ 
                                                            fontWeight: 600,
                                                            fontSize: '1rem',
                                                            color: 'var(--color-text)'
                                                        }}>
                                                            {option.label}
                                                            {option.value === 'light' && theme === 'light' && ' (Default)'}
                                                        </div>
                                                        {theme === option.value && (
                                                            <div style={{
                                                                width: '18px',
                                                                height: '18px',
                                                                borderRadius: '50%',
                                                                background: 'var(--color-primary)',
                                                                display: 'flex',
                                                                alignItems: 'center',
                                                                justifyContent: 'center',
                                                                color: 'white',
                                                                fontSize: '12px',
                                                                fontWeight: 'bold'
                                                            }}>
                                                                ✓
                                                            </div>
                                                        )}
                                                    </div>
                                                    <div style={{ 
                                                        fontSize: '0.85rem',
                                                        color: 'var(--color-text-secondary)',
                                                        lineHeight: '1.3'
                                                    }}>
                                                        {option.description}
                                                    </div>
                                                </div>
                                            </div>
                                        ))}
                                    </div>
                                    
                                    <div style={{ 
                                        marginTop: 'var(--space-lg)',
                                        padding: 'var(--space-md)',
                                        backgroundColor: 'var(--color-background-medium)',
                                        borderRadius: 'var(--border-radius-md)',
                                        fontSize: '0.9rem',
                                        borderLeft: `4px solid var(--color-primary)`
                                    }}>
                                        <strong>Currently active:</strong> {effectiveTheme === 'light' ? '☀️ Light theme' : '🌙 Dark theme'}
                                        {theme === 'system' && ' (automatically set by your device)'}
                                    </div>
                                </FormGroup>
                            </div>
                        </div>
            </ProfileContentLayout>
        </Layout>
    );
};

export default ProfilePreferences;
